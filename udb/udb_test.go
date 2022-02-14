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

package udb

import (
	"encoding/json"
	"testing"

	"github.com/pdupub/go-pdu/newv/core"
)

// url & api token
const (
	testUrl   = "https://blue-surf-570228.us-east-1.aws.cloud.dgraph.io/graphql"
	testToken = "OTNiNjMyMTM1ODU3ZWE2NWU3YWQyNmM2MzgwZmViODg="
)

func TestUDB(t *testing.T) {
	udb, err := New(testUrl, testToken)
	if err != nil {
		t.Error(err)
	}

	defer udb.Close()

	if err := udb.initSchema(); err != nil {
		t.Error(err)
	}

	if err := udb.dropData(); err != nil {
		t.Error(err)
	}

	sig1 := []byte("0x001")
	sig2 := []byte("0x002")
	sig3 := []byte("0x003")
	sender := "0xa01"

	// test1 : quantum with empty refs
	quantum := &core.Quantum{
		Signature: sig3,
		UnsignedQuantum: core.UnsignedQuantum{
			Type:       core.QuantumTypeInfo,
			References: []core.Sig{sig1, sig2},
			Contents: []*core.QContent{{
				Format: core.QCFmtStringTEXT,
				Data:   []byte("Hello World!"),
			}},
		},
	}

	if err := sgcheck(udb, t, quantum, sender); err != nil {
		t.Error(err)
	}

	// test2 : quantum fill empty quantum
	quantum = &core.Quantum{
		Signature: sig2,
		UnsignedQuantum: core.UnsignedQuantum{
			Type:       core.QuantumTypeProfile,
			References: []core.Sig{sig1},
			Contents: []*core.QContent{{
				Format: core.QCFmtStringJSON,
				Data:   []byte("{\"nickname\":[\"Hello\",\"World\"]}"),
			}},
		},
	}

	if err := sgcheck(udb, t, quantum, sender); err != nil {
		t.Error(err)
	}

}

func sgcheck(udb *UDB, t *testing.T, quantum *core.Quantum, sender string) error {

	if uid, err := udb.SetQuantum(quantum, sender); err != nil {
		return err
	} else {
		res, _ := json.Marshal(quantum)
		t.Log("SetQuantum", "uid", uid, "sender", sender, "quantum", string(res))
	}

	if newQ, dbQ, err := udb.GetQuantum(quantum.Signature); err != nil {
		return err
	} else {
		res, _ := json.Marshal(newQ)
		t.Log("GetQuantum", "uid", dbQ.UID, "sender", dbQ.Sender.Address, "quantum", string(res))
	}
	return nil
}
