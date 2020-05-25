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

package pdu

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pborman/uuid"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/crypto"
)

// PEngine is the engine of pdu
type PEngine struct {
	name string
}

// New is used to create PEngine
func New() *PEngine {
	return &PEngine{name: crypto.PDU}
}

// Name return name of pdu (PDU)
func (e PEngine) Name() string {
	return e.name
}

// GenKey generate the private and public key pair
func (e PEngine) GenKey(params ...interface{}) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	return crypto.GenKey(e.name, genKey, params...)
}

// parsePriKey parse the private key
func parsePriKey(priKey interface{}) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	switch priKey.(type) {
	case *ecdsa.PrivateKey:
		pk = priKey.(*ecdsa.PrivateKey)
	case ecdsa.PrivateKey:
		*pk = priKey.(ecdsa.PrivateKey)
	case []byte:
		pk.PublicKey.Curve = elliptic.P256()
		pk.D = new(big.Int).SetBytes(priKey.([]byte))
		pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	case *big.Int:
		pk.PublicKey.Curve = elliptic.P256()
		pk.D = new(big.Int).Set(priKey.(*big.Int))
		pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

// parsePubKey parse the public key
func parsePubKey(pubKey interface{}) (*ecdsa.PublicKey, error) {
	pk := new(ecdsa.PublicKey)
	switch pubKey.(type) {
	case *ecdsa.PublicKey:
		pk = pubKey.(*ecdsa.PublicKey)
	case ecdsa.PublicKey:
		*pk = pubKey.(ecdsa.PublicKey)
	case []byte:
		pk.Curve = elliptic.P256()
		pk.X = new(big.Int).SetBytes(pubKey.([]byte)[:32])
		pk.Y = new(big.Int).SetBytes(pubKey.([]byte)[32:])
	case *big.Int:
		pk.Curve = elliptic.P256()
		pk.X = new(big.Int).SetBytes(pubKey.(*big.Int).Bytes()[:32])
		pk.Y = new(big.Int).SetBytes(pubKey.(*big.Int).Bytes()[32:])
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

func sign(hash []byte, privKey interface{}) ([]byte, *ecdsa.PublicKey, error) {

	pk, err := parsePriKey(privKey)
	if err != nil {
		return nil, nil, err
	}
	r, s, err := ecdsa.Sign(rand.Reader, pk, hash[:])
	if err != nil {
		return nil, nil, err
	}

	rb := common.Bytes2Hash(r.Bytes())
	sb := common.Bytes2Hash(s.Bytes())
	signature := append(rb[:], sb[:]...)

	return signature, &pk.PublicKey, nil
}

// Sign is used to create signature of content by private key
func (e PEngine) Sign(hash []byte, priKey *crypto.PrivateKey) (*crypto.Signature, error) {
	return crypto.Sign(e.name, hash, priKey, sign)
}

func verify(hash []byte, pubKey interface{}, signature []byte) (bool, error) {
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return false, err
	}
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:])
	return ecdsa.Verify(pk, hash, r, s), nil
}

func parseMulSig(signature []byte) [][]byte {
	size := len(signature) / 64
	sigs := make([][]byte, size)
	for i := 0; i < size; i++ {
		sigs[i] = signature[i*64 : (i+1)*64]
	}
	return sigs
}

// Verify is used to verify the signature
func (e PEngine) Verify(hash []byte, sig *crypto.Signature) (bool, error) {
	return crypto.Verify(e.name, hash, sig, verify, parseMulSig)
}

// Unmarshal unmarshal private & public key from json
func (e PEngine) Unmarshal(privKeyBytes, pubKeyBytes []byte) (privKey *crypto.PrivateKey, pubKey *crypto.PublicKey, err error) {

	if len(pubKeyBytes) > 0 {
		if pubKey, err = e.unmarshalPubKey(pubKeyBytes); err != nil {
			return
		}
	}
	if len(privKeyBytes) > 0 {
		if privKey, err = e.unmarshalPrivKey(privKeyBytes); err != nil {
			return
		}
	}
	return
}

