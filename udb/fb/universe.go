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

package fb

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	errDocumentLoadDataFail   = errors.New("document load data fail")
	errReceivedQuantumMissing = errors.New("received quantum missing")
	errUnmarshalQuantumFail   = errors.New("unmarshal quantum fail")
	errQuantumIsReject        = errors.New("quantum is reject")
	errQuantumIsWaiting       = errors.New("quantum is waiting")
	errIndividualLevelUnknown = errors.New("individual level is unknown")
)

var (
	firstRef = core.Sig2Hex(core.FirstQuantumReference)
)

type FBSig struct {
	SigHex string `json:"refs"`
}

type UniverseStatus struct {
	Sequence        int64  `json:"lastSequence,omitempty"`
	SigHex          string `json:"lastSigHex,omitempty"`
	UpdateTimestamp int64  `json:"updateTime"`
}

type FBUniverse struct {
	ctx         context.Context
	app         *firebase.App
	client      *firestore.Client
	status      *UniverseStatus
	quantumC    *firestore.CollectionRef
	communityC  *firestore.CollectionRef
	individualC *firestore.CollectionRef
	universeC   *firestore.CollectionRef
}

type PlatformConfig struct {
	Platform string `json:"platform"`
	Version  int    `json:"version"`
	Action   string `json:"action"`
	Params   string `json:"params,omitempty"`
}

const (
	universeStatusDocID   = "status"
	platformIOS           = "ios"
	platformVersion       = 1
	platformActionPost    = "post"
	platformActionReply   = "reply"
	platformActionComment = "comment"
)

func NewFBUniverse(ctx context.Context, keyFilename string, projectID string) (*FBUniverse, error) {
	fbu := &FBUniverse{ctx: ctx, status: &UniverseStatus{}}
	opt := option.WithCredentialsFile(keyFilename)
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, err
	}
	fbu.app = app

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	fbu.client = client
	fbu.quantumC = fbu.client.Collection("quantum")
	fbu.communityC = fbu.client.Collection("community")
	fbu.individualC = fbu.client.Collection("individual")
	fbu.universeC = fbu.client.Collection("universe")

	// init config
	if err = fbu.loadUniverse(); err != nil {
		return nil, err
	}

	return fbu, nil
}

func (fbu *FBUniverse) Close() error {
	return fbu.client.Close()
}

func (fbu *FBUniverse) loadUniverse() error {
	docRef := fbu.universeC.Doc(universeStatusDocID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return err
	}
	dMap := docSnapshot.Data()
	if sequence, ok := dMap["lastSequence"]; ok {
		fbu.status.Sequence = sequence.(int64)
	}
	return nil
}

func (fbu *FBUniverse) increaseUniverseSequence(newSigHex string) error {
	docRef := fbu.universeC.Doc(universeStatusDocID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return err
	}
	dMap := docSnapshot.Data()
	if sequence, ok := dMap["lastSequence"]; ok {
		fbu.status.Sequence = sequence.(int64)
		fbu.status.Sequence += 1
	} else {
		fbu.status.Sequence = 1
	}
	dMap["lastSequence"] = fbu.status.Sequence
	dMap["lastSigHex"] = newSigHex
	dMap["updateTime"] = time.Now().UnixMilli()
	docRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"lastSequence"}, []string{"lastSigHex"}, []string{"updateTime"}))

	return nil
}

func (fbu *FBUniverse) loadUnprocessedQuantum(sig core.Sig) (*core.Quantum, error) {
	docSnapshot, err := fbu.quantumC.Doc(core.Sig2Hex(sig)).Get(fbu.ctx)
	if err != nil {
		return nil, err
	}
	if docSnapshot == nil {
		return nil, errDocumentLoadDataFail
	}

	qRes := &core.Quantum{}
	if qBytes, ok := docSnapshot.Data()["recv"]; !ok || qBytes == nil {
		return nil, errReceivedQuantumMissing
	} else {
		if err := json.Unmarshal(qBytes.([]byte), qRes); err != nil {
			docSnapshot.Ref.Delete(fbu.ctx)
			return nil, errUnmarshalQuantumFail
		}
	}
	return qRes, nil
}

