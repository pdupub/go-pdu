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
	"encoding/json"

	"github.com/pdupub/go-pdu/core"
)

func Quantum2FBQuantum(q *core.Quantum) (string, *FBQuantum) {
	docID := core.Sig2Hex(q.Signature)
	fbq := &FBQuantum{Contents: q.Contents, Type: q.Type}
	for _, ref := range q.References {
		fbq.FBRef = append(fbq.FBRef, &FBSig{SigHex: core.Sig2Hex(ref)})
	}
	return docID, fbq
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

func Data2FBIndividual(d map[string]interface{}) (*FBIndividual, error) {
	dataBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	fbq := new(FBIndividual)
	err = json.Unmarshal(dataBytes, fbq)
	if err != nil {
		return nil, err
	}
	return fbq, nil
}

func FBQuantum2Quantum(uid string, fbq *FBQuantum) (*core.Quantum, error) {
	q := core.Quantum{}
	q.Contents = fbq.Contents
	q.Type = fbq.Type
	q.Signature = core.Hex2Sig(uid)
	for _, ref := range fbq.FBRef {
		q.References = append(q.References, core.Hex2Sig(ref.SigHex))
	}
	return &q, nil
}

func FBStruct2Data(fbstruct interface{}) (map[string]interface{}, error) {
	qBytes, err := json.Marshal(fbstruct)
	if err != nil {
		return nil, err
	}
	aMap := make(map[string]interface{})
	err = json.Unmarshal(qBytes, &aMap)
	if err != nil {
		return nil, err
	}
	return aMap, nil
}
