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

package user

import (
	"encoding/json"
	"fmt"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/pdu"
	"testing"
)

func TestCreateRootUsers(t *testing.T) {

	pk, err := pdu.GenerateKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}
	pubKey := crypto.PublicKey{Source: pdu.SourceName, SigType: pdu.Signature2PublicKey, PubKey: pk.PublicKey}

	users, err := CreateRootUsers(pubKey)
	for i, user := range users {
		if user != nil {
			if crypto.Byte2String(user.ID()) == crypto.Byte2String(encode(user).ID()) {
				fmt.Println("User:", i, "ID:", crypto.Byte2String(user.ID()))
			} else {
				t.Errorf("%s : %s json Encode & Decode fail ", pdu.SourceName, pdu.Signature2PublicKey)
			}
		} else {
			fmt.Println("User:", i, "No user being created")
		}
	}

	pk2, err := pdu.GenerateKey()
	if err != nil {
		t.Errorf("generate key pair fail, err : %s", err)
	}
	pubKey = crypto.PublicKey{Source: pdu.SourceName, SigType: pdu.MultipleSignatures, PubKey: append(append([]interface{}{}, pk.PublicKey), pk2.PublicKey)}

	users, err = CreateRootUsers(pubKey)
	for i, user := range users {
		if user != nil {
			if crypto.Byte2String(user.ID()) == crypto.Byte2String(encode(user).ID()) {
				fmt.Println("User:", i, "ID:", crypto.Byte2String(user.ID()))
			} else {
				t.Errorf("%s : %s json Encode & Decode fail ", pdu.SourceName, pdu.MultipleSignatures)
			}
		} else {
			fmt.Println("User:", i, "No user being created")
		}
	}

}

func encode(u *User) *User {
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

func TestCreateNewUser(t *testing.T) {

}
