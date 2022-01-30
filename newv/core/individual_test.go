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

package core

import (
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestNewIndividual(t *testing.T) {
	did, _ := identity.New()
	did.UnlockWallet("../../"+params.TestKeystore(0), params.TestPassword)

	newInd := NewIndividual(did.GetAddress())

	if did.GetAddress() != newInd.GetAddress() {
		t.Error("address not match")
	}

	k1, _ := NewContent(QCFmtStringTEXT, []byte("nickname"))
	v1, _ := NewContent(QCFmtStringTEXT, []byte("pdu"))

	if err := newInd.UpsertProfile([]*QContent{k1, v1}); err != nil {
		t.Error(err)
	}

	for k := range newInd.Profile {
		t.Log(k)
	}
}
