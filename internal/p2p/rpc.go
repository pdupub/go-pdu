package p2p

import (
	"fmt"
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
