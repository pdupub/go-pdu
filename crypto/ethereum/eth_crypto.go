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

package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/math"
	eth "github.com/ethereum/go-ethereum/crypto"

	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pborman/uuid"
	"github.com/pdupub/go-pdu/crypto"
)

// EEngine is the engine of ETH
type EEngine struct {
	name string
}

// New create new EEngine
func New() *EEngine {
	return &EEngine{name: crypto.ETH}
}

// Name return the name of this engine (ETH)
func (e EEngine) Name() string {
	return e.name
}

// GenKey generate the private and public key pair
func (e EEngine) GenKey(params ...interface{}) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	if len(params) == 0 {
		return nil, nil, crypto.ErrSigTypeNotSupport
	}
	sigType := params[0].(string)
	switch sigType {
	case crypto.Signature2PublicKey:
		pk, err := genKey()
		if err != nil {
			return nil, nil, err
		}
		return &crypto.PrivateKey{Source: e.name, SigType: crypto.Signature2PublicKey, PriKey: pk}, &crypto.PublicKey{Source: e.name, SigType: crypto.Signature2PublicKey, PubKey: &pk.PublicKey}, nil

	case crypto.MultipleSignatures:
		if len(params) == 1 {
			return nil, nil, crypto.ErrParamsMissing
		}
		var privKeys, pubKeys []interface{}
		for i := 0; i < params[1].(int); i++ {
			pk, err := genKey()
			if err != nil {
				return nil, nil, err
			}
			privKeys = append(privKeys, pk)
			pubKeys = append(pubKeys, &pk.PublicKey)
		}
		return &crypto.PrivateKey{Source: e.name, SigType: crypto.MultipleSignatures, PriKey: privKeys}, &crypto.PublicKey{Source: e.name, SigType: crypto.MultipleSignatures, PubKey: pubKeys}, nil
	default:
		return nil, nil, crypto.ErrSigTypeNotSupport
	}
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
		return eth.ToECDSA(priKey.([]byte))
	case *big.Int:
		return eth.ToECDSA(priKey.(*big.Int).Bytes())
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
		return eth.UnmarshalPubkey(pubKey.([]byte))
	case *big.Int:
		return eth.UnmarshalPubkey(pubKey.(*big.Int).Bytes())
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

// Sign is used to create signature of content by private key
func (e EEngine) Sign(hash []byte, priKey *crypto.PrivateKey) (*crypto.Signature, error) {
	if priKey.Source != e.name {
		return nil, crypto.ErrSourceNotMatch
	}
	switch priKey.SigType {
	case crypto.Signature2PublicKey:
		pk, err := parsePriKey(priKey.PriKey)
		if err != nil {
			return nil, err
		}
		signature, err := eth.Sign(signHash(hash), pk)
		if err != nil {
			return nil, err
		}
		return &crypto.Signature{
			PublicKey: crypto.PublicKey{Source: e.name, SigType: priKey.SigType, PubKey: pk.PublicKey},
			Signature: signature,
		}, nil
	case crypto.MultipleSignatures:
		pks := priKey.PriKey.([]interface{})
		var pubKeys []interface{}
		var signatures []byte
		for _, item := range pks {
			pk, err := parsePriKey(item)
			if err != nil {
				return nil, err
			}
			signature, err := eth.Sign(signHash(hash), pk)
			if err != nil {
				return nil, err
			}
			signatures = append(signatures, signature...)
			pubKeys = append(pubKeys, pk.PublicKey)
		}
		return &crypto.Signature{
			PublicKey: crypto.PublicKey{Source: e.name, SigType: priKey.SigType, PubKey: pubKeys},
			Signature: signatures,
		}, nil
	default:
		return nil, crypto.ErrSigTypeNotSupport
	}
}

func (e EEngine) verify(hash []byte, pubKey interface{}, signature []byte) (bool, error) {
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return false, err
	}
	recoveredPubKey, err := eth.SigToPub(signHash(hash), signature)
	if err != nil || recoveredPubKey == nil {
		return false, err
	}

	if recoveredPubKey.Y.Cmp(pk.Y) != 0 || recoveredPubKey.X.Cmp(pk.X) != 0 {
		return false, nil
	}
	return true, nil
}

// Verify is used to verify the signature
func (e EEngine) Verify(hash []byte, sig *crypto.Signature) (bool, error) {
	if sig.Source != e.name {
		return false, crypto.ErrSourceNotMatch
	}
	switch sig.SigType {
	case crypto.Signature2PublicKey:
		return e.verify(hash, sig.PubKey, sig.Signature)
	case crypto.MultipleSignatures:
		pks := sig.PubKey.([]interface{})
		if len(pks) != len(sig.Signature)/65 {
			return false, crypto.ErrSigPubKeyNotMatch
		}
		for i, pubkey := range pks {
			if verify, err := e.verify(hash, pubkey, sig.Signature[i*65:(i+1)*65]); err != nil || !verify {
				return verify, err
			}
		}
		return true, nil
	default:
		return false, crypto.ErrSigTypeNotSupport
	}
}

