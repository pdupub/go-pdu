package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

const (
	protocolName    = "example-protocol"
	protocolVersion = 1
	protocolID      = 0x01
)

func main() {
	// 启动 Node1
	go func() {
		enodeURL := startNode("Node1", 30301, "")
		time.Sleep(2 * time.Second) // 确保 Node1 启动完成
		// 启动 Node2 并连接到 Node1
		startNode("Node2", 30302, enodeURL)
	}()

	// 保持程序运行
	select {}
}

// startNode 创建一个 P2P 节点并启动
func startNode(name string, port int, connectToEnode string) string {
	// 生成私钥作为节点标识
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("[%s] Failed to generate private key: %v", name, err)
	}

	// 配置节点
	cfg := p2p.Config{
		PrivateKey: privateKey,
		MaxPeers:   10,
		Name:       name,
		Protocols: []p2p.Protocol{
			{
				Name:    protocolName,
				Version: protocolVersion,
				Length:  1,
				Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
					log.Printf("[%s] Connected to peer: %s", name, peer.ID())
					// 发送 "hello" 消息
					if err := sendMessage(rw, "hello"); err != nil {
						log.Printf("[%s] Failed to send message: %v", name, err)
						return err
					}
					// 监听收到的消息
					for {
						msg, err := rw.ReadMsg()
						if err != nil {
							log.Printf("[%s] Error reading message: %v", name, err)
							return err
						}
						// 读取消息内容
						payload := make([]byte, msg.Size)
						if _, err := msg.Payload.Read(payload); err != nil {
							log.Printf("[%s] Error reading payload: %v", name, err)
							return err
						}
						log.Printf("[%s] Received message: %s", name, string(payload))

						// 如果消息是 "hello"，回复 "hi"
						if string(payload) == "hello" {
							if err := sendMessage(rw, "hi"); err != nil {
								log.Printf("[%s] Failed to send reply: %v", name, err)
								return err
							}
						}
					}
				},
			},
		},
		ListenAddr: fmt.Sprintf("127.0.0.1:%d", port),
	}

	// 创建并启动服务器
	server := &p2p.Server{
		Config: cfg,
	}

	if err := server.Start(); err != nil {
		log.Fatalf("[%s] Failed to start server: %v", name, err)
	}
	defer server.Stop()

	// 打印节点的 Enode URL
	enodeURL := server.Self().URLv4()
	log.Printf("[%s] Node started on port %d", name, port)
	log.Printf("[%s] Enode URL: %s", name, enodeURL)

	// 如果提供了目标 Enode URL，则连接到目标节点
	if connectToEnode != "" {
		node := enode.MustParse(connectToEnode)
		server.AddPeer(node)
		log.Printf("[%s] Connected to peer: %s", name, connectToEnode)
	}

	// 返回当前节点的 Enode URL
	return enodeURL
}

// sendMessage 发送消息
func sendMessage(rw p2p.MsgReadWriter, message string) error {
	return p2p.Send(rw, protocolID, []byte(message))
}
