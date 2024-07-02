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

package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type NodeDB struct {
	db *sql.DB
}

// Peer represents a peer in the network.
type Peer struct {
	ID            string
	Address       string
	Status        string
	LastConnected time.Time
}

// NewNodeDB creates a new NodeDB instance and initializes the peer table.
func NewNodeDB(dbName string) (*NodeDB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	// Create peer table if it doesn't exist.
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS peer (
		id TEXT PRIMARY KEY,
		address TEXT,
		status TEXT,
		last_connected DATETIME
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &NodeDB{db: db}, nil
}

// AddPeer adds a new peer to the database.
func (db *NodeDB) AddPeer(peer Peer) error {
	insertQuery := `INSERT OR REPLACE INTO peer (id, address, status, last_connected) VALUES (?, ?, ?, ?)`
	_, err := db.db.Exec(insertQuery, peer.ID, peer.Address, peer.Status, peer.LastConnected)
	return err
}

// GetPeers retrieves peers from the database with pagination.
func (db *NodeDB) GetPeers(skip, limit int) ([]Peer, error) {
	query := `SELECT id, address, status, last_connected FROM peer ORDER BY last_connected DESC LIMIT ? OFFSET ?`
	rows, err := db.db.Query(query, limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []Peer
	for rows.Next() {
		var peer Peer
		var lastConnectedStr string
		err := rows.Scan(&peer.ID, &peer.Address, &peer.Status, &lastConnectedStr)
		if err != nil {
			return nil, err
		}
		peer.LastConnected, err = time.Parse(time.RFC3339, lastConnectedStr)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}
	return peers, nil
}

// CloseDB closes the SQLite database.
func (db *NodeDB) CloseDB() {
	db.db.Close()
}
