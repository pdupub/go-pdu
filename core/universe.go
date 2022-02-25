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
	"github.com/pdupub/go-pdu/udb"
)

// Universe is struct contain all quantums which be received, select and accept by yourself.
// Your universe may same or not with other's, usually your universe only contains part of whole
// exist quantums (not conflict). By methods in Universe, communities be created by quantum and individuals
// be invited into community can be found. Universse also have some aggregate infomation on quantums.
type Universe struct {
	// `json:"address"`
	db udb.UDB
	// database connection
}

// NewUniverse is
func NewUniverse(db udb.UDB) (*Universe, error) {
	universe := Universe{
		db: db,
	}

	return &universe, nil
}

func (u *Universe) RecvQuantum(quantum *Quantum) error {
	signer, err := quantum.Ecrecover()
	if err != nil {
		return err
	}

	dbQuantum := ToUDBQuantum(quantum, "", "")
	_, _, err = u.db.NewQuantum(dbQuantum)
	if err != nil {
		return err
	}

	if quantum.Type == QuantumTypeCommunity {
		community, err := NewCommunity(quantum)
		if err != nil {
			return err
		}
		dbCommunity := ToUDBCommunity(community, "")
		_, err = u.db.NewCommunity(dbCommunity)
		if err != nil {
			return err
		}

		// creator join community
		u.JoinCommunity(community.Define, signer)
		// init members join community
		for _, v := range dbCommunity.InitMembers {
			u.JoinCommunity(community.Define, identity.HexToAddress(v.Address))
		}
	}
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
	community, _ := u.db.GetCommunity(Sig2Hex(defineSig))
	individual, _ := u.db.GetIndividual(address.Hex())
	individual.Communities = append(individual.Communities, community)

	return u.db.Update(individual)
}

func (u *Universe) QueryQuantum(address identity.Address, qType int, pageIndex int, pageSize int, desc bool) []*Quantum {
	dbQuantums, err := u.db.QueryQuantum(address.Hex(), qType, pageIndex, pageSize, desc)
	if err != nil {
		return nil
	}
	quantum := []*Quantum{}
	for _, v := range dbQuantums {
		q, _, _ := FromUDBQuantum(v)
		quantum = append(quantum, q)
	}
	return quantum
}

func (u *Universe) QueryIndividual(community *Community) []*Individual {
	return nil
}

func (u *Universe) QueryCommunity(community *Community) []*Community {
	return nil
}

func (u *Universe) GetIndividual(address identity.Address) *Individual {
	dbIndividual, err := u.db.GetIndividual(address.Hex())
	if err != nil {
		return nil
	}
	individual, _ := FromUDBIndividual(dbIndividual)
	return individual
}
func (u *Universe) GetQuantum(sig Sig) *Quantum {
	dbQuantum, err := u.db.GetQuantum(Sig2Hex(sig))
	if err != nil {
		return nil
	}
	quantum, _, _ := FromUDBQuantum(dbQuantum)
	return quantum
}
