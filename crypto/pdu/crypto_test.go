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

package pdu

import (
	"github.com/pdupub/go-pdu/crypto"
	"testing"
)

func TestPDUCrypto_Sign(t *testing.T) {

	pk, err := GenerateKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}

	pdu := New()
	content1 := "hello world"
	sig1, err := pdu.Sign([]byte(content1), crypto.PrivateKey{Source: SourceName, SigType: Signature2PublicKey, PriKey: pk.D})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	sig2, err := pdu.Sign([]byte(content1), crypto.PrivateKey{Source: SourceName, SigType: Signature2PublicKey, PriKey: pk.D.Bytes()})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	sig3, err := pdu.Sign([]byte(content1), crypto.PrivateKey{Source: SourceName, SigType: Signature2PublicKey, PriKey: pk})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	if sig1.Source != sig2.Source || sig1.Source != sig3.Source || sig1.Source != SourceName {
		t.Errorf("signature source should be %s", SourceName)
	}

	if sig1.SigType != sig2.SigType || sig1.SigType != sig3.SigType || sig1.SigType != Signature2PublicKey {
		t.Errorf("signature type should be %s", Signature2PublicKey)
	}

	verify1, err := pdu.Verify([]byte(content1), crypto.Signature{Source: SourceName,
		SigType: Signature2PublicKey, Signature: sig1.Signature, PubKey: pk.PublicKey})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify1 {
		t.Errorf("verify fail")
	}

	verify2, err := pdu.Verify([]byte(content1), crypto.Signature{Source: SourceName,
		SigType: Signature2PublicKey, Signature: sig1.Signature, PubKey: &pk.PublicKey})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify2 {
		t.Errorf("verify fail")
	}

	pubKeyBytes := append(pk.PublicKey.X.Bytes(), pk.PublicKey.Y.Bytes()...)
	verify3, err := pdu.Verify([]byte(content1), crypto.Signature{Source: SourceName,
		SigType: Signature2PublicKey, Signature: sig1.Signature, PubKey: pubKeyBytes})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify3 {
		t.Errorf("verify fail")
	}

}
