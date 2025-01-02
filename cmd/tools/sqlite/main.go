package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=rwc")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 建表：存储一个 JSON 字符串
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user_data (
            id   INTEGER PRIMARY KEY AUTOINCREMENT,
            info TEXT
        )
    `)
	if err != nil {
		panic(err)
	}

	// 插入示例数据：JSON 格式
	_, err = db.Exec(`
        INSERT INTO user_data (info) VALUES
        ('{"name":"Alice","age":25}'),
        ('{"name":"Bob","age":30}')
    `)
	if err != nil {
		panic(err)
	}

	// 查询：从 JSON 中提取 name 并筛选 age > 25
	rows, err := db.Query(`
        SELECT json_extract(info, '$.name'), json_extract(info, '$.age')
        FROM user_data
        WHERE json_extract(info, '$.age') > 25
    `)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var age int
		err = rows.Scan(&name, &age)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Name: %s, Age: %d\n", name, age)
	}
}
