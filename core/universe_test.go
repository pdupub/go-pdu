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
	"github.com/pdupub/go-pdu/udb"
	"github.com/pdupub/go-pdu/udb/dgraph"
)

func TestNewUniverse(t *testing.T) {
	t.Log("new universe")
	sender, _ := identity.New()
	sender.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)

	var db udb.UDB
	cdgd, err := dgraph.New(params.TestDGraphURL, params.TestDGraphToken)
	if err != nil {
		t.Error(err)
	}
	db = cdgd

	defer db.Close()

	if err := db.SetSchema(); err != nil {
		t.Error(err)
	}

	if err := db.DropData(); err != nil {
		t.Error(err)
	}

	universe, err := NewUniverse(db)
	if err != nil {
		t.Error(err)
	}

	// build and receive first message
	qc1, err := NewContent(QCFmtStringTEXT, []byte("Hello World!"))
	if err != nil {
		t.Error(err)
	}
	qc2, err := NewContent(QCFmtStringTEXT, []byte("Second Content"))
	if err != nil {
		t.Error(err)
	}
	quantum1, err := NewQuantum(QuantumTypeInfo, []*QContent{qc1, qc2})
	if err != nil {
		t.Error(err)
	}

	if err := quantum1.Sign(sender); err != nil {
		t.Error(err)
	}
	if err := universe.RecvQuantum(quantum1); err != nil {
		t.Error(err)
	}

	// load first message from database, check if equal
	quantum2 := universe.GetQuantum(quantum1.Signature)
	if quantum2 == nil {
		t.Error("get quantum2  fail")
	}

	if len(quantum1.References) != len(quantum2.References) {
		t.Error("references length not match")
	}
	for i := 0; i < len(quantum1.References); i++ {
		if Sig2Hex(quantum1.References[i]) != Sig2Hex(quantum2.References[i]) {
			t.Error("reference", i, "not match")
		}
	}

	individual := universe.GetIndividual(sender.GetAddress())
	showIndividual(t, individual)

	if individual.Address.Hex() != sender.GetAddress().Hex() {
		t.Error("individual not match")
	}

	// build and receive second message
	qc3, err := NewContent(QCFmtStringTEXT, []byte("Second Message"))
	if err != nil {
		t.Error(err)
	}
	quantum3, err := NewQuantum(QuantumTypeInfo, []*QContent{qc3}, quantum1.Signature)
	if err != nil {
		t.Error(err)
	}
	if err := quantum3.Sign(sender); err != nil {
		t.Error(err)
	}
	if err := universe.RecvQuantum(quantum3); err != nil {
		t.Error(err)
	}
	individual = universe.GetIndividual(sender.GetAddress())
	showIndividual(t, individual)

	// build and receive profile
	qc4, _ := NewContent(QCFmtStringTEXT, []byte("nickname"))
	qc5, _ := NewContent(QCFmtStringTEXT, []byte("PDU"))

	quantum4, err := NewQuantum(QuantumTypeProfile, []*QContent{qc4, qc5}, quantum3.Signature, quantum1.Signature)
	if err != nil {
		t.Error(err)
	}
	if err := quantum4.Sign(sender); err != nil {
		t.Error(err)
	}
	if err := universe.RecvQuantum(quantum4); err != nil {
		t.Error(err)
	}
	individual = universe.GetIndividual(sender.GetAddress())
	showIndividual(t, individual)

	// query quantums
	quantumQueryResult := universe.QueryQuantum(sender.GetAddress(), 0, 1, 10, false)
	for _, v := range quantumQueryResult {
		t.Log(Sig2Hex(v.Signature))
	}
}

func showIndividual(t *testing.T, individual *Individual) {
	t.Log("individual")
	t.Log("address", individual.Address.Hex())
	t.Log("quantums:")
	for _, v := range individual.Quantums {
		t.Log("uid", Sig2Hex(v.Signature))
	}
	t.Log("profile:")
	for k, v := range individual.Profile {
		t.Log("key", k, "value", v)
	}
	t.Log("communities")
	for _, v := range individual.Communities {
		t.Log("note", v.Note)
	}
}
