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
	MappingKey(*PrivateKey, *PublicKey) (map[string]interface{}, map[string]interface{}, error)
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
type funcSign func([]byte, interface{}) ([]byte, *ecdsa.PublicKey, error)
type funcParseKey func(interface{}) (interface{}, interface{}, error)
type funcVerify func(hash []byte, pubKey interface{}, signature []byte) (bool, error)
type funcPrivKeyToKeyBytes func(interface{}) ([]byte, []byte, error)
type funcParsePubKey func(interface{}) (*ecdsa.PublicKey, error)
type funcParseKeyToString func(interface{}) (string, string, error)
type funcParsePubKeyToString func(interface{}) (string, error)
type funcParseMulSig func(signature []byte) [][]byte

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

// Unmarshal unmarshal private & public key from json
func Unmarshal(source string, privKeyBytes, pubKeyBytes []byte, parseKey funcParseKey, parsePubKey funcParsePubKey) (privKey *PrivateKey, pubKey *PublicKey, err error) {
	if len(pubKeyBytes) > 0 {
		if pubKey, err = unmarshalPubKey(source, pubKeyBytes, parsePubKey); err != nil {
			return
		}
	}
	if len(privKeyBytes) > 0 {
		if privKey, err = unmarshalPrivKey(source, privKeyBytes, parseKey); err != nil {
			return
		}
	}
	return
}

func unmarshalPrivKey(source string, input []byte, parseKey funcParseKey) (*PrivateKey, error) {
	p := PrivateKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	}
	p.Source = aMap["source"].(string)
	p.SigType = aMap["sigType"].(string)

	if p.Source == source {
		if p.SigType == Signature2PublicKey {
			pk := aMap["privKey"].(interface{})
			d, err := hex.DecodeString(pk.(string))
			if err != nil {
				return nil, err
			}
			privKey, _, err := parseKey(d)
			if err != nil {
				return nil, err
			}
			p.PriKey = privKey
		} else if p.SigType == MultipleSignatures {
			pk := aMap["privKey"].([]interface{})
			var privKeys []interface{}
			for i := 0; i < len(pk); i++ {
				d, err := hex.DecodeString(pk[i].(string))
				if err != nil {
					return nil, err
				}
				privKey, _, err := parseKey(d)
				if err != nil {
					return nil, err
				}
				privKeys = append(privKeys, privKey)
			}
			p.PriKey = privKeys
		} else {
			return nil, ErrSigTypeNotSupport
		}
	} else {
		return nil, ErrSourceNotMatch
	}
	return &p, nil
}

func unmarshalPubKey(source string, input []byte, parsePubKey funcParsePubKey) (*PublicKey, error) {
	p := PublicKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	}
	p.Source = aMap["source"].(string)
	p.SigType = aMap["sigType"].(string)

	if p.Source == source {
		if p.SigType == Signature2PublicKey {
			pk, err := hex.DecodeString(aMap["pubKey"].(string))
			if err != nil {
				return nil, err
			}
			pubKey, err := parsePubKey(pk)
			if err != nil {
				return nil, err
			}
			p.PubKey = pubKey
		} else if p.SigType == MultipleSignatures {
			pks := aMap["pubKey"].([]interface{})
			var pubKeys []interface{}
			for _, v := range pks {
				pk, err := hex.DecodeString(v.(string))
				if err != nil {
					return nil, err
				}
				pubKey, err := parsePubKey(pk)
				if err != nil {
					return nil, err
				}
				pubKeys = append(pubKeys, pubKey)
			}
			p.PubKey = pubKeys
		} else {
			return nil, ErrSigTypeNotSupport
		}
	} else {
		return nil, ErrSourceNotMatch
	}

	return &p, nil
}

// Marshal marshal private & public key to json
func Marshal(source string, privKey *PrivateKey, pubKey *PublicKey, parseKeyToString funcParseKeyToString, parsePubKeyToString funcParsePubKeyToString) (privKeyBytes []byte, pubKeyBytes []byte, err error) {
	m := make(map[string]interface{})
	if privKey != nil {
		m, err = MappingPrivKey(source, privKey, parseKeyToString)
		if err != nil {
			return
		}
		if privKeyBytes, err = json.Marshal(m); err != nil {
			return
		}
	}
	if pubKey != nil {
		m, err = MappingPubKey(source, pubKey, parsePubKeyToString)
		if err != nil {
			return
		}
		if pubKeyBytes, err = json.Marshal(m); err != nil {
			return
		}
	}
	return
}

// MappingKey build private & public key content into map for display or marshal
func MappingKey(source string, privKey *PrivateKey, pubKey *PublicKey, parseKeyToString funcParseKeyToString, parsePubKeyToString funcParsePubKeyToString) (privKeyM, pubKeyM map[string]interface{}, err error) {
	if privKey != nil {
		privKeyM, err = MappingPrivKey(source, privKey, parseKeyToString)
		if err != nil {
			return
		}
	}
	if pubKey != nil {
		pubKeyM, err = MappingPubKey(source, pubKey, parsePubKeyToString)
		if err != nil {
			return
		}
	}
	return
}

// MappingPrivKey display the content of private key
func MappingPrivKey(source string, a *PrivateKey, parseKeyToString funcParseKeyToString) (map[string]interface{}, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == source {
		if a.SigType == Signature2PublicKey {
			pk, _, err := parseKeyToString(a.PriKey)
			if err != nil {
				return nil, err
			}
			aMap["privKey"] = pk
		} else if a.SigType == MultipleSignatures {
			switch a.PriKey.(type) {
			case []interface{}:
				pks := a.PriKey.([]interface{})
				privKey := make([]string, len(pks))
				for i, v := range pks {
					pk, _, err := parseKeyToString(v)
					if err != nil {
						return nil, err
					}
					privKey[i] = pk
				}
				aMap["privKey"] = privKey
			default:
				return nil, ErrSigTypeNotSupport
			}
		} else {
			return nil, ErrSigTypeNotSupport
		}
	} else {
		return nil, ErrSourceNotMatch
	}
	return aMap, nil
}

// MappingPubKey display the content of public key
func MappingPubKey(source string, a *PublicKey, parsePubKeyToString funcParsePubKeyToString) (map[string]interface{}, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == source {
		if a.SigType == Signature2PublicKey {
			pk, err := parsePubKeyToString(a.PubKey)
			if err != nil {
				return nil, err
			}
			aMap["pubKey"] = pk
		} else if a.SigType == MultipleSignatures {
			switch a.PubKey.(type) {
			case []interface{}:
				pks := a.PubKey.([]interface{})
				pubKey := make([]string, len(pks))
				for i, v := range pks {
					pk, err := parsePubKeyToString(v)
					if err != nil {
						return nil, err
					}
					pubKey[i] = pk
				}
				aMap["pubKey"] = pubKey
			default:
				return nil, ErrSigTypeNotSupport
			}
		} else {
			return nil, ErrSigTypeNotSupport
		}
	} else {
		return nil, ErrSourceNotMatch
	}
	return aMap, nil
}