// Unmarshal unmarshal private & public key from json
func (e EEngine) Unmarshal(privKeyBytes, pubKeyBytes []byte) (privKey *crypto.PrivateKey, pubKey *crypto.PublicKey, err error) {

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

func (e EEngine) unmarshalPrivKey(input []byte) (*crypto.PrivateKey, error) {
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
			pk, err := hex.DecodeString(aMap["privKey"].(string))
			if err != nil {
				return nil, err
			}
			privKey, err := parsePriKey(pk)
			if err != nil {
				return nil, err
			}
			p.PriKey = privKey
		} else if p.SigType == crypto.MultipleSignatures {
			pks := aMap["privKey"].([]interface{})
			var privKeys []interface{}
			for _, v := range pks {
				pk, err := hex.DecodeString(v.(string))
				if err != nil {
					return nil, err
				}
				privKey, err := parsePriKey(pk)
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

func (e EEngine) unmarshalPubKey(input []byte) (*crypto.PublicKey, error) {

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
			pk, err := hex.DecodeString(aMap["pubKey"].(string))
			if err != nil {
				return nil, err
			}
			pubKey, err := parsePubKey(pk)
			if err != nil {
				return nil, err
			}
			p.PubKey = pubKey
		} else if p.SigType == crypto.MultipleSignatures {
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
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}

	return &p, nil
}

// Marshal marshal private & public key to json
func (e EEngine) Marshal(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) (privKeyBytes []byte, pubKeyBytes []byte, err error) {
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

func (e EEngine) marshalPrivKey(a *crypto.PrivateKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == e.name {
		if a.SigType == crypto.Signature2PublicKey {
			pk, err := parsePriKey(a.PriKey)
			if err != nil {
				return nil, err
			}
			aMap["privKey"] = hex.EncodeToString(eth.FromECDSA(pk))
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
					privKey[i] = hex.EncodeToString(eth.FromECDSA(pk))
				}
				aMap["privKey"] = privKey
			}

		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}
	return json.Marshal(aMap)
}

func (e EEngine) marshalPubKey(a *crypto.PublicKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == e.name {
		if a.SigType == crypto.Signature2PublicKey {
			pk, err := parsePubKey(a.PubKey)
			if err != nil {
				return nil, err
			}
			aMap["pubKey"] = hex.EncodeToString(eth.FromECDSAPub(pk))
		} else if a.SigType == crypto.MultipleSignatures {
			switch a.PubKey.(type) {
			case []interface{}:
				pks := a.PubKey.([]interface{})
				pubKey := make([]string, len(pks))
				for i, v := range pks {
					pk, err := parsePubKey(v)
					if err != nil {
						return nil, err
					}
					pubKey[i] = hex.EncodeToString(eth.FromECDSAPub(pk))
				}
				aMap["pubKey"] = pubKey
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
func (e EEngine) EncryptKey(priKey *crypto.PrivateKey, pass string) ([]byte, error) {

	if priKey.Source != crypto.ETH {
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
	return json.Marshal(crypto.EncryptedPrivateKey{Source: crypto.ETH, SigType: priKey.SigType, EPK: ekl})
}

func (e EEngine) encryptKey(priKey *ecdsa.PrivateKey, pass string) (*crypto.EncryptedKeyJSONV3, error) {
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    eth.PubkeyToAddress(priKey.PublicKey),
		PrivateKey: priKey,
	}

	keyBytes := math.PaddedBigBytes(key.PrivateKey.D, 32)
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
func (e EEngine) DecryptKey(keyJSON []byte, pass string) (*crypto.PrivateKey, *crypto.PublicKey, error) {

	var k crypto.EncryptedPrivateKey
	if err := json.Unmarshal(keyJSON, &k); err != nil {
		return nil, nil, err
	} else if k.Source != crypto.ETH {
		return nil, nil, crypto.ErrSourceNotMatch
	}
	var priKeys, pubKeys []interface{}
	for _, v := range k.EPK {
		keyBytes, err := keystore.DecryptDataV3(v.Crypto, pass)
		if err != nil {
			return nil, nil, err
		}
		pk := eth.ToECDSAUnsafe(keyBytes)

		pk.PublicKey.Curve = eth.S256()
		pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())

		priKeys = append(priKeys, pk)
		pubKeys = append(pubKeys, pk.PublicKey)
		if k.SigType == crypto.Signature2PublicKey {
			return &crypto.PrivateKey{Source: crypto.ETH, SigType: crypto.Signature2PublicKey, PriKey: pk}, &crypto.PublicKey{Source: crypto.ETH, SigType: crypto.Signature2PublicKey, PubKey: pk.PublicKey}, nil
		}
	}
	return &crypto.PrivateKey{Source: crypto.ETH, SigType: crypto.MultipleSignatures, PriKey: priKeys}, &crypto.PublicKey{Source: crypto.ETH, SigType: crypto.MultipleSignatures, PubKey: pubKeys}, nil

}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return eth.Keccak256([]byte(msg))
}

func genKey() (*ecdsa.PrivateKey, error) {
	return eth.GenerateKey()
}
