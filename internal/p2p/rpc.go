package p2p

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
)

// 定义一个对外提供的 API
type PDUAPI struct {
	node *Node
}

func NewPDUAPI(node *Node) *PDUAPI {
	return &PDUAPI{
		node: node,
	}
}

func (p *PDUAPI) Chat(msg string) string {
	fmt.Println("Received message: ", msg)
	return fmt.Sprintf("You said: %s", msg)
}

func (p *PDUAPI) Message(peerID, msg string) string {
	if len(p.node.streams) == 0 {
		return "Connect to no peer"
	}

	pID, err := peer.Decode(peerID)
	if peerID == "" || err != nil {
		return "PeerID is missing"
	}

	if err := p.node.SendMessage(pID, msg); err != nil {
		return fmt.Sprintf("Send message err : %s", err)
	}

	return fmt.Sprintf("Send %s to %s", msg, peerID)
}
