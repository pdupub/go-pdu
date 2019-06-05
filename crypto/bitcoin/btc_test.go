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

import "testing"

func TestBTC(t *testing.T) {

	btc := New()
	priv, pub, err := btc.CrateNewKeyPair()
	if err != nil {
		t.Errorf("create key pari fail, err : %s", err)
	}

	content := "hello world"
	r, s, err := btc.Sign(priv, []byte(content))
	if err != nil {
		t.Errorf("sign content fail, err : %s", err)
	}

	v := btc.Verify(pub, []byte(content), r, s)
	if !v {
		t.Errorf("verify fail")
	}
}
