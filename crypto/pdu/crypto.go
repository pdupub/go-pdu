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
	"encoding/json"
	"math/big"

	"github.com/pdupub/go-pdu/crypto"
)

const (
	SourceName = "PDU"

	MultipleSignatures  = "MS"
	Signature2PublicKey = "S2PK"
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

type PDUCrypto struct {
}

func New() *PDUCrypto {
	return &PDUCrypto{}
}

func GetKey(priKey interface{}) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	switch priKey.(type) {
	case *ecdsa.PrivateKey:
		pk = priKey.(*ecdsa.PrivateKey)
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

func GetPubKey(pubKey interface{}) (*ecdsa.PublicKey, error) {
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
	default:
		return nil, crypto.ErrKeyTypeNotSupport
	}
	return pk, nil
}

func (pc *PDUCrypto) Sign(hash []byte, priKey crypto.PrivateKey) (*crypto.Signature, error) {
	if priKey.Source != SourceName {
		return nil, crypto.ErrSourceNotMatch
	}
	switch priKey.SigType {
	case Signature2PublicKey:
		pk, err := GetKey(priKey.PriKey)
		if err != nil {
			return nil, err
		}
		r, s, err := ecdsa.Sign(rand.Reader, pk, hash[:])
		if err != nil {
			return nil, err
		}
		return &crypto.Signature{
			PublicKey: crypto.PublicKey{Source: SourceName, SigType: priKey.SigType, PubKey: pk.PublicKey},
			Signature: append(r.Bytes(), s.Bytes()...),
		}, nil
	case MultipleSignatures:
		pks := priKey.PriKey.([]interface{})
		var pubKeys []ecdsa.PublicKey
		var signature []byte
		for _, item := range pks {
			pk, err := GetKey(item)
			if err != nil {
				return nil, err
			}
			r, s, err := ecdsa.Sign(rand.Reader, pk, hash[:])
			if err != nil {
				return nil, err
			}
			signature = append(signature, append(r.Bytes(), s.Bytes()...)...)
			pubKeys = append(pubKeys, pk.PublicKey)
		}
		return &crypto.Signature{
			PublicKey: crypto.PublicKey{Source: SourceName, SigType: priKey.SigType, PubKey: pubKeys},
			Signature: signature,
		}, nil
	default:
		return nil, crypto.ErrSigTypeNotSupport
	}
}

func (pc *PDUCrypto) Verify(hash []byte, sig crypto.Signature) (bool, error) {
	if sig.Source != SourceName {
		return false, crypto.ErrSourceNotMatch
	}
	switch sig.SigType {
	case Signature2PublicKey:
		pk, err := GetPubKey(sig.PubKey)
		if err != nil {
			return false, err
		}
		r := new(big.Int).SetBytes(sig.Signature[:32])
		s := new(big.Int).SetBytes(sig.Signature[32:])
		return ecdsa.Verify(pk, hash, r, s), nil
	case MultipleSignatures:
		pks := sig.PubKey.([]interface{})
		if len(pks) != len(sig.Signature)/64 {
			return false, crypto.ErrSigPubKeyNotMatch
		}
		for i, pubkey := range pks {
			pk, err := GetPubKey(pubkey)
			if err != nil {
				return false, err
			}
			r := new(big.Int).SetBytes(sig.Signature[i*64 : i*64+32])
			s := new(big.Int).SetBytes(sig.Signature[i*64+32 : i*64+64])
			if !ecdsa.Verify(pk, hash, r, s) {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, crypto.ErrSigTypeNotSupport
	}
}

func UnmarshalJSON(input []byte) (*crypto.PublicKey, error) {
	p := crypto.PublicKey{}
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return nil, err
	} else {
		p.Source = aMap["source"].(string)
		p.SigType = aMap["sigType"].(string)

		if p.Source == SourceName {
			if p.SigType == Signature2PublicKey {
				pubKey := new(ecdsa.PublicKey)
				pubKey.Curve = elliptic.P256()
				pubKey.X, pubKey.Y = big.NewInt(0), big.NewInt(0)
				pk := aMap["pubKey"].([]interface{})
				pubKey.X.UnmarshalText([]byte(pk[0].(string)))
				pubKey.Y.UnmarshalText([]byte(pk[1].(string)))
				p.PubKey = *pubKey
			} else if p.SigType == MultipleSignatures {
				pk := aMap["pubKey"].([]interface{})
				var pubKeys []ecdsa.PublicKey
				for i := 0; i < len(pk)/2; i++ {
					pubKey := new(ecdsa.PublicKey)
					pubKey.Curve = elliptic.P256()
					pubKey.X, pubKey.Y = big.NewInt(0), big.NewInt(0)
					pubKey.X.UnmarshalText([]byte(pk[i*2].(string)))
					pubKey.Y.UnmarshalText([]byte(pk[i*2+1].(string)))
					pubKeys = append(pubKeys, *pubKey)
				}
				p.PubKey = pubKeys
			} else {
				return nil, crypto.ErrSigTypeNotSupport
			}
		} else {
			return nil, crypto.ErrSourceNotMatch
		}
	}
	return &p, nil
}

func MarshalJSON(a crypto.PublicKey) ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == SourceName {
		if a.SigType == Signature2PublicKey {
			pk := a.PubKey.(ecdsa.PublicKey)
			pubKey := make([]string, 2)
			pubKey[0] = pk.X.String()
			pubKey[1] = pk.Y.String()
			aMap["pubKey"] = pubKey
		} else if a.SigType == MultipleSignatures {
			pks := a.PubKey.([]interface{})
			pubKey := make([]string, len(pks)*2)
			for i, v := range pks {
				pk := v.(ecdsa.PublicKey)
				pubKey[i*2] = pk.X.String()
				pubKey[i*2+1] = pk.Y.String()
			}
			aMap["pubKey"] = pubKey
		} else {
			return nil, crypto.ErrSigTypeNotSupport
		}
	} else {
		return nil, crypto.ErrSourceNotMatch
	}
	return json.Marshal(aMap)
}
