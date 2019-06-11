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

package types

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/pdupub/go-pdu/common"
)

type User struct {
	*Vertex
}

func NewUser(msg Message) (*User, common.Hash, error) {

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, common.Hash{}, err
	}

	var hashKey common.Hash
	hashKey = sha256.Sum256(msgBytes)

	user := &User{
		NewVertex(hashKey, Gene{}, msg.parents),
	}
	return user, hashKey, nil
}

func RootUser(gender bool) (*User, common.Hash, error) {

	return &User{}, common.Hash{}, nil
}
