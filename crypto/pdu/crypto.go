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

	MultipleSignatures = "MS"
	Signature2PublicKey = "S2PK"
)

var (
	errSourceNotMatch = errors.New("signature source not match")
	errTypeNotSupport = errors.New("signature type not support")
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

type PDUCrypto struct {

}


func (pc *PDUCrypto)Sign(hash []byte,priKey crypto.PrivateKey) (*crypto.Signature, error) {
	if priKey.Source != SourceName {
		return nil, errSourceNotMatch
	}
	if priKey.SigType != MultipleSignatures && priKey.SigType != Signature2PublicKey {
		return nil, errTypeNotSupport
	}
	switch priKey.SigType {
	case Signature2PublicKey:

		pk := new(ecdsa.PrivateKey)
		switch priKey.PriKey.(type){
		case *ecdsa.PrivateKey:
			pk = priKey.PriKey.(*ecdsa.PrivateKey)
		case []byte:
			pk.PublicKey.Curve = elliptic.P256()
			pk.D = new(big.Int).SetBytes(priKey.PriKey.([]byte))
			pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
		case *big.Int:
			pk.PublicKey.Curve = elliptic.P256()
			pk.D = new(big.Int).Set(priKey.PriKey.(*big.Int))
			pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())

		}
		r, s, err := ecdsa.Sign(rand.Reader, pk, hash[:])
		if err != nil {
			return nil, err
		}
		return &crypto.Signature{
			Source:SourceName,
			SigType:priKey.SigType,
			Signature: append(r.Bytes(),s.Bytes()...),
			PubKey:pk.PublicKey,
		},nil
	case MultipleSignatures:
		pks := priKey.PriKey.([]*ecdsa.PrivateKey)
		for _, pk := range pks {
			r, s, err := ecdsa.Sign(rand.Reader, pk, hash[:])
			if err != nil {
				return nil,err
			}
		}
	}
	return &crypto.Signature{},nil
}

func (pc *PDUCrypto)Verify(hash []byte, sig crypto.Signature) bool {
	return true
}