func (fbu *FBUniverse) loadUnprocessedQuantums(limit, skip int) ([]*core.Quantum, error) {
	var receivedQuantums []*core.Quantum
	iter := fbu.quantumC.Where("recv", "!=", []byte{}).Offset(skip).Limit(limit).Documents(fbu.ctx)
	for docSnapshot, err := iter.Next(); err != iterator.Done; docSnapshot, err = iter.Next() {
		if docSnapshot == nil {
			return nil, errDocumentLoadDataFail
		}

		if qBytes, ok := docSnapshot.Data()["recv"]; ok {
			qRes := &core.Quantum{}
			if err := json.Unmarshal(qBytes.([]byte), qRes); err != nil {
				// delete current quantums, (or set recv to "")
				docSnapshot.Ref.Delete(fbu.ctx)
			} else {
				receivedQuantums = append(receivedQuantums, qRes)
			}
		}
	}
	return receivedQuantums, nil
}

func (fbu *FBUniverse) preprocessQuantums(quantums []*core.Quantum) (signatureQuantumMap map[string]*core.Quantum, referenceSignatureMap map[string]string, addressStatusMap map[string]int) {
	// signatureQuantumMap is used for find quantum by signature
	signatureQuantumMap = make(map[string]*core.Quantum) // sig:quantum
	// referenceSignatureMap is used for from individual last find next quantum by same individual
	referenceSignatureMap = make(map[string]string) // first_ref_sig:sig
	// addressStatusMap is used for get attitude for each of quantum signer
	addressStatusMap = make(map[string]int) // address:core.Attitude...

	// fill data struct
	for _, qRes := range quantums {
		// ecrecover the author address
		addr, err := qRes.Ecrecover()
		if err != nil {
			continue
		}
		// update individual attitude
		if _, ok := addressStatusMap[addr.Hex()]; !ok {
			addressStatusMap[addr.Hex()] = fbu.getStatusLevelByAddressHex(addr.Hex())
		}

		sigHex := core.Sig2Hex(qRes.Signature)
		selfRef := core.Sig2Hex(qRes.References[0])
		signatureQuantumMap[sigHex] = qRes
		if selfRef != firstRef {
			referenceSignatureMap[selfRef] = sigHex
		}
	}
	return
}

func (fbu *FBUniverse) getStatusLevelByAddressHex(addrHex string) int {
	statusLevel := core.AttitudeAccept // default accept
	iDocRef := fbu.individualC.Doc(addrHex)
	snapshot, err := iDocRef.Get(fbu.ctx)
	if err == nil && snapshot.Exists() {
		attitude, err := snapshot.DataAt("attitude")
		if err == nil {
			at := attitude.(map[string]interface{})
			statusLevel = int(at["level"].(float64))
		}
	}
	return statusLevel
}

func (fbu *FBUniverse) ProcessQuantums(limit, skip int) (accept []core.Sig, wait []core.Sig, reject []core.Sig, err error) {
	// load unprocessed quantums
	unprocessedQuantums, err := fbu.loadUnprocessedQuantums(limit, skip)
	if err != nil || len(unprocessedQuantums) == 0 {
		return
	}
	return fbu.proccessQuantums(unprocessedQuantums)
}

