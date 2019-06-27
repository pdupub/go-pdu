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
	"crypto/sha256"
	"fmt"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/msg"
	"math/big"
)

const (
	rootMName     = "Adam"
	rootMDOBExtra = "Hello World!"
	rootFName     = "Eve"
	rootFDOBExtra = ";-)"
	male          = true
	female        = false
)

type User struct {
	name     string      `json:"name"`
	dobExtra []byte      `json:"extra"`
	auth     Auth        `json:"auth"`
	dobMsg   msg.Message `json:"dobMsg"`
}

// CreateRootUser try to create two root users by public key
// One Male user and one female user,
func CreateRootUsers(key crypto.PublicKey) ([2]*User, error) {
	rootUsers := [2]*User{nil, nil}
	rootFUser := User{name: rootFName, dobExtra: []byte(rootFDOBExtra), auth: Auth{key}, dobMsg: msg.Message{}}
	if rootFUser.Gender() == female {
		rootUsers[0] = &rootFUser
	}
	rootMUser := User{name: rootMName, dobExtra: []byte(rootMDOBExtra), auth: Auth{key}, dobMsg: msg.Message{}}
	if rootMUser.Gender() == male {
		rootUsers[1] = &rootMUser
	}
	return rootUsers, nil
}

// CreateNewUser create new user by cosign message
func CreateNewUser(msg *msg.Message) (*User, error) {
	newUser := User{}

	return &newUser, nil
}

// ID return the vertex.id, related to parents and value of the vertex
// ID cloud use as address of user account
func (u User) ID() []byte {
	hash := sha256.New()
	hash.Reset()
	auth := fmt.Sprintf("%v", u.auth)
	dobMsg := fmt.Sprintf("%v", u.dobMsg)
	hash.Write(append(append(append([]byte(u.name), u.dobExtra...), auth...), dobMsg...))
	return hash.Sum(nil)
}

// Gender return the gender of user, true = male = end of ID is odd
func (u User) Gender() bool {
	if uid := new(big.Int).SetBytes(u.ID()); uid.Mod(uid, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
		return true
	}
	return false
}

// Value return the vertex.value
func (u User) Value() interface{} {
	return nil
}

// ParentsID return the ID of user parents,
// res[0] should be the female parent (id end by even)
// res[1] should be the male parent (id end by odd)
func (u User) ParentsID() [2][]byte {
	var PID [2][]byte

	PID[0] = []byte{}
	PID[1] = []byte{}
	return PID
}
