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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//HashFile return sha256 value of file
func HashFile(filePath string) (content []byte, hashStr string, err error) {
	var reader io.Reader
	if strings.HasPrefix(filePath, "http") {
		resp, err := http.Get(filePath)
		if err != nil {
			return content, hashStr, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return content, hashStr, http.ErrMissingFile
		}

		content, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return content, hashStr, err
		}
		reader = bytes.NewReader(content)
	} else {
		file, err := os.Open(filePath)
		if err != nil {
			return content, hashStr, err
		}
		defer file.Close()

		content, err = ioutil.ReadAll(file)
		if err != nil {
			return content, hashStr, err
		}
		reader = bytes.NewReader(content)
	}
	hash := sha256.New()
	if _, err = io.Copy(hash, reader); err != nil {
		return
	}
	hashStr = hex.EncodeToString(hash.Sum(nil))
	return
}
