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

package types

import (
	"math/big"
	"testing"
)

func TestNewMessage(t *testing.T) {
	// do not check signature right now
	sig := MsgSig{
		R: big.NewInt(1),
		S: big.NewInt(2),
	}

	body1 := MsgBody{
		Title:    "title1",
		Category: 100,
		Nonce:    big.NewInt(0),
		Author:   "author1",
	}
	// build
	_, key1, err := RootMessage(body1, sig)
	if err != nil {
		t.Errorf("create root msg fail, err : %s", err)
	}

	body2 := MsgBody{
		Title:    "title2",
		Category: 3,
		Nonce:    big.NewInt(1),
		Author:   "author2",
	}
	_, key2, err := RootMessage(body2, sig)
	if err != nil {
		t.Errorf("create root msg fail, err : %s", err)
	}

	body3 := MsgBody{
		Title:    "title3",
		Category: 3,
		Nonce:    big.NewInt(1),
		Author:   "author3",
	}

	ref1 := MsgRef{hash: key1}
	ref2 := MsgRef{hash: key2}
	_, _, err = NewMessage(body3, sig, ref1, ref2)
	if err != nil {
		t.Errorf("create new msg fail, err : %s", err)
	}

}
