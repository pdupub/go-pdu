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
)

const (
	TypeText = iota
	TypeDOB
)

type MsgValue struct {
	ContentType int
	Content     []byte
}

type DOBMsgContent struct {
	User    User
	Parents [2]ParentSig
}

type ParentSig struct {
	PID crypto.Hash
	Sig []byte
}

func CreateDOBMsgContent(name string, extra string, auth *Auth) (*DOBMsgContent, error) {
	user := User{Name: name, DOBExtra: extra, Auth: auth}
	return &DOBMsgContent{User: user}, nil

}

func (mv *DOBMsgContent) SignByParent(user *User, privKey crypto.PrivateKey) error {

	jsonByte, err := json.Marshal(mv.User)
	if err != nil {
		return err
	}
	var signature *crypto.Signature
	switch privKey.Source {
	case pdu.SourceName:
		signature, err = pdu.Sign(jsonByte, privKey)
		if err != nil {
			return err
		}
	}
	if user.Gender() == male {
		mv.Parents[1] = ParentSig{PID: user.ID(), Sig: signature.Signature}
	} else {
		mv.Parents[0] = ParentSig{PID: user.ID(), Sig: signature.Signature}
	}
	return nil
}
