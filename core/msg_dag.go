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
	msgD *dag.DAG
	ids  []common.Hash
	tpD  *dag.DAG //map[common.Hash]*TimeProof
	ugD  map[common.Hash]*Group
}

// NewMsgDag create MsgDAG
// the msg will also be used to create time proof as msg.SenderID
func NewMsgDag(userDAG *Group, msg *Message) (*MsgDAG, error) {
	// check msg sender from valid user
	if nil == userDAG.GetUserByID(msg.SenderID) {
		return nil, ErrMsgFromInvalidUser
	}
	// build msg dag
	msgVertex, err := dag.NewVertex(msg.ID(), msg)
	if err != nil {
		return nil, err
	}
	msgD, err := dag.NewDAG(msgVertex)
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
	tpVertex, err := dag.NewVertex(msg.SenderID, tp)
	if err != nil {
		return nil, err
	}
	tpD, err := dag.NewDAG(tpVertex)
	if err != nil {
		return nil, err
	}

	msgDAG := MsgDAG{msgD: msgD,
		ids: ids,
		tpD: tpD,
		ugD: map[common.Hash]*Group{msg.SenderID: userDAG}}
	return &msgDAG, nil
}

// CheckUserValid check if the user valid in this MsgDAG
func (md *MsgDAG) CheckUserValid(userID common.Hash) bool {
	for k, v := range md.ugD {
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
	if nil != md.tpD.GetVertex(msg.SenderID) {
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
				// todo : tpVertex should add parents
				// todo : parent contain all current tpVertex which this msg is valid
				tpVertex, err := dag.NewVertex(msg.SenderID, tp)
				if err != nil {
					return err
				}
				if err = md.tpD.AddVertex(tpVertex); err != nil {
					return err
				}
				initialize = false
			} else {
				if err := md.updateTimeProof(msgTP); err != nil {
					return err
				}
			}
		}
	}

	md.ugD[msg.SenderID] = md.createUserMap(msg.SenderID)

	return nil
}

//
func (md *MsgDAG) createUserMap(userID common.Hash) *Group {
	// todo : the new UserDAG should contain all parent users
	// todo : in all userDag which this userID is valid
	// todo : need deep copy
	return nil
}

// GetUserDAG return userDAG by time proof userID
func (md *MsgDAG) GetUserDAG(userID common.Hash) *Group {
	if userDag, ok := md.ugD[userID]; ok {
		return userDag
	}
	return nil
}

// GetMsgByID will return the msg by msg.ID()
// nil will be return if msg not exist
func (md *MsgDAG) GetMsgByID(mid common.Hash) *Message {
	if v := md.msgD.GetVertex(mid); v != nil {
		return v.Value().(*Message)
	} else {
		return nil
	}
}

// Add will check if the msg from valid user,
// add new msg into MsgDAG, and update time proof if
// msg.SenderID is belong to time proof
func (md *MsgDAG) Add(msg *Message) error {
	// check
	if md.GetMsgByID(msg.ID()) != nil {
		return ErrMsgAlreadyExist
	}
	if !md.CheckUserValid(msg.SenderID) {
		return ErrMsgFromInvalidUser
	}
	// update dag
	var refs []interface{}
	for _, r := range msg.Reference {
		refs = append(refs, r.MsgID)
	}
	msgVertex, err := dag.NewVertex(msg.ID(), msg, refs...)
	if err != nil {
		return err
	}
	err = md.msgD.AddVertex(msgVertex)
	if err != nil {
		return err
	}
	// ids
	md.ids = append(md.ids, msg.ID())
	// update tp
	err = md.updateTimeProof(msg)
	if err != nil {
		return err
	}
	// process the msg
	err = md.processMsg(msg)
	if err != nil {
		return err
	}
	return nil
}

func (md *MsgDAG) processMsg(msg *Message) error {
	switch msg.Value.ContentType {
	case TypeText:
		return nil
	case TypeDOB:
		user, err := CreateNewUser(msg)
		if err != nil {
			return err
		}
		// todo :check the valid time proof for parents in each timeproof
		// user may not can be add to all userMap
		for _, v := range md.ugD {
			err = v.Add(user)
			if err != nil {
				return err
			}
		}
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
	if vertex := md.tpD.GetVertex(msg.SenderID); vertex != nil {

		tp := vertex.Value().(*TimeProof)
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
	if vertex := md.tpD.GetVertex(userID); vertex != nil {
		return vertex.Value().(*TimeProof).maxSeq
	} else {
		return 0
	}
}
