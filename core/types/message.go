// Copyright 2019 The PDU Authors
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

package types

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/pdupub/go-pdu/common"
	"math/big"
)

// Message
type Message struct {
	*Vertex
}

// NewMessage
func NewMessage(body MsgBody, sig MsgSig, refs ...MsgRef) (*Message, common.Hash, error) {
	// todo: verify signature
	var parents []interface{}
	for _, ref := range refs {
		parents = append(parents, ref.hash)
	}
	// build id  from refs & body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, common.Hash{}, err
	}
	refsBytes, err := json.Marshal(refs)
	if err != nil {
		return nil, common.Hash{}, err
	}
	sigBytes, err := json.Marshal(sig)
	if err != nil {
		return nil, common.Hash{}, err
	}
	var hashKey common.Hash
	hashKey = sha256.Sum256(append(bodyBytes, append(refsBytes, sigBytes...)...))

	msg := &Message{
		NewVertex(hashKey, body, parents...),
	}
	return msg, hashKey, nil
}

// RootMessage
func RootMessage(body MsgBody, sig MsgSig) (*Message, common.Hash, error) {

	// build id  from refs & body
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, common.Hash{}, err
	}
	sigBytes, err := json.Marshal(sig)
	if err != nil {
		return nil, common.Hash{}, err
	}
	var hashKey common.Hash
	hashKey = sha256.Sum256(append(bodyBytes, sigBytes...))
	msg := &Message{
		NewVertex(hashKey, body),
	}
	return msg, hashKey, nil
}

// MsgBody
// todo : change to interface
type MsgBody struct {
	Nonce    *big.Int
	Category uint
	Title    string
	Author   string
}

// MsgRef is the parents of this vertex
type MsgRef struct {
	gene Gene
	id   string
	hash common.Hash
}

// MsgSig
type MsgSig struct {
	R *big.Int
	S *big.Int
}
