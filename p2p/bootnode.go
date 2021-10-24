// Copyright 2021 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/pdupub/go-dag"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/msg"
	"go.etcd.io/bbolt"
)

const BucketMsg = "message"
const UploadFileSizeLimit = 1024 * 1024
const PrefixMessage = "msg"    // msg-signature => message
const PrefixIndividual = "ind" // ind-address-pcnt => signature
const PrefixLast = "lst"       // lst-address => pcnt
const PrefixCount = "cnt"      // cnt-signature-address-pcnt => scnt
const PrefixSystem = "sys"     // sys-scnt-address-pcnt => signature

type UploadResponse struct {
	Path    string `json:"path"`
	Hash    string `json:"hash"`
	Expired int64  `json:"expired"`
}

type BootNode struct {
	Uptime              time.Time      `json:"uptime"`
	RequestCount        uint64         `json:"requestCount"`
	Statuses            map[string]int `json:"statuses"`
	NextMsgID           *big.Int       `json:"next"`
	Peers               []string       `json:"peers"`
	IgnoreUnknownSource bool           `json:"ignoreUnknownSource"`
	mutex               sync.RWMutex
	db                  *bbolt.DB
	universe            *core.Universe
	logger              echo.Logger
}

func NewBootNode(db *bbolt.DB, logger echo.Logger) *BootNode {
	return &BootNode{
		Uptime:              time.Now(),
		Statuses:            map[string]int{},
		db:                  db,
		NextMsgID:           big.NewInt(1),
		logger:              logger,
		IgnoreUnknownSource: true,
	}
}

// SetIgnoreUnknownSource set if or not ignore unknow source
func (bn *BootNode) SetIgnoreUnknownSource(state bool) {
	bn.IgnoreUnknownSource = state
}

// SetUniverse init the universe
func (bn *BootNode) SetUniverse(universe *core.Universe) error {
	bn.universe = universe
	return nil
}

func (bn *BootNode) AddPeers(peers []string) error {
	bn.Peers = append(bn.Peers, peers...)
	return nil
}

func joinKey(keys ...interface{}) (resKey []byte) {
	if len(keys) == 0 {
		return
	}
	for _, key := range keys {
		switch key := key.(type) {
		case string:
			resKey = append(resKey, []byte(key)...)
		case []byte:
			resKey = append(resKey, key...)
		case common.Address:
			resKey = append(resKey, key.Bytes()...)
		case *big.Int:
			resKey = append(resKey, key.Bytes()...)
		default:
			break
		}
		resKey = append(resKey, []byte("-")...)
	}
	resKey = resKey[:len(resKey)-1]
	return resKey
}

// UpsertDBMsgs insert msgs into db if not exist on db
func (bn *BootNode) UpsertDBMsgs(msgs []*msg.SignedMsg) error {
	err := bn.db.Batch(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(BucketMsg))
		if b == nil {
			b, err = tx.CreateBucket([]byte(BucketMsg))
			if err != nil {
				return err
			}
		}
		for _, m := range msgs {
			mBytes, err := json.Marshal(m)
			if err != nil {
				return err
			}
			// msg-signature => message
			msgKey := joinKey(PrefixMessage, m.Signature)
			mOrigin := b.Get(msgKey)

			personalCnt := big.NewInt(0)
			author, err := m.Ecrecover()
			if err != nil {
				return err
			}

			if mOrigin != nil {
				// already save, read personalCnt
				lstKey := joinKey(PrefixLast, author)
				if pcnt := b.Get(lstKey); pcnt == nil {
					return ErrMessageNotRecord
				} else {
					lastCnt := new(big.Int).SetBytes(pcnt)
					last := joinKey(PrefixIndividual, author, lastCnt)
					cursor := b.Cursor()
					prefix := joinKey(PrefixIndividual, author)
					for k, sig := cursor.Seek(last); k != nil && bytes.HasPrefix(k, prefix); k, sig = cursor.Prev() {
						if common.Bytes2Hex(sig) == common.Bytes2Hex(m.Signature) {
							personalCnt.SetBytes(k[len(prefix)+1:])
						}
					}
					if personalCnt.Cmp(big.NewInt(0)) <= 0 {
						return ErrMessageNotRecord
					}
				}

			} else {
				// save msg for personal, 3 records
				// msg-signature => message
				if err := b.Put(msgKey, mBytes); err != nil {
					return err
				}
				// lst-address => pcnt
				lstKey := joinKey(PrefixLast, author)
				if pcnt := b.Get(lstKey); pcnt != nil {
					personalCnt = new(big.Int).SetBytes(pcnt)
					personalCnt.Add(personalCnt, big.NewInt(1))
				}
				if err := b.Put(lstKey, personalCnt.Bytes()); err != nil {
					return err
				}

				// ind-address-pcnt => signature
				indKey := joinKey(PrefixIndividual, author, personalCnt)
				if err := b.Put(indKey, m.Signature); err != nil {
					return err
				}
			}

			cntKey := joinKey(PrefixCount, m.Signature, author, personalCnt)
			// check if the author in society before save scnt, sys and update NextMsgID
			if _, err := bn.universe.GetSociety().GetIndividual(author); err == nil && b.Get(cntKey) == nil {
				// cnt-signature-address-pcnt => scnt
				if err := b.Put(cntKey, bn.NextMsgID.Bytes()); err != nil {
					return err
				}

				// sys-scnt-address-pcnt => signature
				sysKey := joinKey(PrefixSystem, bn.NextMsgID, author, personalCnt)
				if err := b.Put(sysKey, m.Signature); err != nil {
					return err
				}

				bn.mutex.Lock()
				bn.NextMsgID.Add(bn.NextMsgID, big.NewInt(1))
				bn.mutex.Unlock()
			}

		}

		return nil
	})
	if err != nil {
		bn.logger.Error(err)
	}
	return err
}