func (fbu *FBUniverse) proccessQuantums(unprocessedQuantums []*core.Quantum) (accept []core.Sig, wait []core.Sig, reject []core.Sig, err error) {
	// format quantums
	signatureQuantumMap, referenceSignatureMap, addressStatusMap := fbu.preprocessQuantums(unprocessedQuantums)

	// process first quantums
	for sigHex, quantum := range signatureQuantumMap {
		// check individual
		qDocRef := fbu.quantumC.Doc(sigHex)

		addr, _ := quantum.Ecrecover()
		// update address info for quantum
		dMap, _ := FBStruct2Data(&FBQuantum{AddrHex: addr.Hex()})
		dMap["createTime"] = time.Now().UnixMilli()

		// deal by individual status
		if addressStatusMap[addr.Hex()] < core.AttitudeIgnoreContent {
			// clear recv
			recvMoveTo := "trash"
			docSnapshot, _ := qDocRef.Get(fbu.ctx)
			if recv, ok := docSnapshot.Data()["recv"]; ok {
				dMap := map[string]interface{}{"recv": []byte{}, recvMoveTo: recv}
				mergeKeys := []firestore.FieldPath{[]string{"recv"}, []string{recvMoveTo}}
				qDocRef.Set(fbu.ctx, dMap, firestore.Merge(mergeKeys...))
			}
			reject = append(reject, quantum.Signature)
			continue
		}
		// quantum created when verfy signature & got user address, not receive quantum
		qDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"address"}, []string{"createTime"}))

		iDocRef := fbu.individualC.Doc(addr.Hex())
		iDocSnapshot, _ := iDocRef.Get(fbu.ctx)
		// TODO: filter unkown source quantums
		if !iDocSnapshot.Exists() && core.Sig2Hex(quantum.References[0]) == core.Sig2Hex(core.FirstQuantumReference) {
			// checked first quantums, can be accepted.
			if err = fbu.increaseUniverseSequence(sigHex); err != nil {
				// reject = append(reject, core.Hex2Sig(sigHex))
				// continue
				return
			}
			// add sequence to quantum
			dMap, _ := FBStruct2Data(&FBQuantum{Sequence: fbu.status.Sequence, SelfSeq: int64(1)})
			qDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"seq"}, []string{"sseq"}))

			// add new individual
			newIndividual := &FBIndividual{AddrHex: addr.Hex(), LastSigHex: sigHex, LastSelfSeq: int64(1), Attitude: &core.Attitude{Level: core.AttitudeAccept}}
			dMap, _ = FBStruct2Data(newIndividual)
			dMap["createTime"] = time.Now().UnixMilli()
			dMap["updateTime"] = time.Now().UnixMilli()
			iDocRef.Set(fbu.ctx, dMap)

			accept = append(accept, core.Hex2Sig(sigHex))
		}
	}

	// process quantums
	for addrHex := range addressStatusMap {
		if addressStatusMap[addrHex] < core.AttitudeIgnoreContent {
			continue
		}

		iDocRef := fbu.individualC.Doc(addrHex)
		iDocSnapshot, _ := iDocRef.Get(fbu.ctx)
		if iDocSnapshot.Exists() {
			var individual *FBIndividual
			individual, err = Data2FBIndividual(iDocSnapshot.Data())
			if err != nil {
				continue
			}

			for {
				if sigHex, ok := referenceSignatureMap[individual.LastSigHex]; ok {
					// accept the quantum
					if _, ok := signatureQuantumMap[sigHex]; ok {

						// checked first quantums, can be accepted.
						if err = fbu.increaseUniverseSequence(sigHex); err != nil {
							// reject = append(reject, core.Hex2Sig(sigHex))
							// continue
							return
						}

						individual.LastSigHex = sigHex
						individual.LastSelfSeq += 1

						qDocRef := fbu.quantumC.Doc(sigHex)
						// add sequence to quantum
						dMap, _ := FBStruct2Data(&FBQuantum{Sequence: fbu.status.Sequence, SelfSeq: individual.LastSelfSeq})
						qDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"seq"}, []string{"sseq"}))

						// add new individual
						dMap, _ = FBStruct2Data(individual)
						dMap["createTime"] = time.Now().UnixMilli()
						dMap["updateTime"] = time.Now().UnixMilli()
						iDocRef.Set(fbu.ctx, dMap)

						if addressStatusMap[addrHex] == core.AttitudeIgnoreContent {
							// clear recv
							recvMoveTo := "ignore"
							docSnapshot, _ := qDocRef.Get(fbu.ctx)
							if recv, ok := docSnapshot.Data()["recv"]; ok {
								dMap := map[string]interface{}{"recv": []byte{}, recvMoveTo: recv}
								mergeKeys := []firestore.FieldPath{[]string{"recv"}, []string{recvMoveTo}}
								qDocRef.Set(fbu.ctx, dMap, firestore.Merge(mergeKeys...))
							}
							// not realy reject, just not executeQuantum in next step
							reject = append(reject, core.Hex2Sig(sigHex))
						} else if addressStatusMap[addrHex] > core.AttitudeIgnoreContent {
							accept = append(accept, core.Hex2Sig(sigHex))
						}

					}
				} else {
					break
				}
			}
		}
	}

	for _, sig := range accept {
		sigHex := core.Sig2Hex(sig)
		qDocRef := fbu.quantumC.Doc(sigHex)
		if quantum, ok := signatureQuantumMap[sigHex]; ok {
			fbu.executeQuantumFunc(quantum, qDocRef)
		}
	}

	for _, k := range accept {
		delete(signatureQuantumMap, core.Sig2Hex(k))
	}
	for _, k := range reject {
		delete(signatureQuantumMap, core.Sig2Hex(k))
	}
	for k := range signatureQuantumMap {
		wait = append(wait, core.Hex2Sig(k))
	}

	return
}

