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
	"errors"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	errDocumentLoadDataFail = errors.New("document load data fail")
)

type FBSig struct {
	SigHex string `json:"refs"`
}

type FBContent struct {
	Data   interface{} `json:"data,omitempty"`
	Format int         `json:"fmt"`
}

type FBQuantum struct {
	Contents   []*core.QContent `json:"cs,omitempty"`
	Type       int              `json:"type"`
	FBRef      []*FBSig         `json:"refs"`
	Sequence   int64            `json:"seq,omitempty"`
	SelfSeq    int64            `json:"sseq,omitempty"`
	AddrHex    string           `json:"address,omitempty"`
	ReadableCS []*FBContent     `json:"rcs,omitempty"`
}

type FBIndividual struct {
	LastSigHex  string                    `json:"last"` // last sig of quantum
	LastSelfSeq int64                     `json:"lseq"` // last self sequance
	Profile     map[string]*core.QContent `json:"profile,omitempty"`
	Attitude    *core.Attitude            `json:"attitude"`
}

type FBCommunity struct {
	Note           *core.QContent  `json:"note"`
	CreatorAddrHex string          `json:"creator"`
	MinCosignCnt   int             `json:"minCosignCnt"`
	MaxInviteCnt   int             `json:"maxInviteCnt"`
	InitMembersHex []string        `json:"initMembers,omitempty"`
	Members        map[string]bool `json:"members,omitempty"`
	InviteCnt      map[string]int  `json:"inviteCnt,omitempty"`
}

type UniverseStatus struct {
	Sequence int64 `json:"universeSequence,omitempty"`
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

const (
	universeStatusDocID = "status"
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
	fbu.loadUniverse()

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
	if sequence, ok := dMap["universeSequence"]; ok {
		fbu.status.Sequence = sequence.(int64)
	}
	return nil
}

func (fbu *FBUniverse) increaseUniverseSequence() error {
	docRef := fbu.universeC.Doc(universeStatusDocID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return err
	}
	dMap := docSnapshot.Data()
	if sequence, ok := dMap["universeSequence"]; ok {
		fbu.status.Sequence = sequence.(int64)
		fbu.status.Sequence += 1
	} else {
		fbu.status.Sequence = 1
	}
	dMap["universeSequence"] = fbu.status.Sequence
	docRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"universeSequence"}))

	return nil
}

