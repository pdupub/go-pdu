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

type Species struct {
	Note         *QContent          `json:"note"`
	Define       Sig                `json:"define"`
	Creator      identity.Address   `json:"creator"`
	MinCosignCnt int                `json:"minCosignCnt"`
	MaxInviteCnt int                `json:"maxInviteCnt"`
	InitMembers  []identity.Address `json:"initMembers"`
}

func NewSpecies(quantum *Quantum) (*Species, error) {
	if quantum.Type != QuantumTypeSpeciation {
		return nil, errQuantumTypeNotFit
	}
	creator, err := quantum.Ecrecover()
	if err != nil {
		return nil, err
	}
	species := Species{
		Define:  quantum.Signature,
		Creator: creator,
	}

	for i, content := range quantum.Contents {
		if i == 0 {
			species.Note = content
		}

		if i == 1 {
			// default min cosign count is 1
			species.MinCosignCnt = 1
			if content.Format == QCFmtStringInt {
				species.MinCosignCnt, _ = strconv.Atoi(string(content.Data))
			}
		}

		if i == 2 {
			// default max invite count is 0
			species.MaxInviteCnt = 0
			if content.Format == QCFmtStringInt {
				species.MaxInviteCnt, _ = strconv.Atoi(string(content.Data))
			}
		}

		if i >= 3 && content.Format == QCFmtStringAddressHex {
			species.InitMembers = append(species.InitMembers, identity.HexToAddress(string(content.Data)))
		}
	}

	return &species, nil
}
