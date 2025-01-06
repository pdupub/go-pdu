package p2p

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/pdupub/go-pdu/internal/config"

	"github.com/ethereum/go-ethereum/rpc"
)

type Node struct {
	Host       host.Host
	DHT        *dht.IpfsDHT
	ctx        context.Context
	cancel     context.CancelFunc
	protocolID protocol.ID
	streams    map[peer.ID]network.Stream
	streamsMux sync.Mutex
	listener   net.Listener
}

// 创建新节点
func NewNode(ctx context.Context) (*Node, error) {
	ctx, cancel := context.WithCancel(ctx)

	// 创建libp2p主机
	h, err := libp2p.New()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	// 创建DHT用于节点发现
	kadDHT, err := dht.New(ctx, h)
	if err != nil {
		h.Close()
		cancel()
		return nil, fmt.Errorf("failed to create DHT: %w", err)
	}

	// 构造协议ID
	protocolID := protocol.ID(fmt.Sprintf("/%s/%s", config.ProtocolName, config.ProtocolVersion))

	node := &Node{
		Host:       h,
		DHT:        kadDHT,
		ctx:        ctx,
		cancel:     cancel,
		protocolID: protocolID,
		streams:    make(map[peer.ID]network.Stream),
	}

	// 设置流处理器
	h.SetStreamHandler(protocolID, node.handleStream)

	// 启动本地节点发现
	if err := node.setupDiscovery(); err != nil {
		return nil, err
	}

	// 创建RPC客户端
	rpcServer := rpc.NewServer()
	if err := rpcServer.RegisterName("pdu", NewPDUAPI(node)); err != nil {
		return nil, fmt.Errorf("failed to register PDU: %w", err)
	}

	http.Handle("/", rpcServer)

	addr := "127.0.0.1:8545"

	fmt.Println("HTTP RPC server listening on", addr)
	go http.ListenAndServe(addr, nil)

	return node, nil
}

func (n *Node) handleStream(stream network.Stream) {
	peerID := stream.Conn().RemotePeer()

	// 将这个新来的 stream 缓存到 map 中
	n.streamsMux.Lock()
	n.streams[peerID] = stream
	n.streamsMux.Unlock()

	// 在单独的 goroutine 中处理“读消息”逻辑
	go func() {
		defer func() {
			// 一旦退出读循环（出现错误或对端关闭等），需要清理
			n.streamsMux.Lock()
			delete(n.streams, peerID)
			n.streamsMux.Unlock()

			// 最后关闭这个 stream
			stream.Close()
		}()

		buf := make([]byte, 1024)
		for {
			// 不断从 stream 中读取数据
			length, err := stream.Read(buf)
			if err != nil {
				// 读出错，说明对端可能断开了或出现其他错误，结束循环
				fmt.Printf("Error reading from %s: %v\n", peerID, err)
				break
			}

			msg := string(buf[:length])
			fmt.Printf("Received message from %s: %s\n", peerID, msg)

			// 简单演示：如果收到 "Hello!"，则回复一句 "How are you"
			// if msg == "Hello!" {
			// 	if _, werr := stream.Write([]byte("How are you")); werr != nil {
			// 		fmt.Printf("Error sending response to %s: %v\n", peerID, werr)
			// 	}
			// }
		}
	}()
}

// 添加获取本地地址的方法
func (n *Node) GetLocalAddress() string {
	// 获取第一个本地地址
	for _, addr := range n.Host.Addrs() {
		// 优先返回本地地址
		if strings.Contains(addr.String(), "127.0.0.1") || strings.Contains(addr.String(), "localhost") {
			return fmt.Sprintf("%s/p2p/%s", addr, n.Host.ID())
		}
	}
	// 如果没有找到本地地址，返回第一个可用地址
	if len(n.Host.Addrs()) > 0 {
		return fmt.Sprintf("%s/p2p/%s", n.Host.Addrs()[0], n.Host.ID())
	}
	return ""
}

// 获取或创建 Stream
func (n *Node) getOrCreateStream(peerID peer.ID) (network.Stream, error) {
	n.streamsMux.Lock()
	defer n.streamsMux.Unlock()

	// // 检查是否已存在活跃的 stream
	// if stream, exists := n.streams[peerID]; exists {
	// 	// 修改验证方式：尝试写入一个空消息来测试流是否有效
	// 	if _, err := stream.Write([]byte{}); err == nil {
	// 		return stream, nil
	// 	}
	// 	// stream 已失效，删除它
	// 	delete(n.streams, peerID)
	// }

	// 创建新的 stream
	stream, err := n.Host.NewStream(n.ctx, peerID, n.protocolID)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	// 保存新创建的 stream
	n.streams[peerID] = stream
	return stream, nil
}

// 发送消息
func (n *Node) SendMessage(peerID peer.ID, message string) error {
	stream, err := n.getOrCreateStream(peerID)
	if err != nil {
		return err
	}

	// 发送消息
	_, err = stream.Write([]byte(message))
	if err != nil {
		// 如果发送失败，移除失效的 stream
		n.streamsMux.Lock()
		delete(n.streams, peerID)
		n.streamsMux.Unlock()
		return fmt.Errorf("failed to send message: %w", err)
	}

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
	err := n.node.SendMessage(pi.ID, "Hi")
	if err != nil {
		fmt.Printf("Failed to send hello message to %s: %s\n", pi.ID, err)
		return
	}
}

// 关闭时清理所有 streams
func (n *Node) Close() error {
	n.streamsMux.Lock()
	for _, stream := range n.streams {
		stream.Close()
	}
	n.streams = nil
	n.streamsMux.Unlock()

	if err := n.DHT.Close(); err != nil {
		return err
	}
	return n.Host.Close()
}
