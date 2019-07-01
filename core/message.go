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

type Message struct {
	MsgReference []*MsgReference
	MsgValue     *MsgValue
}

type MsgReference struct {
	Sender *User
	MsgID  []byte
}

func (msg Message) ID() []byte {
	return []byte{}
}

func (msg Message) Value() *MsgValue {

	return nil
}

// ParentsID return the parents id
// Parents are the message referenced by this Message
func (msg Message) ParentsID() [][]byte {

	return nil
}

func (msg Message) ChildrenID() [][]byte {
	return nil
}
