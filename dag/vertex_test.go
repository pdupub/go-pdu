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
	"testing"
)

func findTargetTest(v *Vertex, args ...interface{}) (bool, error) {
	if len(args) != 1 {
		return false, errors.New("argument is missing")
	}

	if v.ID() == args[0] {
		return true, nil
	}
	return false, nil
}

func TestVertex(t *testing.T) {
	vertex1, _ := NewVertex("id-1", "hello world")
	vertex2, _ := NewVertex("id-2", "hello world again")
	vertex3, _ := NewVertex("id-3", "hello world again")
	vertex4, _ := NewVertex("id-4", "hello world again")
	vertex5, _ := NewVertex("id-5", "hello world again")
	vertex6, _ := NewVertex("id-6", "hello world again")

	vertex7, _ := NewVertex("id-7", "hello world again")
	vertex8, _ := NewVertex("id-8", "hello world again")
	vertex9, _ := NewVertex("id-9", "hello world again")
	vertex10, _ := NewVertex("id-10", "hello world again")

	vertex1.AddChild(vertex2)
	if !vertex1.HasChild(vertex2) {
		t.Errorf("vertex2 should be child ")
	}
	vertex2.AddChild(vertex3)
	vertex3.AddChild(vertex4)
	vertex4.AddChild(vertex5)
	vertex5.AddChild(vertex6)

	// add some noise
	vertex3.AddChild(vertex7)
	vertex3.AddChild(vertex8)
	vertex7.AddChild(vertex9)
	vertex4.AddChild(vertex10)

	pathRes := vertex1.Seek(findTargetTest, 5, SeekForward, "id-6")
	if len(pathRes) != 5 || pathRes[0] != vertex6.ID() {
		t.Error("path can not be found")
	}

	pathRes = vertex1.Seek(findTargetTest, 8, SeekForward, "id-6")
	if len(pathRes) != 5 || pathRes[0] != vertex6.ID() {
		t.Error("path can not be found")
	}

	pathRes = vertex1.Seek(findTargetTest, 4, SeekForward, "id-6")
	if len(pathRes) != 0 {
		t.Error("path should not be found")
	}

	pathRes = vertex3.Seek(findTargetTest, 3, SeekForward, "id-6")
	if len(pathRes) != 3 || pathRes[0] != vertex6.ID() {
		t.Error("path can not be found")
	}

	pathRes = vertex6.Seek(findTargetTest, 5, SeekBackward, "id-1")
	if len(pathRes) != 5 || pathRes[0] != vertex1.ID() {
		t.Error("path can not be found")
	}

	vertex2.AddChild(vertex1)
	if vertex2.HasChild(vertex2) {
		t.Errorf("vertex1 should be child ")
	}

	vertex2.DelChild(vertex1)
	if vertex2.HasChild(vertex1) {
		t.Errorf("vertex1 should be removed ")
	}

	vertex1.SetValue("nihao")
	if vertex1.Value() != "nihao" {
		t.Errorf("vertex1 set value fail")
	}
}
