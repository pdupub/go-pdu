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

package core

import (
	"testing"
)

func TestNewContent(t *testing.T) {
	qc, err := NewContent(QCFmtStringTEXT, []byte("Hello World!"))
	if err != nil {
		t.Error(err)
	}
	t.Log(qc.Format)

	qc = CreateTextContent("Hello")
	if d, err := qc.GetData(); err != nil {
		t.Error(err)
	} else if d.(string) != "Hello" {
		t.Error(errContentParseFail)
	} else {
		t.Log(d.(string))
	}

	qc = CreateEmptyContent()
	if d, err := qc.GetData(); err != nil {
		t.Error(err)
	} else if d.(string) != "" {
		t.Error(errContentParseFail)
	} else {
		t.Log(d.(string))
	}

	qc = CreateIntContent(12345)
	if d, err := qc.GetData(); err != nil {
		t.Error(err)
	} else if d.(int64) != 12345 {
		t.Error(errContentParseFail)
	} else {
		t.Log(d.(int64))
	}

	qc = CreateFloatContent(123.45)
	if d, err := qc.GetData(); err != nil {
		t.Error(err)
	} else if d.(float64) != 123.45 {
		t.Error(errContentParseFail)
	} else {
		t.Log(d.(float64))
	}
}
