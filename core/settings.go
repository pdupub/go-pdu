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

import "github.com/ethereum/go-ethereum/common"

const (
	FamilyTraceLevel = 5
)

var (
	// same with params.TestAddrs()[0], TODO : reset this before online
	GenesisRoots = []common.Address{common.HexToAddress("0xaf040ed5498f9808550402ebb6c193e2a73b860a")}

	// DefaultGLimit = []*GenerationLimit{
	// 	{ParentsMinSize: 0, ChildrenMaxSize: 3},
	// 	{ParentsMinSize: 1, ChildrenMaxSize: 1},
	// 	{ParentsMinSize: 1, ChildrenMaxSize: 4},
	// 	{ParentsMinSize: 2, ChildrenMaxSize: 10},
	// 	{ParentsMinSize: 2, ChildrenMaxSize: 8},
	// }

	DefaultGLimit = []*GenerationLimit{
		{ParentsMinSize: 0, ChildrenMaxSize: 3},
		{ParentsMinSize: 1, ChildrenMaxSize: 1},
		{ParentsMinSize: 1, ChildrenMaxSize: 2},
		{ParentsMinSize: 2, ChildrenMaxSize: 2},
		{ParentsMinSize: 2, ChildrenMaxSize: 8},
	}
)
