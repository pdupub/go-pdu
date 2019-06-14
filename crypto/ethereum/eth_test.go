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

import "testing"

func TestETH(t *testing.T) {
	eth := New()
	priv, pub, err := eth.CrateNewKeyPair()
	if err != nil {
		t.Errorf("\nprivateKey : %s \npublicKey : %s \nerr : %v", byteString(priv), byteString(pub), err)
	}
}
