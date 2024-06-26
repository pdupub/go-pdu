package node

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/pdupub/go-pdu/core"
)

const protocolID = "/p2p/1.0.0"

type Node struct {
	Host     host.Host
	Universe *core.Universe
	Ctx      context.Context
}

func NewNode(listenPort int, dbName string) (*Node, error) {
	ctx := context.Background()
	h, err := libp2p.New(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
	)
	if err != nil {
		return nil, err
	}

	universe, err := core.NewUniverse(dbName)
	if err != nil {
		return nil, err
	}

	return &Node{
		Host:     h,
		Universe: universe,
		Ctx:      ctx,
	}, nil
}

func (n *Node) handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down...")
		if err := n.Host.Close(); err != nil {
			log.Fatal(err)
		}
		n.Universe.DB.CloseDB()
		os.Exit(0)
	}()
}

func (n *Node) handleStream(s network.Stream) {
	fmt.Println("Got a new stream!")
	s.Close()
}

func (n *Node) connectToPeer(peerAddr string) {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		log.Fatalf("Invalid multiaddress: %s", err)
	}

	peerinfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalf("Failed to get peer info: %s", err)
	}

	if err := n.Host.Connect(n.Ctx, *peerinfo); err != nil {
		log.Fatalf("Failed to connect to peer: %s", err)
	}

	fmt.Printf("Connected to %s\n", peerinfo.ID.String())
}

func (n *Node) handleWebsite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("node/static", "index.html"))
}

func (n *Node) startWebServer(port int) {
	http.HandleFunc("/", n.handleWebsite)
	fmt.Printf("Starting Website server on port %d... \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (n *Node) startRPCServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", n.handleRPCRequest)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: withCORS(mux),
	}

	fmt.Printf("Starting RPC server on port %d...\n", port)
	log.Fatal(server.ListenAndServe())
}

func (n *Node) Run(webPort, rpcPort int) {
	n.handleInterrupt()

	n.Host.SetStreamHandler(protocol.ID(protocolID), n.handleStream)

	fmt.Printf("Node ID: %s\n", n.Host.ID().String())
	for _, addr := range n.Host.Addrs() {
		fmt.Printf("Node Address: %s\n", addr.String())
	}

	go n.startWebServer(webPort)
	go n.startRPCServer(rpcPort)

	fmt.Println("Enter the multiaddr of a peer to connect to (empty to skip):")
	peerAddr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	peerAddr = strings.TrimSpace(peerAddr)

	if peerAddr != "" {
		n.connectToPeer(peerAddr)
	}

	<-n.Ctx.Done() // 保持程序运行
}
