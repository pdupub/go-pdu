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
		if _, pubKey, err := pdu.GenKey(pdu.Signature2PublicKey); err != nil {
			t.Errorf("generate key fail, err :%s", err)
		} else {
			if users, err := CreateRootUsers(*pubKey); err != nil {
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
		if _, pubKey, err := pdu.GenKey(pdu.MultipleSignatures, 3); err != nil {

		} else {
			if users, err := CreateRootUsers(*pubKey); err != nil {
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

}

func TestCreateNewUser(t *testing.T) {
	Adam, privKeyAdam, Eve, privKeyEve := createRootUsers()
	value := MsgValue{
		ContentType: TypeDOB,
	}

	_, pubKey, err := pdu.GenKey(pdu.MultipleSignatures, 5)
	if err != nil {
		t.Errorf("generate key fail, err:%s", err)
	}
	// build auth for new user
	auth := Auth{PublicKey: *pubKey}
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
	var privKeyRes crypto.PrivateKey
	for i := 0; i < retryCnt; i++ {

		privKey, pubKey, err := pdu.GenKey(pdu.MultipleSignatures, 3)
		if err != nil {
			continue
		}
		if users, err := CreateRootUsers(*pubKey); err != nil {
			if err != nil {
				continue
			}
		} else {
			if users[0] != nil && users[1] != nil {
				Adam = users[0]
				Eve = users[1]
				privKeyRes = *privKey
				break
			}
		}
	}
	return Adam, privKeyRes, Eve, privKeyRes
}
