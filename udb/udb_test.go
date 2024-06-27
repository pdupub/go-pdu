package udb

import (
	"os"
	"testing"
)

func TestUDB(t *testing.T) {
	// 设置数据库文件名
	dbName := "test.db"
	// 在测试结束时删除数据库文件
	defer os.Remove(dbName)

	// 初始化数据库
	db, err := InitDB(dbName)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// 定义测试数据
	sig1 := "testsig123"
	sig2 := "testsig456"
	qtype := 1
	contents := "test contents"
	references := "0xf09ec9d2fd43cfad1f0c93859e5678450c05d26a33c9298673ed991497e4e01c6a125379618b2fb2eb70fba5cc2ae1c946fa76bb0eeadc4723ded6176743ebab1c,0xd2ad5214bf586da7aaa829e460c07b7864d93f7edfbd6fc9e1fc1fc67df448de2bd52d1f1423b0c89b43b93cbe5e761abce559da564f84680e24deb1bf5e1ce61c"
	address := "0xC604E94a66bCE26FdffcE7dE309d0c9Af26fb33F"

	// 插入测试数据
	err = db.PutQuantum(sig1, contents, address, references, qtype)
	if err != nil {
		t.Fatalf("Failed to put quantum: %v", err)
	}

	err = db.PutQuantum(sig2, contents, address, references, qtype)
	if err != nil {
		t.Fatalf("Failed to put quantum: %v", err)
	}

	// 检索测试数据
	retContents, retReferences, retAddress, retQtype, err := db.GetQuantum(sig1)
	if err != nil {
		t.Fatalf("Failed to get quantum: %v", err)
	}

	// 验证插入的数据是否正确
	if retContents != contents {
		t.Errorf("Expected contents %v, got %v", contents, retContents)
	}
	if len(retReferences) != 2 {
		t.Errorf("Expected 2 references, got %d", len(retReferences))
	}
	if retAddress != address {
		t.Errorf("Expected address %v, got %v", address, retAddress)
	}
	if retQtype != qtype {
		t.Errorf("Expected qtype %v, got %v", qtype, retQtype)
	}

	// 根据地址检索量子数据
	quantums, err := db.GetQuantumsByAddress(address, 10, 0, true)
	if err != nil {
		t.Fatalf("Failed to get quantums by address: %v", err)
	}

	// 验证检索到的数据是否正确
	if len(quantums) != 2 {
		t.Errorf("Expected 2 quantums, got %d", len(quantums))
	}

	for _, quantum := range quantums {
		if quantum["contents"] != contents {
			t.Errorf("Expected contents %v, got %v", contents, quantum["contents"])
		}
		if quantum["qtype"] != qtype {
			t.Errorf("Expected qtype %v, got %v", qtype, quantum["qtype"])
		}
		if quantum["sig"] != sig1 && quantum["sig"] != sig2 {
			t.Errorf("Unexpected sig %v", quantum["sig"])
		}
		if len(quantum["references"].([]string)) != 2 {
			t.Errorf("Expected 2 references, got %d", len(quantum["references"].([]string)))
		}
	}
}
