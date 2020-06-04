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

package dag

import (
	"fmt"
)

const (
	// SeekForward is the flag used to find the target from children
	SeekForward = iota

	// SeekBackward is the flag used to find the target from parents
	SeekBackward
)

type funcJudge func(*Vertex, ...interface{}) (bool, error)

// Vertex is a node in DAG
type Vertex struct {
	id       interface{}
	value    interface{}
	parents  map[interface{}]*Vertex
	children map[interface{}]*Vertex
}

// NewVertex create vertex, id, value and parents must be set and is immutable
func NewVertex(id interface{}, value interface{}, parents ...interface{}) (*Vertex, error) {
	v := &Vertex{
		id:       id,
		value:    value,
		parents:  make(map[interface{}]*Vertex),
		children: make(map[interface{}]*Vertex),
	}
	for _, parent := range parents {
		pk := getItemID(parent)
		v.parents[pk] = nil
	}
	return v, nil
}

// ID is the id of vertex
func (v Vertex) ID() interface{} {
	return v.id
}

// ParentIDs is the vertexes which current vertex reference
func (v Vertex) ParentIDs() []interface{} {
	var pks []interface{}
	for k := range v.parents {
		pks = append(pks, k)
	}
	return pks
}

// Parents return the parent vertexes
func (v Vertex) Parents() []*Vertex {
	var pvs []*Vertex
	for _, parent := range v.parents {
		pvs = append(pvs, parent)
	}
	return pvs
}

// Children is the vertexes which reference this vertex
func (v Vertex) Children() []*Vertex {
	var cvs []*Vertex
	for _, child := range v.children {
		cvs = append(cvs, child)
	}
	return cvs
}

// Value is the content of vertex
func (v Vertex) Value() interface{} {
	return v.value
}

// SetValue set the content of vertex
func (v *Vertex) SetValue(value interface{}) {
	v.value = value
}

// AddChild just add the child for this vertex (usually the key or point of child object)
// not add this vertex as parent of the child vertex or check their parents at the same time
func (v *Vertex) AddChild(children ...*Vertex) {
	for _, child := range children {
		v.children[child.ID()] = child
		if parent, ok := child.parents[v.ID()]; !ok || parent == nil {
			child.parents[v.ID()] = v
		}
	}
}

// DelChild remove the children vertexes
// param children is Vertex, *Vertex or ID
func (v *Vertex) DelChild(items ...interface{}) {
	for _, child := range items {
		ck := getItemID(child)
		if cv, ok := v.children[ck]; ok {
			if _, ok := cv.parents[v.ID()]; ok {
				delete(cv.parents, v.ID())
			}
		}
		delete(v.children, ck)
	}
}

// HasParent return true if this vertex have parents
// param children is Vertex, *Vertex or ID
func (v Vertex) HasParent(item interface{}) bool {
	if _, ok := v.parents[getItemID(item)]; !ok {
		return false
	}
	return true
}

// HasChild return true if this vertex have children
func (v Vertex) HasChild(item interface{}) bool {
	if _, ok := v.children[getItemID(item)]; !ok {
		return false
	}
	return true
}

// Seek return ID slice if target is found, empty slice if target can be found
func (v Vertex) Seek(judge funcJudge, maxSteps, direction int, seekArgs ...interface{}) []interface{} {
	var fullPath []interface{}
	if maxSteps <= 0 {
		return fullPath
	}
	scope := v.children
	if direction == SeekBackward {
		scope = v.parents
	}

	for _, v := range scope {
		if found, err := judge(v, seekArgs...); err == nil && found {
			fullPath = []interface{}{v.ID()}
			break
		}
		if path := v.Seek(judge, maxSteps-1, direction, seekArgs...); len(path) > 0 {
			fullPath = append(path, v.ID())
			break
		}
	}

	return fullPath
}

// String used to print the content of vertex
func (v Vertex) String() string {
	result := fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value: %v\n", v.id, len(v.ParentIDs()), len(v.Children()), v.value)
	return result
}

func getItemID(item interface{}) interface{} {
	switch item.(type) {
	case *Vertex:
		return item.(*Vertex).ID()
	case Vertex:
		return item.(Vertex).ID()
	default:
		return item
	}
}
