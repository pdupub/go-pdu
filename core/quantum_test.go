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
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
	"golang.org/x/crypto/sha3"
)

func TestUtils(t *testing.T) {
	msg := "hello123"
	t.Log("message ", msg)
	t.Log("msg to bytes", []byte(msg))

	d := sha3.NewLegacyKeccak256()
	d.Write([]byte(msg))
	hashBytes := d.Sum(nil)
	t.Log("hash bytes", hashBytes)
	t.Log("hash hex", hex.EncodeToString(hashBytes))
}

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

	refs := []Sig{Hex2Sig("0x3e34d7ba1ed979e0b5c5cb0507837554bf16798e5bde0b4b55df1f33a1e12fa22dff93ae2f2b362eabe866e6394d56cff4625f22efab18540769e827053a574c00"),
		Hex2Sig("0x8bcb61fd8d0e280b7eaa92de0821bfedf62795ddfb1d8f6f6cf6ee6fb7974dd947fa9c73cf3ef1ac47003da28ab3f4c44eb58a619920a4d6fe8604ff4aa10c4e00")}

	qc, err := NewContent(QCFmtStringTEXT, []byte("Hello World!"))
	if err != nil {
		t.Error(err)
	}
	t.Log("content fmt", QCFmtStringTEXT)
	t.Log("content data", string(qc.Data))

	q, err := NewQuantum(QuantumTypeInfo, []*QContent{qc}, refs...)
	if err != nil {
		t.Error(err)
	}

	b, err := json.Marshal(q.UnsignedQuantum)
	if err != nil {
		t.Error(err)
	}
	t.Log("unsigned q", string(b))
	t.Log("###", b)

	sig, err := did.Sign(b)
	if err != nil {
		t.Error(err)
	}
	t.Log("signature", Sig2Hex(sig))

	if err := q.Sign(did); err != nil {
		t.Error(err)
	}

	signedQ, err := json.Marshal(q)
	if err != nil {
		t.Error(err)
	}
	t.Log("signed q", string(signedQ))

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
