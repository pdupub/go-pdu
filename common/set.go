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

import (
	"fmt"
	"sync"
)

var Exists = struct{}{}

type Set struct {
	m  map[interface{}]struct{}
	mu sync.RWMutex
}

func NewSet(items ...interface{}) *Set {
	s := &Set{
		m: make(map[interface{}]struct{}),
	}
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s *Set) Add(item interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[item] = Exists
}

func (s *Set) Remove(item interface{}) {
	s.mu.Lock()
	s.mu.Unlock()
	delete(s.m, item)
}

func (s Set) Has(item interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s Set) Size() int {
	return len(s.m)
}

func (s *Set) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[interface{}]struct{})
}

func (s Set) IsEmpty() bool {
	if s.Size() == 0 {
		return true
	}
	return false
}

func (s Set) List() []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []interface{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s Set) String() string {
	var str string
	for item := range s.m {
		str += fmt.Sprintf("%v", item)
	}
	return str
}
