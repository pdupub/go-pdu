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
	"encoding/json"
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestInfoQuantum(t *testing.T) {
	did, _ := identity.New()
	did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)

	c0 := NewTextC("Hello!")
	c1 := NewIntC(100)
	c2 := NewTextC(">")
	c3 := NewFloatC(99.9)

	q, err := NewQuantum(QuantumTypeInfo, []*QContent{c0, c1, c2, c3}, FirstQuantumReference)
	if err != nil {
		t.Error(err)
	}
	q.Sign(did)
	if j, err := json.Marshal(q); err != nil {
		t.Error(err)
	} else {
		t.Log("###########################################################")
		t.Log("signed q is:")
		t.Log(string(j))
		t.Log("hex signature is:", len(q.Signature))
		t.Log(Sig2Hex(q.Signature))
		t.Log("[]byte(sig) is:")
		t.Log(Hex2Sig(string([]byte(Sig2Hex(q.Signature)))))
		t.Log("q.Signature is:")
		t.Log(q.Signature)
		t.Log("###########################################################")
	}

	if d, err := q.Contents[0].GetData(); err != nil || d.(string) != "Hello!" {
		t.Error(err)
	}
	if d, err := q.Contents[1].GetData(); err != nil || d.(int) != 100 {
		t.Error(err)
	}
	if d, err := q.Contents[2].GetData(); err != nil || d.(string) != ">" {
		t.Error(err)
	}
	if d, err := q.Contents[3].GetData(); err != nil || d.(float64) != 99.9 {
		t.Error(err)
	}

}

func TestSignAndVerify(t *testing.T) {
	did, _ := identity.New()
	did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)

	refs := []Sig{[]byte("0x070d15041083041b48d0f2297357ce59ad18f6c608d70a1e6e04bcf494e366db"),
		[]byte("0x08fd3282eecbf25d31a9a5e51ed2d79a806f14281fbb583a5ee4024589b959d9")}

	qc, err := NewContent(QCFmtStringTEXT, []byte("Hello World!"))
	if err != nil {
		t.Error(err)
	}

	q, err := NewQuantum(QuantumTypeInfo, []*QContent{qc}, refs...)
	if err != nil {
		t.Error(err)
	}

	if err := q.Sign(did); err != nil {
		t.Error(err)
	}

	addr, err := q.Ecrecover()
	if err != nil {
		t.Error(err)
	}
	if addr != did.GetAddress() {
		t.Error("address not match")
	}
	t.Log("ecrecover address", addr.Hex())

	if jsonUnsignedQuantumBytes, err := json.Marshal(q.UnsignedQuantum); err != nil {
		t.Error(err)
	} else {
		t.Log("unsigned quantum json", string(jsonUnsignedQuantumBytes))
	}

	if jsonQuantumBytes, err := json.Marshal(q); err != nil {
		t.Error(err)
	} else {
		t.Log("quantum json", string(jsonQuantumBytes))
	}
}
