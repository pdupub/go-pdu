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
	"errors"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/dag"
)

var (
	ErrMsgFromInvalidUser = errors.New("msg from invalid user")
	ErrMsgAlreadyExist    = errors.New("msg already exist")
	ErrMsgNotFound        = errors.New("msg not found")
	ErrTPAlreadyExist     = errors.New("time proof already exist")
)

type TimeProof struct {
	maxSeq uint64
	dag    *dag.DAG
}

type MsgDAG struct {
	dag     *dag.DAG
	ids     []common.Hash
	tpMap   map[common.Hash]*TimeProof
	userMap map[common.Hash]*UserDAG
}

// NewMsgDag create MsgDAG
// the msg will also be used to create time proof as msg.SenderID
func NewMsgDag(userDAG *UserDAG, msg *Message) (*MsgDAG, error) {
	// check msg sender from valid user
	if nil == userDAG.GetUserByID(msg.SenderID) {
		return nil, ErrMsgFromInvalidUser
	}
	// build msg dag
	msgVertex, err := dag.NewVertex(msg.ID(), msg)
	if err != nil {
		return nil, err
	}
	dag, err := dag.NewDAG(msgVertex)
	if err != nil {
		return nil, err
	}
	// init ids
	ids := []common.Hash{msg.ID()}
	// build time proof
	tp, err := createTimeProof(msg)
	if err != nil {
		return nil, err
	}

	msgDAG := MsgDAG{dag: dag,
		ids:     ids,
		tpMap:   map[common.Hash]*TimeProof{msg.SenderID: tp},
		userMap: map[common.Hash]*UserDAG{msg.SenderID: userDAG}}
	return &msgDAG, nil
}

// CheckUserValid check if the user valid in this MsgDAG
func (md *MsgDAG) CheckUserValid(userID common.Hash) bool {
	for k, v := range md.userMap {
		if k == userID {
			return true
		}
		if nil != v.GetUserByID(userID) {
			return true
		}
	}
	return false
}

// AddTimeProof will get all messages save in MsgDAG with same msg.SenderID
// and build the time proof by those messages
func (md *MsgDAG) AddTimeProof(msg *Message) error {
	if md.GetMsgByID(msg.ID()) == nil {
		return ErrMsgNotFound
	}
	if _, ok := md.tpMap[msg.SenderID]; ok {
		return ErrTPAlreadyExist
	}
	if !md.CheckUserValid(msg.SenderID) {
		return ErrMsgFromInvalidUser
	}

	initialize := true
	for _, id := range md.ids {
		if msgTP := md.GetMsgByID(id); msgTP != nil && msgTP.SenderID == msg.SenderID {
			if initialize {
				tp, err := createTimeProof(msgTP)
				if err != nil {
					return err
				}
				md.tpMap[msg.SenderID] = tp
				initialize = false
			} else {
				if err := md.updateTimeProof(msgTP); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// GetMsgByID will return the msg by msg.ID()
// nil will be return if msg not exist
func (md *MsgDAG) GetMsgByID(mid common.Hash) *Message {
	if v := md.dag.GetVertex(mid); v != nil {
		return v.Value().(*Message)
	} else {
		return nil
	}
}

// Add will check if the msg from valid user,
// add new msg into MsgDAG, and update time proof if
// msg.SenderID is belong to time proof
func (md *MsgDAG) Add(msg *Message) error {
	if md.GetMsgByID(msg.ID()) != nil {
		return ErrMsgAlreadyExist
	}

	if !md.CheckUserValid(msg.SenderID) {
		return ErrMsgFromInvalidUser
	}
	var refs []interface{}
	for _, r := range msg.Reference {
		refs = append(refs, r.MsgID)
	}

	msgVertex, err := dag.NewVertex(msg.ID(), msg, refs...)
	if err != nil {
		return err
	}
	err = md.dag.AddVertex(msgVertex)
	if err != nil {
		return err
	}

	md.ids = append(md.ids, msg.ID())
	err = md.updateTimeProof(msg)
	if err != nil {
		return err
	}
	return nil
}

func createTimeProof(msg *Message) (*TimeProof, error) {
	timeVertex, err := dag.NewVertex(msg.ID(), uint64(1))
	if err != nil {
		return nil, err
	}
	timeDag, err := dag.NewDAG(timeVertex)
	if err != nil {
		return nil, err
	}
	return &TimeProof{maxSeq: timeVertex.Value().(uint64), dag: timeDag}, nil
}

func (md *MsgDAG) updateTimeProof(msg *Message) error {
	if tp, ok := md.tpMap[msg.SenderID]; ok {
		var currentSeq uint64 = 1
		for _, r := range msg.Reference {
			if r.SenderID == msg.SenderID {
				refSeq := tp.dag.GetVertex(r.MsgID).Value().(uint64)
				if currentSeq <= refSeq {
					currentSeq = refSeq + 1
				}
			}
		}
		timeVertex, err := dag.NewVertex(msg.ID(), currentSeq)
		if err != nil {
			return err
		}

		if err := tp.dag.AddVertex(timeVertex); err != nil {
			return err
		} else if currentSeq > tp.maxSeq {
			tp.maxSeq = currentSeq
		}
	}
	return nil
}

// GetMaxSeq will return the max time proof sequence for
// time proof by the userID
func (md *MsgDAG) GetMaxSeq(userID common.Hash) uint64 {
	if tp, ok := md.tpMap[userID]; ok {
		return tp.maxSeq
	} else {
		return 0
	}
}
