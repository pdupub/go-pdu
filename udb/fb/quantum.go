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

	"cloud.google.com/go/firestore"
	"github.com/pdupub/go-pdu/core"
)

type FBQuantum struct {
	Received        []byte           `json:"recv,omitempty"`
	Origin          []byte           `json:"origin,omitempty"`
	Contents        []*core.QContent `json:"cs,omitempty"`
	Type            int              `json:"type"`
	FBRef           []*FBSig         `json:"refs"`
	Sequence        int64            `json:"seq,omitempty"`
	SelfSeq         int64            `json:"sseq,omitempty"`
	AddrHex         string           `json:"address,omitempty"`
	SigHex          string           `json:"sig,omitempty"`
	ReadableCS      []*FBContent     `json:"rcs,omitempty"`
	CreateTimestamp int64            `json:"createTime"`
}

func NewFBQuantumFromSnap(docSnapshot *firestore.DocumentSnapshot) (*FBQuantum, error) {
	if docSnapshot == nil {
		return nil, errDocumentLoadDataFail
	}
	return Data2FBQuantum(docSnapshot.Data())
}

func NewFBQuantumFromDB(sig core.Sig, collect *firestore.CollectionRef, ctx context.Context) (*FBQuantum, error) {
	docID := core.Sig2Hex(sig)
	docRef := collect.Doc(docID)
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	return NewFBQuantumFromSnap(docSnapshot)
}

func NewFBQuantum(q *core.Quantum) (*FBQuantum, error) {
	sigHex := core.Sig2Hex(q.Signature)
	address, err := q.Ecrecover()
	if err != nil {
		return nil, err
	}
	fbq := &FBQuantum{Contents: q.Contents, Type: q.Type, AddrHex: address.Hex(), SigHex: sigHex}
	for _, ref := range q.References {
		fbq.FBRef = append(fbq.FBRef, &FBSig{SigHex: core.Sig2Hex(ref)})
	}
	return fbq, nil
}

func Data2FBQuantum(d map[string]interface{}) (*FBQuantum, error) {
	dataBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	fbq := new(FBQuantum)
	err = json.Unmarshal(dataBytes, fbq)
	if err != nil {
		return nil, err
	}
	return fbq, nil
}

func (fbq *FBQuantum) GetOriginQuantum() (*core.Quantum, error) {
	resQuantum := core.Quantum{}
	err := json.Unmarshal(fbq.Origin, &resQuantum)
	return &resQuantum, err
}
