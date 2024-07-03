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
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/pdupub/go-pdu/node/db"
)

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

func (n *Node) connectToPeer(peerAddr string) (string, error) {
	maddr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return "", fmt.Errorf("invalid multiaddress: %s", err)
	}

	peerinfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return "", fmt.Errorf("failed to get peer info: %s", err)
	}

	if err := n.Host.Connect(n.Ctx, *peerinfo); err != nil {
		return "", fmt.Errorf("failed to connect to peer: %s", err)
	}

	peer := db.Peer{
		ID:            peerinfo.ID.String(),
		Address:       peerAddr,
		Status:        "connected",
		LastConnected: time.Now(),
	}
	if err := n.ndb.AddPeer(peer); err != nil {
		return "", fmt.Errorf("failed to add peer to database: %s", err)
	}

	return peerinfo.ID.String(), nil
}

func (n *Node) chatWithPeer(peerAddr string) {

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

func (n *Node) connectPeers() {

	peers, err := n.ndb.GetPeers(0, 10)
	if err == nil && len(peers) > 0 {

		for _, peer := range peers {
			// if peer.Status == "connected" {
			// 	continue
			// }
			if targetID, err := n.connectToPeer(peer.Address); err != nil {
				fmt.Printf("Failed to connect to peer: %s\n", err)
			} else {
				fmt.Printf("Connected to peer: %s\n", targetID)
				n.peerAddrs = append(n.peerAddrs, peer.Address)
				n.chatWithPeer(peer.Address)
			}
		}
	}

	if len(n.peerAddrs) == 0 {
		fmt.Printf("Failed to get peers from database: %s\n", err)
		fmt.Println("Enter the multiaddr of a peer to connect to (empty to skip):")
		peerAddr, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		peerAddr = strings.TrimSpace(peerAddr)

		if peerAddr != "" {
			if targetID, err := n.connectToPeer(peerAddr); err != nil {
				fmt.Printf("Failed to connect to peer: %s\n", err)
			} else {
				fmt.Printf("Connected to peer: %s\n", targetID)
				n.peerAddrs = append(n.peerAddrs, peerAddr)
				n.chatWithPeer(peerAddr)

			}
		}
	}

}
