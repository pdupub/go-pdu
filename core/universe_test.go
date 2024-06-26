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
	"os"
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestNewUniverse(t *testing.T) {
	const testDBName = "universe_test.db"

	// 删除测试数据库文件，以确保测试从干净的状态开始
	os.Remove(testDBName)

	// 初始化 Universe 并获取 UDB 实例
	universe, err := NewUniverse(testDBName)
	if err != nil {
		t.Fatalf("NewUniverse failed: %v", err)
	}
	defer universe.DB.CloseDB()
	defer os.Remove(testDBName)

	// 检查数据库是否初始化正确
	if universe.DB == nil {
		t.Fatalf("Expected non-nil DB, got nil")
	}
}

func TestRecv(t *testing.T) {
	const testDBName = "universe_test.db"

	// 删除测试数据库文件，以确保测试从干净的状态开始
	os.Remove(testDBName)

	// 初始化 Universe 并获取 UDB 实例
	universe, err := NewUniverse(testDBName)
	if err != nil {
		t.Fatalf("NewUniverse failed: %v", err)
	}
	defer universe.DB.CloseDB()
	defer os.Remove(testDBName)

	// // 创建测试 Quantum 对象
	// quantum := &Quantum{
	// 	Signature: []byte("test-sig"),
	// 	Contents:  QCS{&QContent{Data: []byte("test-content"), Format: "txt", Zipped: false}},

	// 	// Contents:   QCS{&QContent{Data: []byte("test-content"), Format: "txt", Zipped: false}},
	// 	References: []Sig{[]byte("ref1"), []byte("ref2")},
	// }

	cs := []*QContent{
		{Data: []byte("content1"), Format: "txt", Zipped: false},
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

	// 调用 Recv 方法
	err = universe.RecvQuantum(quantum)
	if err != nil {
		t.Fatalf("Recv failed: %v", err)
	}

	// 检查 Quantum 表中是否有对应的记录
	contents, _, _, _, err := universe.DB.GetQuantum(quantum.Signature.toHex())
	if err != nil {
		t.Fatalf("GetQuantum failed: %v", err)
	}
	expectedContent, _ := quantum.Contents.String()
	if contents != expectedContent {
		t.Fatalf("GetQuantum: expected %s, got %s", expectedContent, contents)
	}

	// 检查 Reference 表中是否有对应的记录
	refsRes, err := universe.DB.GetReferencesBySig(quantum.Signature.toHex())
	if err != nil {
		t.Fatalf("GetReferencesBySig failed: %v", err)
	}
	if len(refsRes) != len(quantum.References) {
		t.Fatalf("GetReferencesBySig: expected %d refs, got %d", len(quantum.References), len(refs))
	}
	for i, ref := range refsRes {
		if ref != quantum.References[i].toHex() {
			t.Fatalf("GetReferencesBySig: expected ref %s, got %s", quantum.References[i], ref)
		}
	}

	// 测试重复插入相同的 Quantum
	err = universe.RecvQuantum(quantum)
	if err == nil {
		t.Fatalf("Expected error for duplicate quantum, got nil")
	}
	if err.Error() != "quantum already exists" {
		t.Fatalf("Expected 'quantum already exists' error, got %v", err)
	}
}
