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

import "github.com/pdupub/go-pdu/udb"

// Universe is struct contain all quantums which be received, select and accept by yourself.
// Your universe may same or not with other's, usually your universe only contains part of whole
// exist quantums (not conflict). By methods in Universe, communities be created by quantum and individuals
// be invited into community can be found. Universse also have some aggregate infomation on quantums.
type Universe struct {
	// `json:"address"`
	db udb.UDB
	// database connection
}

func NewUniverse(db udb.UDB) (*Universe, error) {
	universe := Universe{
		db: db,
	}

	universe.db.NewIndividual("0xabc")

	return &universe, nil
}
