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
			fmt.Println("User:", i, "ID:", crypto.Byte2String(user.ID()))
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
			fmt.Println("User:", i, "ID:", crypto.Byte2String(user.ID()))
		} else {
			fmt.Println("User:", i, "No user being created")
		}
	}

}

func TestCreateNewUser(t *testing.T) {

}
