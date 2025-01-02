package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// ----------------------------
// 你的结构定义
// ----------------------------
type QContent struct {
	Data   []byte `json:"data,omitempty"` // Store as []byte
	Format string `json:"fmt"`
	Zipped bool   `json:"zipped,omitempty"`
}

type UnsignedQuantum struct {
	Contents []*QContent `json:"cs,omitempty"`
	Nonce    int         `json:"nonce"`
	Type     int         `json:"type,omitempty"`
}

// ----------------------------
// 常量或函数
// ----------------------------
const createTableSQL = `
CREATE TABLE IF NOT EXISTS quantum (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quantum_data TEXT NOT NULL
);
`

const insertSQL = `
INSERT INTO quantum (quantum_data) VALUES (?);
`

const selectSQL = `
SELECT id, quantum_data FROM quantum;
`

func main() {
	// 1) 打开（或创建）SQLite 数据库文件
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("open db error: %v", err)
	}
	defer db.Close()

	// 2) 建表（若不存在则创建）
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("create table error: %v", err)
	}

	// 3) 构造示例数据
	uq := UnsignedQuantum{
		Contents: []*QContent{
			{
				Data:   []byte("Hello Bytes"),
				Format: "raw",
				Zipped: false,
			},
			{
				Data:   []byte{0x01, 0x02, 0x03},
				Format: "binary",
				Zipped: true,
			},
		},
		Nonce: 42,
		Type:  1,
	}

	// 4) 序列化 UnsignedQuantum 为 JSON
	quantumJSON, err := json.Marshal(uq)
	if err != nil {
		log.Fatalf("json marshal error: %v", err)
	}

	// 5) 插入到数据库
	result, err := db.Exec(insertSQL, string(quantumJSON))
	if err != nil {
		log.Fatalf("insert error: %v", err)
	}
	newID, _ := result.LastInsertId()
	fmt.Printf("Inserted record with ID: %d\n", newID)

	// 6) 查询并读取数据
	rows, err := db.Query(selectSQL)
	if err != nil {
		log.Fatalf("select error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int64
			quantumData string
		)
		err := rows.Scan(&id, &quantumData)
		if err != nil {
			log.Fatalf("scan error: %v", err)
		}

		// 解析回 UnsignedQuantum
		var uqParsed UnsignedQuantum
		err = json.Unmarshal([]byte(quantumData), &uqParsed)
		if err != nil {
			log.Fatalf("json unmarshal error: %v", err)
		}

		// 7) 打印结果
		fmt.Printf("\n--- Record ID: %d ---\n", id)
		fmt.Printf("Nonce: %d\nType: %d\n", uqParsed.Nonce, uqParsed.Type)
		for i, c := range uqParsed.Contents {
			fmt.Printf(" Content[%d]: Data=%v, Format=%s, Zipped=%v\n",
				i, c.Data, c.Format, c.Zipped)
		}
	}
	// 如果 rows.Next() 里返回 false，需要检查 rows.Err()
	if err := rows.Err(); err != nil {
		log.Fatalf("rows error: %v", err)
	}
}