// LoadUniverse load msg from db
func (bn *BootNode) LoadUniverse() error {
	if bn.universe == nil {
		bn.logger.Error(ErrUniverseNotExist)
		return ErrUniverseNotExist
	}

	// load msg by order in db and receive
	return bn.db.View(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(BucketMsg))
		if b != nil {
			cursor := b.Cursor()
			prefix := joinKey(PrefixSystem)
			for k, sig := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, sig = cursor.Next() {
				curMsg := b.Get(joinKey(PrefixMessage, sig))
				if curMsg == nil {
					break
				}
				bn.NextMsgID.Add(bn.NextMsgID, big.NewInt(1))
				m := new(msg.SignedMsg)
				if err := json.Unmarshal(curMsg, m); err != nil {
					bn.logger.Error(err)
					continue
				}

				author, err := m.Ecrecover()
				if err != nil {
					bn.logger.Error(err)
					continue
				}

				if _, err := bn.universe.ReceiveMsg(author, m.Signature, m.Content, m.References...); err != nil && err != core.ErrPhotonAlreadyExist {
					bn.logger.Error(err)
					continue
				}
			}
		}
		return nil
	})
}

// Process is the middleware function used to count request.
func (bn *BootNode) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		bn.mutex.Lock()
		defer bn.mutex.Unlock()
		bn.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		bn.Statuses[status]++
		return nil
	}
}

// welcome is the endpoint to get bootnode information.
func (bn *BootNode) welcome(c echo.Context) error {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()
	return c.JSON(http.StatusOK, bn)
}

