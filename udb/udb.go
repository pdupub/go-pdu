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
	"fmt"
	"log"
	"strings"

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
	CREATE TABLE IF NOT EXISTS Quantum (
		sig VARCHAR(132) PRIMARY KEY,
		qtype INTEGER,
		contents TEXT,
		nonce INTEGER,
		refs TEXT,
		address VARCHAR(42)
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
func (udb *UDB) PutQuantum(sig, contents, address, references string, nonce, qtype int) error {
	_, err := udb.db.Exec("INSERT INTO Quantum (sig, qtype, contents, nonce, refs, address) VALUES (?, ?, ?, ?, ?, ?)", sig, qtype, contents, nonce, references, address)
	return err
}

// GetQuantum retrieves a value by key from the Quantum table
func (udb *UDB) GetQuantum(sig string) (string, int, []string, string, int, error) {
	var nonce, qtype int
	var contents, references, address string
	err := udb.db.QueryRow("SELECT qtype, contents, nonce, refs, address FROM Quantum WHERE sig = ?", sig).Scan(&qtype, &contents, &nonce, &references, &address)
	return contents, nonce, strings.Split(references, ","), address, qtype, err
}

// GetReferencesBySig retrieves all references for a given sig from the Quantum table
func (udb *UDB) GetReferencesBySig(sig string) ([]string, error) {
	_, _, refs, _, _, err := udb.GetQuantum(sig)
	if err != nil {
		return nil, err
	}
	return refs, nil
}

// GetQuantumsByAddress retrieves all quantums for a given address from the Quantum table
func (udb *UDB) GetQuantumsByAddress(address string, limit, skip int, asc bool) ([]map[string]interface{}, error) {
	order := "DESC"
	if asc {
		order = "ASC"
	}
	query := fmt.Sprintf("SELECT sig, qtype, contents, nonce, refs FROM Quantum WHERE address = ? ORDER BY sig %s LIMIT ? OFFSET ?", order)

	rows, err := udb.db.Query(query, address, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quantums []map[string]interface{}
	for rows.Next() {
		var sig, contents, references string
		var nonce, qtype int
		err := rows.Scan(&sig, &qtype, &contents, &nonce, &references)
		if err != nil {
			return nil, err
		}
		quantum := map[string]interface{}{
			"sig":        sig,
			"qtype":      qtype,
			"nonce":      nonce,
			"contents":   contents,
			"references": strings.Split(references, ","),
		}
		quantums = append(quantums, quantum)
	}
	return quantums, nil
}
