package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}

	// 创建表
	statements := []string{
		createQuantumTable,
		createContentTable,
		createReferenceTable,
		createQuantumReferenceTable,
	}
	for _, stmt := range statements {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}

	return db
}
