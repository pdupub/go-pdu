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
	"testing"
)

func TestQContent(t *testing.T) {
	data := []byte("example data")
	format := "text/plain"
	zipped := true

	content, err := NewQContent(data, format, zipped)
	if err != nil {
		t.Errorf("error creating QContent: %v", err)
	}

	decodedData, err := content.GetData()
	if err != nil {
		t.Errorf("error getting data: %v", err)
	}

	if zipped {
		var buf bytes.Buffer
		writer := gzip.NewWriter(&buf)
		_, err := writer.Write(data)
		if err != nil {
			t.Errorf("error writing gzip: %v", err)
		}
		writer.Close()

		compressedData := buf.Bytes()

		if !bytes.Equal(content.Data, compressedData) {
			t.Errorf("expected %v, got %v", compressedData, content.Data)
		}
	} else if !bytes.Equal(decodedData, data) {
		t.Errorf("expected %s, got %s", string(data), string(decodedData))
	}

	if content.GetFormat() != format {
		t.Errorf("expected %s, got %s", format, content.GetFormat())
	}

	if content.IsZipped() != zipped {
		t.Errorf("expected %t, got %t", zipped, content.IsZipped())
	}

	jsonStr, err := content.ToJSON()
	if err != nil {
		t.Errorf("error converting to JSON: %v", err)
	}

	newContent := &QContent{}
	err = newContent.FromJSON(jsonStr)
	if err != nil {
		t.Errorf("error converting from JSON: %v", err)
	}

	newDecodedData, err := newContent.GetData()
	if err != nil {
		t.Errorf("error getting data: %v", err)
	}

	if !bytes.Equal(newDecodedData, data) {
		t.Errorf("expected %s, got %s", string(data), string(newDecodedData))
	}

	if newContent.GetFormat() != format {
		t.Errorf("expected %s, got %s", format, newContent.GetFormat())
	}

	if newContent.IsZipped() != zipped {
		t.Errorf("expected %t, got %t", zipped, newContent.IsZipped())
	}
}

func TestNewTXTQContent(t *testing.T) {
	data := []byte("example text data")

	content, err := NewTXTQContent(data)
	if err != nil {
		t.Errorf("error creating NewTXTQContent: %v", err)
	}

	if !bytes.Equal(content.Data, data) {
		t.Errorf("expected %s, got %s", string(data), string(content.Data))
	}

	if content.GetFormat() != "txt" {
		t.Errorf("expected txt, got %s", content.GetFormat())
	}

	if content.IsZipped() {
		t.Errorf("expected false, got %t", content.IsZipped())
	}
}
