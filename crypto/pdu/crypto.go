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
	errPKTypeNotSupport  = errors.New("privatekey type not support")
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
		return nil, errPKTypeNotSupport
	}
	return pk, nil
}

func (pc *PDUCrypto) Sign(hash []byte, priKey crypto.PrivateKey) (*crypto.Signature, error) {
	if priKey.Source != SourceName {
		return nil, errSourceNotMatch
	}
	if priKey.SigType != MultipleSignatures && priKey.SigType != Signature2PublicKey {
		return nil, errSigTypeNotSupport
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
			Source:    SourceName,
			SigType:   priKey.SigType,
			Signature: append(r.Bytes(), s.Bytes()...),
			PubKey:    pk.PublicKey,
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
			Source:    SourceName,
			SigType:   priKey.SigType,
			Signature: signature,
			PubKey:    pubKeys,
		}, nil
	}
	return &crypto.Signature{}, nil
}

func (pc *PDUCrypto) Verify(hash []byte, sig crypto.Signature) bool {
	return true
}
