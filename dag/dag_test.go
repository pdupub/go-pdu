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
	"testing"
)

func TestDAG_AddVertex(t *testing.T) {
	v1, _ := NewVertex("id-1", "hello world")
	v2, _ := NewVertex("id-2", "hello you")
	dag, err := NewDAG(2, v1, v2)
	if err != nil {
		t.Errorf("create DAG fail , err : %s", err)
	}

	if len(dag.GetIDs()) != 2 {
		t.Errorf("id number not match, should be %d, dag getIDs is %d", 2, dag.GetIDs())
	}

	v3, _ := NewVertex("id-3", "hello you", v1, v2)
	if err := dag.AddVertex(v3); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	if len(dag.GetIDs()) != 3 {
		t.Errorf("id number not match, should be %d, dag getIDs is %d", 3, dag.GetIDs())
	}

	v3b, _ := NewVertex("id-3", "hello", v1, v2)
	if err := dag.AddVertex(v3b); err != ErrVertexAlreadyExist {
		t.Errorf("add vertex should not be success, becasuse id depulicate")
	}

	v0, _ := NewVertex("id-0", "hello you")
	v4, _ := NewVertex("id-4", "hello", v0)
	if err := dag.AddVertex(v4); err != ErrVertexParentNotExist {
		t.Errorf("add vertex should not be success, becasuse not parent exist")
	}

	v5, _ := NewVertex("id-5", "hello", v0, v1)
	if err := dag.AddVertex(v5); err != ErrVertexParentNotExist {
		t.Errorf("add vertex should not be success, becasuse not all parents exist")
	}

	if len(dag.GetIDs()) != 3 {
		t.Errorf("id number not match, should be %d, dag getIDs is %d", 3, dag.GetIDs())
	}

	v6, _ := NewVertex("id-4", "hello", v1, v3)
	if err := dag.AddVertex(v6); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}
}

func TestDAG_AddVertex_RootCnt(t *testing.T) {
	v1, _ := NewVertex("id-1", "hello world")
	v2, _ := NewVertex("id-2", "hello you")

	dag, err := NewDAG(3)
	if err != nil {
		t.Errorf("create DAG fail , err : %s", err)
	}

	if len(dag.GetIDs()) != 0 {
		t.Errorf("id number not match, should be %d, dag getIDs is %d", 0, dag.GetIDs())
	}

	if err := dag.AddVertex(v1); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	if err := dag.AddVertex(v2); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	v3, _ := NewVertex("id-3", "hello you", v1, v2)
	if err := dag.AddVertex(v3); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	v4, _ := NewVertex("id-4", "hello you", v1, v3)
	if err := dag.AddVertex(v4); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	v5, _ := NewVertex("id-5", "hello you too")
	if err := dag.AddVertex(v5); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	v6, _ := NewVertex("id-6", "hello you too")
	if err := dag.AddVertex(v6); err != ErrRootNumberOutOfRange {
		t.Errorf("add vertex should fail, err should be %s not %s", ErrRootNumberOutOfRange, err)
	}
}

func TestDAG_AddVertex_Strict(t *testing.T) {
	v1, _ := NewVertex("id-1", "hello world")
	v2, _ := NewVertex("id-2", "hello you")

	dag, err := NewDAG(2, v1, v2)
	if err != nil {
		t.Errorf("create DAG fail , err : %s", err)
	}

	v3, _ := NewVertex("id-3", "hello you", v1, v2)
	v4, _ := NewVertex("id-4", "hello you", v1, v3)
	v5, _ := NewVertex("id-5", "hello you too", v2, v3)

	if err := dag.AddVertex(v4); err != ErrVertexParentNotExist {
		t.Errorf("add vertex should fail with err %s ,but: %s", ErrVertexParentNotExist, err)
	}
	dag.RemoveStrict()
	if err := dag.AddVertex(v4); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}
	if children, ok := dag.awcf[v3.ID()]; !ok {
		t.Errorf("awaiting confirmation should have key %s", v3.ID())
	} else if len(children) != 1 || children[0] != v4.ID() {
		t.Errorf("children should contain 1 child ,which ID is %s", v4.ID())
	}

	if err := dag.AddVertex(v5); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}
	if children, ok := dag.awcf[v3.ID()]; !ok {
		t.Errorf("awaiting confirmation should have key %s", v3.ID())
	} else if len(children) != 2 || children[0] != v4.ID() || children[1] != v5.ID() {
		t.Errorf("children should contain 2 children ,which IDs are %s and %s", v4.ID(), v5.ID())
	}

	if err := dag.AddVertex(v3); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}
	if _, ok := dag.awcf[v3.ID()]; ok {
		t.Errorf("v3.ID %s should be deleted from awaiting confirmation", v3.ID())
	}

}

func TestDAG_DelVertex(t *testing.T) {
	v1, _ := NewVertex("id-1", "hello world")
	v2, _ := NewVertex("id-2", "hello you")
	dag, err := NewDAG(2, v1, v2)
	if err != nil {
		t.Errorf("create DAG fail , err : %s", err)
	}

	v3, _ := NewVertex("id-3", "hello you", v1, v2)
	if err := dag.AddVertex(v3); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	v4, _ := NewVertex("id-4", "hello", v1, v3)
	if err := dag.AddVertex(v4); err != nil {
		t.Errorf("add vertex fail, err : %s", err)
	}

	if err := dag.DelVertex("id-1"); err != ErrVertexHasChildren {
		t.Errorf("del vertex should not be success, because child exist")
	}

	if err := dag.DelVertex("id-5"); err != ErrVertexNotExist {
		t.Error("del vertex should not be success, because id not exist")
	}

	if err := dag.DelVertex("id-4"); err != nil {
		t.Errorf("del vertex fail, err : %s", err)
	}

	if err := dag.DelVertex("id-3"); err != nil {
		t.Errorf("del vertex fail, err : %s", err)
	}

	if err := dag.DelVertex("id-2"); err != nil {
		t.Errorf("del vertex fail, err : %s", err)
	}
	if err := dag.DelVertex("id-2"); err != ErrVertexNotExist {
		t.Errorf("del vertex fail should fail, because this key already being removed, but err is : %s", err)
	}
}
