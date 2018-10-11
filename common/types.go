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

package common

import (
	"fmt"
	"math/big"

	"github.com/TATAUFO/PDU/common/hexutil"
)

const (
	HashLength    = 32 // HashLength is ...
	AddressLength = 20 // AddressLength is ...
)

// Hash is fixed length []byte
type Hash [HashLength]byte

// BytesToHash convert []byte to Hash
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash convert *big.Int to Hash
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash convert string to Hash
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Str get the string representation of the underlying hash
func (h Hash) Str() string { return string(h[:]) }

// Bytes get []byte representation of the underlying hash
func (h Hash) Bytes() []byte { return h[:] }

// Big get big.Int representation of the underlying hash
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex get the Hex string representation of the underlying hash
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}
	copy(h[HashLength-len(b):], b)
}

// SetString set string `s` to h. If s is larger than len(h) s will be cropped (from left) to fit.
func (h *Hash) SetString(s string) { h.SetBytes([]byte(s)) }

// Set h to other
func (h *Hash) Set(other Hash) {
	for i, v := range other {
		h[i] = v
	}
}

// Address is fixed length []byte
type Address [AddressLength]byte

// BytesToAddress
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// BigToAddress
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsHexAddress
func IsHexAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Str get the string representation of the underlying address
func (a Address) Str() string { return string(a[:]) }

// Bytes get []byte ...
func (a Address) Bytes() []byte { return a[:] }

// Big get big.Int ...
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash get Hash ...
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex get Hex string ...
func (a Address) Hex() string { return hexutil.Encode(a[:]) }

// Format address string
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b. If b is larger than len(a) it will panic
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// SetString set string `s` to a. If s is larger than len(a) it will panic
func (a *Address) SetString(s string) { a.SetBytes([]byte(s)) }

// Set a to other
func (a *Address) Set(other Address) {
	for i, v := range other {
		a[i] = v
	}
}
