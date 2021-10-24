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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pdupub/go-dag"
)

type Entropy struct {
	*dag.DAG
	bornMsg map[common.Address]string
	lastMsg map[common.Address]string
}

func entropyIDFunc(v *dag.Vertex) (string, error) {
	return common.Bytes2Hex(v.Value().(Event).Key), nil
}

func NewEntropy() (*Entropy, error) {
	return &Entropy{DAG: new(dag.DAG), bornMsg: make(map[common.Address]string), lastMsg: make(map[common.Address]string)}, nil
}

func (e Entropy) GetLastEventID(author common.Address) []byte {
	if id, ok := e.lastMsg[author]; ok {
		return common.Hex2Bytes(id)
	}
	return nil
}

func (e *Entropy) AddEvent(idByte []byte, author common.Address, p *Photon, pidBytes ...[]byte) error {
	id := common.Bytes2Hex(idByte)

	if _, ok := e.DAG.GetVertex(id); ok {
		return ErrPhotonAlreadyExist
	}

	if len(e.DAG.Roots()) == 0 {
		if len(pidBytes) != 0 {
			return ErrPhotonReferenceMissing
		}
		entD, err := dag.New(entropyIDFunc, dag.NewVertex(Event{Key: idByte, Nonce: big.NewInt(0), Author: author, P: p}))
		if err != nil {
			return err
		}
		e.DAG = entD
		// first msg should come from root
		e.bornMsg[author] = id

	} else {
		if err := e.CheckRefs(author, pidBytes...); err != nil {
			return err
		}

		var ps []*dag.Vertex
		maxParentNonce := big.NewInt(0)
		for _, pid := range pidBytes {
			p, ok := e.DAG.GetVertex(common.Bytes2Hex(pid))
			if !ok {
				return ErrPhotonReferenceNotExist
			}
			if maxParentNonce.Cmp(p.Value().(Event).Nonce) < 0 {
				maxParentNonce.Set(p.Value().(Event).Nonce)
			}
			ps = append(ps, p)
		}
		maxParentNonce.Add(maxParentNonce, big.NewInt(1))
		newP := dag.NewVertex(Event{Key: idByte, Nonce: maxParentNonce, Author: author, P: p}, ps...)
		if err := e.DAG.Upsert(newP, 0); err != nil {
			return err
		}
	}

	if p != nil && p.Type == PhotonTypeBorn {
		newID, err := p.GetNewBorn()
		if err != nil {
			return err
		}
		e.bornMsg[newID] = id
	}
	e.lastMsg[author] = id
	return nil
}

func (e *Entropy) IsExist(idByte []byte) bool {
	_, ok := e.DAG.GetVertex(common.Bytes2Hex(idByte))
	if !ok {
		return false
	}
	return true
}

func (e *Entropy) CheckRefs(author common.Address, idBytes ...[]byte) error {
	if len(idBytes) == 0 {
		return ErrPhotonReferenceMissing
	}
	for _, idByte := range idBytes {
		if !e.IsExist(idByte) {
			return ErrPhotonReferenceNotExist
		}
	}
	v, ok := e.DAG.GetVertex(common.Bytes2Hex(idBytes[0]))
	// is not by same author , must be the create ID msg
	if ok && v.Value().(Event).Author != author {
		if v.Value().(Event).P.Type != PhotonTypeBorn {
			return ErrPhotonMissingReferenceByAuthor
		}
		newID, err := v.Value().(Event).P.GetNewBorn()
		if err != nil {
			return err
		}
		if newID != author {
			return ErrPhotonMissingReferenceByAuthor
		}
	}
	// first reference should not have child vertex by same author
	for _, child := range v.Children() {
		if child.Value().(Event).Author == author {
			return ErrPhotonReferenceNotCorrect
		}
	}

	return nil
}

func (e Entropy) marshalPhotonFunc(v *dag.Vertex) (interface{}, error) {
	id := common.Bytes2Hex(v.Value().(Event).Key)
	label := id[0:4] + "..." + id[len(id)-4:]
	return label, nil
}

func (e *Entropy) Dump(keys []string, parentLimit, childLimit int) (*dag.DAGData, error) {
	return e.DAG.Dump(e.marshalPhotonFunc, keys, parentLimit, childLimit)
}

func (e Entropy) DumpByAuthor(author common.Address, keys []string, parentLimit, childLimit int) (*dag.DAGData, error) {
	auD, err := e.GetEvents(author)
	if err != nil {
		return nil, err
	}
	return auD.Dump(e.marshalPhotonFunc, keys, parentLimit, childLimit)
}

func (e Entropy) GetEvents(author common.Address) (*dag.DAG, error) {
	id, ok := e.bornMsg[author]
	if !ok {
		return nil, ErrPhotonNotExist
	}

	auD, err := dag.New(entropyIDFunc, dag.NewVertex(Event{Key: common.Hex2Bytes(id), Nonce: big.NewInt(0), Author: author, P: nil}))
	if err != nil {
		return nil, err
	}
	e.walkFill(auD, id, author)
	return auD, nil
}

func (e Entropy) walkFill(auD *dag.DAG, id string, author common.Address) {
	// return if not find idByte on e.DAG
	v, ok := e.DAG.GetVertex(id)
	if !ok {
		return
	}

	p, ok := auD.GetVertex(id)
	if !ok {
		return
	}
	nextNonce := new(big.Int).Add(big.NewInt(1), p.Value().(Event).Nonce)
	for _, c := range v.Children() {
		if c.Value().(Event).Author == author {
			// save this event , and walFill this key
			val := c.Value().(Event)
			val.Nonce = nextNonce

			newP := dag.NewVertex(val, p)
			cID, err := entropyIDFunc(newP)
			if err != nil {
				continue
			}
			// continue if already exist
			_, ok = auD.GetVertex(cID)
			if ok {
				continue
			}

			if err := auD.Upsert(newP, 0); err != nil {
				continue
			}
			e.walkFill(auD, cID, author)
		}
	}

	return
}
