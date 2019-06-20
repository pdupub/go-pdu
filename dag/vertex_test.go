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

import "testing"

func TestVertex_AddChild(t *testing.T) {
	vertex1 := NewVertex("id-1", "hello world")
	vertex2 := NewVertex("id-2", "hello world again")
	vertex1.AddChild(vertex2)
	if _, ok := vertex1.Children()[vertex2]; !ok {
		t.Errorf("vertex2 should be child ")
	}

	vertex2.AddChild(vertex1)
	if _, ok := vertex2.Children()[vertex1]; !ok {
		t.Errorf("vertex1 should be child ")
	}
}
