// Copyright 2024 The PDU Authors
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
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

type QContent struct {
	Data   []byte `json:"data,omitempty"` // Store as []byte
	Format string `json:"fmt"`
	Zipped bool   `json:"zipped,omitempty"`
}

// NewQContent creates a new QContent object
func NewQContent(data []byte, format string, zipped bool) (*QContent, error) {
	if zipped {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err := writer.Write(data)
		if err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
		data = buf.Bytes()
	}

	return &QContent{
		Data:   data,
		Format: format,
		Zipped: zipped,
	}, nil
}

// NewTXTQContent creates a new textQContent object
func NewTXTQContent(data []byte) (*QContent, error) {
	return NewQContent(data, "txt", false)
}

// GetData returns the binary data of the QContent object
func (q *QContent) GetData() ([]byte, error) {
	data := q.Data
	if q.Zipped {
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		data, err = io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// GetFormat returns the format of the QContent object
func (q *QContent) GetFormat() string {
	return q.Format
}

// IsZipped returns whether the QContent object is zipped
func (q *QContent) IsZipped() bool {
	return q.Zipped
}

// ToJSON converts the QContent object to a JSON string
func (q *QContent) ToJSON() (string, error) {
	data, err := json.Marshal(q)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON populates the QContent object from a JSON string
func (q *QContent) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), q)
}
