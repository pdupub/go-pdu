package node

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = "/p2p/1.0.0"

// createNode creates a new libp2p host
func createNode(listenPort int) (host.Host, context.Context) {
	ctx := context.Background()
	h, err := libp2p.New(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
	)
	if err != nil {
		log.Fatal(err)
	}
	return h, ctx
}

// handleInterrupt handles OS interrupts to gracefully shut down the host
func handleInterrupt(_ context.Context, h host.Host) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down...")
		if err := h.Close(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
}

// handleStream is the function to handle incoming streams
func handleStream(s network.Stream) {
	fmt.Println("Got a new stream!")
	// 处理流的代码
	s.Close()
}

// connectToPeer connects to another peer using its multiaddress
func connectToPeer(h host.Host, peerAddr string) {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		log.Fatalf("Invalid multiaddress: %s", err)
	}

	peerinfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalf("Failed to get peer info: %s", err)
	}

	if err := h.Connect(context.Background(), *peerinfo); err != nil {
		log.Fatalf("Failed to connect to peer: %s", err)
	}

	fmt.Printf("Connected to %s\n", peerinfo.ID.String())
}

// startWebServer starts a simple web server
func startWebServer(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is the web server!")
	})

	addr := fmt.Sprintf(":%d", port)
	go func() {
		fmt.Printf("Starting web server on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start web server: %s", err)
		}
	}()
}

// startRPCServer starts a JSON-RPC server compatible with MetaMask
func startRPCServer(port int) {
	server := rpc.NewServer()
	// 注册需要的RPC服务，这里可以自定义自己的RPC方法
	// server.RegisterName("YourServiceName", &YourService{})

	addr := fmt.Sprintf(":%d", port)
	go func() {
		fmt.Printf("Starting RPC server on %s\n", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("Failed to start RPC server: %s", err)
		}
		if err := server.ServeListener(listener); err != nil {
			log.Fatalf("Failed to serve RPC requests: %s", err)
		}
	}()
}

// Run starts the libp2p node and listens for incoming connections
func Run(listenPort, webPort, rpcPort int) {
	h, ctx := createNode(listenPort)
	handleInterrupt(ctx, h)

	h.SetStreamHandler(protocol.ID(protocolID), handleStream)

	fmt.Printf("Node ID: %s\n", h.ID().String())
	for _, addr := range h.Addrs() {
		fmt.Printf("Node Address: %s\n", addr.String())
	}

	startWebServer(webPort)
	startRPCServer(rpcPort)

	fmt.Println("Enter the multiaddr of a peer to connect to (empty to skip):")
	peerAddr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	peerAddr = strings.TrimSpace(peerAddr)

	if peerAddr != "" {
		connectToPeer(h, peerAddr)
	}

	<-ctx.Done() // 保持程序运行
}
