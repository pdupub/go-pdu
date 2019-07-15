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
	"github.com/pdupub/go-pdu/core/rule"
	"github.com/pdupub/go-pdu/dag"
)

var (
	ErrUserNotExist     = errors.New("user not exist")
	ErrMsgAlreadyExist  = errors.New("msg already exist")
	ErrMsgNotFound      = errors.New("msg not found")
	ErrTPAlreadyExist   = errors.New("time proof already exist")
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrNotSupportYet    = errors.New("not error, just not support ye")
	ErrNewUserAddFail   = errors.New("new user add fail")
)

const (
	UserStatusNormal = iota
)

// UserState contain the information except pass by DOBMsg
// the state related to nature rule is start by nature
// the other state start by local
type UserInfo struct {
	natureState      int    // validation state depend on nature rule
	natureLastCosign uint64 // last DOB cosign
	localNickname    string
}

// SpaceTime contain time proof of this space time and the user info who is valid in this space time
type SpaceTime struct {
	maxTimeSequence uint64
	timeProofD      *dag.DAG // msg.id  : time sequence
	userStateD      *dag.DAG // user.id : state of user
}

// Universe
type Universe struct {
	msgD  *dag.DAG `json:"messageDAG"`   // contain all messages valid in any universe (time proof)
	userD *dag.DAG `json:"userDAG"`      // contain all users valid in any universe (time proof)
	stD   *dag.DAG `json:"spaceTimeDAG"` // contain all space time, which is the origin thought of PDU
}

// NewUniverse create Universe from two user with diff gender
func NewUniverse(Eve, Adam *User) (*Universe, error) {
	if Eve.Gender() == Adam.Gender() {
		return nil, ErrNotSupportYet
	}
	EveVertex, err := dag.NewVertex(Eve.ID(), Eve)
	if err != nil {
		return nil, err
	}
	AdamVertex, err := dag.NewVertex(Adam.ID(), Adam)
	if err != nil {
		return nil, err
	}
	userD, err := dag.NewDAG(EveVertex, AdamVertex)
	if err != nil {
		return nil, err
	}
	userD.SetMaxParentsCount(2)
	return &Universe{userD: userD}, nil
}

