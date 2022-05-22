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

package core

import (
	"github.com/pdupub/go-pdu/identity"
)

// Universe is struct contain all quantums which be received, select and accept by yourself.
// Your universe may same or not with other's, usually your universe only contains part of whole
// exist quantums (not conflict). By methods in Universe, communities be created by quantum and individuals
// be invited into community can be found. Universse also have some aggregate infomation on quantums.
type Universe struct {
	// `json:"address"`

}

// NewUniverse is
func NewUniverse(db *UDB) (*Universe, error) {
	universe := Universe{}
	return &universe, nil
}

func (u *Universe) RecvQuantum(quantum *Quantum) error {
	return nil
}

func (u *Universe) SetAttitude(address identity.Address, level int, judgment string, evidence ...[]Sig) error {

	return nil
}

func (u *Universe) GetAttitude(address identity.Address) (*Attitude, error) {

	return nil, nil
}

func (u *Universe) JoinCommunity(defineSig Sig, address identity.Address) error {
	// update individual.community of creator & initMembers
	return nil
}

func (u *Universe) QueryQuantum(address identity.Address, qType int, pageIndex int, pageSize int, desc bool) []*Quantum {
	return nil
}

func (u *Universe) QueryIndividual(community *Community) []*Individual {
	return nil
}

func (u *Universe) QueryCommunity(community *Community) []*Community {
	return nil
}

func (u *Universe) GetIndividual(address identity.Address) *Individual {
	return nil
}
func (u *Universe) GetQuantum(sig Sig) *Quantum {
	return nil
}
