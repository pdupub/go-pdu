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

package core

import (
	"strconv"

	"github.com/pdupub/go-pdu/identity"
)

type Community struct {
	Note          *QContent          `json:"note"`
	RuleSig       Sig                `json:"ruleSig"`
	Creator       identity.Address   `json:"creator"`
	BaseCommunity Sig                `json:"baseCommunity"`
	MinCosignCnt  int                `json:"minCosignCnt"`
	MaxInviteCnt  int                `json:"maxInviteCnt"`
	Members       []identity.Address `json:"members"`
}

func NewCommunity(quantum *Quantum) (*Community, error) {
	if quantum.Type != QuantumTypeCommunity {
		return nil, errQuantumTypeNotFit
	}
	creator, err := quantum.Ecrecover()
	if err != nil {
		return nil, err
	}
	community := Community{
		RuleSig: quantum.Signature,
		Creator: creator,
	}
	community.Members = append(community.Members, creator)

	for i, content := range quantum.Contents {
		if i == 0 {
			community.Note = content
		}

		if i == 1 && content.Format == QCFmtBytesSignature {
			// default base community is nil
			community.BaseCommunity = Sig(content.Data)
		}

		if i == 2 {
			// default min cosign count is 1
			community.MinCosignCnt = 1
			if content.Format == QCFmtStringInt {
				community.MinCosignCnt, _ = strconv.Atoi(string(content.Data))
			}
		}

		if i == 3 {
			// default max invite count is 0
			community.MaxInviteCnt = 0
			if content.Format == QCFmtStringInt {
				community.MaxInviteCnt, _ = strconv.Atoi(string(content.Data))
			}
		}

		if i >= 4 && content.Format == QCFmtStringHexAddress {
			community.Members = append(community.Members, identity.HexToAddress(string(content.Data)))
		}
	}

	return &community, nil
}