// AddMsg will check if the msg from valid user,
// add new msg into Universe, and update time proof if
// msg.SenderID is belong to time proof
func (u *Universe) AddMsg(msg *Message) error {
	if !u.CheckUserExist(msg.SenderID) {
		return ErrUserNotExist
	}
	if u.msgD == nil {
		if err := u.initializeMsgD(msg); err != nil {
			return err
		}
		if err := u.AddSpaceTime(msg); err != nil {
			return err
		}
	} else {

		// check
		if u.GetMsgByID(msg.ID()) != nil {
			return ErrMsgAlreadyExist
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
		err = u.msgD.AddVertex(msgVertex)
		if err != nil {
			return err
		}
		// update tp
		err = u.updateTimeProof(msg)
		if err != nil {
			return err
		}
		// process the msg
		err = u.processMsg(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddSpaceTime will get all messages save in Universe with same msg.SenderID
func (u *Universe) AddSpaceTime(msg *Message) error {
	if u.GetMsgByID(msg.ID()) == nil {
		return ErrMsgNotFound
	}
	if u.stD != nil && nil != u.stD.GetVertex(msg.SenderID) {
		return ErrTPAlreadyExist
	}
	if !u.CheckUserExist(msg.SenderID) {
		return ErrUserNotExist
	}
	// update time proof
	initialize := true
	for _, id := range u.msgD.GetIDs() {
		if msgSpaceTime := u.GetMsgByID(id); msgSpaceTime != nil && msgSpaceTime.SenderID == msg.SenderID {
			if initialize {
				st, err := u.createSpaceTime(msgSpaceTime)
				if err != nil {
					return err
				}
				stVertex, err := dag.NewVertex(msg.SenderID, st, u.findValidSpaceTime(msg.SenderID)...)
				if err != nil {
					return err
				}
				if u.stD == nil {
					stD, err := dag.NewDAG(stVertex)
					if err != nil {
						return err
					}
					u.stD = stD
				} else {
					if err = u.stD.AddVertex(stVertex); err != nil {
						return err
					}
				}
				initialize = false
			} else {
				if err := u.updateTimeProof(msgSpaceTime); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// CheckUserExist check if the user valid in this Universe
func (u Universe) CheckUserExist(userID common.Hash) bool {
	if nil != u.GetUserByID(userID) {
		return true
	}
	return false
}

// GetUserByID return the user from userD, not userInfo by space time
func (u Universe) GetUserByID(uid common.Hash) *User {
	if v := u.userD.GetVertex(uid); v != nil {
		return v.Value().(*User)
	} else {
		return nil
	}
}

// GetMaxSeq will return the max time proof sequence for
// time proof by the userID
func (u Universe) GetMaxSeq(userID common.Hash) uint64 {
	if vertex := u.stD.GetVertex(userID); vertex != nil {
		return vertex.Value().(*SpaceTime).maxTimeSequence
	} else {
		return 0
	}
}

// GetMsgByID will return the msg by msg.ID()
// nil will be return if msg not exist
func (u Universe) GetMsgByID(mid interface{}) *Message {
	if v := u.msgD.GetVertex(mid); v != nil {
		return v.Value().(*Message)
	} else {
		return nil
	}
}

func (u *Universe) initializeMsgD(msg *Message) error {
	// build msg dag
	msgVertex, err := dag.NewVertex(msg.ID(), msg)
	if err != nil {
		return err
	}
	msgD, err := dag.NewDAG(msgVertex)
	if err != nil {
		return err
	}
	u.msgD = msgD
	return nil
}

func (u *Universe) processMsg(msg *Message) error {
	switch msg.Value.ContentType {
	case TypeText:
		return nil
	case TypeDOB:
		err := u.addUserByMsg(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Universe) createSpaceTime(msg *Message) (*SpaceTime, error) {
	timeVertex, err := dag.NewVertex(msg.ID(), uint64(1))
	if err != nil {
		return nil, err
	}
	timeProofDag, err := dag.NewDAG(timeVertex)
	if err != nil {
		return nil, err
	}
	userStateD, err := dag.NewDAG()
	if err != nil {
		return nil, err
	}

	if nil != u.stD {
		//
		for _, k := range u.stD.GetIDs() {
			stV := u.stD.GetVertex(k)
			usD := stV.Value().(*SpaceTime).userStateD
			if nil != usD.GetVertex(msg.SenderID) {
				// msg.SenderID is valid in this space time
				for _, userID := range usD.GetIDs() {
					if nil == userStateD.GetVertex(userID) {
						userStateVertex, err := dag.NewVertex(userID, &UserInfo{natureState: UserStatusNormal, natureLastCosign: 0}, usD.GetVertex(userID).Parents()...)
						if err != nil {
							return nil, err
						}
						userStateD.AddVertex(userStateVertex)
					}
				}
			}
		}
	} else {
		// create the userState for the first space time
		for _, k := range u.userD.GetIDs() {
			if nil == userStateD.GetVertex(k) {
				userStateVertex, err := dag.NewVertex(k, &UserInfo{natureState: UserStatusNormal, natureLastCosign: 0}, u.userD.GetVertex(k).Parents()...)
				if err != nil {
					return nil, err
				}
				userStateD.AddVertex(userStateVertex)
			}
		}
	}
	return &SpaceTime{maxTimeSequence: timeVertex.Value().(uint64), timeProofD: timeProofDag, userStateD: userStateD}, nil
}

func (u *Universe) updateTimeProof(msg *Message) error {
	if vertex := u.stD.GetVertex(msg.SenderID); vertex != nil {

		st := vertex.Value().(*SpaceTime)
		var currentSeq uint64 = 1
		for _, r := range msg.Reference {
			if r.SenderID == msg.SenderID {
				refSeq := st.timeProofD.GetVertex(r.MsgID).Value().(uint64)
				if currentSeq <= refSeq {
					currentSeq = refSeq + 1
				}
			}
		}
		timeVertex, err := dag.NewVertex(msg.ID(), currentSeq)
		if err != nil {
			return err
		}

		if err := st.timeProofD.AddVertex(timeVertex); err != nil {
			return err
		} else if currentSeq > st.maxTimeSequence {
			st.maxTimeSequence = currentSeq
		}
	}
	return nil
}

// findValidSpaceTime is used when create new space time
func (u *Universe) findValidSpaceTime(senderID common.Hash) []interface{} {
	var ugs []interface{}
	// todo: find valid space time for reference
	return ugs
}

// addUser user to u.userD
// update info of u.stD need other func
func (u *Universe) addUserByMsg(msg *Message) error {
	user, err := CreateNewUser(msg)
	if err != nil {
		return err
	}
	if u.GetUserByID(user.ID()) != nil {
		return ErrUserAlreadyExist
	}

	var dobContent DOBMsgContent
	err = json.Unmarshal(user.DOBMsg.Value.Content, &dobContent)
	if err != nil {
		return err
	}

	var validST []interface{}
	for _, ref := range msg.Reference {
		if stV := u.stD.GetVertex(ref.SenderID); stV != nil {
			st := stV.Value().(*SpaceTime)

			if tp := st.timeProofD.GetVertex(ref.MsgID); tp != nil {
				msgSeq := tp.Value().(uint64)
				p0 := st.userStateD.GetVertex(dobContent.Parents[0].PID)

				if p0 == nil {
					continue
				}
				p1 := st.userStateD.GetVertex(dobContent.Parents[0].PID)
				if p1 == nil {
					continue
				}

				if msgSeq-p0.Value().(*UserInfo).natureLastCosign > rule.REPRODUCTION_INTERVAL &&
					msgSeq-p1.Value().(*UserInfo).natureLastCosign > rule.REPRODUCTION_INTERVAL {
					// update nature last cosign number as msgSeq
					p0.Value().(*UserInfo).natureLastCosign = msgSeq
					p1.Value().(*UserInfo).natureLastCosign = msgSeq
					// append validST
					validST = append(validST, ref.SenderID)
					// add user in this st
					userVertex, err := dag.NewVertex(user.ID(), &UserInfo{natureState: UserStatusNormal, natureLastCosign: 0}, p0, p1)
					if err != nil {
						return err
					}
					st.userStateD.AddVertex(userVertex)
				}
			}
		}
	}

	if len(validST) == 0 {
		return ErrNewUserAddFail
	}
	userVertex, err := dag.NewVertex(user.ID(), user, dobContent.Parents[0].PID, dobContent.Parents[1].PID)
	if err != nil {
		return err
	}
	err = u.userD.AddVertex(userVertex)
	if err != nil {
		return err
	}
	return nil
}
