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

import "testing"

func TestDAG_AddVertex(t *testing.T) {
	v1 := NewVertex("id-1", "hello world")
	v2 := NewVertex("id-2", "hello you")
	dag, err := NewDAG(v1, v2)
	if err != nil {
		t.Errorf("create DAG fail , err : %s", err)
	}
	v3 := NewVertex("id-3", "hello you", "id-1", "id-2")
	if err := dag.AddVertex(v3); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}
	v3_ := NewVertex("id-3", "hello", "id-1", "id-2")
	if err := dag.AddVertex(v3_); err == nil {
		t.Errorf("add vertex should not be success, becasuse id depulicate")
	}
	v4 := NewVertex("id-4", "hello", "id-0")
	if err := dag.AddVertex(v4); err == nil {
		t.Errorf("add vertex should not be success, becasuse not parent exist")
	}
	v5 := NewVertex("id-5", "hello", "id-0", "id-1")
	if err := dag.AddVertex(v5); err == nil {
		t.Errorf("add vertex should not be success, becasuse not all parents exist")
	}
	v6 := NewVertex("id-4", "hello", "id-1", "id-3")
	if err := dag.AddVertex(v6); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

}