func (fbu *FBUniverse) ProcessQuantum(skip, limit int) (accept []core.Sig, wait []core.Sig, reject []core.Sig, err error) {

	signatureQuantumMap := make(map[string]*core.Quantum) // sig:quantum
	signatureAddressMap := make(map[string]string)        // sig:address
	referenceSignatureMap := make(map[string]string)      // first_ref_sig:sig
	addressExistMap := make(map[string]bool)              // address:struct{} 	// address exist

	// load all undeal quantums
	iter := fbu.quantumC.Where("type", ">", 0).Offset(skip).Limit(limit).Documents(fbu.ctx)
	for docSnapshot, err := iter.Next(); err != iterator.Done; docSnapshot, err = iter.Next() {

		if docSnapshot == nil {
			return nil, nil, nil, errDocumentLoadDataFail
		}

		// get data of snapshot
		fbqRes, err := Data2FBQuantum(docSnapshot.Data())
		if err != nil {
			continue
		}
		qRes, err := FBQuantum2Quantum(docSnapshot.Ref.ID, fbqRes)
		if err != nil {
			continue
		}

		// ecrecover the author address
		addr, err := qRes.Ecrecover()
		if err != nil {
			continue
		}

		signatureQuantumMap[docSnapshot.Ref.ID] = qRes
		signatureAddressMap[docSnapshot.Ref.ID] = addr.Hex()

		if core.Sig2Hex(qRes.References[0]) != core.Sig2Hex(core.FirstQuantumReference) {
			referenceSignatureMap[core.Sig2Hex(qRes.References[0])] = docSnapshot.Ref.ID
		}

		// update individual attitude
		if _, ok := addressExistMap[addr.Hex()]; !ok {
			addressExistMap[addr.Hex()] = true

			iDocRef := fbu.individualC.Doc(addr.Hex())
			snapshot, err := iDocRef.Get(fbu.ctx)
			if err == nil {
				attitude, err := snapshot.DataAt("attitude")
				if err == nil {
					at := attitude.(map[string]interface{})
					if int(at["level"].(float64)) < core.AttitudeIgnoreContent {
						addressExistMap[addr.Hex()] = false
					}
				}
			}

		}
	}

	// set address for quantums
	for sigHex := range signatureQuantumMap {
		// update quantums with address
		qDocRef := fbu.quantumC.Doc(sigHex)
		// update address info for quantum
		dMap, _ := FBStruct2Data(&FBQuantum{AddrHex: signatureAddressMap[sigHex]})
		qDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"address"}))
	}

	// process first quantums
	for sigHex, quantum := range signatureQuantumMap {
		// check individual
		qDocRef := fbu.quantumC.Doc(sigHex)
		iDocRef := fbu.individualC.Doc(signatureAddressMap[sigHex])
		iDocSnapshot, _ := iDocRef.Get(fbu.ctx)
		if !iDocSnapshot.Exists() && core.Sig2Hex(quantum.References[0]) == core.Sig2Hex(core.FirstQuantumReference) {
			// checked first quantums, can be accepted.
			if err := fbu.increaseUniverseSequence(); err != nil {
				// reject
				reject = append(reject, core.Hex2Sig(sigHex))
				continue
			}
			// add sequence to quantum
			dMap, _ := FBStruct2Data(&FBQuantum{Sequence: fbu.status.Sequence, SelfSeq: int64(1)})
			qDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"seq"}, []string{"sseq"}))

			// add new individual
			newIndividual := &FBIndividual{LastSigHex: sigHex, LastSelfSeq: int64(1), Attitude: &core.Attitude{Level: core.AttitudeAccept}}
			dMap, _ = FBStruct2Data(newIndividual)
			iDocRef.Set(fbu.ctx, dMap)

			accept = append(accept, core.Hex2Sig(sigHex))
		}
	}

	// process quantums
	for addrHex := range addressExistMap {
		iDocRef := fbu.individualC.Doc(addrHex)
		iDocSnapshot, _ := iDocRef.Get(fbu.ctx)
		if iDocSnapshot.Exists() {
			individual, err := Data2FBIndividual(iDocSnapshot.Data())
			if err != nil {
				continue
			}

			for {
				if sigHex, ok := referenceSignatureMap[individual.LastSigHex]; ok {
					// accept the quantum
					if _, ok := signatureQuantumMap[sigHex]; ok {

						// checked first quantums, can be accepted.
						if err := fbu.increaseUniverseSequence(); err != nil {
							reject = append(reject, core.Hex2Sig(sigHex))
							continue
						}

						individual.LastSigHex = sigHex
						individual.LastSelfSeq += 1

						qDocRef := fbu.quantumC.Doc(sigHex)
						// add sequence to quantum
						dMap, _ := FBStruct2Data(&FBQuantum{Sequence: fbu.status.Sequence, SelfSeq: individual.LastSelfSeq})
						qDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"seq"}, []string{"sseq"}))

						// add new individual
						dMap, _ = FBStruct2Data(individual)
						iDocRef.Set(fbu.ctx, dMap)

						accept = append(accept, core.Hex2Sig(sigHex))

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
			addrHex := signatureAddressMap[core.Sig2Hex(quantum.Signature)]

			// resave into db as readable data
			if readableCS, err := CS2Readable(quantum.Contents); err == nil {
				readableRecord := make(map[string]interface{})
				readableRecord["rcs"] = readableCS
				qDocRef.Set(fbu.ctx, readableRecord, firestore.Merge([]string{"rcs"}))
			}
			switch quantum.Type {
			case core.QuantumTypeProfile:
				profileMap := make(map[string]*core.QContent)
				var mergeKeys []firestore.FieldPath
				for i := 0; i < len(quantum.Contents); i += 2 {
					k := string(quantum.Contents[i].Data)
					mergeKeys = append(mergeKeys, []string{"profile", k})
					profileMap[k] = quantum.Contents[i+1]
				}

				iDocRef := fbu.individualC.Doc(addrHex)
				dMap, _ := FBStruct2Data(&FBIndividual{Profile: profileMap})
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
					CreatorAddrHex: addrHex,
					MinCosignCnt:   minCosignCnt,
					MaxInviteCnt:   maxInviteCnt,
					InitMembersHex: initMembers,
					Members:        members,
					InviteCnt:      inviteCnt,
				})
				cDocRef := fbu.communityC.Doc(sigHex)
				cDocRef.Set(fbu.ctx, dMap)

			case core.QuantumTypeInvitation:
				communtiyHex := core.Sig2Hex(quantum.Contents[0].Data)
				targets := make(map[string]struct{})
				for i := 1; i < len(quantum.Contents); i++ {
					targets[string(quantum.Contents[i].Data)] = struct{}{}
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
									newCommunity.InviteCnt[target] = cnt.(int) + 1
								} else {
									newCommunity.InviteCnt[target] = 1
								}
								mergeKeys = append(mergeKeys, []string{"inviteCnt", target})
								if newCommunity.InviteCnt[target] >= int(dMap["minCosignCnt"].(float64)) {
									newCommunity.Members[addrHex] = true
									mergeKeys = append(mergeKeys, []string{"members", target})
								}
							}
							newMap, _ := FBStruct2Data(newCommunity)
							cDocRef.Set(fbu.ctx, newMap, firestore.Merge(mergeKeys...))
						}
					}
				}
			case core.QuantumTypeEnd:
				iDocRef := fbu.individualC.Doc(addrHex)
				dMap, _ := FBStruct2Data(&FBIndividual{Attitude: &core.Attitude{Level: core.AttitudeReject}})
				iDocRef.Set(fbu.ctx, dMap, firestore.Merge([]string{"attitude", "level"}))
			default:
				// core.QuantumTypeInfo or unknown

			}
			// reset quantum type, so this quantum has been deal
			qDocRef.Set(fbu.ctx, map[string]int64{"type": int64(-quantum.Type)}, firestore.Merge([]string{"type"}))
		}
	}

	for _, k := range accept {
		delete(signatureAddressMap, core.Sig2Hex(k))
	}
	for _, k := range reject {
		delete(signatureAddressMap, core.Sig2Hex(k))
	}
	for k, _ := range signatureAddressMap {
		wait = append(wait, core.Hex2Sig(k))
	}

	return
}

