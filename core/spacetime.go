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
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core/rule"
	"github.com/pdupub/go-pdu/dag"
)

// SpaceTime contain time proof of this space time and the user info who is valid in this space time
type SpaceTime struct {
	maxTimeSequence uint64
	timeProofD      *dag.DAG // msg.id  : time sequence
	userStateD      *dag.DAG // user.id : user info (strict)
}

func NewSpaceTime(u *Universe, msg *Message, ref *MsgReference) (*SpaceTime, error) {
	timeVertex, err := dag.NewVertex(msg.ID(), uint64(1))
	if err != nil {
		return nil, err
	}
	timeProofDag, err := dag.NewDAG(1, timeVertex)
	if err != nil {
		return nil, err
	}
	timeProofDag.RemoveStrict()

	spaceTime := &SpaceTime{maxTimeSequence: timeVertex.Value().(uint64), timeProofD: timeProofDag, userStateD: nil}
	if nil != u.stD {
		if ref == nil {
			return nil, ErrCreateSpaceTimeFail
		}
		refExist := false
		for _, v := range msg.Reference {
			if v.SenderID == ref.SenderID && v.MsgID == ref.MsgID {
				refExist = true
			}
		}
		if !refExist {
			return nil, ErrCreateSpaceTimeFail
		}
		stVertex := u.stD.GetVertex(ref.SenderID)
		if stVertex == nil {
			return nil, ErrCreateSpaceTimeFail
		}
		st := stVertex.Value().(*SpaceTime)
		if err := spaceTime.createUserStateD(st, ref, u.userD); err != nil {
			return nil, err
		}
	} else {
		// create the userState for the first space time
		if err := spaceTime.createUserStateD(nil, ref, u.userD); err != nil {
			return nil, err
		}
	}
	return spaceTime, nil
}

func (s *SpaceTime) createUserStateD(st *SpaceTime, ref *MsgReference, userD *dag.DAG) error {
	newUserStateD, err := dag.NewDAG(2)
	newUserStateD.SetMaxParentsCount(2)
	if err != nil {
		return err
	}
	refSeq := uint64(0)
	lifeMaxSeq := rule.MaxLifeTime
	userStateD := userD
	if st != nil {
		refSeq = st.timeProofD.GetVertex(ref.MsgID).Value().(uint64)
		userStateD = st.userStateD
	}
	for _, k := range userStateD.GetIDs() {
		if st != nil {
			lifeMaxSeq = userStateD.GetVertex(k).Value().(*UserInfo).natureLifeMaxSeq - refSeq
		}

		userStateVertex, err := dag.NewVertex(k, &UserInfo{natureState: UserStatusNormal, natureLastCosign: 0, natureLifeMaxSeq: lifeMaxSeq, natureDOBSeq: 0, localNickname: userD.GetVertex(k).Value().(*User).Name}, userD.GetVertex(k).Parents()...)
		if err != nil {
			return err
		}

		newUserStateD.AddVertex(userStateVertex)
	}
	s.userStateD = newUserStateD
	return nil
}

// GetUserIDs returns users ID of space-time
func (s SpaceTime) GetUserIDs() (userIDs []common.Hash) {
	for _, id := range s.userStateD.GetIDs() {
		userIDs = append(userIDs, id.(common.Hash))
	}
	return userIDs
}

// GetUserInfo returns userInfo in current space-time by userID
func (s SpaceTime) GetUserInfo(userID common.Hash) *UserInfo {
	if uVertex := s.userStateD.GetVertex(userID); uVertex != nil {
		return uVertex.Value().(*UserInfo)
	}
	return nil
}
