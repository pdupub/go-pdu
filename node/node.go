// Copyright 2019 The PDU Authors
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
	"github.com/pdupub/go-pdu/common/log"
	"os"
	"time"
)

// Node is struct of node
type Node struct {
	tpEnable   bool
	tpInterval uint64
}

// New is used to create new node
func New() (node *Node, err error) {
	node = &Node{
		tpInterval: uint64(1),
	}

	return node, nil
}

func (n *Node) SetTPEnable(enable bool, val uint64) error {

	if enable {
		n.tpEnable = true
		n.tpInterval = val
	} else {
		n.tpEnable = false
		n.tpInterval = 0
	}

	return nil
}

// Run the node
func (n *Node) Run(c <-chan os.Signal) {
	sigN, waitN := make(chan struct{}), make(chan struct{})
	sigTP, waitTP := make(chan struct{}), make(chan struct{})
	go n.runNode(sigN, waitN)
	log.Info("Start node server")

	if n.tpEnable {
		go n.runTimeProof(sigTP, waitTP)
		log.Info("Start time proof server")
	}

	for {
		select {
		case <-c:
			close(sigN)
			close(sigTP)

			if n.tpEnable {
				<-waitTP
			}

			<-waitN
			log.Info("Stop node")
			return

		}
	}
}

func (n *Node) runNode(sig <-chan struct{}, wait chan<- struct{}) {
	for {
		select {
		case <-sig:
			log.Info("Stop server")
			close(wait)
			return
		}
	}
}

func (n *Node) runTimeProof(sig <-chan struct{}, wait chan<- struct{}) {
	for {
		select {
		case <-sig:
			log.Info("Stop time proof server")
			close(wait)
			return
		case <-time.After(time.Second * time.Duration(n.tpInterval)):
			log.Info("Broadcast a new message")
		}
	}
}