func (fbu *FBUniverse) executeInfoPlatformCustom(quantum *core.Quantum, qDocRef *firestore.DocumentRef) {
	if quantum.Type != core.QuantumTypeInfo {
		return
	}

	if len(quantum.Contents) == 0 {
		return
	}

	// platform config should be the last content of quantum
	// and encode by string JSON
	configContent := quantum.Contents[len(quantum.Contents)-1]
	if configContent.Format != core.QCFmtStringJSON {
		return
	}

	var config PlatformConfig
	if err := json.Unmarshal(configContent.Data, &config); err != nil {
		return
	}

	// only process ios:1: ...
	if config.Platform != platformIOS || config.Version != platformVersion {
		return
	}

	platformSetting := make(map[string]string)
	platformSetting["action"] = config.Action
	if config.Action == platformActionComment || config.Action == platformActionReply {
		platformSetting["param"] = config.Params
	}
	data := make(map[string]interface{})
	data[config.Platform] = platformSetting
	qDocRef.Set(fbu.ctx, data, firestore.Merge([]string{config.Platform}))
}

func (fbu *FBUniverse) executeQuantumFunc(quantum *core.Quantum, qDocRef *firestore.DocumentRef) {
	qid, _ := quantum.Ecrecover()
	addrHex := qid.Hex()
	// resave into db as readable data
	if readableCS, err := CS2Readable(quantum.Contents); err == nil {
		readableRecord := make(map[string]interface{})
		readableRecord["rcs"] = readableCS
		qDocRef.Set(fbu.ctx, readableRecord, firestore.Merge([]string{"rcs"}))
	}
	switch quantum.Type {
	case core.QuantumTypeInfo:
		fbu.executeInfoPlatformCustom(quantum, qDocRef)
	case core.QuantumTypeProfile:
		profileMap := make(map[string]*core.QContent)
		readableProfileMap := make(map[string]interface{})
		var mergeKeys []firestore.FieldPath
		for i := 0; i < len(quantum.Contents); i += 2 {
			k := string(quantum.Contents[i].Data)
			mergeKeys = append(mergeKeys, []string{"profile", k})
			mergeKeys = append(mergeKeys, []string{"rp", k})

			profileMap[k] = quantum.Contents[i+1]
			readableProfileMap[k], _ = Content2Readable(quantum.Contents[i+1])
		}

		iDocRef := fbu.individualC.Doc(addrHex)
		dMap, _ := FBStruct2Data(&FBIndividual{Profile: profileMap})
		dMap["rp"] = readableProfileMap
		dMap["updateTime"] = time.Now().UnixMilli()
		mergeKeys = append(mergeKeys, []string{"updateTime"})
		iDocRef.Set(fbu.ctx, dMap, firestore.Merge(mergeKeys...))
	case core.QuantumTypeCommunity:
		minCosignCnt, err := strconv.Atoi(string(quantum.Contents[1].Data))
		if err != nil {
			minCosignCnt = 1
		}
		maxInviteCnt, err := strconv.Atoi(string(quantum.Contents[2].Data))
		if err != nil {
			maxInviteCnt = 0
		}

		initMembers := []string{}
		members := map[string]bool{addrHex: true}
		inviteCnt := map[string]int{addrHex: minCosignCnt}
		for i := 3; i < len(quantum.Contents) && i < 16; i++ {
			memberHex := string(quantum.Contents[i].Data)
			initMembers = append(initMembers, memberHex)
			members[memberHex] = true
			inviteCnt[memberHex] = minCosignCnt
		}

		dMap, _ := FBStruct2Data(&FBCommunity{
			Note:           quantum.Contents[0],
			DefineSigHex:   core.Sig2Hex(quantum.Signature),
			CreatorAddrHex: addrHex,
			MinCosignCnt:   minCosignCnt,
			MaxInviteCnt:   maxInviteCnt,
			InitMembersHex: initMembers,
			Members:        members,
			InviteCnt:      inviteCnt,
		})

		dMap["createTime"] = time.Now().UnixMilli()
		dMap["updateTime"] = time.Now().UnixMilli()
		cDocRef := fbu.communityC.Doc(core.Sig2Hex(quantum.Signature))
		cDocRef.Set(fbu.ctx, dMap)

	case core.QuantumTypeInvitation:
		communtiyHex := core.Sig2Hex(quantum.Contents[0].Data)          // if QCFmtBytesSignature = 33
		if quantum.Contents[0].Format == core.QCFmtStringSignatureHex { // QCFmtStringSignatureHex = 7
			communtiyHex = string(quantum.Contents[0].Data)
		}

		targets := make(map[string]struct{})
		for i := 1; i < len(quantum.Contents); i++ {
			if quantum.Contents[i].Format == core.QCFmtStringAddressHex {
				targets[string(quantum.Contents[i].Data)] = struct{}{}
			} else if quantum.Contents[i].Format == core.QCFmtBytesAddress {
				targets[identity.BytesToAddress(quantum.Contents[i].Data).Hex()] = struct{}{}
			}
		}

		cDocRef := fbu.communityC.Doc(communtiyHex)
		if snapshot, err := cDocRef.Get(fbu.ctx); err == nil {
			dMap := snapshot.Data()

			if members, ok := dMap["members"]; ok {
				if _, ok := members.(map[string]interface{})[addrHex]; ok {
					inviteCnt := dMap["inviteCnt"].(map[string]interface{})

					// TODO : count if out of max sign limit
					var mergeKeys []firestore.FieldPath
					newCommunity := &FBCommunity{Members: make(map[string]bool), InviteCnt: make(map[string]int)}
					for target := range targets {
						if cnt, ok := inviteCnt[target]; ok {
							newCommunity.InviteCnt[target] = int(cnt.(float64)) + 1
						} else {
							newCommunity.InviteCnt[target] = 1
						}
						mergeKeys = append(mergeKeys, []string{"inviteCnt", target})
						if newCommunity.InviteCnt[target] >= int(dMap["minCosignCnt"].(float64)) {
							newCommunity.Members[target] = true
							mergeKeys = append(mergeKeys, []string{"members", target})
						}
					}
					dMap, _ := FBStruct2Data(newCommunity)
					dMap["updateTime"] = time.Now().UnixMilli()
					mergeKeys = append(mergeKeys, []string{"updateTime"})
					cDocRef.Set(fbu.ctx, dMap, firestore.Merge(mergeKeys...))
				}
			}
		}
	case core.QuantumTypeEnd:
		iDocRef := fbu.individualC.Doc(addrHex)
		dMap, _ := FBStruct2Data(&FBIndividual{Attitude: &core.Attitude{Level: core.AttitudeReject}})
		dMap["updateTime"] = time.Now().UnixMilli()
		iDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"attitude", "level"}, []string{"updateTime"}))
	default:
		// core.QuantumTypeInfo or unknown

	}
	// reset copy the bytes from recv to origin, and clear recv
	docSnapshot, _ := qDocRef.Get(fbu.ctx)
	if recv, ok := docSnapshot.Data()["recv"]; ok {
		var refs []*FBSig
		for _, ref := range quantum.References {
			refs = append(refs, &FBSig{SigHex: core.Sig2Hex(ref)})
		}
		dMap := map[string]interface{}{"recv": []byte{}, "origin": recv, "type": quantum.Type, "refs": refs}
		mergeKeys := []firestore.FieldPath{[]string{"recv"}, []string{"origin"}, []string{"type"}, []string{"refs"}}
		qDocRef.Set(fbu.ctx, dMap, firestore.Merge(mergeKeys...))
	}
}

