// node_test.go
package node

import (
	"os"
	"testing"
)

func TestNewNode(t *testing.T) {
	peerPort := 4001
	nodeKey := "node_test.key"
	dbName := "pdu_test.db"

	defer os.Remove(nodeKey)
	defer os.Remove(dbName)

	n, err := NewNode(peerPort, nodeKey, dbName)
	if err != nil {
		t.Fatalf("Failed to create node: %s", err)
	}

	if n.Host == nil {
		t.Fatalf("Node host is nil")
	}

	if n.Universe == nil {
		t.Fatalf("Node universe is nil")
	}
}

func TestRunNode(t *testing.T) {
	peerPort := 4001
	webPort := 8546
	rpcPort := 8545
	nodeKey := "node_test.key"
	dbName := "pdu_test.db"

	defer os.Remove(nodeKey)
	defer os.Remove(dbName)

	n, err := NewNode(peerPort, nodeKey, dbName)
	if err != nil {
		t.Fatalf("Failed to create node: %s", err)
	}

	go n.Run(webPort, rpcPort)

	// 在此添加更多检查和断言，以确保节点正确运行
	// 例如，你可以检查 Web 和 RPC 服务器是否可访问

	// 清理
	defer n.Host.Close()
	defer n.Universe.DB.CloseDB()

}
