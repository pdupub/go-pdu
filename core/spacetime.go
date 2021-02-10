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
	dag "github.com/pdupub/go-dag"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core/rule"
)

// SpaceTime contain time proof of this space time and the user info who is valid in this space time
type SpaceTime struct {
	maxTimeSequence uint64
	timeProofD      *dag.DAG // msg.id  : time sequence
	userStateD      *dag.DAG // user.id : user info (strict)
}

// NewSpaceTime create the new space-time
func NewSpaceTime(u *Universe, msg *Message, ref *MsgReference) (*SpaceTime, error) {
	spaceTime := &SpaceTime{}
	// create time proof and set max time sequence
	if err := spaceTime.createTimeProofD(msg); err != nil {
		return nil, err
	}

	// create the user state
	if nil != u.stD {
		// base on parent space time
		st, err := spaceTime.getParentSpaceTime(u, msg, ref)
		if err != nil {
			return nil, err
		}
		if err := spaceTime.createUserStateD(st, ref, u.userD); err != nil {
			return nil, err
		}
	} else {
		// create the userState for the first space time
		if err := spaceTime.createFirstUserStateD(ref, u.userD); err != nil {
			return nil, err
		}
	}
	return spaceTime, nil
}

// getParentSpaceTime return the parent space time. ref must be one of msg.ref. and SpaceTime of ref.SenderID must already be exist in universe.
func (s SpaceTime) getParentSpaceTime(u *Universe, msg *Message, ref *MsgReference) (*SpaceTime, error) {
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
	return stVertex.Value().(*SpaceTime), nil
}

// createTimeProofD create time proof DAG start from sequence 1
func (s *SpaceTime) createTimeProofD(msg *Message) error {
	timeSequence := uint64(1)
	timeVertex, err := dag.NewVertex(msg.ID(), timeSequence)
	if err != nil {
		return err
	}
	timeProofDag, err := dag.NewDAG(1, timeVertex)
	if err != nil {
		return err
	}
	timeProofDag.RemoveStrict()
	s.timeProofD = timeProofDag
	s.maxTimeSequence = timeSequence
	return nil
}

// createUserStateD create user state DAG base on userD
func (s *SpaceTime) createUserStateD(st *SpaceTime, ref *MsgReference, userD *dag.DAG) error {
	userStateD, err := dag.NewDAG(2)
	if err != nil {
		return err
	}
	userStateD.SetMaxParentsCount(2)

	refSeq := st.timeProofD.GetVertex(ref.MsgID).Value().(uint64)
	refUserStateD := st.userStateD
	for _, k := range refUserStateD.GetIDs() {
		lifeMaxSeq := refUserStateD.GetVertex(k).Value().(*UserInfo).natureLifeMaxSeq - refSeq
		userStateVertex, err := dag.NewVertex(k, NewUserInfo(userD.GetVertex(k).Value().(*User).Name, lifeMaxSeq, 0), userD.GetVertex(k).ParentIDs()...)
		if err != nil {
			return err
		}
		userStateD.AddVertex(userStateVertex)
	}
	s.userStateD = userStateD
	return nil
}

// createFirstUserStateD create user state DAG, this space time is from one of roots, who send the first message on universe
func (s *SpaceTime) createFirstUserStateD(ref *MsgReference, userD *dag.DAG) error {
	userStateD, err := dag.NewDAG(2)
	if err != nil {
		return err
	}
	userStateD.SetMaxParentsCount(2)
	for _, k := range userD.GetIDs() {
		userStateVertex, err := dag.NewVertex(k, NewUserInfo(userD.GetVertex(k).Value().(*User).Name, rule.MaxLifeTime, 0))
		if err != nil {
			return err
		}
		userStateD.AddVertex(userStateVertex)
	}
	s.userStateD = userStateD
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

// UpdateTimeProof update the tp info
func (s *SpaceTime) UpdateTimeProof(msg *Message) error {
	var currentSeq uint64 = 1
	var ref interface{}
	for _, r := range msg.Reference {
		if r.SenderID == msg.SenderID {
			refVertex := s.timeProofD.GetVertex(r.MsgID)
			if refVertex != nil {
				refSeq := refVertex.Value().(uint64)
				if currentSeq <= refSeq {
					currentSeq = refSeq + 1
					ref = r.MsgID
				}
			}
		}
	}
	timeVertex, err := dag.NewVertex(msg.ID(), currentSeq, ref)
	if err != nil {
		return err
	}

	if err := s.timeProofD.AddVertex(timeVertex); err != nil {
		return err
	} else if currentSeq > s.maxTimeSequence {
		s.maxTimeSequence = currentSeq
	}
	return nil
}

// AddUser add user info to this space time
func (s *SpaceTime) AddUser(ref *MsgReference, contentBirth ContentBirth, user *User) error {
	if tp := s.timeProofD.GetVertex(ref.MsgID); tp != nil {
		msgSeq := tp.Value().(uint64)
		p0 := s.userStateD.GetVertex(contentBirth.Parents[0].UserID)
		if p0 == nil {
			return ErrAddUserToSpaceTimeFail
		}
		userInfo0 := p0.Value().(*UserInfo)
		p1 := s.userStateD.GetVertex(contentBirth.Parents[1].UserID)
		if p1 == nil {
			return ErrAddUserToSpaceTimeFail
		}
		userInfo1 := p1.Value().(*UserInfo)
		if userInfo0.natureBirthSeq+userInfo0.natureLifeMaxSeq > msgSeq &&
			userInfo1.natureBirthSeq+userInfo1.natureLifeMaxSeq > msgSeq &&
			msgSeq-userInfo0.natureLastCosign > rule.ReproductionInterval &&
			msgSeq-userInfo1.natureLastCosign > rule.ReproductionInterval {
			// update nature last cosign number as msgSeq
			userInfo0.natureLastCosign = msgSeq
			userInfo1.natureLastCosign = msgSeq
			// add user in this st
			userVertex, err := dag.NewVertex(user.ID(), NewUserInfo(user.Name, user.LifeTime, msgSeq), p0, p1)
			if err != nil {
				return err
			}
			if err := s.userStateD.AddVertex(userVertex); err != nil {
				return err
			}
			return nil
		}
	}
	return ErrAddUserToSpaceTimeFail
}
