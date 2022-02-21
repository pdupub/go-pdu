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
	"github.com/pdupub/go-pdu/identity"
)

// Individual is the user in pdu system
type Individual struct {
	Address     identity.Address     `json:"address"`
	Profile     map[string]*QContent `json:"profile"`
	Communities []*Community         `json:"communities"`
	Quantums    []*Quantum           `json:"quantums"`
}

func NewIndividual(address identity.Address) *Individual {
	return &Individual{Address: address, Profile: make(map[string]*QContent)}
}

func (ind Individual) GetAddress() identity.Address {
	return ind.Address
}

func (ind *Individual) UpsertProfile(cs []*QContent) error {
	// upsert profile
	for i := 0; i < len(cs); i += 2 {
		if cs[i].Format == QCFmtStringTEXT {
			k := string(cs[i].Data)
			ind.Profile[k] = cs[i+1]
		}
	}
	return nil
}
