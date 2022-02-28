// Copyright 2021 The PDU Authors
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

package udb

// Value of DType is same with the name of pdu.type in schema (used for expand)
const (
	DTypeQuantum    = "quantum"
	DTypeContent    = "content"
	DTypeIndividual = "individual"
	DTypeAttitude   = "attitude"
	DTypeCommunity  = "community"
)

const (
	DefaultAttitudeLevel = 2 // core.AttitudeAccept
)

// Quantum is similar to type quantum in schema, with all the predicate of quantum and two more item UID & DType
// All fields in Quantum is omitempy for JSON Marshal, not means this field can be omit in record, Only used for
// database operation more efficient.
type Quantum struct {
	UID       string      `json:"uid,omitempty"`
	Sig       string      `json:"quantum.sig,omitempty"`
	Type      int         `json:"quantum.type,omitempty"`
	Refs      []*Quantum  `json:"quantum.refs,omitempty"`
	Contents  []*Content  `json:"quantum.contents,omitempty"`
	Sender    *Individual `json:"quantum.sender,omitempty"`
	Timestamp int         `json:"quantum.timestamp,omitempty"`
	DType     string      `json:"pdu.type,omitempty"`
}

// Content only set one for each quantum, will never update in any condition.
type Content struct {
	UID   string `json:"uid,omitempty"`
	Fmt   int    `json:"content.fmt,omitempty"`
	Data  string `json:"content.data,omitempty"`
	DType string `json:"pdu.type,omitempty"`
}

// Individual record be created only with Address field, community & quantums can be add when new quantum be accept by system.
type Individual struct {
	UID         string       `json:"uid,omitempty"`
	Address     string       `json:"individual.address,omitempty"`
	Communities []*Community `json:"individual.communities,omitempty"`
	Attitude    *Attitude    `json:"individual.attitude,omitempty"`
	Timestamp   int          `json:"individual.timestamp,omitempty"`
	DType       string       `json:"pdu.type,omitempty"`
}

type Attitude struct {
	UID      string     `json:"uid,omitempty"`
	Level    int        `json:"attitude.level,omitempty"`
	Judgment string     `json:"attitude.judgment,omitempty"`
	Evidence []*Quantum `json:"attitude.evidence,omitempty"`
	DType    string     `json:"pdu.type,omitempty"`
}

// Community record be created when rule quantum be accept by system, Invitations & Memebers can be updated.
type Community struct {
	UID          string        `json:"uid,omitempty"`
	Note         *Content      `json:"community.note,omitempty"`
	Base         *Community    `json:"community.base,omitempty"`
	MaxInviteCnt int           `json:"community.maxInviteCnt,omitempty"`
	MinCosignCnt int           `json:"community.minCosignCnt,omitempty"`
	Define       *Quantum      `json:"community.define,omitempty"` // creator is sender of Define
	InitMembers  []*Individual `json:"community.initMembers,omitempty"`
	DType        string        `json:"pdu.type,omitempty"`
}

// UDB is ...
type UDB interface {
	SetSchema() error
	NewQuantum(quantum *Quantum) (qid string, sid string, err error)
	QueryQuantum(address string, qType int, pageIndex int, pageSize int, desc bool) ([]*Quantum, error)
	GetQuantum(sig string) (*Quantum, error)
	GetIndividual(address string) (*Individual, error)
	NewCommunity(community *Community) (cid string, err error)
	GetCommunity(sig string) (*Community, error)
	Update(interface{}) error
	DropData() error
	Close() error
}
