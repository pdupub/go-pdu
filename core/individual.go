// Copyright 2024 The PDU Authors
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

// The information in Individual is intended for the information recipient or consumer.
// For different information consumers, the information in Individual, except for the
// address, may vary. Consumer does not necessarily have the identity represented by
// the private key and address (as Individual). All Individual information within each
// information consumer also does not need to be publicly disclosed.

// 在PDU中，用户的身份被分割成信息的发布者和使用者，同一用户往往同时具备这两个身份，但用户的这两个身份间
// 并不存在必然联系。Individual中所包含的信息，使用者本地对于自身可见的信息发布者的记录，这些不同的使用者
// 对于相同的发布者其记录内容除Address之外，都可能不同。使用者无义务公开自己本地的Individual信息，即无
// 义务表达对于其他用户的态度。

// All Attitude State below is used to show my current subjective attitude to
// that individual for new quantums, not influence the quantums already be accepted.
const (
	AttitudeRejectOnRef   = -2 // reject the quantum which use any quantum from this address as reference
	AttitudeReject        = -1 // reject any quantum from this address
	AttitudeIgnoreContent = 1  // accept the quantum from this address, but not eval the content, such as identify ...
	AttitudeAccept        = 2  // accept the quantum from this address as normal
	AttitudeBroadcast     = 3  // accept the quantum from this address, broadcast them and used them as reference
)

// Attitude show the subjective attitude to individual & reasons
type Attitude struct {
	Level    int        `json:"level"`              // Attitude const level
	Judgment string     `json:"judgment,omitempty"` // my subjective judgment
	Evidence []*Quantum `json:"evidence,omitempty"` // evidence of my judgment, which can be omit but all quantum should come from current individual
}

// Individual is the publisher in pdu system
type Individual struct {
	Address  identity.Address     `json:"address"`
	Profile  map[string]*QContent `json:"profile,omitempty"` // profile info base on quantums user accept from this address
	Species  []*Species           `json:"species,omitempty"`
	Attitude *Attitude            `json:"attitude"`
	LastSig  Sig                  `json:"lastSignature,omitempty"`
	LastSeq  int64                `json:"lastSequence,omitempty"`
}

func NewIndividual(address identity.Address) *Individual {
	return &Individual{Address: address, Profile: make(map[string]*QContent), Attitude: &Attitude{Level: AttitudeAccept}}
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

func (ind *Individual) UpdateAttitude(na *Attitude) error {
	ind.Attitude = na
	return nil
}
