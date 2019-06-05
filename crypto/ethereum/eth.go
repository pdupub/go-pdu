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

package ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	privKeyBytesLen = 32
	pubKeyBytesLen  = 64
)

type ETH struct {
}

func New() *ETH {
	return &ETH{}
}

func (eth *ETH) CrateNewKeyPair() ([]byte, []byte, error) {
	private, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	d := private.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	privateKey := paddedAppend(privKeyBytesLen, b, d)
	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return privateKey, publicKey, nil
}

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