func (fbu *FBUniverse) QueryQuantums(address identity.Address, qType int, skip int, limit int, desc bool) ([]*core.Quantum, error) {
	var qs []*core.Quantum
	quantumQuery := fbu.quantumC.Query

	// TODO: filter by params
	emptyAddress := identity.Address{}
	if address != emptyAddress {
		quantumQuery = quantumQuery.Where("address", "==", address.Hex())
	}

	quantumQuery = quantumQuery.Where("type", "==", qType)

	iter := quantumQuery.Offset(skip).Limit(limit).Documents(fbu.ctx)

	// load all undeal quantums
	for docSnapshot, err := iter.Next(); err != iterator.Done; docSnapshot, err = iter.Next() {
		qRes, err := NewFBQuantumFromSnap(docSnapshot)
		if err != nil {
			return nil, err
		}
		q, err := qRes.GetOriginQuantum()
		if err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}

	return qs, nil
}

func (fbu *FBUniverse) ReceiveQuantums(quantums []*core.Quantum) (accept []core.Sig, wait []core.Sig, reject []core.Sig, err error) {
	for _, q := range quantums {
		qBytes, err := json.Marshal(q)
		if err != nil {
			reject = append(reject, q.Signature)
		} else {
			sigHex := core.Sig2Hex(q.Signature)
			docRef := fbu.quantumC.Doc(sigHex)
			// check exist
			if snapshot, _ := docRef.Get(fbu.ctx); snapshot.Exists() {
				continue
			}

			dMap := make(map[string]interface{})
			dMap["recv"] = qBytes
			if _, err := docRef.Set(fbu.ctx, dMap); err != nil {
				reject = append(reject, q.Signature)
			}
		}
	}
	limit := len(quantums) - len(reject) + 100 // 100 is random number
	accept, wait, r, err := fbu.ProcessQuantums(limit, 0)
	reject = append(reject, r...)
	return
}

