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
	"errors"
	"strconv"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
)

var (
	errContentFmtMissing = errors.New("content fmt is missing")
	errContentsMissing   = errors.New("contents is missing")
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

func Data2FBCommunity(d map[string]interface{}) (*FBCommunity, error) {
	dataBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	fbq := new(FBCommunity)
	err = json.Unmarshal(dataBytes, fbq)
	if err != nil {
		return nil, err
	}
	return fbq, nil
}

func FBQuantum2Quantum(uid string, fbq *FBQuantum) (*core.Quantum, error) {
	q := core.Quantum{}
	q.Contents = fbq.Contents
	if fbq.Type > 0 {
		q.Type = fbq.Type
	} else {
		q.Type = -fbq.Type
	}
	q.Signature = core.Hex2Sig(uid)
	for _, ref := range fbq.FBRef {
		q.References = append(q.References, core.Hex2Sig(ref.SigHex))
	}
	return &q, nil
}

func FBIndividual2Individual(uid string, fbi *FBIndividual) (*core.Individual, error) {
	i := core.Individual{}
	i.Address = identity.HexToAddress(uid)
	i.Profile = fbi.Profile
	// i.Communities
	i.Attitude = fbi.Attitude
	i.LastSig = core.Hex2Sig(fbi.LastSigHex)
	i.LastSeq = fbi.LastSelfSeq
	return &i, nil
}

func FBCommunity2Community(uid string, fbc *FBCommunity) (*core.Community, error) {
	c := core.Community{}
	c.Note = fbc.Note
	c.Define = core.Hex2Sig(uid)
	c.Creator = identity.HexToAddress(fbc.CreatorAddrHex)
	c.MinCosignCnt = fbc.MinCosignCnt
	c.MaxInviteCnt = fbc.MaxInviteCnt

	for _, addrHex := range fbc.InitMembersHex {
		c.InitMembers = append(c.InitMembers, identity.HexToAddress(addrHex))
	}

	return &c, nil
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

func CS2Readable(contents interface{}) (interface{}, error) {

	switch contents := contents.(type) {
	case interface{}:
		readableCS := []interface{}{}
		for _, v := range contents.([]*core.QContent) {
			cc, err := Content2Readable(v)
			if err != nil {
				return nil, err
			}
			readableCS = append(readableCS, cc)
		}
		return readableCS, nil
	}
	return nil, errContentsMissing
}

func Content2Readable(content *core.QContent) (map[string]interface{}, error) {
	cc := make(map[string]interface{})
	switch content.Format {
	case core.QCFmtStringTEXT, core.QCFmtStringAddressHex, core.QCFmtStringSignatureHex:
		cc["data"] = string(content.Data)
	case core.QCFmtStringInt, core.QCFmtStringFloat:
		dataFloat, err := strconv.ParseFloat(string(content.Data), 64)
		if err != nil {
			return nil, err
		}
		cc["data"] = dataFloat
	case core.QCFmtBytesAddress:
		cc["data"] = identity.BytesToAddress(content.Data).Hex()
	case core.QCFmtBytesSignature:
		cc["data"] = core.Sig2Hex(content.Data)
	default:
		return nil, errContentFmtMissing
	}
	cc["fmt"] = content.Format

	return cc, nil
}
