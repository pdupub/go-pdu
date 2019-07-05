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

package core

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/pdu"
	"testing"
)

const (
	retryCnt = 100
)

func TestCreateRootUsers(t *testing.T) {

	for i := 0; i < retryCnt; i++ {
		if pk, err := pdu.GenerateKey(); err != nil {
			t.Errorf("generate key pair fail, err : %s", err)
		} else {
			pubKey := crypto.PublicKey{Source: pdu.SourceName, SigType: pdu.Signature2PublicKey, PubKey: pk.PublicKey}
			if users, err := CreateRootUsers(pubKey); err != nil {
				t.Errorf("create root users fail, err : %s", err)
			} else {
				for _, user := range users {
					if user != nil && user.ID() != copy(user).ID() {
						t.Errorf("%s : %s json Encode & Decode fail ", pdu.SourceName, pdu.Signature2PublicKey)
					}
				}
				if users[0] != nil && users[1] != nil {
					break
				}
			}
		}
	}

	for i := 0; i < retryCnt; i++ {
		pk1, err := pdu.GenerateKey()
		if err != nil {
			t.Errorf("generate key pair fail, err : %s", err)
		}
		pk2, err := pdu.GenerateKey()
		if err != nil {
			t.Errorf("generate key pair fail, err : %s", err)
		}
		pk3, err := pdu.GenerateKey()
		if err != nil {
			t.Errorf("generate key pair fail, err : %s", err)
		}
		pubKey := crypto.PublicKey{Source: pdu.SourceName, SigType: pdu.MultipleSignatures, PubKey: append(append(append([]interface{}{}, pk1.PublicKey), pk2.PublicKey), pk3.PublicKey)}
		if users, err := CreateRootUsers(pubKey); err != nil {
			t.Errorf("create root users fail, err : %s", err)
		} else {
			for _, user := range users {
				if user != nil && user.ID() != copy(user).ID() {
					t.Errorf("%s : %s json Encode & Decode fail ", pdu.SourceName, pdu.MultipleSignatures)
				}
			}
			if users[0] != nil && users[1] != nil {
				break
			}
		}
	}

}

func TestCreateNewUser(t *testing.T) {
	Adam, privKeyAdam, Eve, privKeyEve := createRootUsers()
	value := MsgValue{
		ContentType: TypeDOB,
	}
	// build public key
	var privKeyA2Group []*ecdsa.PrivateKey
	for i := 0; i < 5; i++ {
		pk, _ := pdu.GenerateKey()
		privKeyA2Group = append(privKeyA2Group, pk)
	}
	var pubKeyA2Group []interface{}
	for _, v := range privKeyA2Group {
		pubKeyA2Group = append(pubKeyA2Group, v.PublicKey)
	}
	// build auth for new user
	auth := Auth{PublicKey: crypto.PublicKey{Source: pdu.SourceName, SigType: pdu.MultipleSignatures, PubKey: pubKeyA2Group}}
	// build dob msg content
	content, err := CreateDOBMsgContent("A2", "1234", &auth)
	if err != nil {
		t.Errorf("create bod content fail, err: %s", err)
	}
	content.SignByParent(Adam, privKeyAdam)
	content.SignByParent(Eve, privKeyEve)
	value.Content, err = json.Marshal(content)
	if err != nil {
		t.Errorf("content marshal fail , err: %s", err)
	}
	// build dob msg
	dobMsg, err := CreateMsg(Eve, &value, &privKeyEve)
	if err != nil {
		t.Errorf("create msg fail, err :%s", err)
	}
	// create new user by dob msg
	newUser, err := CreateNewUser(dobMsg)
	if err != nil {
		t.Errorf("create new user fail, err:%s", err)
	} else {
		if newUser.ID() != copy(newUser).ID() {
			t.Errorf("json Encode & Decode fail ")
		}
	}
}

func copy(u *User) *User {
	res, err := json.Marshal(u)
	if err != nil {
		return nil
	}
	var user User
	err = json.Unmarshal(res, &user)
	if err != nil {
		return nil
	}
	return &user
}

func createRootUsers() (*User, crypto.PrivateKey, *User, crypto.PrivateKey) {
	var Adam, Eve *User
	var privKeyGroup []*ecdsa.PrivateKey
	for i := 0; i < retryCnt; i++ {
		pk1, err := pdu.GenerateKey()
		if err != nil {
			continue
		}
		pk2, err := pdu.GenerateKey()
		if err != nil {
			continue
		}
		pk3, err := pdu.GenerateKey()
		if err != nil {
			continue
		}
		pubKey := crypto.PublicKey{Source: pdu.SourceName, SigType: pdu.MultipleSignatures, PubKey: append(append(append([]interface{}{}, pk1.PublicKey), pk2.PublicKey), pk3.PublicKey)}
		if users, err := CreateRootUsers(pubKey); err != nil {
			if err != nil {
				continue
			}
		} else {
			if users[0] != nil && users[1] != nil {
				Adam = users[0]
				Eve = users[1]
				privKeyGroup = append(append(append(privKeyGroup, pk1), pk2), pk3)
				break
			}
		}
	}
	return Adam, buildPrivateKey(privKeyGroup), Eve, buildPrivateKey(privKeyGroup)
}

func buildPrivateKey(privKeyGroup []*ecdsa.PrivateKey) crypto.PrivateKey {
	var privKeys []interface{}
	for _, k := range privKeyGroup {
		privKeys = append(privKeys, k)
	}
	return crypto.PrivateKey{
		Source:  pdu.SourceName,
		SigType: pdu.MultipleSignatures,
		PriKey:  privKeys,
	}
}
