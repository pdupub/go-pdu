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
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/pdupub/go-pdu/crypto"
)

func tGenKey() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {

	privKey, pubKey, err := genKey()
	if err != nil {
		return nil, nil, err
	}
	return privKey.(*ecdsa.PrivateKey), pubKey.(*ecdsa.PublicKey), nil
}

func TestS2PKSignature(t *testing.T) {
	E := New()
	pk, _, err := tGenKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}

	content := "hello world"
	sig1, err := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PriKey: pk.D})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	sig2, err := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PriKey: pk.D.Bytes()})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	sig3, err := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PriKey: pk})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	if sig1.Source != sig2.Source || sig1.Source != sig3.Source || sig1.Source != crypto.PDU {
		t.Errorf("signature source should be %s", crypto.PDU)
	}

	if sig1.SigType != sig2.SigType || sig1.SigType != sig3.SigType || sig1.SigType != crypto.Signature2PublicKey {
		t.Errorf("signature type should be %s", crypto.Signature2PublicKey)
	}

}

func TestS2PKVerify(t *testing.T) {
	E := New()

	pk, _, err := tGenKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}

	content := "hello world"
	sig, _ := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PriKey: pk.D})

	verify1, err := E.Verify([]byte(content), &crypto.Signature{
		PublicKey: crypto.PublicKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PubKey: pk.PublicKey}, Signature: sig.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify1 {
		t.Errorf("verify fail")
	}

	verify2, err := E.Verify([]byte(content), &crypto.Signature{PublicKey: crypto.PublicKey{Source: crypto.PDU,
		SigType: crypto.Signature2PublicKey, PubKey: &pk.PublicKey}, Signature: sig.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify2 {
		t.Errorf("verify fail")
	}

	pubKeyBytes := fromECDSAPub(&pk.PublicKey)
	verify3, err := E.Verify([]byte(content), &crypto.Signature{PublicKey: crypto.PublicKey{Source: crypto.PDU,
		SigType: crypto.Signature2PublicKey, PubKey: pubKeyBytes}, Signature: sig.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify3 {
		t.Errorf("verify fail")
	}

}

func TestMSSignature(t *testing.T) {
	E := New()

	pk1, _, err := tGenKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}
	pk2, _, err := tGenKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}
	pk3, _, err := tGenKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}

	content := "hello world"
	var pks []interface{}
	sig1, err := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.MultipleSignatures, PriKey: append(pks, pk1.D, pk2.D, pk3.D)})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}
	pks = []interface{}{}
	sig2, err := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.MultipleSignatures, PriKey: append(pks, pk1.D.Bytes(), pk2.D.Bytes(), pk3.D.Bytes())})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}
	pks = []interface{}{}
	sig3, err := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.MultipleSignatures, PriKey: append(pks, pk1, pk2, pk3)})
	if err != nil {
		t.Errorf("sign fail, err : %s", err)
	}

	if sig1.Source != sig2.Source || sig1.Source != sig3.Source || sig1.Source != crypto.PDU {
		t.Errorf("signature source should be %s", crypto.PDU)
	}

	if sig1.SigType != sig2.SigType || sig1.SigType != sig3.SigType || sig1.SigType != crypto.MultipleSignatures {
		t.Errorf("signature type should be %s", crypto.MultipleSignatures)
	}
}

