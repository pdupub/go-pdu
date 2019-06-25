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

package dag

import (
	"errors"
	"fmt"
	"sync"
)

var (
	errRootVertexParentsExist       = errors.New("root vertex parents exist")
	errVertexAlreadyExist           = errors.New("vertex already exist")
	errVertexNotExist               = errors.New("vertex not exist")
	errVertexHasChildren            = errors.New("vertex has children")
	errVertexParentNotExist         = errors.New("parent not exist")
	errVertexParentNumberOutOfRange = errors.New("parent number is out of range")
)

type DAG struct {
	mu              sync.Mutex
	maxParentsCount int // 0 = unlimited
	store           map[interface{}]*Vertex
}

// NewDAG
func NewDAG(rootVertex ...*Vertex) (*DAG, error) {
	dag := &DAG{
		store: make(map[interface{}]*Vertex),
	}
	for _, vertex := range rootVertex {
		if len(vertex.Parents()) == 0 {
			dag.store[vertex.ID()] = vertex
		} else {
			return nil, errRootVertexParentsExist
		}
	}
	return dag, nil
}

// SetMaxParentsCount
func (d *DAG) SetMaxParentsCount(maxCount int) {
	d.maxParentsCount = maxCount
}

// GetMaxParentsCount
func (d DAG) GetMaxParentsCount() int {
	return d.maxParentsCount
}

// AddVertex
func (d *DAG) AddVertex(vertex *Vertex) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// check the vertex if exist or not
	if _, ok := d.store[vertex.ID()]; ok {
		return errVertexAlreadyExist
	}

	if d.maxParentsCount != 0 && len(vertex.Parents()) > d.maxParentsCount {
		return errVertexParentNumberOutOfRange
	}
	// check parents cloud be found
	for _, pid := range vertex.Parents() {
		if _, ok := d.store[pid]; !ok {
			return errVertexParentNotExist
		}
	}
	// add vertex into store
	d.store[vertex.ID()] = vertex
	// update the parent vertex children
	for _, pid := range vertex.Parents() {
		d.store[pid].AddChild(vertex.ID())
	}
	return nil
}

// DelVertex
func (d *DAG) DelVertex(id interface{}) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// check the key exist and no children
	if v, ok := d.store[id]; !ok {
		return errVertexNotExist
	} else if len(v.Children()) > 0 {
		return errVertexHasChildren
	} else {
		// remove this child vertex from parents
		for _, pid := range v.Parents() {
			if parent, ok := d.store[pid]; ok {
				parent.DelChild(id)
			}
		}
	}
	delete(d.store, id)
	return nil
}

func (d DAG) String() string {
	result := fmt.Sprintf("maxParentsCount : %d - storeSize : %d \n", d.maxParentsCount, len(d.store))
	for k, v := range d.store {
		result += fmt.Sprintf("k = %v \n", k)
		result += fmt.Sprintf("v = %v \n", v)
	}
	return result
}
