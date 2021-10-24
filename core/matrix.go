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
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type Matrix struct {
	society *Society // id
	entropy *Entropy // msg
}

func NewMatrix() (*Matrix, error) {
	return &Matrix{society: new(Society), entropy: new(Entropy)}, nil
}

func (m *Matrix) SetSociety(soc *Society) {
	m.society = soc
}

func (m *Matrix) GetSociety() *Society {
	return m.society
}

func (m *Matrix) SetEntropy(ent *Entropy) {
	m.entropy = ent
}

func (m *Matrix) GetEntropy() *Entropy {
	return m.entropy
}

func (m *Matrix) ReceiveMsg(author common.Address, key, content []byte, refs ...[]byte) (*Photon, error) {
	// check photon if exist
	if m.entropy.IsExist(key) {
		return nil, ErrPhotonAlreadyExist
	}

	// check auth if exist
	ap, err := m.society.GetIndividual(author)
	if err != nil {
		return nil, err
	}

	if ap.Addr != author {
		return nil, ErrSocietyIDConflict
	}

	photon := new(Photon)
	if err := json.Unmarshal(content, photon); err != nil {
		return nil, err
	}

	// store photon into entropy
	if err := m.entropy.AddEvent(key, author, photon, refs...); err != nil {
		return nil, err
	}

	// do photon func
	if photon.Type == PhotonTypeBorn {

		newID, err := photon.GetNewBorn()
		if err != nil {
			return nil, err
		}
		parents, err := photon.GetParents()
		if err != nil {
			return nil, err
		}
		if err := m.society.AddIndividual(newID, parents...); err != nil {
			return nil, err
		}
	} else if photon.Type == PhotonTypeProfile {
		profile, err := photon.GetProfile()
		if err != nil {
			return nil, err
		}
		// update profile, author->profile
		if err := m.society.UpdateIndividualProfile(author, profile); err != nil {
			return nil, err
		}
	}

	return photon, nil
}
