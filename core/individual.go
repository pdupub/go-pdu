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

// All Attitude State below is used to show my current subjective attitude to
// that individual for new quantums, not influence the quantums already be accepted before.
const (
	AttitudeRejectOnRef   = -2 // reject the quantum which use any quantum from this address as reference
	AttitudeReject        = -1 // reject any quantum from this address
	AttitudeIgnoreContent = 1  // accept the quantum from this address, but not eval the content, such as invite ...
	AttitudeAccept        = 2  // accept the quantum from this address as normal
	AttitudeBroadcast     = 3  // accept the quantum from this address, broadcast them and used them as reference
)

// Attitude show the subjective attitude to individual & reasons
type Attitude struct {
	Level    int        `json:"level"`              // Attitude const level
	Judgment string     `json:"judgment,omitempty"` // my subjective judgment
	Evidence []*Quantum `json:"evidence,omitempty"` // evidence of my judgment, can be omit but all quantum should come from current individual
}

// Individual is the user in pdu system
type Individual struct {
	Address     identity.Address     `json:"address"`
	Profile     map[string]*QContent `json:"profile"`
	Communities []*Community         `json:"communities"`
	Attitude    *Attitude            `json:"attitude"`
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