func (fbu *FBUniverse) ProcessSingleQuantum(sig core.Sig) error {
	// TODO: check if exist
	// TODO: check if unprocessed
	quantum, err := fbu.loadUnprocessedQuantum(sig)
	if err != nil {
		return err
	}

	// accept []core.Sig, wait []core.Sig, reject []core.Sig, err error
	_, wait, reject, err := fbu.proccessQuantums([]*core.Quantum{quantum})
	if err != nil {
		return err
	} else if len(reject) > 0 {
		return errQuantumIsReject
	} else if len(wait) > 0 {
		return errQuantumIsWaiting
	}

	return nil
}

// HideProcessedQuantum is not belong to interface of PDU Universe. In PDU system, each quantum is
// public and can not be completely removed from system, this function is just for convenience.
func (fbu *FBUniverse) HideProcessedQuantum(sig core.Sig) error {
	pDocRef := fbu.quantumC.Doc(core.Sig2Hex(sig))

	snapshot, err := pDocRef.Get(fbu.ctx)
	if err != nil {
		return err
	}
	if !snapshot.Exists() {
		return errDocumentLoadDataFail
	}
	origin, err := snapshot.DataAt("origin")

	if _, err := pDocRef.Update(fbu.ctx, []firestore.Update{
		{
			Path:  "rcs",
			Value: firestore.Delete,
		},
		{
			Path:  "refs",
			Value: firestore.Delete,
		},
		{
			Path:  "type",
			Value: firestore.Delete,
		},
		{
			Path:  "origin",
			Value: firestore.Delete,
		},
		{
			Path:  "manualIgnore",
			Value: origin,
		},
	}); err != nil {
		return err
	}

	return nil
}

// GetFBDataByTable is not belong to interface of core/universe
func (fbu *FBUniverse) GetFBDataByTable(collection string) ([]byte, error) {
	var res []map[string]interface{}
	iter := fbu.client.Collection(collection).Documents(fbu.ctx)
	for docS, err := iter.Next(); err != iterator.Done; docS, err = iter.Next() {
		res = append(res, docS.Data())
	}
	return json.Marshal(res)
}

