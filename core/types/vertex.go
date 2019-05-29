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

package types

import (
	"fmt"
	"github.com/pdupub/go-pdu/common"
)

type Vertex struct {
	id      	interface{}
	value  		interface{}
	parents		*common.Set
	children	*common.Set
}

// NewVertex create vertex, id, value and parents must be set and is immutable
// parents cloud be Vertex or just key
func NewVertex(id interface{},value interface{}, parents ...interface{}) *Vertex {
	v := &Vertex{
		id:id,
		value:value,
		parents:common.NewSet(),
		children:common.NewSet(),
	}
	for _,parent := range parents { v.parents.Add(parent)}
	return v
}

func (v Vertex) ID() interface{} {
	return v.id
}

func (v Vertex) Parents() common.Set {
	return *v.parents
}

func (v Vertex) Children() common.Set {
	return *v.children
}

func (v Vertex) Value() interface{} {
	return v.value
}

func (v *Vertex) AddChild(children ... interface{}) {
	for _, child := range children{
		v.children.Add(child)
	}
}

func (v Vertex) String() string {
	result := fmt.Sprintf("ID: %s - Parents: %d - Children: %d - Value: %v\n", v.id, v.Parents().Size(), v.Children().Size(), v.value)
	return result
}

