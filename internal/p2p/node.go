package p2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

type Node struct {
	Host host.Host
	DHT  *kaddht.IpfsDHT
}

func NewNode(ctx context.Context) (*Node, error) {
	h, err := libp2p.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	dht, err := kaddht.New(ctx, h)
	if err != nil {
		h.Close()
		return nil, fmt.Errorf("failed to create DHT: %w", err)
	}

	return &Node{
		Host: h,
		DHT:  dht,
	}, nil
}

func (n *Node) Close() error {
	if err := n.DHT.Close(); err != nil {
		return err
	}
	return n.Host.Close()
}
