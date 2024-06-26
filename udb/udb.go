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

package udb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type UDB struct {
	db *sql.DB
}

// InitDB initializes the SQLite database with the given name and returns a UDB instance
func InitDB(dbName string) (*UDB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	createTables := `
	CREATE TABLE IF NOT EXISTS Setting (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);

	CREATE TABLE IF NOT EXISTS Quantum (
		sig TEXT PRIMARY KEY,
		contents TEXT,
		address VARCHAR(42)
	);

	CREATE TABLE IF NOT EXISTS Publisher (
		address TEXT PRIMARY KEY,
		value TEXT
	);

	CREATE TABLE IF NOT EXISTS Reference (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sig TEXT,
		ref TEXT
	);
	`
	_, err = db.Exec(createTables)
	if err != nil {
		return nil, err
	}
	return &UDB{db: db}, nil
}

// CloseDB closes the SQLite database
func (udb *UDB) CloseDB() {
	udb.db.Close()
}

// PutQuantum stores a key-value pair in the Quantum table
func (udb *UDB) PutQuantum(sig, contents, address string) error {
	_, err := udb.db.Exec("INSERT OR REPLACE INTO Quantum (sig, contents, address) VALUES (?, ?, ?)", sig, contents, address)
	return err
}

// GetQuantum retrieves a value by key from the Quantum table
func (udb *UDB) GetQuantum(sig string) (string, string, error) {
	var contents, address string
	err := udb.db.QueryRow("SELECT contents, address FROM Quantum WHERE sig = ?", sig).Scan(&contents, &address)
	return contents, address, err
}

// PutPublisher stores a key-value pair in the Publisher table
func (udb *UDB) PutPublisher(address, value string) error {
	_, err := udb.db.Exec("INSERT OR REPLACE INTO Publisher (address, value) VALUES (?, ?)", address, value)
	return err
}

// GetPublisher retrieves a value by key from the Publisher table
func (udb *UDB) GetPublisher(address string) (string, error) {
	var value string
	err := udb.db.QueryRow("SELECT value FROM Publisher WHERE address = ?", address).Scan(&value)
	return value, err
}

// PutReference stores a reference in the Reference table
func (udb *UDB) PutReference(sig, ref string) error {
	_, err := udb.db.Exec("INSERT INTO Reference (sig, ref) VALUES (?, ?)", sig, ref)
	return err
}

// GetReference retrieves a reference by ID from the Reference table
func (udb *UDB) GetReference(id int64) (string, string, error) {
	var sig, ref string
	err := udb.db.QueryRow("SELECT sig, ref FROM Reference WHERE id = ?", id).Scan(&sig, &ref)
	return sig, ref, err
}

// GetReferencesBySig retrieves all refs for a given sig from the Reference table
func (udb *UDB) GetReferencesBySig(sig string) ([]string, error) {
	rows, err := udb.db.Query("SELECT ref FROM Reference WHERE sig = ?", sig)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []string
	for rows.Next() {
		var ref string
		if err := rows.Scan(&ref); err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	return refs, nil
}
