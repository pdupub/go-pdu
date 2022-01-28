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
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/pdupub/go-pdu/identity"
)

const (
	maxReferencesCnt = 16
	maxContentsCnt   = 16
	maxContentSize   = 1024 * 512
)

var (
	errQuantumRefsCntOutOfLimit     = errors.New("quantum references count out of limit")
	errQuantumContentsCntOutOfLimit = errors.New("quantum contents count out of limit")
	errQuantumContentSizeOutOfLimit = errors.New("quantum content size out of limit")
	errQuantumSignedAlready         = errors.New("quantum have been signed already")
)

const (
	// QuantumTypeInfo specifies quantum which user just want to share information.
	QuantumTypeInfo = 1

	// QuantumTypeProfile specifies the quantum to update user's (signer's) profile.
	// contents = [key1(QCFmtStringTEXT), value1, key2, value2 ...]
	// if user want to update key1, send new QuantumTypeProfile quantum with
	// contents = [key1, newValue ...]
	// if user want to delete key1, send new QuantumTypeProfile quantum with
	// contents = [key1, newValue which content with empty data]
	QuantumTypeProfile = 2

	// QuantumTypeRule specifies the quantum of rule to build new group
	// contents[0] is base group of current group
	// {fmt:QCFmtBytesSignature, data:signature of base group rule quantum}
	// contents[1] is the number of co-signature from users in current group
	// {fmt: QCFmtStringInt, data:1} number should at least 1
	// contents[2] is the max number of create by one user
	// {fmt: QCFmtStringInt, data:1} -1 means no limit, 0 means not allowed
	// contents[3] ~ contents[15] is the initial users in this group
	// {fmt: QCFmtBytesAddress, data:0x1232...}
	// signer of this group is also the initial user in this group
	QuantumTypeRule = 3
	// QuantumTypeInvite specifies
	QuantumTypeInvite = 4
	// QuantumTypeQuit specifies
	QuantumTypeQuit = 5
)

// UnsignedQuantum defines the single message from user without signature,
// all variables should be in alphabetical order.
type UnsignedQuantum struct {
	Contents   []*QContent `json:"cs"`
	References [][]byte    `json:"refs"`
	Type       int         `json:"type"`
}

// Quantum defines the single message signed by user.
type Quantum struct {
	UnsignedQuantum
	Signature []byte `json:"sig,omitempty"`
}

const (
	QCFmtStringTEXT  = 1
	QCFmtStringURL   = 2
	QCFmtStringJSON  = 3
	QCFmtStringInt   = 4
	QCFmtStringFloat = 5

	QCFmtBytesAddress   = 33
	QCFmtBytesSignature = 34

	QCFmtImagePNG = 65
	QCFmtImageJPG = 66
	QCFmtImageBMP = 67

	QCFmtAudioWAV = 97
	QCFmtAudioMP3 = 98

	QCFmtVideoMP4 = 129
)

// QContent is one piece of data in Quantum,
// all variables should be in alphabetical order.
type QContent struct {
	Data   []byte `json:"data,omitempty"`
	Format int    `json:"fmt"`
}

// NewQuantum try to build Quantum without signature
func NewQuantum(t int, cs []*QContent, refs ...[]byte) (*Quantum, error) {
	if len(cs) > maxContentsCnt {
		return nil, errQuantumContentsCntOutOfLimit
	}
	if len(refs) > maxReferencesCnt {
		return nil, errQuantumRefsCntOutOfLimit
	}
	for _, v := range cs {
		if len(v.Data) > maxContentSize {
			return nil, errQuantumContentSizeOutOfLimit
		}
		// TODO: check fmt of each content
	}

	uq := UnsignedQuantum{
		Contents:   cs,
		References: refs,
		Type:       t}

	return &Quantum{UnsignedQuantum: uq}, nil
}

func NewContent(fmt int, data []byte) (*QContent, error) {
	return &QContent{Format: fmt, Data: data}, nil
}

// Sign try to add signature to Quantum
func (q *Quantum) Sign(did *identity.DID) error {
	if q.Signature != nil {
		return errQuantumSignedAlready
	}

	b, err := json.Marshal(q.UnsignedQuantum)
	if err != nil {
		return err
	}
	hash := crypto.Keccak256(b)
	sig, err := crypto.Sign(hash, did.GetKey().PrivateKey)
	if err != nil {
		return err
	}
	q.Signature = sig
	return nil
}

// Ecrecover recover
func (q *Quantum) Ecrecover() (common.Address, error) {
	b, err := json.Marshal(q.UnsignedQuantum)
	if err != nil {
		return common.Address{}, err
	}
	hash := crypto.Keccak256(b)
	pk, err := crypto.Ecrecover(hash, q.Signature)
	if err != nil {
		return common.Address{}, err
	}
	signer := common.Address{}
	copy(signer[:], crypto.Keccak256(pk[1:])[12:])

	return signer, nil
}
