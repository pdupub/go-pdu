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

import "errors"

var (
	// ErrRootVertexParentsExist returns if create new dag with vertex has parent
	ErrRootVertexParentsExist = errors.New("root vertex parents exist")

	// ErrRootNumberOutOfRange returns if numbers of no parents vertex add to a dag is more then the number of roots
	ErrRootNumberOutOfRange = errors.New("root number is out of range")

	// ErrVertexAlreadyExist returns try to add a vertex which already exist in the dag
	ErrVertexAlreadyExist = errors.New("vertex already exist")

	// ErrVertexNotExist returns when try to get vertex by ID, but that vertex not exist
	ErrVertexNotExist = errors.New("vertex not exist")

	// ErrVertexHasChildren returns when try to delete a vertex, but that vertex have children vertex
	ErrVertexHasChildren = errors.New("vertex has children")

	// ErrVertexParentNotExist returns when dag is under the strict rule(default), try to add new vertex,
	// but at least one parent is not exist in the dag
	ErrVertexParentNotExist = errors.New("parent not exist")

	// ErrVertexParentNumberOutOfRange returns parents number if more than setting parents number (default 255)
	ErrVertexParentNumberOutOfRange = errors.New("parent number is out of range")

	// ErrVertexIDInvalid returns try to user a invalid type as id of vertex
	ErrVertexIDInvalid = errors.New("vertex ID invalid")
)
