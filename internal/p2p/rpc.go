package p2p

import (
	"fmt"
	"strings"

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

func (p *PDUAPI) List() []string {
	_, files, err := p.node.ListKeystoreFiles()
	if err != nil {
		return nil
	}
	// 创建一个新的字符串切片用于存储结果
	result := make([]string, len(files))

	// 遍历原字符串切片，逐个去除 .json 后缀
	for i, str := range files {
		result[i] = strings.TrimSuffix(str, ".json")
	}

	return result
}

func (p *PDUAPI) Unlock(addr, password string) string {
	if err := p.node.UnlockPrivKey(addr, password); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Success unlock %s", addr)
}

func (p *PDUAPI) lock(addr string) string {
	p.node.ClearPrivKey(addr)
	return fmt.Sprintf("Success lock %s", addr)
}
