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

const (
	AttitudeRejectOnRef   = -2 // reject the quantum which use any quantum from this address as reference
	AttitudeReject        = -1 // reject any quantum from this address
	AttitudeIgnoreContent = 1  // accept the quantum from this address, but not eval the content, such as identify ...
	AttitudeAccept        = 2  // accept the quantum from this address as normal
	AttitudeBroadcast     = 3  // accept the quantum from this address, broadcast them and used them as reference
)

// Attitude show the subjective attitude to publisher & reasons
type Attitude struct {
	Level    int        `json:"level"`              // Attitude const level
	Judgment string     `json:"judgment,omitempty"` // my subjective judgment
	Evidence []*Quantum `json:"evidence,omitempty"` // evidence of my judgment, which can be omit but all quantum should come from current publisher
}

// Publisher is the publisher in pdu system
type Publisher struct {
	Address  identity.Address     `json:"address"`
	Profile  map[string]*QContent `json:"profile,omitempty"` // profile info base on quantums user accept from this address
	Attitude *Attitude            `json:"attitude"`
	LastSig  Sig                  `json:"lastSignature,omitempty"`
	LastSeq  int64                `json:"lastSequence,omitempty"`
}

func NewPublisher(address identity.Address) *Publisher {
	return &Publisher{Address: address, Profile: make(map[string]*QContent), Attitude: &Attitude{Level: AttitudeAccept}}
}

func (publisher Publisher) GetAddress() identity.Address {
	return publisher.Address
}

func (publisher *Publisher) UpsertProfile(k string, v *QContent) error {
	if publisher.Profile == nil {
		publisher.Profile = make(map[string]*QContent)
	}
	if publisher.Profile[k] != nil && v == nil {
		delete(publisher.Profile, k)
		return nil
	}

	publisher.Profile[k] = v
	return nil
}

func (publisher *Publisher) UpdateAttitude(na *Attitude) error {
	publisher.Attitude = na
	return nil
}
