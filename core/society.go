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
	"github.com/pdupub/go-dag"
)

// GenerationLimit is limit of seek relatives of one user
type GenerationLimit struct {
	ParentsMinSize  int // minimize parents cnt need to sign for creating ID on current generation
	ChildrenMaxSize int // maximize children cnt current generation can be envolve in creating
}

// Society is main struct store all users of system
type Society struct {
	GLimit []*GenerationLimit
	*dag.DAG
	profiles map[common.Address]*PProfile
}

func societyIDFunc(v *dag.Vertex) (string, error) { return v.Value().(Individual).Addr.Hex(), nil }

// NewSociety is used to create new user system
func NewSociety(roots ...common.Address) (*Society, error) {
	var vs []*dag.Vertex
	if len(roots) == 0 {
		return nil, dag.ErrRootNumberOutOfRange
	}
	for _, root := range roots {
		vs = append(vs, dag.NewVertex(Individual{Addr: root, Level: 0}))
	}
	socD, err := dag.New(societyIDFunc, vs...)
	if err != nil {
		return nil, err
	}

	society := &Society{
		GLimit:   DefaultGLimit,
		DAG:      socD,
		profiles: make(map[common.Address]*PProfile),
	}

	return society, nil
}

// UpdateIndividualProfile update profile information which come from Quantum of QuantumProfileType
func (s *Society) UpdateIndividualProfile(author common.Address, profile *PProfile) error {
	s.profiles[author] = profile
	return nil
}

// GetIndividualProfile get profile by author
func (s *Society) GetIndividualProfile(author common.Address) *PProfile {
	if profile, ok := s.profiles[author]; ok {
		return profile
	}
	return new(PProfile)
}

// AddIndividual create new user
func (s *Society) AddIndividual(id common.Address, pids ...common.Address) error {
	if _, ok := s.DAG.GetVertex(id.Hex()); ok {
		return dag.ErrVertexAlreadyExist
	}

	if len(pids) == 0 {
		return dag.ErrVertexParentNotExist
	}

	parents := make(map[common.Address]struct{})
	for _, pid := range pids {
		parents[pid] = struct{}{}
	}

	var ps []*dag.Vertex
	maxParentLevel := uint(0)
	for pid := range parents {
		p, ok := s.DAG.GetVertex(pid.Hex())
		if !ok {
			return dag.ErrVertexParentNotExist
		}
		if len(p.Children()) >= s.GLimit[p.Value().(Individual).Level].ChildrenMaxSize {
			return ErrAddIndividualBeyondChildrenMaxLimit
		}
		if maxParentLevel < p.Value().(Individual).Level {
			maxParentLevel = p.Value().(Individual).Level
		}
		ps = append(ps, p)
	}

	currentLevel := maxParentLevel + 1
	if len(parents) < s.GLimit[currentLevel].ParentsMinSize {
		return ErrAddIndividualWithoutEnoughParents
	}

	newID := dag.NewVertex(Individual{Addr: id, Level: currentLevel}, ps...)
	if err := s.DAG.Upsert(newID, 0); err != nil {
		return err
	}

	return nil
}

// GetIndividual return individual if exist
func (s *Society) GetIndividual(id common.Address) (*Individual, error) {
	v, ok := s.DAG.GetVertex(id.Hex())
	if !ok {
		return nil, ErrIndividualNotExistInSociety
	}

	p := v.Value().(Individual)
	return &p, nil
}

func (s Society) marshalIndividualFunc(v *dag.Vertex) (interface{}, error) {
	addr := v.Value().(Individual).Addr
	label := addr.Hex()[2:6] + "..." + addr.Hex()[38:42]
	return label, nil
}

// Dump dump Society to JSON
func (s *Society) Dump(keys []string, parentLimit, childLimit int) (*dag.DAGData, error) {
	return s.DAG.Dump(s.marshalIndividualFunc, keys, parentLimit, childLimit)
}

// MarshalJSON encode Society to JSON
func (s *Society) MarshalJSON() ([]byte, error) {
	fmtData, err := s.Dump(nil, -1, -1)
	if err != nil {
		return nil, nil
	}
	return json.Marshal(fmtData)
}

// UnmarshalJSON decode JSON to Society (unfinished)
func (s *Society) UnmarshalJSON(b []byte) error {
	fmtData := make(map[string]interface{})
	if err := json.Unmarshal(b, &fmtData); err != nil {
		return err
	}
	// TODO :
	// rebuild society

	return nil
}
