// Copyright 2021 The PDU Authors
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

package msg

import (
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestMessage(t *testing.T) {
	did, _ := identity.New()
	did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)

	refs := [][]byte{[]byte("0x070d15041083041b48d0f2297357ce59ad18f6c608d70a1e6e04bcf494e366db"),
		[]byte("0x08fd3282eecbf25d31a9a5e51ed2d79a806f14281fbb583a5ee4024589b959d9")}
	content := did.GetKey().Address.Hash().Bytes()

	testMsg := New(content, refs...)

	sig, err := testMsg.Sign(did)
	if err != nil {
		t.Error(err)
	}
	err = testMsg.Verify(sig, did.GetKey().Address)
	if err != nil {
		t.Error(err)
	}

}
