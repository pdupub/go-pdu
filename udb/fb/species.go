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

type FBSpecies struct {
	Note            *core.QContent  `json:"note"`
	DefineSigHex    string          `json:"define"`
	CreatorAddrHex  string          `json:"creator"`
	MinCosignCnt    int             `json:"minCosignCnt"`
	MaxIdentifyCnt  int             `json:"maxIdentifyCnt"`
	InitMembersHex  []string        `json:"initMembers"`
	Members         map[string]bool `json:"members"`
	IdentifyCnt     map[string]int  `json:"identifyCnt"`
	CreateTimestamp int64           `json:"createTime"`
	UpdateTimestamp int64           `json:"updateTime"`
}

func FBSpecies2Species(uid string, fbc *FBSpecies) (*core.Species, error) {
	c := core.Species{}
	c.Note = fbc.Note
	c.Define = core.Hex2Sig(uid)
	c.Creator = identity.HexToAddress(fbc.CreatorAddrHex)
	c.MinCosignCnt = fbc.MinCosignCnt
	c.MaxIdentifyCnt = fbc.MaxIdentifyCnt

	for _, addrHex := range fbc.InitMembersHex {
		c.InitMembers = append(c.InitMembers, identity.HexToAddress(addrHex))
	}

	return &c, nil
}

func Data2FBSpecies(d map[string]interface{}) (*FBSpecies, error) {
	dataBytes, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	fbq := new(FBSpecies)
	err = json.Unmarshal(dataBytes, fbq)
	if err != nil {
		return nil, err
	}
	return fbq, nil
}
