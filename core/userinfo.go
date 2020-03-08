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

import "fmt"

const (
	// UserStatusNormal is the status of user, will be add more later, like punished...
	UserStatusNormal = iota
)

// UserInfo contain the information except pass by DOBMsg
// the state related to nature rule is start by nature
// the other state start by local
type UserInfo struct {
	natureState      int    // validation state depend on nature rule
	natureLastCosign uint64 // last DOB cosign
	natureLifeMaxSeq uint64 // max time sequence this use can use as reference in this space time
	natureDOBSeq     uint64 // sequence of dob in this space time
	localNickname    string
}

// String used to print user info
func (ui UserInfo) String() string {
	return fmt.Sprintf("localNickname:\t%s\tnatureState:\t%d\tnatureLastCosign:\t%d\tnatureLifeMaxSeq:\t%d\tnatureDOBSeq:\t%d\t", ui.localNickname, ui.natureState, ui.natureLastCosign, ui.natureLifeMaxSeq, ui.natureDOBSeq)
}