// GetQuantums is not belong to interface of core/universe
// this function can be used to backup data by node owner or download data from node by user
func (fbu *FBUniverse) GetQuantums(limit int, skip int, desc bool) ([]*core.Quantum, error) {
	var qs []*core.Quantum
	quantumQuery := fbu.quantumC.Query.Where("seq", ">", 0)

	fDirection := firestore.Desc
	if !desc {
		fDirection = firestore.Asc
	}
	iter := quantumQuery.Offset(skip).Limit(limit).OrderBy("seq", fDirection).Documents(fbu.ctx)

	// load all undeal quantums
	for docSnapshot, err := iter.Next(); err != iterator.Done; docSnapshot, err = iter.Next() {
		qRes, err := NewFBQuantumFromSnap(docSnapshot)
		if err != nil {
			return nil, err
		}
		q, err := qRes.GetOriginQuantum()
		if err != nil {
			return nil, err
		}
		qs = append(qs, q)
	}

	return qs, nil
}

func (fbu *FBUniverse) JudgeIndividual(address identity.Address, level int, judgment string, evidence ...[]core.Sig) error {
	// defult status should be accept & broadcast
	// judge can be done even the individual not exist, use like blacklist.
	iDocRef := fbu.individualC.Doc(address.Hex())

	// check the level
	if level > core.AttitudeBroadcast || level < core.AttitudeRejectOnRef {
		return errIndividualLevelUnknown
	}

	// ignore judgment & evidence
	dMap, _ := FBStruct2Data(&FBIndividual{Attitude: &core.Attitude{Level: level}})
	dMap["updateTime"] = time.Now().UnixMilli()
	iDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"attitude", "level"}, []string{"updateTime"}))

	return nil
}

func (fbu *FBUniverse) JudgeCommunity(sig core.Sig, level int, statement string) error {
	// defulat community should be not follow
	return nil
}

func (fbu *FBUniverse) QueryIndividuals(sig core.Sig, skip int, limit int, desc bool) ([]*core.Individual, error) {
	docID := core.Sig2Hex(sig)
	docRef := fbu.communityC.Doc(docID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return nil, err
	}
	fbCommunity, err := Data2FBCommunity(docSnapshot.Data())
	if err != nil {
		return nil, err
	}
	index := 0
	count := 0
	individuals := []*core.Individual{}
	for addrHex := range fbCommunity.Members {
		if skip <= index {
			if ind, err := fbu.GetIndividual(identity.HexToAddress(addrHex)); err == nil {
				individuals = append(individuals, ind)
				count += 1
				if count >= limit {
					break
				}
			}
		}
		index += 1
	}
	return individuals, nil
}

func (fbu *FBUniverse) GetCommunity(sig core.Sig) (*core.Community, error) {
	docID := core.Sig2Hex(sig)
	docRef := fbu.communityC.Doc(docID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return nil, err
	}

	fbCommunity, err := Data2FBCommunity(docSnapshot.Data())
	if err != nil {
		return nil, err
	}

	community, err := FBCommunity2Community(docID, fbCommunity)
	if err != nil {
		return nil, err
	}
	return community, nil
}

func (fbu *FBUniverse) GetIndividual(address identity.Address) (*core.Individual, error) {
	docID := address.Hex()
	docRef := fbu.individualC.Doc(docID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return nil, err
	}

	fbIndividual, err := Data2FBIndividual(docSnapshot.Data())
	if err != nil {
		return nil, err
	}

	individual, err := FBIndividual2Individual(docID, fbIndividual)
	if err != nil {
		return nil, err
	}
	return individual, nil
}

func (fbu *FBUniverse) GetQuantum(sig core.Sig) (*core.Quantum, error) {

	fbQuantum, err := NewFBQuantumFromDB(sig, fbu.quantumC, fbu.ctx)
	if err != nil {
		return nil, err
	}

	quantum, err := fbQuantum.GetOriginQuantum()
	if err != nil {
		return nil, err
	}

	return quantum, nil
}
