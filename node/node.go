package node

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// JSON-RPC request structure
type JSONRPCRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params"`
	ID      json.RawMessage `json:"id,omitempty"`
}

// JSON-RPC response structure
type JSONRPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

// RPCError structure
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

func handleWebsite(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("node/static", "index.html"))
}

// startWebServer starts a simple web server
func startWebServer(port int) {
	http.HandleFunc("/", handleWebsite)
	fmt.Printf("Starting Website server on port %d... \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleCustomJSONRequest(w http.ResponseWriter, body []byte) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the JSON data received
	log.Printf("Received custom JSON data: %+v\n", jsonData)

	quantum, err := core.JsonToQuantum(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the quantum
	log.Printf("Quantum: %+v\n", quantum)
	log.Printf("Quantum signature: %s\n", core.Sig2Hex(quantum.Signature))
	log.Printf("Quantum content[0] data: %+s\n", quantum.Contents[0].Data)
	log.Printf("Quantum content[0] format: %+s\n", quantum.Contents[0].Format)
	log.Printf("Quantum ref[0]: %s\n", core.Sig2Hex(quantum.References[0]))

	addr, err := quantum.Ecrecover()
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// return
		log.Printf("Failed to recover address: %s\n", err.Error())
	}
	log.Printf("Quantum address: %s\n", addr.Hex())

	// Respond with a success message
	response := map[string]string{"status": "success"}
	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func handleRPCRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the incoming request
	log.Printf("Received request: %s\n", string(body))

	var req JSONRPCRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Method == "" {
		// If the method field is missing, treat it as a custom JSON request
		handleCustomJSONRequest(w, body)
		return
	}

	// Handle standard JSON-RPC requests
	var result interface{}
	var rpcErr *RPCError

	switch req.Method {
	case "eth_chainId":
		result = "0x2304" // Example chain ID, replace with actual chain ID
	case "net_version":
		result = "1" // Example network version, replace with actual network version
	case "eth_blockNumber":
		result = "0xBC614E" // Example block number, replace with actual block number
	case "eth_getBlockByNumber":
		result = "0xBC614E" // Example block data, replace with actual block data
	case "eth_gasPrice":
		result = "0x09184e72a000" // Example gas price, replace with actual gas price
	case "eth_getBalance":
		result = "0x8AC7230489E80000" // Example balance, replace with actual logic to fetch balance
	case "eth_getTransactionCount":
		result = "0x1" // Example transaction count, replace with actual logic to fetch transaction count
	// case "eth_call":
	// 	result, rpcErr = handleEthCall(req.Params)
	default:
		rpcErr = &RPCError{
			Code:    -32601,
			Message: "Method not found",
		}
	}

	resp := JSONRPCResponse{
		Jsonrpc: "2.0",
		Result:  result,
		Error:   rpcErr,
		ID:      req.ID,
	}

	responseBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the outgoing response
	log.Printf("Sending response: %s\n", string(responseBody))

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

// withCORS adds CORS headers to a handler
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

// startRPCServer starts a JSON-RPC server compatible with MetaMask
func startRPCServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", handleRPCRequest)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: withCORS(mux),
	}

	fmt.Printf("Starting RPC server on port %d...\n", port)
	log.Fatal(server.ListenAndServe())
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

	go startWebServer(webPort)
	go startRPCServer(rpcPort)

	fmt.Println("Enter the multiaddr of a peer to connect to (empty to skip):")
	peerAddr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	peerAddr = strings.TrimSpace(peerAddr)

	if peerAddr != "" {
		connectToPeer(h, peerAddr)
	}

	<-ctx.Done() // 保持程序运行
}
