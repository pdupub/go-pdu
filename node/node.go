// Copyright 2024 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package node

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ed25519"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/multiformats/go-multiaddr"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/node/db"
)

// staticFiles holds our static web server content.
//
//go:embed static/*
var staticFiles embed.FS

const protocolID = "/p2p/1.0.0"

type Node struct {
	Host     host.Host
	Universe *core.Universe
	Ctx      context.Context
	ndb      *db.NodeDB
}

func NewNode(listenPort int, nodeKey, dbName string) (*Node, error) {
	ctx := context.Background()

	privKey, err := loadOrCreatePrivateKey(nodeKey)
	if err != nil {
		return nil, err
	}

	h, err := libp2p.New(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.Identity(privKey),
	)
	if err != nil {
		return nil, err
	}

	universe, err := core.NewUniverse(dbName)
	if err != nil {
		return nil, err
	}

	nodeDB, err := db.NewNodeDB("node_" + dbName)
	if err != nil {
		return nil, err
	}

	return &Node{
		Host:     h,
		Universe: universe,
		Ctx:      ctx,
		ndb:      nodeDB,
	}, nil
}

func loadOrCreatePrivateKey(filename string) (crypto.PrivKey, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File does not exist, create a new private key and save it
		_, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(filename, []byte(hex.EncodeToString(priv)), 0600)
		if err != nil {
			return nil, err
		}
		fmt.Println("New private key generated and saved to", filename)
		return crypto.UnmarshalEd25519PrivateKey(priv)
	}

	// File exists, load the private key
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	priv, err := hex.DecodeString(string(keyBytes))
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalEd25519PrivateKey(priv)
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
		n.ndb.CloseDB()
		os.Exit(0)
	}()
}

func (n *Node) handleStream(s network.Stream) {
	fmt.Println("Got a new stream!")
	defer s.Close()

	remotePeerID := s.Conn().RemotePeer()
	remotePeerAddr := s.Conn().RemoteMultiaddr()

	remoteFullAddr := fmt.Sprintf("%s/p2p/%s", remotePeerAddr, remotePeerID)
	fmt.Printf("Connected to peer: %s\n", remoteFullAddr)

	buf := new(bytes.Buffer)
	buf.ReadFrom(s)
	message := buf.String()
	fmt.Printf("Received message: %s\n", message)

	// save peer info into ndb

	info := strings.Split(remoteFullAddr, "/")
	peer := db.Peer{
		ID:            info[len(info)-1],
		Address:       remoteFullAddr,
		Status:        "connected",
		LastConnected: time.Now(),
	}
	if err := n.ndb.AddPeer(peer); err != nil {
		log.Printf("Failed to add peer to database: %s", err)
	}

}

func (n *Node) sendMessage(peerID peer.ID, message string) error {
	s, err := n.Host.NewStream(n.Ctx, peerID, protocol.ID(protocolID))
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = s.Write([]byte(message))
	return err
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

	peer := db.Peer{
		ID:            peerinfo.ID.String(),
		Address:       peerAddr,
		Status:        "connected",
		LastConnected: time.Now(),
	}
	if err := n.ndb.AddPeer(peer); err != nil {
		log.Printf("Failed to add peer to database: %s", err)
	}

	fmt.Printf("Connected to %s\n", peerinfo.ID.String())
}

func (n *Node) startWebServer(port int) {
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(subFS))))

	// http.HandleFunc("/", n.handleWebsite)
	fmt.Printf("Starting Website server on port %d... \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
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

func (n *Node) connectPeers() {
	fmt.Println("Enter the multiaddr of a peer to connect to (empty to skip):")
	peerAddr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	peerAddr = strings.TrimSpace(peerAddr)

	if peerAddr != "" {
		n.connectToPeer(peerAddr)

		fmt.Println("Enter message to send (empty to skip):")
		message, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		message = strings.TrimSpace(message)

		if message != "" {
			peerinfo, err := peer.AddrInfoFromP2pAddr(multiaddr.StringCast(peerAddr))
			if err == nil {
				err = n.sendMessage(peerinfo.ID, message)
				if err != nil {
					fmt.Printf("Failed to send message: %s\n", err)
				} else {
					fmt.Println("Message sent successfully")
				}
			} else {
				fmt.Printf("Failed to get peer info: %s\n", err)
			}
		}
	}
}

func (n *Node) Run(webPort, rpcPort int) {
	n.handleInterrupt()

	n.Host.SetStreamHandler(protocol.ID(protocolID), n.handleStream)

	fmt.Printf("Node ID: %s\n", n.Host.ID().String())
	for _, addr := range n.Host.Addrs() {
		fmt.Printf("Node Address: %s\n", addr.String())
		fmt.Printf("Node Fully Address: %s/p2p/%s\n", addr.String(), n.Host.ID().String())

	}

	go n.startWebServer(webPort)
	go n.startRPCServer(rpcPort)
	go n.connectPeers()

	<-n.Ctx.Done() // 保持程序运行
}
