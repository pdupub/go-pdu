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

// Quantum is similar to type quantum in schema, with all the predicate of quantum and two more item UID & DType
// All fields in Quantum is omitempy for JSON Marshal, not means this field can be omit in record, Only used for
// database operation more efficient.
type Quantum struct {
	UID      string      `json:"uid,omitempty"`
	Sig      string      `json:"quantum.sig,omitempty"`
	Type     int         `json:"quantum.type,omitempty"`
	Refs     []*Quantum  `json:"quantum.refs,omitempty"`
	Contents []*Content  `json:"quantum.contents,omitempty"`
	Sender   *Individual `json:"quantum.sender,omitempty"`
	DType    []string    `json:"dgraph.type,omitempty"`
}

// Content only set one for each quantum, will never update in any condition.
type Content struct {
	UID   string   `json:"uid,omitempty"`
	Fmt   int      `json:"content.fmt,omitempty"`
	Data  string   `json:"content.data,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

// Individual record be created only with Address field, community & quantums can be add when new quantum be accept by system.
type Individual struct {
	UID         string       `json:"uid,omitempty"`
	Address     string       `json:"individual.address,omitempty"`
	Communities []*Community `json:"individual.communities,omitempty"`
	Quantums    []*Quantum   `json:"individual.quantums,omitempty"`
	DType       []string     `json:"dgraph.type,omitempty"`
}

// Community record be created when rule quantum be accept by system, Invitations & Memebers can be updated.
type Community struct {
	UID          string        `json:"uid,omitempty"`
	Base         *Quantum      `json:"community.base,omitempty"`
	Invitations  []*Quantum    `json:"community.invitations,omitempty"`
	MaxInviteCnt int           `json:"community.maxInviteCnt,omitempty"`
	MinCosignCnt int           `json:"community.minCosignCnt,omitempty"`
	Members      []*Individual `json:"community.members,omitempty"`
	Rule         *Quantum      `json:"community.rule,omitempty"`
	DType        []string      `json:"dgraph.type,omitempty"`
}

// UDB is ...
type UDB interface {
	SetQuantum(quantum *Quantum, address string) (uid string, err error)
	GetQuantum(sig string) (*Quantum, error)
	NewIndividual(address string) (uid string, err error)
	GetIndividual(address string) (*Individual, error)
	Close() error
}
