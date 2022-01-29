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

type Group struct {
	ruleSig      Sig
	creator      identity.Address
	baseGroupSig Sig
	minCosignCnt int
	maxInviteCnt int
	individuals  []identity.Address
}

func NewGroup(quantum *Quantum) (*Group, error) {
	if quantum.Type != QuantumTypeRule {
		return nil, errQuantumTypeNotFit
	}
	creator, err := quantum.Ecrecover()
	if err != nil {
		return nil, err
	}
	group := Group{
		ruleSig: quantum.Signature,
		creator: creator,
	}
	group.individuals = append(group.individuals, creator)

	for i, content := range quantum.Contents {
		if i == 0 && content.Format == QCFmtBytesSignature {
			group.baseGroupSig = Sig(content.Data)
		}

		if i == 1 && content.Format == QCFmtStringInt {
			group.minCosignCnt, err = strconv.Atoi(string(content.Data))
			if err != nil {
				return nil, err
			}
		}

		if i == 2 && content.Format == QCFmtStringInt {
			group.maxInviteCnt, err = strconv.Atoi(string(content.Data))
			if err != nil {
				return nil, err
			}
		}

		if i >= 3 && content.Format == QCFmtStringHexAddress {
			group.individuals = append(group.individuals, identity.HexToAddress(string(content.Data)))
		}
	}

	return &group, nil
}
