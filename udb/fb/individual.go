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
	"github.com/pdupub/go-pdu/identity"
)

type FBIndividual struct {
	AddrHex         string                    `json:"address"`
	LastSigHex      string                    `json:"last"` // last sig of quantum
	LastSelfSeq     int64                     `json:"lseq"` // last self sequance
	Profile         map[string]*core.QContent `json:"profile,omitempty"`
	ReadableProfile map[string]*FBContent     `json:"rp,omitempty"`
	Attitude        *core.Attitude            `json:"attitude"`
	CreateTimestamp int64                     `json:"createTime"`
	UpdateTimestamp int64                     `json:"updateTime"`
}

func FBIndividual2Individual(uid string, fbi *FBIndividual) (*core.Individual, error) {
	i := core.Individual{}
	i.Address = identity.HexToAddress(uid)
	i.Profile = fbi.Profile
	// i.Species
	i.Attitude = fbi.Attitude
	i.LastSig = core.Hex2Sig(fbi.LastSigHex)
	i.LastSeq = fbi.LastSelfSeq
	return &i, nil
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