func TestMSVerify(t *testing.T) {
	E := New()

	pk1, _, _ := tGenKey()
	pk2, _, _ := tGenKey()
	pk3, _, _ := tGenKey()
	content := "hello world"
	var pks []interface{}
	sig, _ := E.Sign([]byte(content), &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.MultipleSignatures, PriKey: append(pks, pk1.D, pk2.D, pk3.D)})

	var pubks []interface{}
	verify1, err := E.Verify([]byte(content), &crypto.Signature{PublicKey: crypto.PublicKey{Source: crypto.PDU,
		SigType: crypto.MultipleSignatures, PubKey: append(pubks, pk1.PublicKey, pk2.PublicKey, pk3.PublicKey)}, Signature: sig.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify1 {
		t.Errorf("verify fail")
	}

	pubks = []interface{}{}
	verify2, err := E.Verify([]byte(content), &crypto.Signature{PublicKey: crypto.PublicKey{Source: crypto.PDU,
		SigType: crypto.MultipleSignatures, PubKey: append(pubks, &pk1.PublicKey, &pk2.PublicKey, &pk3.PublicKey)}, Signature: sig.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify2 {
		t.Errorf("verify fail")
	}

	pubks = []interface{}{}
	pubKeyBytes1 := fromECDSAPub(&pk1.PublicKey)
	pubKeyBytes2 := fromECDSAPub(&pk2.PublicKey)
	pubKeyBytes3 := fromECDSAPub(&pk3.PublicKey)
	verify3, err := E.Verify([]byte(content), &crypto.Signature{PublicKey: crypto.PublicKey{Source: crypto.PDU,
		SigType: crypto.MultipleSignatures, PubKey: append(pubks, pubKeyBytes1, pubKeyBytes2, pubKeyBytes3)}, Signature: sig.Signature})
	if err != nil {
		t.Errorf("verify fail, err : %s", err)
	}
	if !verify3 {
		t.Errorf("verify fail")
	}

	pubks = []interface{}{}
	pubKeyBytes1 = fromECDSAPub(&pk1.PublicKey)
	pubKeyBytes3 = fromECDSAPub(&pk3.PublicKey)
	_, err = E.Verify([]byte(content), &crypto.Signature{PublicKey: crypto.PublicKey{Source: crypto.PDU,
		SigType: crypto.MultipleSignatures, PubKey: append(pubks, pubKeyBytes1, pubKeyBytes3)}, Signature: sig.Signature})
	if err != crypto.ErrSigPubKeyNotMatch {
		t.Errorf("verify should fail with err : %s", crypto.ErrSigPubKeyNotMatch)
	}

	pubks = []interface{}{}
	pubKeyBytes1 = fromECDSAPub(&pk1.PublicKey)
	pubKeyBytes2 = fromECDSAPub(&pk2.PublicKey)
	pubKeyBytes3 = fromECDSAPub(&pk3.PublicKey)
	verify4, err := E.Verify([]byte(content), &crypto.Signature{PublicKey: crypto.PublicKey{Source: crypto.PDU,
		SigType: crypto.MultipleSignatures, PubKey: append(pubks, pubKeyBytes1, pubKeyBytes3, pubKeyBytes2)}, Signature: sig.Signature})
	if err != nil {
		t.Errorf("verify should fail with no err : %s", err)
	}
	if verify4 {
		t.Errorf("verify should fail")
	}

}

func TestParsePriKey(t *testing.T) {
	E := New()

	priKey, _, err := E.GenKey(crypto.Signature2PublicKey)
	if err != nil {
		t.Error(err)
	}
	pkTarget := priKey.PriKey.(*ecdsa.PrivateKey)

	if pk, err := parsePriKey(pkTarget); err != nil {
		t.Error(err)
	} else if pk.D.Cmp(pkTarget.D) != 0 {
		t.Error("private key not equal")
	}

	if pk, err := parsePriKey(*pkTarget); err != nil {
		t.Error(err)
	} else if pk.D.Cmp(pkTarget.D) != 0 {
		t.Error("private key not equal")
	}

	if pk, err := parsePriKey(pkTarget.D); err != nil {
		t.Error(err)
	} else if pk.D.Cmp(pkTarget.D) != 0 {
		t.Error("private key not equal")
	}

	if pk, err := parsePriKey(pkTarget.D.Bytes()); err != nil {
		t.Error(err)
	} else if pk.D.Cmp(pkTarget.D) != 0 {
		t.Error("private key not equal")
	}

}

func TestParsePubKey(t *testing.T) {
	E := New()

	priKey, _, err := E.GenKey(crypto.Signature2PublicKey)
	if err != nil {
		t.Error(err)
	}
	pkTarget := priKey.PriKey.(*ecdsa.PrivateKey).PublicKey

	if pk, err := parsePubKey(pkTarget); err != nil {
		t.Error(err)
	} else if pk.X.Cmp(pkTarget.X) != 0 || pk.Y.Cmp(pkTarget.Y) != 0 {
		t.Error("private key not equal")
	}

	if pk, err := parsePubKey(&pkTarget); err != nil {
		t.Error(err)
	} else if pk.X.Cmp(pkTarget.X) != 0 || pk.Y.Cmp(pkTarget.Y) != 0 {
		t.Error("private key not equal")
	}

	if pk, err := parsePubKey(fromECDSAPub(&pkTarget)); err != nil {
		t.Error(err)
	} else if pk.X.Cmp(pkTarget.X) != 0 || pk.Y.Cmp(pkTarget.Y) != 0 {
		t.Error("private key not equal")
	}

	if pk, err := parsePubKey(new(big.Int).SetBytes(fromECDSAPub(&pkTarget))); err != nil {
		t.Error(err)
	} else if pk.X.Cmp(pkTarget.X) != 0 || pk.Y.Cmp(pkTarget.Y) != 0 {
		t.Error("private key not equal")
	}

}

func TestSign(t *testing.T) {
	E := New()

	priKey, _, err := E.GenKey(crypto.Signature2PublicKey)
	if err != nil {
		t.Error(err)
	}
	content := "hello world, 今天天气不错 🍊"
	if signature, err := E.Sign([]byte(content), priKey); err != nil {
		t.Error(err)
	} else {
		if v, err := E.Verify([]byte(content), signature); err != nil {
			t.Error(err)
		} else if v == false {
			t.Error("verify fail")
		}
	}

	priKey2, _, err := E.GenKey(crypto.MultipleSignatures, 3)
	if err != nil {
		t.Error(err)
	}

	if signature, err := E.Sign([]byte(content), priKey2); err != nil {
		t.Error(err)
	} else {
		if v, err := E.Verify([]byte(content), signature); err != nil {
			t.Error(err)
		} else if v == false {
			t.Error("verify fail")
		}
	}

}

func TestEEngine_EncryptKey(t *testing.T) {
	E := New()

	pk, _, err := tGenKey()
	if err != nil {
		t.Error("generate key pair fail", err)
	}

	privateKey := &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.Signature2PublicKey, PriKey: pk}

	keyJSON, err := E.EncryptKey(privateKey, "123")
	if err != nil {
		t.Error("encrypt key fail", err)
	}

	newPrivateKey, newPublicKey, err := E.DecryptKey(keyJSON, "123")
	if err != nil {
		t.Error("decrypt key fail")
	}

	if privateKey.Source != newPrivateKey.Source {
		t.Error("source not equal")
	}
	if privateKey.SigType != newPrivateKey.SigType {
		t.Error("sig type not equal")
	}

	if privateKey.PriKey.(*ecdsa.PrivateKey).D.Cmp(newPrivateKey.PriKey.(*ecdsa.PrivateKey).D) != 0 {
		t.Error("private key not equal")
	}

	if pk.X.Cmp(newPublicKey.PubKey.(*ecdsa.PublicKey).X) != 0 ||
		pk.Y.Cmp(newPublicKey.PubKey.(*ecdsa.PublicKey).Y) != 0 {
		t.Error("public key not equal")
	}

}

func TestEEngine_EncryptKeyMS(t *testing.T) {
	E := New()

	pk1, _, err := tGenKey()
	if err != nil {
		t.Error("generate key pair fail", err)
	}
	pk2, _, err := tGenKey()
	if err != nil {
		t.Error("generate key pair fail", err)
	}
	pk3, _, err := tGenKey()
	if err != nil {
		t.Error("generate key pair fail", err)
	}
	var pks []interface{}
	pks = append(pks, pk1, pk2, pk3)
	privateKey := &crypto.PrivateKey{Source: crypto.PDU, SigType: crypto.MultipleSignatures, PriKey: pks}

	keyJSON, err := E.EncryptKey(privateKey, "123")
	if err != nil {
		t.Error("encrypt key fail", err)
	}

	newPrivateKey, _, err := E.DecryptKey(keyJSON, "123")
	if err != nil {
		t.Error("decrypt key fail")
	}

	if privateKey.Source != newPrivateKey.Source {
		t.Error("source not equal")
	}
	if privateKey.SigType != newPrivateKey.SigType {
		t.Error("sig type not equal")
	}

	for k, item := range newPrivateKey.PriKey.([]interface{}) {
		v := item.(*ecdsa.PrivateKey)
		if v.D.Cmp(pks[k].(*ecdsa.PrivateKey).D) != 0 {
			t.Error("private key not equal")
		}
	}

}

func TestMarshalPrivateKey(t *testing.T) {
	E := New()

	privKey, _, err := E.GenKey(crypto.Signature2PublicKey)
	if err != nil {
		t.Error(err)
	}

	privKeyBytes, err := E.marshalPrivKey(privKey)
	if err != nil {
		t.Error(err)
	}

	uPrivKey, err := E.unmarshalPrivKey(privKeyBytes)
	if err != nil {
		t.Error(err)
	}

	if privKey.Source != uPrivKey.Source {
		t.Error("source not match")
	}
	if privKey.SigType != uPrivKey.SigType {
		t.Error("sig type not match")
	}
	if privKey.PriKey.(*ecdsa.PrivateKey).D.Cmp(uPrivKey.PriKey.(*ecdsa.PrivateKey).D) != 0 {
		t.Error("private key not match")
	}

	size := 3
	privKey, _, err = E.GenKey(crypto.MultipleSignatures, size)
	if err != nil {
		t.Error(err)
	}

	privKeyBytes, err = E.marshalPrivKey(privKey)
	if err != nil {
		t.Error(err)
	}

	uPrivKey, err = E.unmarshalPrivKey(privKeyBytes)
	if err != nil {
		t.Error(err)
	}

	if privKey.Source != uPrivKey.Source {
		t.Error("source not match")
	}
	if privKey.SigType != uPrivKey.SigType {
		t.Error("sig type not match")
	}

	for i := 0; i < size; i++ {
		if privKey.PriKey.([]interface{})[i].(*ecdsa.PrivateKey).D.Cmp(uPrivKey.PriKey.([]interface{})[i].(*ecdsa.PrivateKey).D) != 0 {
			t.Error("private key not match")
		}
	}
}

func TestMarshalPublicKey(t *testing.T) {
	E := New()
	size := 3
	_, pubKey, _ := E.GenKey(crypto.Signature2PublicKey)

	pubKeyBytes, err := E.marshalPubKey(pubKey)
	if err != nil {
		t.Error(err)
	}

	uPubKey, err := E.unmarshalPubKey(pubKeyBytes)
	if err != nil {
		t.Error(err)
	}

	if pubKey.Source != uPubKey.Source {
		t.Error("source not match")
	}
	if pubKey.SigType != uPubKey.SigType {
		t.Error("sig type not match")
	}

	if pubKey.PubKey.(*ecdsa.PublicKey).X.Cmp(uPubKey.PubKey.(*ecdsa.PublicKey).X) != 0 ||
		pubKey.PubKey.(*ecdsa.PublicKey).Y.Cmp(uPubKey.PubKey.(*ecdsa.PublicKey).Y) != 0 {
		t.Error("private key not match")
	}

	_, pubKey, err = E.GenKey(crypto.MultipleSignatures, size)
	if err != nil {
		t.Error(err)
	}

	pubKeyBytes, err = E.marshalPubKey(pubKey)
	if err != nil {
		t.Error(err)
	}

	uPubKey, err = E.unmarshalPubKey(pubKeyBytes)
	if err != nil {
		t.Error(err)
	}

	if pubKey.Source != uPubKey.Source {
		t.Error("source not match")
	}
	if pubKey.SigType != uPubKey.SigType {
		t.Error("sig type not match")
	}

	for i := 0; i < size; i++ {
		if pubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).X.Cmp(uPubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).X) != 0 {
			t.Error("public key not match")
		}
		if pubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).Y.Cmp(uPubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).Y) != 0 {
			t.Error("public key not match")
		}
	}
}

