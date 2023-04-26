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
	"strconv"
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestNewSpecies(t *testing.T) {

	creator, _ := identity.New()
	creator.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)

	did, _ := identity.New()
	did.UnlockWallet("../"+params.TestKeystore(1), params.TestPassword)

	var contents []*QContent

	note, _ := NewContent(QCFmtStringTEXT, []byte("Coder"))
	contents = append(contents, note)

	k1, _ := NewContent(QCFmtStringInt, []byte(strconv.Itoa(1)))
	contents = append(contents, k1)

	k2, _ := NewContent(QCFmtStringInt, []byte(strconv.Itoa(2)))
	contents = append(contents, k2)

	contents = append(contents, &QContent{Format: QCFmtStringAddressHex, Data: []byte(did.GetAddress().Hex())})

	baseSig := Sig([]byte("0x070d15041083041b48d0f2297357ce59ad18f6c608d70a1e6e04bcf494e366db"))

	newQuantum, _ := NewQuantum(QuantumTypeSpecies, contents, baseSig)
	if err := newQuantum.Sign(creator); err != nil {
		t.Error(err)
	}
}
