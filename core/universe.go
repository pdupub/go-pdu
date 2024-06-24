// Copyright 2024 The PDU Authors
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
	"log"
	"sync"

	"github.com/pdupub/go-pdu/udb"
)

type Universe struct {
	mu sync.RWMutex
	DB *udb.UDB
}

func NewUniverse(dbName string) (*Universe, error) {
	db, err := udb.InitDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	return &Universe{
		DB: db,
	}, nil
}

func (u *Universe) Recv(quantum *Quantum) error {

	u.mu.Lock()
	defer u.mu.Unlock()

	// 检查数据库的Quantum表中是否有对应的sig
	_, err := u.DB.GetQuantum(quantum.Signature.toHex())
	if err == nil {
		return errors.New("quantum already exists")
	}

	// 将quantum存到Quantum表
	qcs, err := quantum.Contents.String()
	if err != nil {
		return err
	}
	err = u.DB.PutQuantum(quantum.Signature.toHex(), qcs)
	if err != nil {
		return err
	}

	// 将quantum中的refs列表存到Reference表
	for _, ref := range quantum.References {
		err = u.DB.PutReference(quantum.Signature.toHex(), ref.toHex())
		if err != nil {
			return err
		}
	}

	return nil
}
