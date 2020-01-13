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

package common

import "fmt"

// Hash is fixed length []byte
type Hash [HashLength]byte

// HashLength is the length of Hash
const HashLength = 32

// Hash2String is transform Hash to string
func Hash2String(h Hash) (s string) {
	s = ""
	for i := 0; i < HashLength; i++ {
		s += fmt.Sprintf("%02X", h[i])
	}
	return s
}

// Bytes2String is transform []byte to string
func Bytes2String(b []byte) (s string) {
	s = ""
	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02X", b[i])
	}
	return s
}

// Bytes2Hash is transform []byte to Hash
func Bytes2Hash(b []byte) Hash {
	hash := [HashLength]byte{}
	if len(b) > HashLength {
		b = b[len(b)-HashLength:]
	}
	copy(hash[HashLength-len(b):], b)
	return hash
}
