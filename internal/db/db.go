package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pdupub/go-pdu/internal/core"
)

type DB struct {
	db   *sql.DB
	path string
}

func NewDB(filename string) *DB {
	return &DB{db: initDB(filename),
		path: filename}
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) InsertQuantum(sq *core.SignedQuantum) error {
	return insertQuantum(db.db, sq)
}

func (db *DB) QueryQuantumsByReference(refText string) ([]core.SignedQuantum, error) {
	return queryQuantumsByReference(db.db, refText)
}

func initDB(filename string) *sql.DB {
	db, err := sql.Open("sqlite3", filename)
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
