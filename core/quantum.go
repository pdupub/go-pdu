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
	// QuantumTypeInfo specifies quantum which user just want to share information, such as
	// post articles, reply to others or chat message with others (encrypted by receiver's public key)
	QuantumTypeInfo = 0

	// QuantumTypeProfile specifies the quantum to update user's (signer's) profile.
	// contents = [key1(QCFmtStringTEXT), value1, key2, value2 ...]
	// if user want to update key1, send new QuantumTypeProfile quantum with
	// contents = [key1, newValue ...]
	// if user want to delete key1, send new QuantumTypeProfile quantum with
	// contents = [key1, newValue which content with empty data]
	QuantumTypeProfile = 1

	// QuantumTypeCommunity specifies the quantum of rule to build new community
	// contents[0] is the display information of current community
	// {fmt:QCFmtStringJSON/QCFmtStringTEXT..., data: ...}
	// contents[1] is the number of invitation (co-signature) from users in current community
	// {fmt:QCFmtStringInt, data:1} at least 1. (creater is absolutely same with others)
	// contents[2] is the max number of invitation by one user
	// {fmt:QCFmtStringInt, data:1} -1 means no limit, 0 means not allowed
	// contents[3] ~ contents[15] is the initial users in this community
	// {fmt:QCFmtBytesAddress/QCFmtStringAddressHex, data:0x1232...}
	// signer of this community is also the initial user in this community
	QuantumTypeCommunity = 2

	// QuantumTypeInvitation specifies the quantum of invitation
	// contents[0] is the signature of target community rule quantum
	// {fmt:QCFmtBytesSignature/QCFmtStringSignatureHex, data:signature of target community rule quantum}
	// contents[1] ~ contents[n] is the address of be invited
	// {fmt:QCFmtBytesAddress/QCFmtStringAddressHex, data:0x123...}
	// no matter all invitation send in same quantum of different quantum
	// only first n address (rule.contents[2]) will be accepted
	// user can not quit community, but any user can block any other user (or self) from any community
	// accepted by any community is decided by user in that community feel about u, not opposite.
	// User belong to community, quantum belong to user. (On trandation forum, posts usually belong to
	// one topic and have lots of tag, is just the function easy to implemnt not base struct here)
	QuantumTypeInvitation = 3

	// QuantumTypeEnd specifies the quantum of ending life cycle
	// When receive this quantum, no more quantum from same address should be received or broadcast.
	// The decision of end from any identity should be respected, no matter for security reason or just want to leave.
	QuantumTypeEnd = 4
)

var (
	FirstQuantumReference = Hex2Sig("0x00")
	SameQuantumReference  = Hex2Sig("0x01")
)

// UnsignedQuantum defines the single message from user without signature,
// all variables should be in alphabetical order.
type UnsignedQuantum struct {
	// Contents contain all data in this quantum
	Contents []*QContent `json:"cs,omitempty"`

	// References must from the exist signature or 0x00
	// References[0] is the last signature by user-self, 0x00 if this quantum is the first quantum by user
	// References[1] is the last signature of public or important quantum by user-self, which quantum user
	// want other user to view most. (refs[1] usually contains article and forward, but not comments, like or chat)
	// if References[1] same with References[0], References[1] can be SameQuantumReference(0x01)
	// References[2] ~ References[n] is optional, recommend to use new & valid quantum
	// If two quantums by same user with same References[0], these two quantums will cause conflict, and
	// this user maybe block by others. The reason to do that punishment is user should act like individual,
	// all proactive event from one user should be sequence should be total order (全序关系). References[1~n]
	// do not need follow this restriction, because all other references show the partial order (偏序关系).
	// all the quantums which be reference as References[1] should be total order too.
	References []Sig `json:"refs"`

	// Type specifies the type of this quantum
	Type int `json:"type"`
}

// Quantum defines the single message signed by user.
type Quantum struct {
	UnsignedQuantum
	Signature Sig `json:"sig,omitempty"`
}

// NewQuantum try to build Quantum without signature
func NewQuantum(t int, cs []*QContent, refs ...Sig) (*Quantum, error) {
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
		// TODO: check fmt of each content
	}

	uq := UnsignedQuantum{
		Contents:   cs,
		References: refs,
		Type:       t}

	return &Quantum{UnsignedQuantum: uq}, nil
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

func CreateInfoQuantum(qcs []*QContent, refs ...Sig) (*Quantum, error) {
	return NewQuantum(QuantumTypeInfo, qcs, refs...)
}

func CreateProfileQuantum(profiles map[string]interface{}, refs ...Sig) (*Quantum, error) {
	var qcs []*QContent
	for k, v := range profiles {
		key := CreateTextContent(k)
		var value *QContent
		switch v.(type) {
		case string:
			value = CreateTextContent(v.(string))
		case float64:
			value = CreateFloatContent(v.(float64))
		case int:
			value = CreateIntContent(int64(v.(int)))
		case int64:
			value = CreateIntContent(v.(int64))
		default:
			// img for avator will be deal later
			return nil, errContentFmtNotFit
		}
		qcs = append(qcs, key, value)
	}

	return NewQuantum(QuantumTypeProfile, qcs, refs...)
}

func CreateCommunityQuantum(note string, minCosignCnt int, maxInviteCnt int, initAddrsHex []string, refs ...Sig) (*Quantum, error) {
	var qcs []*QContent

	qcs = append(qcs, CreateTextContent(note))
	qcs = append(qcs, CreateIntContent(int64(minCosignCnt)))
	qcs = append(qcs, CreateIntContent(int64(maxInviteCnt)))
	for _, addrHex := range initAddrsHex {
		qcs = append(qcs, &QContent{Format: QCFmtStringAddressHex, Data: []byte(addrHex)})
	}

	return NewQuantum(QuantumTypeCommunity, qcs, refs...)
}

func CreateInvitationQuantum(target Sig, addrsHex []string, refs ...Sig) (*Quantum, error) {
	var qcs []*QContent
	qcs = append(qcs, &QContent{Format: QCFmtBytesSignature, Data: target})
	for _, addrHex := range addrsHex {
		qcs = append(qcs, &QContent{Format: QCFmtStringAddressHex, Data: []byte(addrHex)})
	}
	return NewQuantum(QuantumTypeInvitation, qcs, refs...)
}

func CreateEndQuantum(refs ...Sig) (*Quantum, error) {
	return NewQuantum(QuantumTypeEnd, []*QContent{}, refs...)
}
