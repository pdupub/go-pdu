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
	"encoding/json"
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
func (u *Universe) QueryQuantums(addressHex string, limit int, skip int, asc bool) []*Quantum {

	qRows, err := u.DB.GetQuantumsByAddress(addressHex, limit, skip, asc)
	if err != nil {
		return nil
	}

	qs := []*Quantum{}
	for _, v := range qRows {
		q := &Quantum{}
		if qType, ok := v["qtype"]; ok {
			q.Type = int(qType.(int))
		}

		if contents, ok := v["contents"]; ok {
			err = json.Unmarshal([]byte(contents.(string)), &q.Contents)
			if err != nil {
				return nil
			}
		}

		if refs, ok := v["references"]; ok {
			for _, r := range refs.([]string) {
				q.References = append(q.References, Hex2Sig(r))
			}
		}

		if sig, ok := v["sig"]; ok {
			q.Signature = Hex2Sig(sig.(string))
		}
		qs = append(qs, q)
	}
	return qs
}

func (u *Universe) RecvQuantum(quantum *Quantum) error {

	u.mu.Lock()
	defer u.mu.Unlock()

	// 检查数据库的Quantum表中是否有对应的sig
	_, _, _, _, err := u.DB.GetQuantum(quantum.Signature.toHex())
	if err == nil {
		return errors.New("quantum already exists")
	}

	// 将quantum存到Quantum表
	qcs, err := quantum.Contents.String()
	if err != nil {
		return err
	}

	addr, err := quantum.Ecrecover()
	if err != nil {
		return err
	}

	refs := ""
	for _, v := range quantum.References {
		refs += v.toHex() + ","
	}
	refs = refs[:len(refs)-1]

	err = u.DB.PutQuantum(quantum.Signature.toHex(), qcs, addr.Hex(), refs, quantum.Type)
	if err != nil {
		return err
	}

	return nil
}