// newMsg
func (bn *BootNode) newMsg(c echo.Context) error {
	m := new(msg.SignedMsg)
	if err := c.Bind(m); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := bn.receiveMsg(m)
	if err != nil {
		bn.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bn.broadcastMsg(m)
	return c.JSON(http.StatusCreated, m)
}

func (bn *BootNode) mergeDBMsgsByAuthor(author common.Address) error {

	return bn.db.Batch(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(BucketMsg))
		if b == nil {
			return nil
		}

		// lst-address => pcnt
		lstKey := joinKey(PrefixLast, author)
		if pcnt := b.Get(lstKey); pcnt == nil {
			// no msgs
			return nil
		} else {
			personalCnt := new(big.Int).SetBytes(pcnt)
			for i := int64(1); i < personalCnt.Int64(); i++ {
				// ind-address-pcnt => signature
				indKey := joinKey(PrefixIndividual, author, big.NewInt(i))
				sig := b.Get(indKey)
				if sig == nil {
					break
				}
				// msg-signature => message
				msgKey := joinKey(PrefixMessage, sig)
				msgBytes := b.Get(msgKey)
				if msgBytes == nil {
					break
				}

				m := new(msg.SignedMsg)
				if err := json.Unmarshal(msgBytes, m); err != nil {
					return err
				}

				if err := bn.receiveMsg(m); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (bn *BootNode) receiveMsg(m *msg.SignedMsg) error {
	author, err := m.Ecrecover()
	if err != nil {
		return err
	}

	if photon, err := bn.universe.ReceiveMsg(author, m.Signature, m.Content, m.References...); err == nil {
		if err := bn.UpsertDBMsgs([]*msg.SignedMsg{m}); err != nil {
			return err
		}
		if photon.Type == core.PhotonTypeBorn && !bn.IgnoreUnknownSource {
			if newbe, err := photon.GetNewBorn(); err != nil {
				bn.logger.Error(err)
			} else {
				if err := bn.mergeDBMsgsByAuthor(newbe); err != nil {
					bn.logger.Error(err)
				}
			}
		}
	} else if err == core.ErrIndividualNotExistInSociety && !bn.IgnoreUnknownSource {
		if err := bn.UpsertDBMsgs([]*msg.SignedMsg{m}); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (bn *BootNode) upload(c echo.Context) error {
	fileHash := c.FormValue("hash")
	author := c.FormValue("author")
	signature := c.FormValue("signature")

	if _, err := bn.universe.GetSociety().GetIndividual(common.HexToAddress(author)); err != nil {
		bn.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	file, err := c.FormFile("file")
	if err != nil {
		bn.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if file.Size > UploadFileSizeLimit {
		bn.logger.Error(ErrFileSizeBeyondLimit, file.Size)
		return echo.NewHTTPError(http.StatusBadRequest, ErrFileSizeBeyondLimit.Error())
	}
	if err := bn.checkHash(file, fileHash); err != nil {
		bn.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := bn.verifyFileSignature(fileHash, signature, author); err != nil {
		bn.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// Destination
	filename := strings.Join([]string{fileHash, file.Filename}, "-")
	if err := bn.storeFile(file, filename); err != nil {
		bn.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	resp := UploadResponse{Path: path.Join("files", filename), Hash: fileHash, Expired: 3600 * 24}
	return c.JSON(http.StatusOK, resp)
}

func (bn *BootNode) verifyFileSignature(fileHash, signature, author string) error {
	m := msg.New([]byte(fileHash))
	sigBytes := common.Hex2Bytes(signature)
	return m.Verify(sigBytes, common.HexToAddress(author))
}

func (bn *BootNode) storeFile(file *multipart.FileHeader, filename string) error {
	dst, err := os.Create(path.Join("tmp/files", filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	src, err := file.Open()
	if err != nil {
		return err
	}
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func (bn *BootNode) checkHash(file *multipart.FileHeader, fileHash string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, src); err != nil {
		return err
	}
	if hex.EncodeToString(hash.Sum(nil)) != fileHash {
		return errors.New("hash not match")
	}
	return nil
}

func (bn *BootNode) broadcastMsg(m *msg.SignedMsg) {
	// broadcast new msg
	for _, node := range bn.Peers {
		_, err := m.Post(node)
		if err != nil {
			bn.logger.Error(err)
		}
	}
}

func (bn *BootNode) getSM(sig []byte) (*msg.SignedMsg, *big.Int, error) {
	sm := new(msg.SignedMsg)
	scnt := big.NewInt(0)
	// load msg by order in db and receive
	err := bn.db.View(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(BucketMsg))
		if b != nil {
			// msg-signature => message
			msgKey := joinKey(PrefixMessage, sig)
			curMsg := b.Get(msgKey)
			if curMsg == nil {
				return ErrMessageNotRecord
			}
			if err := json.Unmarshal(curMsg, sm); err != nil {
				return err
			}

			// seek to find scnt
			cursor := b.Cursor()
			prefix := joinKey(PrefixCount, sig)
			for k, v := cursor.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = cursor.Next() {
				scnt.SetBytes(v)
				break
			}
		}
		return nil
	},
	)

	if err != nil {
		return nil, nil, err
	}
	return sm, scnt, nil
}

func (bn *BootNode) parseSM(sm *msg.SignedMsg) (map[string]interface{}, error) {

	photon := new(core.Photon)
	if err := json.Unmarshal(sm.Content, photon); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["type"] = photon.Type
	result["version"] = photon.Version

	data := make(map[string]interface{})
	if photon.Type == core.PhotonTypeInfo {
		pInfo := new(core.PInfo)
		if err := json.Unmarshal(photon.Data, pInfo); err != nil {
			// TODO:
			data["text"] = string(photon.Data)
			result["data"] = data
			return result, nil
		}
		data["text"] = pInfo.Text
		data["quote"] = common.Bytes2Hex(pInfo.Quote)

		ress := []interface{}{}
		for _, resource := range pInfo.Resources {
			res := make(map[string]interface{})
			res["checksum"] = hex.EncodeToString(resource.Checksum)
			res["data"] = resource.Data
			res["format"] = resource.Format
			res["url"] = resource.URL
			ress = append(ress, res)
		}
		data["resources"] = ress
	} else if photon.Type == core.PhotonTypeBorn {
		pBorn := new(core.PBorn)
		if err := json.Unmarshal(photon.Data, pBorn); err != nil {
			return result, nil
		}
		data["address"] = pBorn.Addr.Hex()
		sigs := []string{}
		for _, sig := range pBorn.Signatures {
			sigs = append(sigs, common.Bytes2Hex(sig))
		}
		data["sigs"] = sigs

	} else if photon.Type == core.PhotonTypeProfile {
		pProfile := new(core.PProfile)
		if err := json.Unmarshal(photon.Data, pProfile); err != nil {
			return result, nil
		}
		data["name"] = pProfile.Name
		data["email"] = pProfile.Email
		data["bio"] = pProfile.Bio
		data["url"] = pProfile.URL
		data["location"] = pProfile.Location
		data["avatar"] = pProfile.Avatar
		data["extra"] = pProfile.Extra
	}

	result["data"] = data

	return result, nil
}

func (bn *BootNode) getProfile(c echo.Context) error {
	author := common.HexToAddress(c.Param("uid"))
	profile := bn.universe.GetSociety().GetIndividualProfile(author)
	exist := true
	_, err := bn.universe.GetSociety().GetIndividual(author)
	if err != nil {
		exist = false
	}
	if c.QueryParam("view") != "" {
		detail := map[string]interface{}{
			"cnt":     0,
			"exist":   exist,
			"profile": profile,
		}
		return c.Render(http.StatusOK, "detail.html", []interface{}{detail})
	}
	return c.JSON(http.StatusOK, profile)
}

func (bn *BootNode) getDetail(c echo.Context) error {
	author := common.HexToAddress(c.Param("uid"))
	exist := true
	_, err := bn.universe.GetSociety().GetIndividual(author)
	if err != nil {
		exist = false
	}
	profile := bn.universe.GetSociety().GetIndividualProfile(author)
	personalCnt := big.NewInt(0)
	bn.db.View(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(BucketMsg))
		if b != nil {
			lstKey := joinKey(PrefixLast, author)
			pcnt := b.Get(lstKey)
			if pcnt != nil {
				personalCnt.SetBytes(pcnt)
			}
		}
		return nil
	})
	detail := map[string]interface{}{
		"cnt":     personalCnt,
		"exist":   exist,
		"profile": profile,
	}

	if c.QueryParam("view") != "" {
		return c.Render(http.StatusOK, "detail.html", []interface{}{detail})
	}
	return c.JSON(http.StatusOK, detail)
}

func (bn *BootNode) readDAGParam(c echo.Context) (keys []string, parentLimit, childLimit int) {
	parentLimit, childLimit = dag.DefaultStepLimit, dag.DefaultStepLimit

	if c.QueryParam("key") != "" {
		keys = append(keys, common.HexToAddress(c.QueryParam("key")).Hex())
	}

	if n, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		parentLimit = n
		childLimit = n
	}
	if n, err := strconv.Atoi(c.QueryParam("plimit")); err == nil {
		parentLimit = n
	}
	if n, err := strconv.Atoi(c.QueryParam("climit")); err == nil {
		childLimit = n
	}
	return
}

// getSociety is sample to display the topo of ID relation.
func (bn *BootNode) getSociety(c echo.Context) error {
	keys, parentLimit, childLimit := bn.readDAGParam(c)
	societyData, err := bn.universe.GetSociety().Dump(keys, parentLimit, childLimit)
	if err != nil {
		bn.logger.Error(err)
		return err
	}

	if c.QueryParam("view") != "" {
		return c.Render(http.StatusOK, "topology.html", map[string]interface{}{
			"data": societyData,
			"name": "Society",
		})
	}
	return c.JSON(http.StatusOK, societyData)

}

// getEntropy is sample to display the topo of Event relation.
func (bn *BootNode) getEntropy(c echo.Context) error {
	keys, parentLimit, childLimit := bn.readDAGParam(c)
	entropyData, err := bn.universe.GetEntropy().Dump(keys, parentLimit, childLimit)
	if err != nil {
		bn.logger.Error(err)
		return err
	}
	if c.QueryParam("view") != "" {
		return c.Render(http.StatusOK, "topology.html", map[string]interface{}{
			"data": entropyData,
			"name": "Entropy",
		})
	}
	return c.JSON(http.StatusOK, entropyData)
}

func (bn *BootNode) getPhotons(c echo.Context) error {
	var sms []*msg.SignedMsg // message slice

	start, limit := 1, 10 // default query limit
	if n, err := strconv.Atoi(c.QueryParam("start")); err == nil {
		start = n
	}
	if n, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
		limit = n
	}

	desc := c.QueryParam("desc") != "" // default query oreder is by order inc

	prefix := joinKey(PrefixSystem)
	if c.QueryParam("author") != "" {
		author := common.HexToAddress(c.QueryParam("author"))
		prefix = joinKey(PrefixIndividual, author)
	}
	first := joinKey(prefix, big.NewInt(int64(start)))

	err := bn.db.View(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(BucketMsg))
		if b != nil {
			cursor := b.Cursor()
			// seek direct
			cursorNext := cursor.Next
			if desc {
				cursorNext = cursor.Prev
			}
			// seek
			for k, sig := cursor.Seek(first); k != nil && bytes.HasPrefix(k, prefix); k, sig = cursorNext() {
				sm := new(msg.SignedMsg)
				msgKey := joinKey(PrefixMessage, sig)
				v := b.Get(msgKey)
				if v == nil {
					break
				}
				if err := json.Unmarshal(v, sm); err != nil {
					return err
				}
				sms = append(sms, sm)
				if len(sms) >= limit {
					break
				}
			}
		}
		return nil
	},
	)
	if err != nil {
		bn.logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if c.QueryParam("view") != "" {
		var photons []interface{}
		for _, sm := range sms {
			photon, err := bn.parseSM(sm)
			if err != nil {
				bn.logger.Error(err)
			}
			photons = append(photons, photon)
		}
		return c.Render(http.StatusOK, "photons.html", photons)
	}
	return c.JSON(http.StatusOK, sms)
}

func (bn *BootNode) parseSig(sig string) ([]byte, error) {
	if len(sig) == 130 {
		return common.Hex2Bytes(sig), nil
	}
	return base64.StdEncoding.DecodeString(sig)
}

func (bn *BootNode) getFullMsg(c echo.Context) error {
	msgID, err := bn.parseSig(c.Param("sig"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	sm, cnt, err := bn.getSM(msgID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	pcnt := big.NewInt(0)
	err = bn.db.View(func(tx *bbolt.Tx) (err error) {
		b := tx.Bucket([]byte(BucketMsg))
		if b != nil {
			author, _ := sm.Ecrecover()
			lstKey := joinKey(PrefixLast, author)
			v := b.Get(lstKey)
			if v != nil {
				pcnt.SetBytes(v)
			}
		}
		return nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if c.QueryParam("view") != "" {
		photonM, err := bn.parseSM(sm)
		if err != nil {
			bn.logger.Error(err)
		}
		photons := []interface{}{photonM}
		return c.Render(http.StatusOK, "photons.html", photons)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": sm, "cnt": cnt, "pcnt": pcnt})
}

func (bn *BootNode) getMessage(c echo.Context) error {
	msgID, err := bn.parseSig(c.Param("sig"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sm, _, err := bn.getSM(msgID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if c.QueryParam("view") != "" {
		photonM, err := bn.parseSM(sm)
		if err != nil {
			bn.logger.Error(err)
		}
		photons := []interface{}{photonM}
		return c.Render(http.StatusOK, "photons.html", photons)
	}

	return c.JSON(http.StatusOK, sm)
}

func (bn *BootNode) getLatest(c echo.Context) error {
	latestID := bn.universe.GetEntropy().GetLastEventID(common.HexToAddress(c.Param("uid")))
	c.SetParamNames("sig")
	c.SetParamValues(common.Bytes2Hex(latestID))
	return bn.getMessage(c)
}
