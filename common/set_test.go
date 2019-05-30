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

package common

import "testing"

func TestNewSet(t *testing.T) {
	set := NewSet(1, 2, 3)
	set.Add(1)
	if set.Size() != 3 {
		t.Errorf("should not keep the duplicate item")
	}
	set.Remove(1)
	if set.Size() != 2 {
		t.Errorf("the size of set is incorrect")
	}
	if set.Has(1) {
		t.Errorf("the item 1 should be removed")
	}
	if !set.Has(2) {
		t.Errorf("the item 2 should be in set")
	}
	if len(set.List()) != 2 {
		t.Errorf("list do not have enough items")
	}
	set.Clear()
	if !set.IsEmpty() {
		t.Errorf("set is not cleared")
	}
}