func (fbu *FBUniverse) QueryQuantum(address identity.Address, qType int, skip int, limit int, desc bool) ([]*core.Quantum, error) {
	var qs []*core.Quantum
	quantumQuery := fbu.quantumC.Query

	// TODO: filter by params
	emptyAddress := identity.Address{}
	if address != emptyAddress {
		quantumQuery = quantumQuery.Where("address", "==", address.Hex())
	}

	if qType != 0 {
		quantumQuery = quantumQuery.Where("type", "==", -(qType))
	}

	iter := quantumQuery.Offset(skip).Limit(limit).Documents(fbu.ctx)

	// load all undeal quantums
	for docSnapshot, err := iter.Next(); err != iterator.Done; docSnapshot, err = iter.Next() {
		if docSnapshot == nil {
			return nil, errDocumentLoadDataFail
		}

		dMap := docSnapshot.Data()
		// get data of snapshot
		fbqRes, err := Data2FBQuantum(dMap)
		if err != nil {
			return nil, err
		}

		qRes, err := FBQuantum2Quantum(docSnapshot.Ref.ID, fbqRes)
		if err != nil {
			return nil, err
		}

		qs = append(qs, qRes)
	}

	return qs, nil
}

func (fbu *FBUniverse) ReceiveQuantum(originQuantums []*core.Quantum) (accept []core.Sig, wait []core.Sig, reject []core.Sig, err error) {
	for _, q := range originQuantums {
		docID, fbq := Quantum2FBQuantum(q)
		dMap, _ := FBStruct2Data(fbq)

		// add document
		docRef := fbu.quantumC.Doc(docID)
		_, err := docRef.Set(fbu.ctx, dMap)

		if err != nil {
			return nil, nil, nil, err
		}
	}
	return fbu.ProcessQuantum(0, len(originQuantums))
}

func (fbu *FBUniverse) ProcessSingleQuantum(sig core.Sig) error {
	// TODO: check if exist
	// TODO: check if unprocessed
	return nil
}

func (fbu *FBUniverse) JudgeIndividual(address identity.Address, level int, judgment string, evidence ...[]core.Sig) error {
	// defult status should be accept & broadcast
	return nil
}

func (fbu *FBUniverse) JudgeCommunity(sig core.Sig, level int, statement string) error {
	// defulat community should be not follow
	return nil
}

func (fbu *FBUniverse) QueryIndividual(sig core.Sig, skip int, limit int, desc bool) ([]*core.Individual, error) {
	docID := core.Sig2Hex(sig)
	docRef := fbu.communityC.Doc(docID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return nil, nil
	}

	fbCommunity, err := Data2FBCommunity(docSnapshot.Data())
	if err != nil {
		return nil, nil
	}
	index := 0
	count := 0
	individuals := []*core.Individual{}
	for addrHex := range fbCommunity.Members {
		if skip <= index {
			ind, err := fbu.GetIndividual(identity.HexToAddress(addrHex))
			if err != nil {
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
	docID := core.Sig2Hex(sig)
	docRef := fbu.quantumC.Doc(docID)
	docSnapshot, err := docRef.Get(fbu.ctx)
	if err != nil {
		return nil, err
	}

	fbQuantum, err := Data2FBQuantum(docSnapshot.Data())
	if err != nil {
		return nil, err
	}

	quantum, err := FBQuantum2Quantum(docID, fbQuantum)
	if err != nil {
		return nil, err
	}

	return quantum, nil
}
