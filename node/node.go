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
	"encoding/json"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"math/big"
	"math/rand"
	"os"
	"time"
)

const (
	displayInterval = 1000
)

// Node is struct of node
type Node struct {
	udb                  db.UDB
	tpEnable             bool
	tpInterval           uint64
	universe             *core.Universe
	tpUnlockedUser       *core.User
	tpUnlockedPrivateKey *crypto.PrivateKey
}

// New is used to create new node
func New(udb db.UDB) (node *Node, err error) {
	node = &Node{
		udb:        udb,
		tpInterval: uint64(1),
	}
	if err := node.initUniverse(); err != nil {
		return nil, err
	}

	if err := node.loadMessage(); err != nil {
		return nil, err
	}
	return node, nil
}

// EnableTP set the time proof settings
func (n *Node) EnableTP(user *core.User, priKey *crypto.PrivateKey, val uint64) error {
	n.tpEnable = true
	n.tpUnlockedUser = user
	n.tpUnlockedPrivateKey = priKey
	n.tpInterval = val

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
			var refs []*core.MsgReference
			// load last msg in universe
			lastMsg, err := n.getLastMsg()
			if err != nil {
				log.Error(err)
				continue
			}
			refs = append(refs, &core.MsgReference{SenderID: lastMsg.SenderID, MsgID: lastMsg.ID()})
			// load last msg from unlock user if exist
			lastMsgByUser, err := n.getLastMsgByUser(n.tpUnlockedUser.ID())
			if err != nil {
				log.Error(err)
				continue
			}
			if lastMsg.ID() != lastMsgByUser.ID() {
				refs = append(refs, &core.MsgReference{SenderID: lastMsgByUser.SenderID, MsgID: lastMsgByUser.ID()})
			}
			// create new msg, use 1.2 as reference
			tpMsgValue := &core.MsgValue{ContentType: core.TypeText, Content: []byte(string(rand.Intn(100000)))}
			tpMsg, err := core.CreateMsg(n.tpUnlockedUser, tpMsgValue, n.tpUnlockedPrivateKey, refs...)
			if err != nil {
				log.Error(err)
				continue
			}
			// save msg into udb,
			if err := n.saveMsg(tpMsg); err != nil {
				log.Error(err)
				continue
			}
			// 5. broadcast the new msg

			log.Info("A new message", common.Hash2String(tpMsg.ID()), "just be created and broadcast")
		}
	}
}

func (n Node) saveMsg(msg *core.Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	countBytes, err := n.udb.Get(db.BucketConfig, db.ConfigMsgCount)
	if err != nil {
		return err
	}
	count := new(big.Int).SetBytes(countBytes)
	err = n.udb.Set(db.BucketMsg, common.Hash2String(msg.ID()), msgBytes)
	if err != nil {
		return err
	}

	err = n.udb.Set(db.BucketMID, count.String(), []byte(common.Hash2String(msg.ID())))
	if err != nil {
		return err
	}
	count = count.Add(count, big.NewInt(1))
	err = n.udb.Set(db.BucketConfig, db.ConfigMsgCount, count.Bytes())
	if err != nil {
		return err
	}

	err = n.udb.Set(db.BucketLastMID, common.Hash2String(msg.SenderID), []byte(common.Hash2String(msg.ID())))
	if err != nil {
		return err
	}
	return nil
}

func (n Node) getLastMsgByUser(userID common.Hash) (*core.Message, error) {
	var msg core.Message
	lastMsgBytes, err := n.udb.Get(db.BucketLastMID, common.Hash2String(userID))
	if err != nil {
		return nil, err
	}
	msgBytes, err := n.udb.Get(db.BucketMsg, string(lastMsgBytes))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (n Node) getLastMsg() (*core.Message, error) {
	var msg core.Message
	countBytes, err := n.udb.Get(db.BucketConfig, db.ConfigMsgCount)
	if err != nil {
		return nil, err
	}
	count := new(big.Int).SetBytes(countBytes)
	mid, err := n.udb.Get(db.BucketMID, count.Sub(count, big.NewInt(1)).String())
	if err != nil {
		return nil, err
	}
	msgBytes, err := n.udb.Get(db.BucketMsg, string(mid))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (n *Node) initUniverse() error {
	var user0, user1 core.User

	root0, err := n.udb.Get(db.BucketConfig, db.ConfigRoot0)
	if err != nil {
		return err
	}
	root1, err := n.udb.Get(db.BucketConfig, db.ConfigRoot1)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(root0, &user0); err != nil {
		return err
	}
	if err := json.Unmarshal(root1, &user1); err != nil {
		return err
	}
	log.Info("root0", common.Hash2String(user0.ID()))
	log.Info("root1", common.Hash2String(user1.ID()))
	n.universe, err = core.NewUniverse(&user0, &user1)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) loadMessage() error {
	cntBytes, err := n.udb.Get(db.BucketConfig, db.ConfigMsgCount)
	if err != nil {
		return err
	}
	msgCount := new(big.Int).SetBytes(cntBytes).Uint64()
	for i := uint64(0); i < msgCount; i++ {
		mid, err := n.udb.Get(db.BucketMID, new(big.Int).SetUint64(i).String())
		if err != nil {
			return err
		}
		msgBytes, err := n.udb.Get(db.BucketMsg, string(mid))
		if err != nil {
			return err
		}
		var msg core.Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			return err
		}
		err = n.universe.AddMsg(&msg)
		if err != nil {
			return err
		}
		if i%displayInterval == 0 {
			log.Info(i+1, "messages be loaded")
		}
	}
	log.Info("All", msgCount, "messages already be loaded")
	return nil
}
