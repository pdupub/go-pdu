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

package bolt

import (
	"time"

	"github.com/boltdb/bolt"
)

// UBoltDB is the db struct by bolt
type UBoltDB struct {
	db *bolt.DB
}

// NewDB initialize the new blot DB, create if no file in given path
func NewDB(path string) (*UBoltDB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	return &UBoltDB{db}, nil
}

// Close the blot DB
func (u *UBoltDB) Close() error {
	return u.db.Close()
}

// CreateBucket create new bucket by name
func (u *UBoltDB) CreateBucket(bucketName []byte) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucketName)
		return err
	})
}

// DeleteBucket delete the bucket by name
func (u *UBoltDB) DeleteBucket(bucketName []byte) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucketName)
	})
}

// Set key/val into bucket
func (u *UBoltDB) Set(bucketName, key, val []byte) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Put(key, val)
	})
}

// Get val by key from bucket
func (u *UBoltDB) Get(bucketName, key []byte) ([]byte, error) {
	var val []byte
	err := u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		val = b.Get(key)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return val, nil
}
