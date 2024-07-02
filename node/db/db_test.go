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
