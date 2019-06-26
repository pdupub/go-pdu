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
	"errors"
	"math/big"

	"github.com/pdupub/go-pdu/crypto"
)

const (
	SourceName = "PDU"

	MultipleSignatures  = "MS"
	Signature2PublicKey = "S2PK"
)

var (
	errSourceNotMatch    = errors.New("signature source not match")
	errSigTypeNotSupport = errors.New("signature type not support")
	errKeyTypeNotSupport = errors.New("key type not support")
	errSigPubKeyNotMatch = errors.New("count of signature and public key not match")
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

type PDUCrypto struct {
}

func New() *PDUCrypto {
	return &PDUCrypto{}
}

func getKey(priKey interface{}) (*ecdsa.PrivateKey, error) {
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
		return nil, errKeyTypeNotSupport
	}
	return pk, nil
}

func getPubKey(pubKey interface{}) (*ecdsa.PublicKey, error) {
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
		return nil, errKeyTypeNotSupport
	}
	return pk, nil
}

func (pc *PDUCrypto) Sign(hash []byte, priKey crypto.PrivateKey) (*crypto.Signature, error) {
	if priKey.Source != SourceName {
		return nil, errSourceNotMatch
	}
	switch priKey.SigType {
	case Signature2PublicKey:
		pk, err := getKey(priKey.PriKey)
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
			pk, err := getKey(item)
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
		return nil, errSigTypeNotSupport
	}
}

func (pc *PDUCrypto) Verify(hash []byte, sig crypto.Signature) (bool, error) {
	if sig.Source != SourceName {
		return false, errSourceNotMatch
	}
	switch sig.SigType {
	case Signature2PublicKey:
		pk, err := getPubKey(sig.PubKey)
		if err != nil {
			return false, err
		}
		r := new(big.Int).SetBytes(sig.Signature[:32])
		s := new(big.Int).SetBytes(sig.Signature[32:])
		return ecdsa.Verify(pk, hash, r, s), nil
	case MultipleSignatures:
		pks := sig.PubKey.([]interface{})
		if len(pks) != len(sig.Signature)/64 {
			return false, errSigPubKeyNotMatch
		}
		for i, pubkey := range pks {
			pk, err := getPubKey(pubkey)
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
		return false, errSigTypeNotSupport
	}
}
