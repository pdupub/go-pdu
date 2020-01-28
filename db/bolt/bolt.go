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
	"bytes"
	"errors"
	"github.com/pdupub/go-pdu/db"
	"time"

	"github.com/boltdb/bolt"
)

var (
	errFindMissingLimit         = errors.New("find operate missing limit")
	errFindArgsNumberNotCorrect = errors.New("find operate number not correct")
	errBucketNotExist           = errors.New("bucket not exist")
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
func (u *UBoltDB) CreateBucket(bucketName string) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucketName))
		return err
	})
}

// DeleteBucket delete the bucket by name
func (u *UBoltDB) DeleteBucket(bucketName string) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}

// Set key/val into bucket
func (u *UBoltDB) Set(bucketName, key string, val []byte) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errBucketNotExist
		}
		return b.Put([]byte(key), val)
	})
}

// Get val by key from bucket
func (u *UBoltDB) Get(bucketName, key string) (val []byte, err error) {
	err = u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errBucketNotExist
		}
		val = b.Get([]byte(key))
		return nil
	})
	return val, err
}

// Find the rows from bucket by prefix
func (u *UBoltDB) Find(bucketName, prefix string, args ...int) (rows []*db.Row, err error) {
	var skip, limit int
	if len(args) == 0 {
		return rows, errFindMissingLimit
	} else if len(args) == 1 {
		skip = 0
		limit = args[0]
	} else if len(args) == 2 {
		skip = args[0]
		limit = args[1]
	} else {
		return rows, errFindArgsNumberNotCorrect
	}

	err = u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errBucketNotExist
		}
		c := b.Cursor()
		prefixBytes := []byte(prefix)
		count := 0
		for k, v := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, v = c.Next() {
			if count >= skip+limit {
				break
			}
			if count >= skip {
				rows = append(rows, &db.Row{K: string(k), V: v})
			}
			count++
		}
		return nil
	})
	return rows, err
}
