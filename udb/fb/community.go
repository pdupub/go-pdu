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

type FBCommunity struct {
	Note           *core.QContent  `json:"note"`
	DefineSigHex   string          `json:"define"`
	CreatorAddrHex string          `json:"creator"`
	MinCosignCnt   int             `json:"minCosignCnt"`
	MaxInviteCnt   int             `json:"maxInviteCnt"`
	InitMembersHex []string        `json:"initMembers,omitempty"`
	Members        map[string]bool `json:"members,omitempty"`
	InviteCnt      map[string]int  `json:"inviteCnt,omitempty"`
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
