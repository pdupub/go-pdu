package db

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pdupub/go-pdu/internal/core"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitDB(t *testing.T) {
	db := NewDB("test.db")
	defer db.Close()

}

func TestInsertQueryQuantum(t *testing.T) {
	db := NewDB("test.db")
	defer db.Close()

	quantum := core.NewUnsignedQuantum([]*core.QContent{
		{
			Data:   "hello world",
			Format: "string",
		},
		{
			Data:   []byte{0x01, 0x02, 0x03},
			Format: "binary",
		},
		{
			Data:   123,
			Format: "number",
		},
	}, core.DefaultLastSig, 1, []string{
		"ref1",
		"ref2",
	})

	// 生成一个私钥（secp256k1）
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("GenerateKey error: %v", err)
	}

	jsonBytes, err := core.GenerateSignedJSON(privateKey, *quantum)
	if err != nil {
		t.Errorf("GenerateSignedJSON error: %v", err)
	}

	var signed core.SignedQuantum

	if err := json.Unmarshal(jsonBytes, &signed); err != nil {
		t.Errorf("json.Unmarshal error: %v", err)
	}

	if err := db.InsertQuantum(&signed); err != nil {
		t.Errorf("insertQuantum error: %v", err)
	}

	if results, err := db.QueryQuantumsByReference("ref1"); err != nil {
		t.Errorf("queryQuantumsByReference error: %v", err)
	} else {
		t.Logf("queryQuantumsByReference results: %v", results)
	}
}
