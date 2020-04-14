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

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/dag"
)

// Universe is the struct contain all the information can be received and validated.
// It is built with two precreate users as root users. Any message should be validated before
// add into msgD, which is a DAG used to store message. If new user create from valid message,
// the user should add into userD, a DAG used to store all users. The universe contain at least
// one spacetime, each spacetime base on one user msg as time line. validation of message mean
// this message should be valid at least in one of spacetime in stD. Information in local universe
// is only part of information in whole decentralized system.
type Universe struct {
	msgD  *dag.DAG // contain all messages valid in at least one spacetime
	userD *dag.DAG // contain all users valid in at least one spacetime (strict)
	stD   *dag.DAG // contain all spacetime, which could be diff by selecting (strict)
}

// NewUniverse create Universe with two user with diff gender as root users
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
	userD, err := dag.NewDAG(2, EveVertex, AdamVertex)
	if err != nil {
		return nil, err
	}
	userD.SetMaxParentsCount(2)
	return &Universe{userD: userD}, nil
}

// AddMsg will check if the message from valid user, who is validated in at least one spacetime
// (in stD). Then new message will be added into Universe and update time proof if msg.SenderID
// is any spacetime based on.
func (u *Universe) AddMsg(msg *Message) error {
	if !u.CheckUserExist(msg.SenderID) {
		return ErrUserNotExist
	}
	if u.msgD == nil {
		if err := u.initializeMsgD(msg); err != nil {
			return err
		}
		if err := u.AddSpaceTime(msg, nil); err != nil {
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

// GetSpaceTimeIDs get ids in of spacetime (list of msg.SenderID of each spacetime)
func (u *Universe) GetSpaceTimeIDs() []common.Hash {
	var ids []common.Hash
	if u.stD != nil {
		for _, id := range u.stD.GetIDs() {
			ids = append(ids, id.(common.Hash))
		}
	}
	return ids
}

// AddSpaceTime will add spacetime in Universe with msg.SenderID, and follow the
// time sequence from ref.
func (u *Universe) AddSpaceTime(msg *Message, ref *MsgReference) error {
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
	startRecord := false
	for _, id := range u.msgD.GetIDs() {
		if id == msg.ID() {
			startRecord = true
		}
		if msgSpaceTime := u.GetMsgByID(id); msgSpaceTime != nil && msgSpaceTime.SenderID == msg.SenderID && startRecord {
			if initialize {
				if err := u.initializeSpaceTime(msgSpaceTime, ref); err != nil {
					return err
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

func (u *Universe) initializeSpaceTime(msgSpaceTime *Message, ref *MsgReference) error {
	st, err := NewSpaceTime(u, msgSpaceTime, ref)
	if err != nil {
		return err
	}
	var stVertex *dag.Vertex
	if ref != nil {
		stVertex, err = dag.NewVertex(msgSpaceTime.SenderID, st, ref.SenderID)
	} else {
		stVertex, err = dag.NewVertex(msgSpaceTime.SenderID, st)
	}
	if err != nil {
		return err
	}
	if u.stD == nil {
		stD, err := dag.NewDAG(1, stVertex)
		if err != nil {
			return err
		}
		u.stD = stD
	} else {
		if err = u.stD.AddVertex(stVertex); err != nil {
			return err
		}
	}
	return nil
}

// CheckUserExist check if the user valid in the Universe
func (u Universe) CheckUserExist(userID common.Hash) bool {
	if nil != u.GetUserByID(userID) {
		return true
	}
	return false
}

// GetUserByID return the user from userD, not userInfo by space time
func (u Universe) GetUserByID(userID common.Hash) *User {
	if v := u.userD.GetVertex(userID); v != nil {
		return v.Value().(*User)
	}
	return nil
}

// GetMaxSeq return the max time proof sequence
func (u Universe) GetMaxSeq(spacetimeID common.Hash) uint64 {
	if vertex := u.stD.GetVertex(spacetimeID); vertex != nil {
		return vertex.Value().(*SpaceTime).maxTimeSequence
	}
	return 0
}

// GetUserIDs return userIDs in this space time
func (u Universe) GetUserIDs(spacetimeID common.Hash) []common.Hash {
	var userIDs []common.Hash
	if u.stD != nil {
		if vertex := u.stD.GetVertex(spacetimeID); vertex != nil {
			userIDs = vertex.Value().(*SpaceTime).GetUserIDs()
		}
	}
	return userIDs
}

// GetUserInfo return the user info in space time
// return nil if not find user
func (u Universe) GetUserInfo(userID common.Hash, spacetimeID common.Hash) *UserInfo {
	if u.stD != nil {
		if stVertex := u.stD.GetVertex(spacetimeID); stVertex != nil {
			return stVertex.Value().(*SpaceTime).GetUserInfo(userID)
		}
	}
	return nil
}

// GetMsgByID will return the msg by msg.ID()
// nil will be return if msg not exist
func (u Universe) GetMsgByID(msgID interface{}) *Message {
	if v := u.msgD.GetVertex(msgID); v != nil {
		return v.Value().(*Message)
	}
	return nil
}

// initializeMsgD only run once to create u.msgD by initial message, and the DAG
// will remove strict rule, so msgD can accept new message if at least one of
// reference exist in whole universe.
func (u *Universe) initializeMsgD(msg *Message) error {
	// build msg dag
	msgVertex, err := dag.NewVertex(msg.ID(), msg)
	if err != nil {
		return err
	}
	msgD, err := dag.NewDAG(1, msgVertex)
	if err != nil {
		return err
	}
	msgD.RemoveStrict()
	u.msgD = msgD
	return nil
}

func (u *Universe) processMsg(msg *Message) error {
	switch msg.Value.ContentType {
	case TypeText:
		return nil
	case TypeBirth:
		err := u.addUserByMsg(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Universe) updateTimeProof(msg *Message) error {
	if vertex := u.stD.GetVertex(msg.SenderID); vertex != nil {
		return vertex.Value().(*SpaceTime).UpdateTimeProof(msg)
	}
	return nil
}

// addUser user to u.userD
// update info of u.stD need other func
func (u *Universe) addUserByMsg(msg *Message) error {
	user, err := CreateNewUser(u, msg)
	if err != nil {
		return err
	}
	if u.GetUserByID(user.ID()) != nil {
		return ErrUserAlreadyExist
	}

	var birthContent BirthMsgContent
	err = json.Unmarshal(user.BirthMsg.Value.Content, &birthContent)
	if err != nil {
		return err
	}

	userAdded := false
	for _, ref := range msg.Reference {
		if err := u.addUserToSpaceTime(ref, birthContent, user); err != nil {
			continue
		}
		// at least add into one space time
		userAdded = true
	}

	if !userAdded {
		return ErrNewUserAddFail
	}

	userVertex, err := dag.NewVertex(user.ID(), user, birthContent.Parents[0].UserID, birthContent.Parents[1].UserID)
	if err != nil {
		return err
	}
	err = u.userD.AddVertex(userVertex)
	if err != nil {
		return err
	}
	return nil
}

// addUserToSpaceTime used to add new user to spacetime base on ref.SenderID, the age of parents in this spacetime
// should fit the nature rule.
// TODO: ref.SenderID not must be spacetime, the new user's life length can be calculated by any ref msg.
func (u *Universe) addUserToSpaceTime(ref *MsgReference, birthContent BirthMsgContent, user *User) error {
	if vertex := u.stD.GetVertex(ref.SenderID); vertex != nil {
		return vertex.Value().(*SpaceTime).AddUser(ref, birthContent, user)
	}
	return ErrAddUserToSpaceTimeFail
}
