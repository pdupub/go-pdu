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
	errQuantumTypeNotFit            = errors.New("quantum type is not fit")
)

type Sig []byte

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
	// contents[1] is the number of invitation (co-signature) from users in current group
	// {fmt:QCFmtStringInt, data:1} at least 1. (creater is absolutely same with others)
	// contents[2] is the max number of invitation by one user
	// {fmt:QCFmtStringInt, data:1} -1 means no limit, 0 means not allowed
	// contents[3] ~ contents[15] is the initial users in this group
	// {fmt:QCFmtBytesAddress, data:0x1232...}
	// signer of this group is also the initial user in this group
	QuantumTypeRule = 3

	// QuantumTypeInvite specifies the quantum of invite
	// contents[0] is the signature of target group rule quantum
	// {fmt:QCFmtBytesSignature, data:signature of target group rule quantum}
	// contents[1] ~ contents[n] is the address of be invited
	// {fmt:QCFmtBytesAddress, data:0x123...}
	// no matter all invitation send in same quantum of different quantum
	// only first n address (rule.contents[2]) will be accepted
	// user can not quit group, but any user can block any other user (or self) from any group
	// accepted by any group is decided by user in that group feel about u, not opposite.
	// User belong to group, quantum belong to user. (On trandation forum, posts usually belong to
	// one topic and have lots of tag, is just the function easy to implemnt not base struct here)
	QuantumTypeInvite = 4
)

// UnsignedQuantum defines the single message from user without signature,
// all variables should be in alphabetical order.
type UnsignedQuantum struct {
	// Contents contain all data in this quantum
	Contents []*QContent `json:"cs"`

	// References must from the exist signature.
	// References[0] is the last signature by user-self, 0x00000... if this quantum is the first quantum
	// References[1] ~ References[n] is optional, recommend to use new & valid quantum
	// If two quantums by same user with same References[0], these two quantums will cause conflict, and
	// this user maybe block by others. The reason to do that punishment is user should act like individual,
	// all proactive event from one user should be sequence should be total order (全序关系). References[1~n]
	// do not need follow this restriction, because all other references show the partial order (偏序关系).
	References []Sig `json:"refs"`

	// Type specifies the type of this quantum
	Type int `json:"type"`
}

// Quantum defines the single message signed by user.
type Quantum struct {
	UnsignedQuantum
	Signature Sig `json:"sig,omitempty"`
}

const (
	QCFmtStringTEXT       = 1
	QCFmtStringURL        = 2
	QCFmtStringJSON       = 3
	QCFmtStringInt        = 4
	QCFmtStringFloat      = 5
	QCFmtStringHexAddress = 6

	QCFmtBytesSignature = 33

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
func NewQuantum(t int, cs []*QContent, refs ...Sig) (*Quantum, error) {
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