func TestMarshal(t *testing.T) {
	E := New()
	size := 3
	privKey, pubKey, err := E.GenKey(crypto.MultipleSignatures, size)
	if err != nil {
		t.Error(err)
	}

	privKeyBytes, pubKeyBytes, err := E.Marshal(privKey, pubKey)
	if err != nil {
		t.Error(err)
	}

	uPrivKey, uPubKey, err := E.Unmarshal(privKeyBytes, pubKeyBytes)
	if err != nil {
		t.Error(err)
	}

	if pubKey.Source != uPubKey.Source {
		t.Error("source not match")
	}
	if pubKey.SigType != uPubKey.SigType {
		t.Error("sig type not match")
	}

	for i := 0; i < size; i++ {
		if pubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).X.Cmp(uPubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).X) != 0 {
			t.Error("public key not match")
		}
		if pubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).Y.Cmp(uPubKey.PubKey.([]interface{})[i].(*ecdsa.PublicKey).Y) != 0 {
			t.Error("public key not match")
		}
	}

	for i := 0; i < size; i++ {
		if privKey.PriKey.([]interface{})[i].(*ecdsa.PrivateKey).D.Cmp(uPrivKey.PriKey.([]interface{})[i].(*ecdsa.PrivateKey).D) != 0 {
			t.Error("private key not match")
		}
	}

}
