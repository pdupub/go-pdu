// Copyright 2021 The PDU Authors
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
	"context"
	"log"
	"os"
	"time"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/udb/fb"
)

// Node
type Node struct {
	interval int64
	univ     core.Universe
}

func New(interval int64, firebaseKeyPath, firebaseProjectID string) (*Node, error) {
	ctx := context.Background()
	fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
	if err != nil {
		return nil, err
	}
	return &Node{interval: interval, univ: fbu}, nil
}

// Run
func (n *Node) Run(c <-chan os.Signal) {
	var runCnt uint64
	busy := false
	for {
		select {
		case <-c:
			log.Println("main closed")
			return
		case <-time.After(time.Second * time.Duration(n.interval)):
			log.Println("execution", runCnt)
			runCnt++
			if !busy {
				busy = true
				n.univ.ProcessQuantum(0, 10)
				busy = false
			}
		}
	}
}
