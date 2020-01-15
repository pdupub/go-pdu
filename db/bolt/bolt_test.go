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
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"testing"
)

func TestNewDB(t *testing.T) {
	bucketName := "testBucket"
	keyPrefix := "key"
	valPrefix := []byte("val")

	dir, _ := os.Getwd()
	filePath := path.Join(dir, "my_test.db")
	// clear file if need
	os.Remove(filePath)

	u, err := NewDB(filePath)
	if err != nil {
		t.Error(err)
	}

	if err := u.CreateBucket(bucketName); err != nil {
		t.Error(err)
	}
	for i := int64(0); i < 10; i++ {
		if err := u.Set(bucketName,
			fmt.Sprintf("%s%s", keyPrefix, big.NewInt(i).Bytes()),
			append(valPrefix, big.NewInt(i).Bytes()...)); err != nil {
			t.Error(err)
		}
	}

	val, err := u.Get(bucketName, fmt.Sprintf("%s%s", keyPrefix, big.NewInt(5).Bytes()))
	if err != nil {
		t.Error(err)
	}
	if string(val) != fmt.Sprintf("%s%s", valPrefix, big.NewInt(5).Bytes()) {
		t.Error("val not equal")
	}

	rows, err := u.Find(bucketName, keyPrefix, 3)
	if err != nil {
		t.Error(err)
	}

	if len(rows) != 3 {
		t.Error(errors.New("result number not match"))
	}

	if err := u.DeleteBucket(bucketName); err != nil {
		t.Error(err)
	}

	if err := u.Close(); err != nil {
		t.Error(err)
	}
	// clear test file
	os.Remove(filePath)
}
