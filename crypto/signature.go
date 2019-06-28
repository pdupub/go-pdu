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

package crypto

import "errors"

var (
	ErrSourceNotMatch    = errors.New("signature source not match")
	ErrSigTypeNotSupport = errors.New("signature type not support")
	ErrKeyTypeNotSupport = errors.New("key type not support")
	ErrSigPubKeyNotMatch = errors.New("count of signature and public key not match")
)

type PublicKey struct {
	Source  string      `json:"source"`
	SigType string      `json:"sigType"`
	PubKey  interface{} `json:"pubKey"`
}

type Signature struct {
	PublicKey
	Signature []byte `json:"signature"`
}

type PrivateKey struct {
	Source  string      `json:"source"`
	SigType string      `json:"sigType"`
	PriKey  interface{} `json:"priKey"`
}

type Engine interface {
	Sign([]byte, PrivateKey) (*Signature, error)
	Verify([]byte, Signature) bool
}
