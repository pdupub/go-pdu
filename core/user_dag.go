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
	"errors"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/dag"
)

type UserDAG struct {
	dag *dag.DAG
}

var (
	errUserAlreadyExist = errors.New("user already exist")
)

func NewUserDag(Eve, Adam *User) (*UserDAG, error) {

	EveVertex, err := dag.NewVertex(Eve.ID(), Eve)
	if err != nil {
		return nil, err
	}

	AdamVertex, err := dag.NewVertex(Adam.ID(), Adam)
	if err != nil {
		return nil, err
	}

	userDAG, err := dag.NewDAG(EveVertex, AdamVertex)
	if err != nil {
		return nil, err
	}
	userDAG.SetMaxParentsCount(2)
	return &UserDAG{dag: userDAG}, nil

}

func (ud *UserDAG) GetUserByID(uid common.Hash) *User {
	if v := ud.dag.GetVertex(uid); v != nil {
		return v.Value().(*User)
	} else {
		return nil
	}
}

func (ud *UserDAG) Add(user *User) error {
	if ud.GetUserByID(user.ID()) != nil {
		return errUserAlreadyExist
	}
	var dobContent DOBMsgContent
	err := json.Unmarshal(user.DOBMsg.Value.Content, &dobContent)
	if err != nil {
		return err
	}
	userVertex, err := dag.NewVertex(user.ID(), user, dobContent.Parents[0].PID, dobContent.Parents[1].PID)
	if err != nil {
		return err
	}
	err = ud.dag.AddVertex(userVertex)
	if err != nil {
		return err
	}
	return nil
}
