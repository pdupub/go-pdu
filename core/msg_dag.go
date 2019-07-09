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
	ErrMsgAlreadyExist = errors.New("msg already exist")
)

type TimeProof struct {
	maxSeq uint64
	dag    *dag.DAG
}

type MsgDAG struct {
	dag   *dag.DAG
	tpMap map[common.Hash]*TimeProof
}

func NewMsgDag(msg *Message) (*MsgDAG, error) {
	msgVertex, err := dag.NewVertex(msg.ID(), msg)
	if err != nil {
		return nil, err
	}
	msgDAG, err := dag.NewDAG(msgVertex)
	if err != nil {
		return nil, err
	}
	// use this msg as time
	timeVertex, err := dag.NewVertex(msg.ID(), uint64(1))
	if err != nil {
		return nil, err
	}
	timeDag, err := dag.NewDAG(timeVertex)
	if err != nil {
		return nil, err
	}
	tp := &TimeProof{maxSeq: timeVertex.Value().(uint64), dag: timeDag}
	tpMap := make(map[common.Hash]*TimeProof)
	tpMap[msg.SenderID] = tp
	return &MsgDAG{dag: msgDAG, tpMap: tpMap}, nil
}

func (md *MsgDAG) GetMsgByID(mid common.Hash) *Message {
	if v := md.dag.GetVertex(mid); v != nil {
		return v.Value().(*Message)
	} else {
		return nil
	}
}

func (md *MsgDAG) Add(msg *Message) error {
	if md.GetMsgByID(msg.ID()) != nil {
		return ErrMsgAlreadyExist
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

func (md *MsgDAG) GetMaxSeq(userID common.Hash) uint64 {
	if tp, ok := md.tpMap[userID]; ok {
		return tp.maxSeq
	} else {
		return 0
	}
}
