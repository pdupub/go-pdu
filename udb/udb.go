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

package udb

import (
	"log"

	"go.etcd.io/bbolt"
)

var DB *bbolt.DB

func InitDB() {
	var err error
	DB, err = bbolt.Open("udb.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	DB.Close()
}

func Put(key, value string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		return b.Put([]byte(key), []byte(value))
	})
}

func Get(key string) (string, error) {
	var value string
	err := DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte(key))
		if v != nil {
			value = string(v)
		}
		return nil
	})
	return value, err
}
