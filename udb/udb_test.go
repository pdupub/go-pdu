package udb

import (
	"os"
	"testing"
)

func TestUDB(t *testing.T) {
	const testDBName = "udb_test.db"

	// 删除测试数据库文件，以确保测试从干净的状态开始
	os.Remove(testDBName)

	// 初始化数据库并获取 UDB 实例
	db, err := InitDB(testDBName)
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.CloseDB()
	defer os.Remove(testDBName)

	// 测试 PutQuantum 和 GetQuantum
	sig := "test-sig"
	contents := "test-contents"
	err = db.PutQuantum(sig, contents, "test-address")
	if err != nil {
		t.Fatalf("PutQuantum failed: %v", err)
	}

	retContents, address, err := db.GetQuantum(sig)
	if err != nil {
		t.Fatalf("GetQuantum failed: %v", err)
	}
	if retContents != contents {
		t.Fatalf("GetQuantum: expected %s, got %s", contents, retContents)
	}

	// 测试 PutPublisher 和 GetPublisher
	value := "test-value"
	err = db.PutPublisher(address, value)
	if err != nil {
		t.Fatalf("PutPublisher failed: %v", err)
	}

	retValue, err := db.GetPublisher(address)
	if err != nil {
		t.Fatalf("GetPublisher failed: %v", err)
	}
	if retValue != value {
		t.Fatalf("GetPublisher: expected %s, got %s", value, retValue)
	}

	// 测试 PutReference 和 GetReferencesBySig
	ref := "test-ref"
	err = db.PutReference(sig, ref)
	if err != nil {
		t.Fatalf("PutReference failed: %v", err)
	}

	refs, err := db.GetReferencesBySig(sig)
	if err != nil {
		t.Fatalf("GetReferencesBySig failed: %v", err)
	}
	if len(refs) != 1 || refs[0] != ref {
		t.Fatalf("GetReferencesBySig: expected [%s], got %v", ref, refs)
	}

	// 测试 GetReference
	err = db.PutReference(sig, ref)
	if err != nil {
		t.Fatalf("PutReference failed: %v", err)
	}

	// 由于是自增ID，我们需要获取最后一个插入的ID
	var id int64
	err = db.db.QueryRow("SELECT MAX(id) FROM Reference").Scan(&id)
	if err != nil {
		t.Fatalf("Query MAX(id) failed: %v", err)
	}

	retSig, retRef, err := db.GetReference(id)
	if err != nil {
		t.Fatalf("GetReference failed: %v", err)
	}
	if retSig != sig || retRef != ref {
		t.Fatalf("GetReference: expected (%s, %s), got (%s, %s)", sig, ref, retSig, retRef)
	}
}
