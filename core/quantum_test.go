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
	"io"
	"os"
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestNewQuantum(t *testing.T) {
	cs := []*QContent{
		{Data: []byte("content1"), Format: "txt", Zipped: false},
	}
	refs := []Sig{InitialQuantumReference}

	quantum, err := NewQuantum(QuantumTypeInformation, cs, refs...)
	if err != nil {
		t.Errorf("error creating Quantum: %v", err)
	}

	if len(quantum.Contents) != len(cs) {
		t.Errorf("expected %d contents, got %d", len(cs), len(quantum.Contents))
	}

	if len(quantum.References) != len(refs) {
		t.Errorf("expected %d references, got %d", len(refs), len(quantum.References))
	}

	if quantum.Type != QuantumTypeInformation {
		t.Errorf("expected type %d, got %d", QuantumTypeInformation, quantum.Type)
	}
}

func TestQuantumSignAndEcrecover(t *testing.T) {

	file, err := os.Open("testdata/logo.png")
	if err != nil {
		t.Errorf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		t.Errorf("Error getting file info: %v", err)
		return
	}

	// 创建一个字节切片来保存文件内容
	fileSize := fileInfo.Size()
	fileBytes := make([]byte, fileSize)

	// 将文件内容读入字节切片
	_, err = io.ReadFull(file, fileBytes)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
		return
	}

	cs := []*QContent{
		{Data: []byte("content1"), Format: "txt", Zipped: false},
		{Data: fileBytes, Format: "png", Zipped: true},
	}
	refs := []Sig{InitialQuantumReference}

	quantum, err := NewQuantum(QuantumTypeInformation, cs, refs...)
	if err != nil {
		t.Errorf("error creating Quantum: %v", err)
	}

	// Create a new DID for signing
	did, err := identity.New()
	if err != nil {
		t.Errorf("error creating new DID: %v", err)
	}

	// Unlock the DID
	err = did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)
	if err != nil {
		t.Errorf("error unlocking wallet: %v", err)
	}

	err = quantum.Sign(did)
	if err != nil {
		t.Errorf("error signing Quantum: %v", err)
	}

	address, err := quantum.Ecrecover()
	if err != nil {
		t.Errorf("error recovering address: %v", err)
	}

	// Add assertions based on expected address values
	expectedAddress := did.GetAddress()
	if address != expectedAddress {
		t.Errorf("expected address %s, got %s", expectedAddress, address)
	}

	// Marshal converts the UnsignedQuantum to a byte slice
	// q, _ := quantum.Marshal()
	// t.Logf("quantum : %s", q)
}
