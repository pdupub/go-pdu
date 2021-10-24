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

package p2p

import (
	"encoding/base64"
	"errors"
	"testing"
)

func TestUtils(t *testing.T) {
	path1 := "../res/appicons/dark/appstore.png"
	content1, hash1, err := HashFile(path1)
	if err != nil {
		t.Error(err)
	}
	t.Log("hash", hash1)
	path2 := "http://127.0.0.1:1323/files/ab3c1f6bae3afcb9beece409242b2f8b220518325a2eab3ebfb1b58bb15969a1-appstore.png"
	content2, hash2, err := HashFile(path2)
	if err != nil {
		t.Error(err)
	}
	t.Log("hash", hash2)
	if hash1 != hash2 {
		t.Error(errors.New("hash not match"))
	}

	if base64.StdEncoding.EncodeToString(content1) != base64.StdEncoding.EncodeToString(content2) {
		t.Error(errors.New("content not match"))
	}

	base1 := base64.StdEncoding.EncodeToString(content1)
	base2 := base64.StdEncoding.EncodeToString(content2)
	t.Log(base1[:40], "...", base1[len(base1)-30:])
	t.Log(base2[:40], "...", base2[len(base2)-30:])
}
