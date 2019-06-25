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
	"errors"
	"fmt"
)

var (
	errVertexIDInvalid = errors.New("vertex ID invalid")
)

type Vertex struct {
	id       interface{}
	value    interface{}
	parents  map[interface{}]struct{}
	children map[interface{}]struct{}
}

// NewVertex create vertex, id, value and parents must be set and is immutable
// parents cloud be Vertex or just key
func NewVertex(id interface{}, value interface{}, parents ...interface{}) (*Vertex, error) {
	// Vertex & *Vertex can not be used as Vertex.ID
	switch id.(type) {
	case *Vertex:
		return nil, errVertexIDInvalid
	case Vertex:
		return nil, errVertexIDInvalid
	}

	v := &Vertex{
		id:       id,
		value:    value,
		parents:  make(map[interface{}]struct{}),
		children: make(map[interface{}]struct{}),
	}
	for _, parent := range parents {
		pk := getItemID(parent)
		v.parents[pk] = struct{}{}
	}
	return v, nil
}

func (v Vertex) ID() interface{} {
	return v.id
}

func (v Vertex) Parents() []interface{} {
	var pks []interface{}
	for k := range v.parents {
		pks = append(pks, k)
	}
	return pks
}

func (v Vertex) Children() []interface{} {
	var cks []interface{}
	for k := range v.children {
		cks = append(cks, k)
	}
	return cks
}

func (v Vertex) Value() interface{} {
	return v.value
}

func (v *Vertex) SetValue(value interface{}) {
	v.value = value
}

// AddChild just add the child for this vertex (usually the key or point of child object)
// not add this vertex as parent of the child vertex or check their parents at the same time
func (v *Vertex) AddChild(children ...interface{}) {
	for _, child := range children {
		ck := getItemID(child)
		v.children[ck] = struct{}{}
	}
}

func (v *Vertex) DelChild(children ...interface{}) {
	for _, child := range children {
		ck := getItemID(child)
		delete(v.children, ck)
	}
}

func (v Vertex) HasParent(item interface{}) bool {
	if _, ok := v.parents[getItemID(item)]; !ok {
		return false
	}
	return true
}

func (v Vertex) HasChild(item interface{}) bool {
	if _, ok := v.children[getItemID(item)]; !ok {
		return false
	}
	return true
}

func (v Vertex) String() string {
	result := fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value: %v\n", v.id, len(v.Parents()), len(v.Children()), v.value)
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
