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
	errRootNumberOutOfRange         = errors.New("root number is out of range")
	errVertexAlreadyExist           = errors.New("vertex already exist")
	errVertexNotExist               = errors.New("vertex not exist")
	errVertexHasChildren            = errors.New("vertex has children")
	errVertexParentNotExist         = errors.New("parent not exist")
	errVertexParentNumberOutOfRange = errors.New("parent number is out of range")
)

const (
	defaultMaxParentsCount = 255
)

// DAG is directed acyclic graph
type DAG struct {
	mu              sync.Mutex
	maxParentsCount int
	strict          bool
	store           map[interface{}]*Vertex
	ids             []interface{}
	rufd            uint // unfilled root count
}

// NewDAG create new DAG by root vertexes
func NewDAG(rootCnt uint, rootVertex ...*Vertex) (*DAG, error) {
	dag := &DAG{
		maxParentsCount: defaultMaxParentsCount,
		strict:          true,
		store:           make(map[interface{}]*Vertex),
		ids:             []interface{}{},
		rufd:            rootCnt,
	}
	for _, vertex := range rootVertex {
		if dag.rufd == 0 {
			return nil, errRootNumberOutOfRange
		} else if len(vertex.Parents()) == 0 {
			dag.store[vertex.ID()] = vertex
			dag.ids = append(dag.ids, vertex.ID())
			dag.rufd--
		} else {
			return nil, errRootVertexParentsExist
		}
	}
	return dag, nil
}

// IsStrict return if all parents must exist when add vertex
func (d *DAG) IsStrict() bool {
	return d.strict
}

// RemoveStrict set strict to false, mean at least one parents exist in dag,
// the vertex can be added, and the strict rule can not from false to true.
func (d *DAG) RemoveStrict() {
	d.strict = false
}

// SetMaxParentsCount set the max number of parents one vertex can get
func (d *DAG) SetMaxParentsCount(maxCount int) {
	d.maxParentsCount = maxCount
}

// GetMaxParentsCount get the max number of parents
func (d *DAG) GetMaxParentsCount() int {
	return d.maxParentsCount
}

// GetVertex can get vertex by ID
func (d *DAG) GetVertex(id interface{}) *Vertex {
	if _, ok := d.store[id]; !ok {
		return nil
	}
	return d.store[id]
}

// AddVertex is add vertex to DAG
func (d *DAG) AddVertex(vertex *Vertex) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	// check the vertex if exist or not
	if _, ok := d.store[vertex.ID()]; ok {
		return errVertexAlreadyExist
	}

	if len(vertex.Parents()) > d.maxParentsCount {
		return errVertexParentNumberOutOfRange
	}
	// check parents cloud be found
	sequenceExist := false
	for _, pid := range vertex.Parents() {
		if _, ok := d.store[pid]; !ok && d.strict {
			return errVertexParentNotExist
		}
		sequenceExist = true
	}
	if !sequenceExist {
		if d.rufd == 0 {
			return errRootNumberOutOfRange
		}
		d.rufd--
	}
	// add vertex into store
	d.store[vertex.ID()] = vertex
	d.ids = append(d.ids, vertex.ID())
	// update the parent vertex children
	for _, pid := range vertex.Parents() {
		if v, ok := d.store[pid]; ok {
			v.AddChild(vertex.ID())
		}
	}
	return nil
}

// DelVertex is used to remove vertex from DAG
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
	for i := 0; i < len(d.ids); i++ {
		if d.ids[i] == id {
			d.ids = append(d.ids[:i], d.ids[i+1:]...)
			break
		}
	}
	return nil
}

// GetIDs get id list of DAG
func (d *DAG) GetIDs() []interface{} {
	return d.ids
}

// String is used to print the DAG content
func (d *DAG) String() string {
	result := fmt.Sprintf("maxParentsCount : %d - storeSize : %d \n", d.maxParentsCount, len(d.store))
	for k, v := range d.store {
		result += fmt.Sprintf("k = %v \n", k)
		result += fmt.Sprintf("v = %v \n", v)
	}
	return result
}
