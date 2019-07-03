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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/pdupub/go-pdu/crypto"
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
	Name     string   `json:"name"`
	DOBExtra string   `json:"extra"`
	Auth     *Auth    `json:"auth"`
	DOBMsg   *Message `json:"dobMsg"`
}

// CreateRootUser try to create two root users by public key
// One Male user and one female user,
func CreateRootUsers(key crypto.PublicKey) ([2]*User, error) {
	rootUsers := [2]*User{nil, nil}
	rootFUser := User{Name: rootFName, DOBExtra: rootFDOBExtra, Auth: &Auth{key}, DOBMsg: &Message{}}
	if rootFUser.Gender() == female {
		rootUsers[0] = &rootFUser
	}
	rootMUser := User{Name: rootMName, DOBExtra: rootMDOBExtra, Auth: &Auth{key}, DOBMsg: &Message{}}
	if rootMUser.Gender() == male {
		rootUsers[1] = &rootMUser
	}
	return rootUsers, nil
}

// CreateNewUser create new user by cosign message
// The msg must be signed by user in local user dag.
// Both parents must be in the local use dag.
// Both parents fit the nature rules.
// The BOD struct signed by both parents.
func CreateNewUser(msg *Message) (*User, error) {
	newUser := User{}

	return &newUser, nil
}

// ID return the vertex.id, related to parents and value of the vertex
// ID cloud use as address of user account
func (u User) ID() crypto.Hash {
	hash := sha256.New()
	hash.Reset()
	auth := fmt.Sprintf("%v", u.Auth)
	dobMsg := fmt.Sprintf("%v", u.DOBMsg)
	hash.Write(append(append(append([]byte(u.Name), u.DOBExtra...), auth...), dobMsg...))
	return crypto.Bytes2Hash(hash.Sum(nil))
}

// Gender return the gender of user, true = male = end of ID is odd
func (u User) Gender() bool {
	hashID := u.ID()
	if uid := new(big.Int).SetBytes(hashID[:]); uid.Mod(uid, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
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
func (u User) ParentsID() [2]crypto.Hash {
	var parentsID [2]crypto.Hash
	if u.DOBMsg != nil {
		// get parents from dobMsg

	}
	return parentsID
}

func (u *User) UnmarshalJSON(input []byte) error {
	userMap := make(map[string]interface{})
	err := json.Unmarshal(input, &userMap)
	if err != nil {
		return err
	} else {
		u.Name = userMap["name"].(string)
		u.DOBExtra = userMap["dobExtra"].(string)
		json.Unmarshal([]byte(userMap["dobMsg"].(string)), &u.DOBMsg)
		json.Unmarshal([]byte(userMap["auth"].(string)), &u.Auth)
	}
	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	userMap := make(map[string]interface{})
	userMap["name"] = u.Name
	userMap["dobExtra"] = u.DOBExtra

	if auth, err := json.Marshal(&u.Auth); err != nil {
		return []byte{}, err
	} else {
		userMap["auth"] = string(auth)
	}

	if dobMsg, err := json.Marshal(&u.DOBMsg); err != nil {
		return []byte{}, err
	} else {
		userMap["dobMsg"] = string(dobMsg)
	}
	return json.Marshal(userMap)
}
