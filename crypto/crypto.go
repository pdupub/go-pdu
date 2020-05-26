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
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
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

	// ErrInvalidPubkey is returned if the public key is invalid
	ErrInvalidPubkey = errors.New("invalid public key")
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
	Unmarshal([]byte, []byte) (*PrivateKey, *PublicKey, error)
	Marshal(*PrivateKey, *PublicKey) ([]byte, []byte, error)
	EncryptKey(*PrivateKey, string) ([]byte, error)
	DecryptKey([]byte, string) (*PrivateKey, *PublicKey, error)
}

// EncryptedPrivateKey is encrypted private key in json
type EncryptedPrivateKey struct {
	Source  string              `json:"source"`
	SigType string              `json:"sigType"`
	EPK     EncryptedKeyJListV3 `json:"priKey"`
}

// EncryptedKeyJSONV3 is from geth
type EncryptedKeyJSONV3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	ID      string              `json:"id"`
	Version int                 `json:"version"`
}

// EncryptedKeyJListV3 is for ms
type EncryptedKeyJListV3 []*EncryptedKeyJSONV3

// EncryptedVersion is version of EncryptedKeyJSONV3
const EncryptedVersion = 3

type funcGenKey func() (interface{}, interface{}, error)

// GenKey generate the private and public key pair
func GenKey(source string, genKey funcGenKey, params ...interface{}) (*PrivateKey, *PublicKey, error) {
	if len(params) == 0 {
		return nil, nil, ErrSigTypeNotSupport
	}
	sigType := params[0].(string)
	switch sigType {
	case Signature2PublicKey:
		privKey, pubKey, err := genKey()
		if err != nil {
			return nil, nil, err
		}
		return &PrivateKey{Source: source, SigType: Signature2PublicKey, PriKey: privKey}, &PublicKey{Source: source, SigType: Signature2PublicKey, PubKey: pubKey}, nil

	case MultipleSignatures:
		if len(params) == 1 {
			return nil, nil, ErrParamsMissing
		}
		var privKeys, pubKeys []interface{}
		for i := 0; i < params[1].(int); i++ {
			privKey, pubKey, err := genKey()
			if err != nil {
				return nil, nil, err
			}
			privKeys = append(privKeys, privKey)
			pubKeys = append(pubKeys, pubKey)
		}
		return &PrivateKey{Source: source, SigType: MultipleSignatures, PriKey: privKeys}, &PublicKey{Source: source, SigType: MultipleSignatures, PubKey: pubKeys}, nil
	default:
		return nil, nil, ErrSigTypeNotSupport
	}
}

type funcSign func([]byte, interface{}) ([]byte, *ecdsa.PublicKey, error)

// Sign is used to create signature of content by private key
func Sign(source string, hash []byte, priKey *PrivateKey, sign funcSign) (*Signature, error) {
	if priKey.Source != source {
		return nil, ErrSourceNotMatch
	}

	switch priKey.SigType {
	case Signature2PublicKey:
		signature, pubKey, err := sign(hash, priKey.PriKey)
		if err != nil {
			return nil, err
		}
		return &Signature{
			PublicKey: PublicKey{Source: source, SigType: priKey.SigType, PubKey: pubKey},
			Signature: signature,
		}, nil
	case MultipleSignatures:
		pks := priKey.PriKey.([]interface{})
		var pubKeys []interface{}
		var signatures []byte
		for _, item := range pks {
			signature, pubKey, err := sign(hash, item)
			if err != nil {
				return nil, err
			}
			signatures = append(signatures, signature...)
			pubKeys = append(pubKeys, pubKey)
		}
		return &Signature{
			PublicKey: PublicKey{Source: source, SigType: priKey.SigType, PubKey: pubKeys},
			Signature: signatures,
		}, nil
	default:
		return nil, ErrSigTypeNotSupport
	}
}

type funcVerify func(hash []byte, pubKey interface{}, signature []byte) (bool, error)

type funcParseMulSig func(signature []byte) [][]byte

