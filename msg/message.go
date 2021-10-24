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

package msg

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/pdupub/go-pdu/identity"
)

type Message struct {
	Content    []byte   `json:"content"` // content by JSON
	References [][]byte `json:"refs"`    // message hash
}

func New(content []byte, refs ...[]byte) *Message {
	return &Message{Content: content, References: refs}
}

func (m *Message) Sign(did *identity.DID) ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	hash := crypto.Keccak256(bytes)
	return crypto.Sign(hash, did.GetKey().PrivateKey)
}

func (m *Message) Verify(signature []byte, address common.Address) error {
	signer, err := m.Ecrecover(signature)
	if err != nil {
		return err
	}
	if address != signer {
		return errors.New("verify signer fail")
	}
	return nil
}

func (m *Message) Ecrecover(signature []byte) (common.Address, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return common.Address{}, err
	}
	hash := crypto.Keccak256(bytes)
	pubkey, err := crypto.Ecrecover(hash, signature)
	if err != nil {
		return common.Address{}, err
	}
	signer := common.Address{}
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	return signer, nil
}
