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
	"encoding/json"
	"errors"

	"github.com/pdupub/go-pdu/identity"
)

const (
	maxReferencesCnt = 64
	maxContentsCnt   = 16
	maxContentSize   = 1024 * 512
)

var (
	errQuantumRefsCntOutOfLimit     = errors.New("quantum references count out of limit")
	errQuantumContentsCntOutOfLimit = errors.New("quantum contents count out of limit")
	errQuantumContentSizeOutOfLimit = errors.New("quantum content size out of limit")
)

var (
	InitialQuantumReference = Hex2Sig("0x00")
)

const (
	// QuantumTypeInformation specifies the quantum to post information. can be omitted
	QuantumTypeInformation = 0

	// QuantumTypeIntegration specifies the quantum to update user's (signer's) profile.
	QuantumTypeIntegration = 1
)

type QCS []*QContent
type UnsignedQuantum struct {
	// Contents contain all data in this quantum
	Contents QCS `json:"cs,omitempty"`
	// Nonce specifies the nonce of this quantum
	Nonce int `json:"nonce"`
	// References contains all references in this quantum
	References []Sig `json:"refs"`
	// Type specifies the type of this quantum
	Type int `json:"type,omitempty"`
}

func (cs QCS) String() (string, error) {
	b, err := json.Marshal(cs)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Quantum defines the single message signed by user.
type Quantum struct {
	UnsignedQuantum
	Signature Sig `json:"sig,omitempty"`
}

func NewQuantum(t int, cs []*QContent, nonce int, refs ...Sig) (*Quantum, error) {
	if len(cs) > maxContentsCnt {
		return nil, errQuantumContentsCntOutOfLimit
	}
	if len(refs) > maxReferencesCnt || len(refs) < 1 {
		return nil, errQuantumRefsCntOutOfLimit
	}
	for _, v := range cs {
		if len(v.Data) > maxContentSize {
			return nil, errQuantumContentSizeOutOfLimit
		}
	}

	uq := UnsignedQuantum{
		Contents:   cs,
		Nonce:      nonce,
		References: refs,
		Type:       t,
	}

	return &Quantum{UnsignedQuantum: uq}, nil
}

func (q *Quantum) Sign(did *identity.DID) error {

	b, err := json.Marshal(q.UnsignedQuantum)
	if err != nil {
		return err
	}
	sig, err := did.Sign(b)
	if err != nil {
		return err
	}
	q.Signature = sig
	return nil
}

// Ecrecover recover
func (q *Quantum) Ecrecover() (identity.Address, error) {
	b, err := json.Marshal(q.UnsignedQuantum)
	if err != nil {
		return identity.Address{}, err
	}
	return identity.Ecrecover(b, q.Signature)
}

// Marshal converts the UnsignedQuantum to a byte slice
func (uq *UnsignedQuantum) Marshal() ([]byte, error) {
	return json.Marshal(uq)
}

// Unmarshal populates the UnsignedQuantum from a byte slice
func (uq *UnsignedQuantum) Unmarshal(data []byte) error {
	return json.Unmarshal(data, uq)
}
