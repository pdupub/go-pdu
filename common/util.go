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

package common

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
)

// ToHash convert data into []byte by gob
func ToHash(data interface{}) []byte {
	var dao bytes.Buffer
	encoder := gob.NewEncoder(&dao)
	encoder.Encode(data)
	return dao.Bytes()
}

// FromHash convert []byte into struct data
func FromHash(hash []byte, data interface{}) error {
	var dao bytes.Buffer
	decoder := gob.NewDecoder(&dao)
	dao.Write(hash)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

// ToMD5 get hash
func ToMD5(src []byte) []byte {
	ctx := md5.New()
	ctx.Write(src)
	dst := make([]byte, HashLength)
	hex.Encode(dst, ctx.Sum(nil))
	return dst
}
