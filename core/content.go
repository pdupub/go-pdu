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
	"errors"
	"fmt"
	"strconv"
)

const (
	QCFmtStringTEXT         = 1
	QCFmtStringURL          = 2
	QCFmtStringJSON         = 3
	QCFmtStringInt          = 4
	QCFmtStringFloat        = 5
	QCFmtStringAddressHex   = 6
	QCFmtStringSignatureHex = 7

	QCFmtBytesSignature = 33
	QCFmtBytesAddress   = 34

	QCFmtImagePNG = 65
	QCFmtImageJPG = 66
	QCFmtImageBMP = 67

	QCFmtAudioWAV = 97
	QCFmtAudioMP3 = 98

	QCFmtVideoMP4 = 129
)

var (
	errContentParseFail = errors.New("contents parse fail")
	errContentFmtNotFit = errors.New("content format not fit")
)

// QContent is one piece of data in Quantum,
// all variables should be in alphabetical order.
type QContent struct {
	Data   []byte `json:"data,omitempty"`
	Format int    `json:"fmt"`
}

func NewContent(fmt int, data []byte) (*QContent, error) {
	return &QContent{Format: fmt, Data: data}, nil
}

func CreateTextContent(t string) *QContent {
	c, _ := NewContent(QCFmtStringTEXT, []byte(t))
	return c
}

func CreateEmptyContent() *QContent {
	c, _ := NewContent(QCFmtStringTEXT, []byte(""))
	return c
}

func CreateIntContent(num int64) *QContent {
	c, _ := NewContent(QCFmtStringInt, []byte(fmt.Sprintf("%d", num)))
	return c
}

func CreateFloatContent(num float64) *QContent {
	c, _ := NewContent(QCFmtStringFloat, []byte(fmt.Sprintf("%g", num)))
	return c
}

func (c *QContent) GetData() (interface{}, error) {
	switch c.Format {
	case QCFmtStringTEXT, QCFmtStringURL, QCFmtStringJSON, QCFmtStringAddressHex, QCFmtBytesSignature:
		return string(c.Data), nil
	case QCFmtStringInt:
		return strconv.ParseInt(string(c.Data), 10, 64)
	case QCFmtStringFloat:
		return strconv.ParseFloat(string(c.Data), 64)
	default:
		return c.Data, nil
	}
}
