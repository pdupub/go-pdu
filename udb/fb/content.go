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

package fb

import (
	"errors"
	"strconv"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
)

var (
	errContentFmtMissing = errors.New("content fmt is missing")
	errContentsMissing   = errors.New("contents is missing")
)

type FBContent struct {
	Data   interface{} `json:"data,omitempty"`
	Format int         `json:"fmt"`
}

func CS2Readable(contents interface{}) (interface{}, error) {

	switch contents := contents.(type) {
	case interface{}:
		readableCS := []interface{}{}
		for _, v := range contents.([]*core.QContent) {
			cc, err := Content2Readable(v)
			if err != nil {
				return nil, err
			}
			readableCS = append(readableCS, cc)
		}
		return readableCS, nil
	}
	return nil, errContentsMissing
}

func Content2Readable(content *core.QContent) (map[string]interface{}, error) {
	cc := make(map[string]interface{})
	switch content.Format {
	case core.QCFmtStringTEXT, core.QCFmtStringURL, core.QCFmtStringAddressHex, core.QCFmtStringSignatureHex:
		cc["data"] = string(content.Data)
	case core.QCFmtStringInt, core.QCFmtStringFloat:
		dataFloat, err := strconv.ParseFloat(string(content.Data), 64)
		if err != nil {
			return nil, err
		}
		cc["data"] = dataFloat
	case core.QCFmtBytesAddress:
		cc["data"] = identity.BytesToAddress(content.Data).Hex()
	case core.QCFmtBytesSignature:
		cc["data"] = core.Sig2Hex(content.Data)
	default:
		return nil, errContentFmtMissing
	}
	cc["fmt"] = content.Format

	return cc, nil
}
