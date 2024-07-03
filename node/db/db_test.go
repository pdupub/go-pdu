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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNodeDB(t *testing.T) {
	db_name := "test.db"
	db, err := NewNodeDB(db_name)
	assert.NoError(t, err)
	defer db.CloseDB()
	defer os.Remove(db_name)

	peer1 := Peer{
		ID:            "12D3KooWQ1",
		Address:       "/ip4/127.0.0.1/tcp/4001",
		Status:        "connected",
		LastConnected: time.Now(),
	}
	peer2 := Peer{
		ID:            "12D3KooWQ2",
		Address:       "/ip4/127.0.0.1/tcp/4002",
		Status:        "disconnected",
		LastConnected: time.Now().Add(-10 * time.Minute),
	}

	err = db.AddPeer(peer1)
	assert.NoError(t, err)
	err = db.AddPeer(peer2)
	assert.NoError(t, err)

	peers, err := db.GetPeers(0, 10)
	assert.NoError(t, err)
	assert.Len(t, peers, 2)

	assert.Equal(t, peer1.ID, peers[0].ID)
	assert.Equal(t, peer2.ID, peers[1].ID)

	peers, err = db.GetPeers(1, 10)
	assert.NoError(t, err)
	assert.Len(t, peers, 1)
	assert.Equal(t, peer2.ID, peers[0].ID)
}
