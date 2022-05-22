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

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/pdupub/go-pdu/core"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FBSig struct {
	SigHex string `json:"refs"`
}

type FBQuantum struct {
	Contents []*core.QContent `json:"cs,omitempty"`
	Type     int              `json:"type"`
	FBRef    []*FBSig         `json:"refs"`
	Sequence int64            `json:"seq,omitempty"`
	SelfSeq  int64            `json:"sseq,omitempty"`
	AddrHex  string           `json:"address,omitempty"`
}

type FBIndividual struct {
	LastSigHex  string `json:"last"`
	LastSelfSeq int64  `json:"lseq"`
}

type SysConfig struct {
	Sequence int64 `json:"sequence,omitempty"`
}

type FBS struct {
	ctx    context.Context
	app    *firebase.App
	client *firestore.Client
	config *SysConfig
}

const (
	collectionQuantum    = "quantum"
	collectionCommunity  = "community"
	collectionIndividual = "individual"
	collectionConfig     = "config"

	documentConfigID = "system"
)

func NewFBS(ctx context.Context, keyFilename string, projectID string) (*FBS, error) {
	fbs := &FBS{ctx: ctx, config: &SysConfig{}}
	opt := option.WithCredentialsFile(keyFilename)
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, err
	}
	fbs.app = app

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	fbs.client = client

	// init config
	fbs.UpdateSysConfig(false)

	return fbs, nil
}

func (fbs *FBS) Close() error {
	return fbs.client.Close()
}

func (fbs *FBS) UpdateSysConfig(increase bool) error {
	config := fbs.client.Collection(collectionConfig)
	docRef := config.Doc(documentConfigID)
	docSnapshot, err := docRef.Get(fbs.ctx)
	if err != nil {
		return err
	}
	dMap := docSnapshot.Data()
	if sequence, ok := dMap["sequence"]; ok {
		fbs.config.Sequence = sequence.(int64)
		if increase {
			fbs.config.Sequence += 1
			dMap["sequence"] = fbs.config.Sequence
			docRef.Set(fbs.ctx, dMap)
		}
	}
	return nil
}

func (fbs *FBS) DealNewQuantums() error {

	sqMap := make(map[string]*core.Quantum) // sig:quantum
	saMap := make(map[string]string)        // sig:address
	rsMap := make(map[string]string)        // ref:sig
	asMap := make(map[string]struct{})      // address:struct{}
	var quantumSigHexSlice []string
	individualCollection := fbs.client.Collection(collectionIndividual)
	quantumCollection := fbs.client.Collection(collectionQuantum)

	// load all undeal quantums
	iter := quantumCollection.Where("type", ">", 0).Documents(fbs.ctx)
	for docSnapshot, err := iter.Next(); err != iterator.Done; docSnapshot, err = iter.Next() {

		// get data of snapshot
		fbqRes, err := Data2FBQuantum(docSnapshot.Data())
		if err != nil {
			return err
		}
		qRes, err := FBQuantum2Quantum(docSnapshot.Ref.ID, fbqRes)
		if err != nil {
			return err
		}
		sqMap[docSnapshot.Ref.ID] = qRes

		// ecrecover the author address
		addr, err := qRes.Ecrecover()
		if err != nil {
			return err
		}
		saMap[docSnapshot.Ref.ID] = addr.Hex()

		if core.Sig2Hex(qRes.References[0]) != core.Sig2Hex(core.FirstQuantumReference) {
			rsMap[core.Sig2Hex(qRes.References[0])] = docSnapshot.Ref.ID
		}

		if _, ok := asMap[addr.Hex()]; !ok {
			asMap[addr.Hex()] = struct{}{}
		}
	}

	// set address for quantums
	for sigHex := range sqMap {
		// update quantums with address
		qDocRef := quantumCollection.Doc(sigHex)
		// update address info for quantum
		dMap, _ := FBStruct2Data(&FBQuantum{AddrHex: saMap[sigHex]})
		qDocRef.Set(fbs.ctx, dMap, firestore.Merge([]string{"address"}))
	}

	// process first quantums
	for sigHex, quantum := range sqMap {
		// check individual
		qDocRef := quantumCollection.Doc(sigHex)
		iDocRef := individualCollection.Doc(saMap[sigHex])
		iDocSnapshot, _ := iDocRef.Get(fbs.ctx)
		if !iDocSnapshot.Exists() && core.Sig2Hex(quantum.References[0]) == core.Sig2Hex(core.FirstQuantumReference) {
			// checked first quantums, can be accepted.
			if err := fbs.UpdateSysConfig(true); err != nil {
				return err
			}
			// add sequence to quantum
			dMap, _ := FBStruct2Data(&FBQuantum{Sequence: fbs.config.Sequence, SelfSeq: int64(1)})
			qDocRef.Set(fbs.ctx, dMap, firestore.Merge([]string{"seq"}, []string{"sseq"}))

			// add new individual
			newIndividual := &FBIndividual{LastSigHex: sigHex, LastSelfSeq: int64(1)}
			dMap, _ = FBStruct2Data(newIndividual)
			iDocRef.Set(fbs.ctx, dMap)

			quantumSigHexSlice = append(quantumSigHexSlice, sigHex)
		}
	}

	// process quantums
	for addrHex := range asMap {
		iDocRef := individualCollection.Doc(addrHex)
		iDocSnapshot, _ := iDocRef.Get(fbs.ctx)
		if iDocSnapshot.Exists() {
			individual, err := Data2FBIndividual(iDocSnapshot.Data())
			if err != nil {
				return err
			}

			for {
				if sigHex, ok := rsMap[individual.LastSigHex]; ok {
					// accept the quantum
					if _, ok := sqMap[sigHex]; ok {

						// checked first quantums, can be accepted.
						if err := fbs.UpdateSysConfig(true); err != nil {
							return err
						}

						individual.LastSigHex = sigHex
						individual.LastSelfSeq += 1

						qDocRef := quantumCollection.Doc(sigHex)
						// add sequence to quantum
						dMap, _ := FBStruct2Data(&FBQuantum{Sequence: fbs.config.Sequence, SelfSeq: individual.LastSelfSeq})
						qDocRef.Set(fbs.ctx, dMap, firestore.Merge([]string{"seq"}, []string{"sseq"}))

						// add new individual
						dMap, _ = FBStruct2Data(individual)
						iDocRef.Set(fbs.ctx, dMap)

						quantumSigHexSlice = append(quantumSigHexSlice, sigHex)

					}
				} else {
					break
				}
			}
		}
	}

	for _, sigHex := range quantumSigHexSlice {
		qDocRef := quantumCollection.Doc(sigHex)
		if quantum, ok := sqMap[sigHex]; ok {
			if quantum.Type != 1 {
				// TODO: deal funcs quantums here
			}
			qDocRef.Set(fbs.ctx, map[string]int64{"type": int64(-quantum.Type)}, firestore.Merge([]string{"type"}))
		}
	}

	return nil
}
