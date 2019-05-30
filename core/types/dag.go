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

// In mathematics, particularly graph theory, and computer science, a directed
// acyclic graph (DAG /ˈdæɡ/ (About this soundlisten)), is a finite directed graph
// with no directed cycles. That is, it consists of finitely many vertices and
// edges (also called arcs), with each edge directed from one vertex to another,
// such that there is no way to start at any vertex v and follow a consistently-directed
// sequence of edges that eventually loops back to v again. Equivalently, a DAG is
// a directed graph that has a topological ordering, a sequence of the vertices
// such that every edge is directed from earlier to later in the sequence.
// -------------------from  https://en.wikipedia.org/wiki/Directed_acyclic_graph

package types

import (
	"errors"
	"sync"
)

const (
	maxParentsCount = 10
)

var (
	errRootVertexParentsExist       = errors.New("root vertex parents exist")
	errVertexAlreadyExist           = errors.New("vertex already exist")
	errVertexParentNotExist         = errors.New("parent not exist")
	errVertexParentNumberOutOfRange = errors.New("parent number is out of range")
)

type DAG struct {
	mu    sync.Mutex
	store map[interface{}]*Vertex
}

// NewDAG
func NewDAG(rootVertex ...*Vertex) (*DAG, error) {
	dag := &DAG{
		store: make(map[interface{}]*Vertex),
	}
	for _, vertex := range rootVertex {
		if vertex.Parents().Size() == 0 {
			dag.store[vertex.ID()] = vertex
		} else {
			return nil, errRootVertexParentsExist
		}
	}
	return dag, nil
}

func (d *DAG) AddVertex(vertex *Vertex) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// check the vertex if exist or not
	if _, ok := d.store[vertex.id]; ok {
		return errVertexAlreadyExist
	}
	if vertex.Parents().Size() > maxParentsCount {
		return errVertexParentNumberOutOfRange
	}
	// check parents cloud be found
	for _, parent := range vertex.Parents().List() {
		if _, ok := d.store[parent]; !ok {
			return errVertexParentNotExist
		}
	}
	// add vertex into store
	d.store[vertex.id] = vertex
	// update the parent vertex children
	for _, parent := range vertex.Parents().List() {
		d.store[parent].AddChild(vertex.id)
	}
	return nil
}
