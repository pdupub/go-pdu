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

func TestS2PK(t *testing.T) {

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

	verify1, err := pdu.Verify([]byte(content1), crypto.Signature{
		PublicKey: crypto.PublicKey{Source: SourceName, SigType: Signature2PublicKey, PubKey: pk.PublicKey}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify1 {
		t.Errorf("verify fail")
	}

	verify2, err := pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: Signature2PublicKey, PubKey: &pk.PublicKey}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify2 {
		t.Errorf("verify fail")
	}

	pubKeyBytes := append(pk.PublicKey.X.Bytes(), pk.PublicKey.Y.Bytes()...)
	verify3, err := pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: Signature2PublicKey, PubKey: pubKeyBytes}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify3 {
		t.Errorf("verify fail")
	}

	pubKeyBytes = append(pk.PublicKey.Y.Bytes(), pk.PublicKey.X.Bytes()...)
	verify4, err := pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: Signature2PublicKey, PubKey: pubKeyBytes}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if verify4 {
		t.Errorf("verify should fail")
	}

}

func TestMS(t *testing.T) {

	pk1, err := GenerateKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}
	pk2, err := GenerateKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}
	pk3, err := GenerateKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}

	pdu := New()
	content1 := "hello world"
	var pks []interface{}
	sig1, err := pdu.Sign([]byte(content1), crypto.PrivateKey{Source: SourceName, SigType: MultipleSignatures, PriKey: append(pks, pk1.D, pk2.D, pk3.D)})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}
	pks = []interface{}{}
	sig2, err := pdu.Sign([]byte(content1), crypto.PrivateKey{Source: SourceName, SigType: MultipleSignatures, PriKey: append(pks, pk1.D.Bytes(), pk2.D.Bytes(), pk3.D.Bytes())})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}
	pks = []interface{}{}
	sig3, err := pdu.Sign([]byte(content1), crypto.PrivateKey{Source: SourceName, SigType: MultipleSignatures, PriKey: append(pks, pk1, pk2, pk3)})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	if sig1.Source != sig2.Source || sig1.Source != sig3.Source || sig1.Source != SourceName {
		t.Errorf("signature source should be %s", SourceName)
	}

	if sig1.SigType != sig2.SigType || sig1.SigType != sig3.SigType || sig1.SigType != MultipleSignatures {
		t.Errorf("signature type should be %s", MultipleSignatures)
	}

	var pubks []interface{}
	verify1, err := pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: MultipleSignatures, PubKey: append(pubks, pk1.PublicKey, pk2.PublicKey, pk3.PublicKey)}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify1 {
		t.Errorf("verify fail")
	}

	pubks = []interface{}{}
	verify2, err := pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: MultipleSignatures, PubKey: append(pubks, &pk1.PublicKey, &pk2.PublicKey, &pk3.PublicKey)}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify2 {
		t.Errorf("verify fail")
	}

	pubks = []interface{}{}
	pubKeyBytes1 := append(pk1.PublicKey.X.Bytes(), pk1.PublicKey.Y.Bytes()...)
	pubKeyBytes2 := append(pk2.PublicKey.X.Bytes(), pk2.PublicKey.Y.Bytes()...)
	pubKeyBytes3 := append(pk3.PublicKey.X.Bytes(), pk3.PublicKey.Y.Bytes()...)
	verify3, err := pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: MultipleSignatures, PubKey: append(pubks, pubKeyBytes1, pubKeyBytes2, pubKeyBytes3)}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify3 {
		t.Errorf("verify fail")
	}

	pubks = []interface{}{}
	pubKeyBytes1 = append(pk1.PublicKey.X.Bytes(), pk1.PublicKey.Y.Bytes()...)
	pubKeyBytes3 = append(pk3.PublicKey.X.Bytes(), pk3.PublicKey.Y.Bytes()...)
	_, err = pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: MultipleSignatures, PubKey: append(pubks, pubKeyBytes1, pubKeyBytes3)}, Signature: sig1.Signature})
	if err != crypto.ErrSigPubKeyNotMatch {
		t.Errorf("verify should fail with err : %s", crypto.ErrSigPubKeyNotMatch)
	}

	pubks = []interface{}{}
	pubKeyBytes1 = append(pk1.PublicKey.X.Bytes(), pk1.PublicKey.Y.Bytes()...)
	pubKeyBytes2 = append(pk2.PublicKey.X.Bytes(), pk2.PublicKey.Y.Bytes()...)
	pubKeyBytes3 = append(pk3.PublicKey.X.Bytes(), pk3.PublicKey.Y.Bytes()...)
	verify4, err := pdu.Verify([]byte(content1), crypto.Signature{PublicKey: crypto.PublicKey{Source: SourceName,
		SigType: MultipleSignatures, PubKey: append(pubks, pubKeyBytes1, pubKeyBytes3, pubKeyBytes2)}, Signature: sig1.Signature})
	if err != nil {
		t.Errorf("verify should fail with no err : %s", err)
	}
	if verify4 {
		t.Errorf("verify should fail")
	}

}