func (e PEngine) unmarshalPrivKey(input []byte) (*crypto.PrivateKey, error) {
	p := crypto.PrivateKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	}
	p.Source = aMap["source"].(string)
	p.SigType = aMap["sigType"].(string)

	if p.Source == e.name {
		if p.SigType == crypto.Signature2PublicKey {
			pk := aMap["privKey"].(interface{})
			d, err := common.String2Hash(pk.(string))
			if err != nil {
				return nil, err
			}
			privKey, err := parsePriKey(common.Hash2Bytes(d))
			if err != nil {
				return nil, err
			}
			p.PriKey = privKey
		} else if p.SigType == crypto.MultipleSignatures {
			pk := aMap["privKey"].([]interface{})
			var privKeys []interface{}
			for i := 0; i < len(pk); i++ {
				d, err := common.String2Hash(pk[i].(string))
				if err != nil {
					return nil, err
				}
				privKey, err := parsePriKey(common.Hash2Bytes(d))
				if err != nil {
					return nil, err
				}
				privKeys = append(privKeys, privKey)
			}
			p.PriKey = privKeys
		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}
	return &p, nil
}

func (e PEngine) unmarshalPubKey(input []byte) (*crypto.PublicKey, error) {
	p := crypto.PublicKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	}
	p.Source = aMap["source"].(string)
	p.SigType = aMap["sigType"].(string)

	if p.Source == e.name {
		if p.SigType == crypto.Signature2PublicKey {
			pk := aMap["pubKey"].([]interface{})
			x, err := common.String2Hash(pk[0].(string))
			if err != nil {
				return nil, err
			}
			y, err := common.String2Hash(pk[1].(string))
			if err != nil {
				return nil, err
			}
			pubKey, err := parsePubKey(append(common.Hash2Bytes(x), common.Hash2Bytes(y)...))
			if err != nil {
				return nil, err
			}
			p.PubKey = pubKey
		} else if p.SigType == crypto.MultipleSignatures {
			pk := aMap["pubKey"].([]interface{})
			var pubKeys []interface{}
			for i := 0; i < len(pk)/2; i++ {
				x, err := common.String2Hash(pk[i*2].(string))
				if err != nil {
					return nil, err
				}
				y, err := common.String2Hash(pk[i*2+1].(string))
				if err != nil {
					return nil, err
				}
				pubKey, err := parsePubKey(append(common.Hash2Bytes(x), common.Hash2Bytes(y)...))
				if err != nil {
					return nil, err
				}
				pubKeys = append(pubKeys, pubKey)
			}
			p.PubKey = pubKeys
		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}

	return &p, nil
}

// Marshal marshal private & public key to json
func (e PEngine) Marshal(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) (privKeyBytes []byte, pubKeyBytes []byte, err error) {
	if privKey != nil {
		if privKeyBytes, err = e.marshalPrivKey(privKey); err != nil {
			return
		}
	}
	if pubKey != nil {
		if pubKeyBytes, err = e.marshalPubKey(pubKey); err != nil {
			return
		}
	}
	return
}

