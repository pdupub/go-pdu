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

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

var (
	// ErrParamsMissing is returned when params is not enough
	ErrParamsMissing = errors.New("params missing")

	// ErrSourceNotMatch is returned if the source name of signature and key not match
	ErrSourceNotMatch = errors.New("signature source not match")

	// ErrSigTypeNotSupport is returned if the signature is not MS or S2PK
	ErrSigTypeNotSupport = errors.New("signature type not support")

	// ErrGenerateKeyFail is returned when generate key fail
	ErrGenerateKeyFail = errors.New("generate key fail")

	// ErrKeyTypeNotSupport is returned if key type not support
	ErrKeyTypeNotSupport = errors.New("key type not support")

	// ErrSigPubKeyNotMatch is returned if the signature and key not match for MS
	ErrSigPubKeyNotMatch = errors.New("count of signature and public key not match")
)

const (
	// MultipleSignatures is type of signature by more than one key pairs
	MultipleSignatures = "MS"
	// Signature2PublicKey is type of signature by one key pair
	Signature2PublicKey = "S2PK"

	// BTC is symbol of Bitcoin
	BTC = "BTC"
	// ETH is symbol of Ethereum
	ETH = "ETH"
	// PDU is symbol of PDU
	PDU = "PDU"
)

// PublicKey contains the source name, type and public key content
type PublicKey struct {
	Source  string      `json:"source"`
	SigType string      `json:"sigType"`
	PubKey  interface{} `json:"pubKey"`
}

// Signature contain the signature and public key
type Signature struct {
	PublicKey
	Signature []byte `json:"signature"`
}

// PrivateKey contains the source name, type and private key content
type PrivateKey struct {
	Source  string      `json:"source"`
	SigType string      `json:"sigType"`
	PriKey  interface{} `json:"priKey"`
}

// Engine is an crypto algorithm engine
type Engine interface {
	Name() string
	GenKey(params ...interface{}) (*PrivateKey, *PublicKey, error)
	Sign([]byte, *PrivateKey) (*Signature, error)
	Verify([]byte, *Signature) (bool, error)
	UnmarshalJSON([]byte) (*PublicKey, error)
	MarshalJSON(PublicKey) ([]byte, error)
	EncryptKey(*PrivateKey, string) ([]byte, error)
	DecryptKey([]byte, string) (*PrivateKey, error)
}

type EncryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Id      string              `json:"id"`
	Version int                 `json:"version"`
}

type EncryptedKeyJListV3 []*EncryptedKeyJSONV3

const EncryptedVersion = 3
