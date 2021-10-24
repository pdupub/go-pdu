// Copyright 2021 The PDU Authors
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

package core

import (
	"errors"
	"testing"

	"github.com/pdupub/go-pdu/params"
)

func TestEntropy(t *testing.T) {
	ent, err := NewEntropy()
	if err != nil {
		t.Error(err)
	}

	if err = ent.AddEvent([]byte("a1"), params.TestAddrs()[0], nil); err != nil {
		t.Error(err)
	}

	if err = ent.AddEvent([]byte("a2"), params.TestAddrs()[0], nil, []byte("a1")); err != nil {
		t.Error(err)
	}

	if err = ent.AddEvent([]byte("a3"), params.TestAddrs()[0], nil, []byte("a2")); err != nil {
		t.Error(err)
	}

	if err = ent.AddEvent([]byte("a4"), params.TestAddrs()[0], nil, []byte("a5")); err == nil {
		t.Error(errors.New("parent should not exist"))
	}

	if lastID := ent.GetLastEventID(params.TestAddrs()[0]); string(lastID) != "a3" {
		t.Error(errors.New("last id not correct"))
	}

	if lastID := ent.GetLastEventID(params.TestAddrs()[1]); lastID != nil {
		t.Error(errors.New("lastID should be nil"))
	}

	if res, err := ent.DumpByAuthor(params.TestAddrs()[0], nil, -1, -1); err != nil {
		t.Error(err)
	} else {
		t.Log("edges", res.Edges)
		t.Log("nodes", res.Nodes)
	}

	for _, root := range ent.DAG.Roots() {
		t.Log(root)
	}
}