func (e PEngine) marshalPrivKey(a *crypto.PrivateKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == e.name {
		if a.SigType == crypto.Signature2PublicKey {
			pk, err := parsePriKey(a.PriKey)
			if err != nil {
				return nil, err
			}
			aMap["privKey"] = common.Bytes2String(pk.D.Bytes())
		} else if a.SigType == crypto.MultipleSignatures {
			switch a.PriKey.(type) {
			case []interface{}:
				pks := a.PriKey.([]interface{})
				privKey := make([]string, len(pks))
				for i, v := range pks {
					pk, err := parsePriKey(v)
					if err != nil {
						return nil, err
					}
					privKey[i] = common.Bytes2String(pk.D.Bytes())
				}
				aMap["privKey"] = privKey
			default:
				return nil, crypto.ErrSigTypeNotSupport
			}
		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}
	return json.Marshal(aMap)
}

func (e PEngine) marshalPubKey(a *crypto.PublicKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == e.name {
		if a.SigType == crypto.Signature2PublicKey {
			pk, err := parsePubKey(a.PubKey)
			if err != nil {
				return nil, err
			}
			pubKey := make([]string, 2)
			pubKey[0] = common.Bytes2String(pk.X.Bytes())
			pubKey[1] = common.Bytes2String(pk.Y.Bytes())
			aMap["pubKey"] = pubKey
		} else if a.SigType == crypto.MultipleSignatures {
			switch a.PubKey.(type) {
			case []interface{}:
				pks := a.PubKey.([]interface{})
				pubKey := make([]string, len(pks)*2)
				for i, v := range pks {
					pk, err := parsePubKey(v)
					if err != nil {
						return nil, err
					}
					pubKey[i*2] = common.Bytes2String(pk.X.Bytes())
					pubKey[i*2+1] = common.Bytes2String(pk.Y.Bytes())
				}
				aMap["pubKey"] = pubKey
			default:
				return nil, crypto.ErrSigTypeNotSupport
			}
		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}
	return json.Marshal(aMap)
}

// EncryptKey encryptKey into file
func (e PEngine) EncryptKey(priKey *crypto.PrivateKey, pass string) ([]byte, error) {
	if priKey.Source != crypto.PDU {
		return nil, crypto.ErrSourceNotMatch
	}
	var ekl crypto.EncryptedKeyJListV3
	if priKey.SigType == crypto.Signature2PublicKey {
		pk, err := parsePriKey(priKey.PriKey)
		if err != nil {
			return nil, err
		}
		ekj, err := e.encryptKey(pk, pass)
		if err != nil {
			return nil, err
		}
		ekl = append(ekl, ekj)
	} else if priKey.SigType == crypto.MultipleSignatures {
		for _, v := range priKey.PriKey.([]interface{}) {
			pk, err := parsePriKey(v)
			if err != nil {
				return nil, err
			}
			ekj, err := e.encryptKey(pk, pass)
			if err != nil {
				return nil, err
			}
			ekl = append(ekl, ekj)
		}
	}
	return json.Marshal(crypto.EncryptedPrivateKey{Source: crypto.PDU, SigType: priKey.SigType, EPK: ekl})
}

func (e PEngine) encryptKey(priKey *ecdsa.PrivateKey, pass string) (*crypto.EncryptedKeyJSONV3, error) {
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    eth.Address{},
		PrivateKey: priKey,
	}

	keyBytes := key.PrivateKey.D.Bytes()
	cryptoStruct, err := keystore.EncryptDataV3(keyBytes, []byte(pass), keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	encryptedKeyJSONV3 := crypto.EncryptedKeyJSONV3{
		Address: hex.EncodeToString(key.Address[:]),
		Crypto:  cryptoStruct,
		ID:      key.Id.String(),
		Version: crypto.EncryptedVersion,
	}
	return &encryptedKeyJSONV3, nil
}

// DecryptKey decrypt private key from file
func (e PEngine) DecryptKey(keyJSON []byte, pass string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	var k crypto.EncryptedPrivateKey
	if err := json.Unmarshal(keyJSON, &k); err != nil {
		return nil, nil, err
	} else if k.Source != crypto.PDU {
		return nil, nil, crypto.ErrSourceNotMatch
	}
	var priKeys, pubKeys []interface{}
	for _, v := range k.EPK {
		keyBytes, err := keystore.DecryptDataV3(v.Crypto, pass)
		if err != nil {
			return nil, nil, err
		}
		pk, err := parsePriKey(keyBytes)
		if err != nil {
			return nil, nil, err
		}
		pk.PublicKey.Curve = elliptic.P256()
		pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())

		priKeys = append(priKeys, pk)
		pubKeys = append(pubKeys, &pk.PublicKey)
		if k.SigType == crypto.Signature2PublicKey {
			return &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PriKey: pk}, &crypto.PublicKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PubKey: &pk.PublicKey}, nil
		}
	}
	return &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.MultipleSignatures, PriKey: priKeys}, &crypto.PublicKey{Source: crypto.PDU, SigType: crypto.MultipleSignatures, PubKey: pubKeys}, nil

}

func genKey() (interface{}, interface{}, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}