// Verify is used to verify the signature
func Verify(source string, hash []byte, sig *Signature, verify funcVerify, parseMulSig funcParseMulSig) (bool, error) {
	if sig.Source != source {
		return false, ErrSourceNotMatch
	}

	switch sig.SigType {
	case Signature2PublicKey:
		return verify(hash, sig.PubKey, sig.Signature)
	case MultipleSignatures:
		pks := sig.PubKey.([]interface{})
		sigs := parseMulSig(sig.Signature)
		if len(pks) != len(sigs) {
			return false, ErrSigPubKeyNotMatch
		}
		for i, pubkey := range pks {
			if verify, err := verify(hash, pubkey, sigs[i]); err != nil || !verify {
				return verify, err
			}
		}
		return true, nil
	default:
		return false, ErrSigTypeNotSupport
	}
}

type funcParseKey func(interface{}) (interface{}, interface{}, error)

// DecryptKey decrypt private key from file
func DecryptKey(source string, keyJSON []byte, pass string, parseKey funcParseKey) (*PrivateKey, *PublicKey, error) {
	var k EncryptedPrivateKey
	if err := json.Unmarshal(keyJSON, &k); err != nil {
		return nil, nil, err
	} else if k.Source != source {
		return nil, nil, ErrSourceNotMatch
	}
	var priKeys, pubKeys []interface{}
	for _, v := range k.EPK {
		keyBytes, err := keystore.DecryptDataV3(v.Crypto, pass)
		if err != nil {
			return nil, nil, err
		}
		privKey, pubKey, err := parseKey(keyBytes)
		if err != nil {
			return nil, nil, err
		}
		priKeys = append(priKeys, privKey)
		pubKeys = append(pubKeys, pubKey)
		if k.SigType == Signature2PublicKey {
			return &PrivateKey{Source: source, SigType: Signature2PublicKey, PriKey: privKey}, &PublicKey{Source: source, SigType: Signature2PublicKey, PubKey: pubKey}, nil
		}
	}
	return &PrivateKey{Source: source, SigType: MultipleSignatures, PriKey: priKeys}, &PublicKey{Source: source, SigType: MultipleSignatures, PubKey: pubKeys}, nil
}

// EncryptSignleKey encrypt single private key
func EncryptSignleKey(keyBytes, address []byte, pass string) (*EncryptedKeyJSONV3, error) {
	cryptoStruct, err := keystore.EncryptDataV3(keyBytes, []byte(pass), keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	encryptedKeyJSONV3 := EncryptedKeyJSONV3{
		Address: hex.EncodeToString(address),
		Crypto:  cryptoStruct,
		ID:      uuid.String(),
		Version: EncryptedVersion,
	}
	return &encryptedKeyJSONV3, nil
}

type funcPrivKeyToKeyBytes func(interface{}) ([]byte, []byte, error)

// EncryptKey encryptKey into file
func EncryptKey(source string, priKey *PrivateKey, pass string, privKeyToKeyBytes funcPrivKeyToKeyBytes) ([]byte, error) {
	if priKey.Source != source {
		return nil, ErrSourceNotMatch
	}
	var ekl EncryptedKeyJListV3
	if priKey.SigType == Signature2PublicKey {
		keyBytes, address, err := privKeyToKeyBytes(priKey.PriKey)
		if err != nil {
			return nil, err
		}
		ekj, err := EncryptSignleKey(keyBytes, address, pass)
		if err != nil {
			return nil, err
		}
		ekl = append(ekl, ekj)

	} else if priKey.SigType == MultipleSignatures {
		for _, v := range priKey.PriKey.([]interface{}) {
			keyBytes, address, err := privKeyToKeyBytes(v)
			if err != nil {
				return nil, err
			}
			ekj, err := EncryptSignleKey(keyBytes, address, pass)
			if err != nil {
				return nil, err
			}
			ekl = append(ekl, ekj)
		}
	}
	return json.Marshal(EncryptedPrivateKey{Source: source, SigType: priKey.SigType, EPK: ekl})
}
