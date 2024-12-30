package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type Node struct {
	Host       host.Host
	DHT        *dht.IpfsDHT
	ctx        context.Context
	protocolID protocol.ID
}

// 创建新节点
func NewNode(ctx context.Context, protocolName string, protocolVersion string) (*Node, error) {
	// 创建libp2p主机
	h, err := libp2p.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	// 创建DHT用于节点发现
	kadDHT, err := dht.New(ctx, h)
	if err != nil {
		h.Close()
		return nil, fmt.Errorf("failed to create DHT: %w", err)
	}

	// 构造协议ID
	protocolID := protocol.ID(fmt.Sprintf("/%s/%s", protocolName, protocolVersion))

	node := &Node{
		Host:       h,
		DHT:        kadDHT,
		ctx:        ctx,
		protocolID: protocolID,
	}

	// 设置流处理器
	h.SetStreamHandler(protocolID, node.handleStream)

	// 启动本地节点发现
	if err := node.setupDiscovery(); err != nil {
		return nil, err
	}

	return node, nil
}

// 处理接收到的流
func (n *Node) handleStream(stream network.Stream) {
	// 读取消息
	buf := make([]byte, 1024)
	len, err := stream.Read(buf)
	if err != nil {
		fmt.Printf("Error reading from stream: %s\n", err)
		return
	}

	// 打印接收到的消息
	msg := string(buf[:len])
	fmt.Printf("Received message from %s: %s\n", stream.Conn().RemotePeer(), msg)

	// 如果收到 "Hello!"，回复 "How are you"
	if msg == "Hello!" {
		response := "How are you"
		_, err := stream.Write([]byte(response))
		if err != nil {
			fmt.Printf("Error sending response: %s\n", err)
		}
	}

	// 关闭流
	stream.Close()
}

// 发送消息到指定节点并等待回复
func (n *Node) SendMessage(peerID peer.ID, message string) error {
	// 打开到目标节点的流
	stream, err := n.Host.NewStream(n.ctx, peerID, n.protocolID)
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}
	defer stream.Close()

	// 发送消息
	_, err = stream.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// 等待回复
	buf := make([]byte, 1024)
	len, err := stream.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// 打印收到的回复
	response := string(buf[:len])
	fmt.Printf("Received response from %s: %s\n", peerID, response)

	return nil
}

// 设置节点发现
func (n *Node) setupDiscovery() error {
	// 启动mDNS发现服务
	discovery := mdns.NewMdnsService(n.Host, "pdu-network", &discoveryNotifee{node: n})
	return discovery.Start()
}

// mDNS发现回调
type discoveryNotifee struct {
	node *Node
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == n.node.Host.ID() {
		return // 忽略自己
	}

	// 尝试连接到发现的节点
	ctx, cancel := context.WithTimeout(n.node.ctx, 10*time.Second)
	defer cancel()

	if err := n.node.Host.Connect(ctx, pi); err != nil {
		fmt.Printf("Failed to connect to peer %s: %s\n", pi.ID, err)
		return
	}

	fmt.Printf("Connected to peer: %s\n", pi.ID)

	// 主动发起连接的一方发送 Hello 消息
	err := n.node.SendMessage(pi.ID, "Hello!")
	if err != nil {
		fmt.Printf("Failed to send hello message to %s: %s\n", pi.ID, err)
		return
	}
}

func (n *Node) Close() error {
	if err := n.DHT.Close(); err != nil {
		return err
	}
	return n.Host.Close()
}
