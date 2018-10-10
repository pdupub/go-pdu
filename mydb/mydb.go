// Copyright 2018 The PDU Authors
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

package mydb

import (
	"github.com/boltdb/bolt"
	"strconv"
)

const (
	defaultDBName   = "memory"
	bucketTimeline  = "timeline"
	bucketGenealogy = "genealogy"
)

type MyDB struct {
	db *bolt.DB
	bs map[string]*bolt.Bucket
}

func Open() (*MyDB, error) {
	db, err := bolt.Open(defaultDBName, 0600, nil)
	if err != nil {
		db.Close()
		return nil, err
	}

	mydb := MyDB{db, make(map[string]*bolt.Bucket)}
	for _, tableName := range []string{bucketTimeline, bucketGenealogy} {
		err = mydb.createTable(tableName)
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	return &mydb, nil
}

func (m *MyDB) SaveTime(timestamp int64, proof string) error {
	return m.Put(bucketTimeline, strconv.Itoa(int(timestamp)), []byte(proof))
}

func (m *MyDB) GetTime(timestamp int64) (string, error) {
	res, err := m.Get(bucketTimeline, strconv.Itoa(int(timestamp)))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (m *MyDB) createTable(tableName string) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(tableName))
		if err != nil {
			return err
		} else {
			m.bs[tableName] = b
		}
		return nil
	})
}

func (m *MyDB) Close() error {
	return m.db.Close()
}

func (m *MyDB) Put(bucketName string, key string, value []byte) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), value)
		return err
	})

}

func (m *MyDB) Get(bucketName string, key string) (res []byte, err error) {
	err = m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		res = b.Get([]byte(key))
		return nil
	})
	return
}
