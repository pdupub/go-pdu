// Copyright 2019 The TTC Authors
// This file is part of the TTC library.
//
// The TTC library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The TTC library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the TTC library. If not, see <http://www.gnu.org/licenses/>.

package crypto

import "math/big"

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

// Crypto is crypto algorithm interface for different block chain
type Crypto interface {
	Create() (*KeyPair, error)
	Sign(priv []byte, hash []byte) (*big.Int, *big.Int, error)
	Verify(pub []byte, hash []byte, r *big.Int, s *big.Int) bool
}
