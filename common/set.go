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

import "sync"

var Exists = struct {}{}

type Set struct {
	m map[interface{}]struct{}
	sync.RWMutex
}

func NewSet(items ...interface{}) *Set {
	s := &Set{
		m: make(map[interface{}]struct{}),
	}
	s.Add(items)
	return s
}

func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = Exists
}

func (s *Set) Remove(item interface{}) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]struct{}{}
}

func (s *Set) IsEmpty() bool {
	if s.Size() == 0 {
		return true
	}
	return false
}

func (s *Set) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	var list []interface{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}