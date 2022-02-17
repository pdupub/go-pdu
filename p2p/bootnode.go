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
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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
	universe            *core.Universe
	logger              echo.Logger
}

func NewBootNode(db *bbolt.DB, logger echo.Logger) *BootNode {
	return &BootNode{
		Uptime:              time.Now(),
		Statuses:            map[string]int{},
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

// LoadUniverse load msg from db
func (bn *BootNode) LoadUniverse() error {
	if bn.universe == nil {
		bn.logger.Error(ErrUniverseNotExist)
		return ErrUniverseNotExist
	}
	return nil
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
	return nil
}

func (bn *BootNode) receiveMsg(m *msg.SignedMsg) error {
	_, err := m.Ecrecover()
	if err != nil {
		return err
	}
	return nil
}

func (bn *BootNode) upload(c echo.Context) error {
	fileHash := c.FormValue("hash")
	author := c.FormValue("author")
	signature := c.FormValue("signature")

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
	return nil, nil, nil

}

func (bn *BootNode) parseSM(sm *msg.SignedMsg) (map[string]interface{}, error) {
	return nil, nil
}

func (bn *BootNode) getProfile(c echo.Context) error {
	// author := common.HexToAddress(c.Param("uid"))

	// if c.QueryParam("view") != "" {
	// 	detail := map[string]interface{}{
	// 		"cnt":     0,
	// 		"exist":   exist,
	// 		"profile": profile,
	// 	}
	// 	return c.Render(http.StatusOK, "detail.html", []interface{}{detail})
	// }
	// return c.JSON(http.StatusOK, profile)
	return nil
}

func (bn *BootNode) getDetail(c echo.Context) error {
	// author := common.HexToAddress(c.Param("uid"))

	// if c.QueryParam("view") != "" {
	// 	return c.Render(http.StatusOK, "detail.html", []interface{}{detail})
	// }
	// return c.JSON(http.StatusOK, detail)
	return nil
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
	// return c.JSON(http.StatusOK, societyData)
	return nil
}

// getEntropy is sample to display the topo of Event relation.
func (bn *BootNode) getEntropy(c echo.Context) error {
	// return c.JSON(http.StatusOK, entropyData)
	return nil
}

func (bn *BootNode) getQuantums(c echo.Context) error {
	// return c.JSON(http.StatusOK, sms)
	return nil
}

func (bn *BootNode) parseSig(sig string) ([]byte, error) {
	if len(sig) == 130 {
		return common.Hex2Bytes(sig), nil
	}
	return base64.StdEncoding.DecodeString(sig)
}

func (bn *BootNode) getFullMsg(c echo.Context) error {
	return nil
}

func (bn *BootNode) getMessage(c echo.Context) error {
	return nil
}

func (bn *BootNode) getLatest(c echo.Context) error {
	return nil
}
