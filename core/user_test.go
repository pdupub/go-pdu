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
	"testing"

	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/ethereum"
)

const (
	retryCnt = 100
)

var (
	userEngine crypto.Engine
)

func TestCreateRootUsersS2PK(t *testing.T) {
	userEngine = ethereum.New()
	var f, m bool
	for i := 0; i < retryCnt; i++ {
		if _, pubKey, err := userEngine.GenKey(crypto.Signature2PublicKey); err != nil {
			t.Error("generate key fail", err)
		} else {
			user := CreateRootUser(*pubKey, "name", "extra")

			if user.ID() != copy(user).ID() {
				t.Errorf("%s : %s json Encode & Decode fail ", userEngine.Name(), crypto.Signature2PublicKey)
			}
			if user.Gender() {
				m = true
			} else {
				f = true
			}
			// both checked
			if f && m {
				break
			}

		}
	}
}
func TestCreateRootUsersMS(t *testing.T) {
	var f, m bool
	for i := 0; i < retryCnt; i++ {
		if _, pubKey, err := userEngine.GenKey(crypto.MultipleSignatures, 3); err != nil {

		} else {
			user := CreateRootUser(*pubKey, "name", "extra")

			if user.ID() != copy(user).ID() {
				t.Errorf("%s : %s json Encode & Decode fail ", userEngine.Name(), crypto.Signature2PublicKey)
			}
			if user.Gender() {
				m = true
			} else {
				f = true
			}
			// both checked
			if f && m {
				break
			}

		}
	}

}

func TestCreateNewUser(t *testing.T) {
	Adam, privKeyAdam, Eve, privKeyEve := createRootUsers()
	value := MsgValue{
		ContentType: TypeBirth,
	}

	_, pubKey, err := userEngine.GenKey(crypto.MultipleSignatures, 5)
	if err != nil {
		t.Error("generate key fail", err)
	}
	// build auth for new user
	auth := Auth{PublicKey: *pubKey}
	// build birth msg content
	content, err := CreateBirthMsgContent("A2", "1234", &auth)
	if err != nil {
		t.Error("create birth content fail", err)
	}
	content.SignByParent(Adam, privKeyAdam)
	content.SignByParent(Eve, privKeyEve)
	value.Content, err = json.Marshal(content)
	if err != nil {
		t.Error("content marshal fail ", err)
	}
	// build birth msg
	birthMsg, err := CreateMsg(Eve, &value, &privKeyEve)
	if err != nil {
		t.Error("create msg fails", err)
	}

	universe, err := NewUniverse(Eve, Adam)
	if err != nil {
		t.Error("create universe fail", err)
	}
	// create new user by birth msg
	newUser, err := CreateNewUser(universe, birthMsg)
	if err != nil {
		t.Error("create new user fail", err)
	} else {
		if newUser.ID() != copy(newUser).ID() {
			t.Error("json Encode & Decode fail ")
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
	var APK, EPK crypto.PrivateKey
	for i := 0; i < retryCnt; i++ {

		privKey, pubKey, err := userEngine.GenKey(crypto.MultipleSignatures, 3)
		if err != nil {
			continue
		}
		user := CreateRootUser(*pubKey, "name", "extra")

		if user.Gender() {
			Adam = user
			APK = *privKey
		} else {
			Eve = user
			EPK = *privKey
		}
		if Adam != nil && Eve != nil {
			break
		}

	}
	return Adam, APK, Eve, EPK
}
