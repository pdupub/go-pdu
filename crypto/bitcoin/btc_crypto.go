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

package bitcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	btc "github.com/btcsuite/btcd/btcec"
	"github.com/pdupub/go-pdu/crypto"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pborman/uuid"
)

// BEngine is the engine of BTC
type BEngine struct {
	name string
}

// New create new BEngine
func New() *BEngine {
	return &BEngine{name: crypto.BTC}
}

// Name return the name of this engine (BTC)
func (e BEngine) Name() string {
	return e.name
}

// GenKey generate the private and public key pair
func (e BEngine) GenKey(params ...interface{}) (*crypto.PrivateKey, *crypto.PublicKey, error) {
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
		return &crypto.PrivateKey{Source: e.name, SigType: crypto.Signature2PublicKey, PriKey: pk}, &crypto.PublicKey{Source: e.name, SigType: crypto.Signature2PublicKey, PubKey: pk.PublicKey}, nil

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
			pubKeys = append(pubKeys, pk.PublicKey)
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
		privateKey, _ := btc.PrivKeyFromBytes(btc.S256(), priKey.([]byte))
		return privateKey.ToECDSA(), nil
	case *big.Int:
		privateKey, _ := btc.PrivKeyFromBytes(btc.S256(), priKey.(*big.Int).Bytes())
		return privateKey.ToECDSA(), nil
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
		publicKey, err := btc.ParsePubKey(pubKey.([]byte), btc.S256())
		if err != nil {
			return nil, err
		}
		return publicKey.ToECDSA(), nil
	case *big.Int:
		publicKey, err := btc.ParsePubKey(pubKey.(*big.Int).Bytes(), btc.S256())
		if err != nil {
			return nil, err
		}
		return publicKey.ToECDSA(), nil
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

// Sign is used to create signature of content by private key
func (e BEngine) Sign(hash []byte, priKey *crypto.PrivateKey) (*crypto.Signature, error) {
	if priKey.Source != e.name {
		return nil, crypto.ErrSourceNotMatch
	}
	switch priKey.SigType {
	case crypto.Signature2PublicKey:
		pk, err := parsePriKey(priKey.PriKey)
		if err != nil {
			return nil, err
		}
		privateKey, _ := btc.PrivKeyFromBytes(btc.S256(), pk.D.Bytes())
		signature, err := privateKey.Sign(hash)
		if err != nil {
			return nil, err
		}
		return &crypto.Signature{
			PublicKey: crypto.PublicKey{Source: e.name, SigType: priKey.SigType, PubKey: pk.PublicKey},
			Signature: signature.Serialize(),
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
			privateKey, _ := btc.PrivKeyFromBytes(btc.S256(), pk.D.Bytes())
			signature, err := privateKey.Sign(hash)
			if err != nil {
				return nil, err
			}
			signatures = append(signatures, signature.Serialize()...)
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

func (e BEngine) verify(hash []byte, pubKey interface{}, sig []byte) (bool, error) {
	signature, err := btc.ParseSignature(sig, btc.S256())
	if err != nil {
		return false, err
	}
	pk, err := parsePubKey(pubKey)
	if err != nil {
		return false, err
	}
	return signature.Verify(hash, (*btc.PublicKey)(pk)), nil
}

// Verify is used to verify the signature
func (e BEngine) Verify(hash []byte, sig *crypto.Signature) (bool, error) {
	if sig.Source != e.name {
		return false, crypto.ErrSourceNotMatch
	}
	switch sig.SigType {
	case crypto.Signature2PublicKey:
		return e.verify(hash, sig.PubKey, sig.Signature)
	case crypto.MultipleSignatures:
		pks := sig.PubKey.([]interface{})

		currentSigLeft := -1
		var currentSig []byte
		var sigs [][]byte
		for _, v := range sig.Signature {
			if currentSigLeft == -1 {
				currentSigLeft = 0
				currentSig = []byte{}
				currentSig = append(currentSig, v)
			} else if currentSigLeft == 0 {
				currentSigLeft = int(v)
				currentSig = append(currentSig, v)
			} else {
				currentSig = append(currentSig, v)
				currentSigLeft -= 1
				if currentSigLeft == 0 {
					sigs = append(sigs, currentSig)
					currentSigLeft = -1
				}
			}
		}

		if len(pks) != len(sigs) {
			return false, crypto.ErrSigPubKeyNotMatch
		}
		for i, pubkey := range pks {
			if verify, err := e.verify(hash, pubkey, sigs[i]); err != nil || !verify {
				return verify, err
			}
		}
		return true, nil
	default:
		return false, crypto.ErrSigTypeNotSupport
	}
}

// UnmarshalJSON unmarshal public key from json
func (e BEngine) UnmarshalJSON(input []byte) (*crypto.PublicKey, error) {
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
			pubKey, err := btc.ParsePubKey(pk, btc.S256())
			if err != nil {
				return nil, err
			}
			p.PubKey = *pubKey.ToECDSA()
		} else if p.SigType == crypto.MultipleSignatures {
			pks := aMap["pubKey"].([]interface{})
			var pubKeys []ecdsa.PublicKey
			for _, v := range pks {
				pk, err := hex.DecodeString(v.(string))
				if err != nil {
					return nil, err
				}
				pubKey, err := btc.ParsePubKey(pk, btc.S256())
				if err != nil {
					return nil, err
				}
				pubKeys = append(pubKeys, *pubKey.ToECDSA())
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

// MarshalJSON marshal public key to json
func (e BEngine) MarshalJSON(a crypto.PublicKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == e.name {
		if a.SigType == crypto.Signature2PublicKey {
			pk := a.PubKey.(btc.PublicKey)
			aMap["pubKey"] = hex.EncodeToString(pk.SerializeUncompressed())
		} else if a.SigType == crypto.MultipleSignatures {
			switch a.PubKey.(type) {
			case []ecdsa.PublicKey:
				pks := a.PubKey.([]ecdsa.PublicKey)
				pubKey := make([]string, len(pks))
				for i, pk := range pks {
					bpk := (btc.PublicKey)(pk)
					pubKey[i] = hex.EncodeToString(bpk.SerializeUncompressed())
				}
				aMap["pubKey"] = pubKey
			case []interface{}:
				pks := a.PubKey.([]interface{})
				pubKey := make([]string, len(pks))
				for i, v := range pks {
					pk := v.(ecdsa.PublicKey)
					bpk := (btc.PublicKey)(pk)
					pubKey[i] = hex.EncodeToString(bpk.SerializeUncompressed())
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
func (e BEngine) EncryptKey(priKey *crypto.PrivateKey, pass string) ([]byte, error) {
	if priKey.Source != crypto.BTC {
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
	return json.Marshal(crypto.EncryptedPrivateKey{Source: crypto.BTC, SigType: priKey.SigType, EPK: ekl})

}

func (e BEngine) encryptKey(priKey *ecdsa.PrivateKey, pass string) (*crypto.EncryptedKeyJSONV3, error) {
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
		Id:      key.Id.String(),
		Version: crypto.EncryptedVersion,
	}
	return &encryptedKeyJSONV3, nil
}

// DecryptKey decrypt private key from file
func (e BEngine) DecryptKey(keyJson []byte, pass string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	var k crypto.EncryptedPrivateKey
	if err := json.Unmarshal(keyJson, &k); err != nil {
		return nil, nil, err
	} else if k.Source != crypto.BTC {
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

		pk.PublicKey.Curve = btc.S256()
		pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())

		priKeys = append(priKeys, pk)
		pubKeys = append(pubKeys, pk.PublicKey)
		if k.SigType == crypto.Signature2PublicKey {
			return &crypto.PrivateKey{Source: crypto.BTC, SigType: crypto.Signature2PublicKey, PriKey: pk}, &crypto.PublicKey{Source: crypto.BTC, SigType: crypto.Signature2PublicKey, PubKey: pk.PublicKey}, nil
		}
	}
	return &crypto.PrivateKey{Source: crypto.BTC, SigType: crypto.MultipleSignatures, PriKey: priKeys}, &crypto.PublicKey{Source: crypto.BTC, SigType: crypto.Signature2PublicKey, PubKey: pubKeys}, nil

}

func genKey() (*ecdsa.PrivateKey, error) {
	pk, err := btc.NewPrivateKey(btc.S256())
	return (*ecdsa.PrivateKey)(pk), err
}
