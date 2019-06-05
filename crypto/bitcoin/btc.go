// Copyright 2018 The PDU Authors
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
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	privKeyBytesLen = 32
	pubKeyBytesLen  = 64
)

type BTC struct {
}

func New() *BTC {
	return &BTC{}
}

func (btc *BTC) CrateNewKeyPair() ([]byte, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	d := private.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	privateKey := paddedAppend(privKeyBytesLen, b, d)
	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return privateKey, publicKey, nil
}

func (btc *BTC) Sign(pri []byte, hash []byte) (*big.Int, *big.Int, error) {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(pri)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())
	return ecdsa.Sign(rand.Reader, privateKey, hash[:])
}

func (btc *BTC) Verify(pub []byte, hash []byte, r *big.Int, s *big.Int) bool {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.PublicKey.X, privateKey.PublicKey.Y = new(big.Int).SetBytes(pub[:pubKeyBytesLen/2]), new(big.Int).SetBytes(pub[pubKeyBytesLen/2:])
	return ecdsa.Verify(&privateKey.PublicKey, hash[:], r, s)
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

func byteString(b []byte) (s string) {
	s = ""
	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02X", b[i])
	}
	return s
}
