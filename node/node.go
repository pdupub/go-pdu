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
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/peer"
	"golang.org/x/net/websocket"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	displayInterval   = 1000
	maxLoadPeersCount = 1000
)

var (
	errParseNodeAddressFail = errors.New("parse node address fail")
)

// Node is struct of node
type Node struct {
	udb                  db.UDB
	tpEnable             bool
	tpInterval           uint64
	universe             *core.Universe
	tpUnlockedUser       *core.User
	tpUnlockedPrivateKey *crypto.PrivateKey
	localPort            uint64
	localNodeKey         string
	peers                map[common.Hash]*peer.Peer
}

// New is used to create new node
func New(udb db.UDB) (node *Node, err error) {
	node = &Node{
		udb:        udb,
		tpInterval: uint64(1),
		localPort:  DefaultLocalPort,
		peers:      make(map[common.Hash]*peer.Peer),
	}

	if err := node.initUniverse(); err != nil {
		return nil, err
	}

	if err := node.loadMessage(); err != nil {
		return nil, err
	}

	if err := node.initNetwork(); err != nil {
		return nil, err
	}

	return node, nil
}

// SetLocalPort set local listen port
func (n *Node) SetLocalPort(port uint64) {
	n.localPort = port
}

// SetNodes set the target nodes [userid@ip:port/nodeKey]
func (n *Node) SetNodes(nodes string) error {
	for _, nodeStr := range strings.Split(nodes, ",") {
		var userID, ip, nodeKey string
		res := strings.Split(nodeStr, "@")
		if len(res) != 2 {
			return errParseNodeAddressFail
		}
		userID = res[0]
		res = strings.Split(res[1], ":")
		if len(res) != 2 {
			return errParseNodeAddressFail
		}
		ip = res[0]
		res = strings.Split(res[1], "/")
		if len(res) != 2 {
			return errParseNodeAddressFail
		}
		nodeKey = res[1]
		port, err := strconv.ParseUint(res[0], 10, 64)
		if err != nil {
			return err
		}
		currentPeer := peer.Peer{IP: ip, Port: port, NodeKey: nodeKey}
		peerBytes, err := json.Marshal(currentPeer)
		if err != nil {
			return err
		}
		err = n.udb.Set(db.BucketPeer, userID, peerBytes)
		if err != nil {
			return err
		}
		n.peers[common.Bytes2Hash([]byte(userID))] = &currentPeer
	}

	return nil
}

func (n *Node) initNetwork() error {
	// set local node key
	if err := n.setLocalNodeKey(); err != nil {
		return err
	}

	// load peers from db
	if err := n.loadPeers(); err != nil {
		return err
	}

	return nil
}

func (n *Node) loadPeers() error {
	rows, err := n.udb.Find(db.BucketPeer, "", maxLoadPeersCount)
	if err != nil {
		return err
	}

	for _, row := range rows {
		var newPeer peer.Peer
		if err := json.Unmarshal(row.V, &newPeer); err != nil {
			return err
		}
		n.peers[common.Bytes2Hash([]byte(row.K))] = &newPeer
	}
	return nil
}

// setLocalNodeKey set the local node key
func (n *Node) setLocalNodeKey() error {
	nodeKey, err := n.udb.Get(db.BucketConfig, db.ConfigLocalNodeKey)
	if err != nil {
		return err
	}

	if nodeKey == nil {
		h := md5.New()
		currentTimestamp := time.Now().UnixNano()
		h.Write([]byte(fmt.Sprintf("%d", currentTimestamp)))
		newNodeKey := h.Sum(nil)
		n.localNodeKey = common.Bytes2String(newNodeKey)
		n.udb.Set(db.BucketConfig, db.ConfigLocalNodeKey, newNodeKey)
		log.Info("Create new local node key", n.localNodeKey)
	} else {
		n.localNodeKey = common.Bytes2String(nodeKey)
		log.Info("Load local node key", n.localNodeKey)
	}
	return nil
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
	go n.runLocalServe()
	log.Info("Start listen on port", n.localPort)

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

func (n Node) wsHandler(ws *websocket.Conn) {
	var err error
	var msg string
	for {
		if err = websocket.Message.Receive(ws, &msg); err != nil {
			break
		}
		if err = websocket.Message.Send(ws, msg); err != nil {
			break
		}
	}
}

func (n Node) nodeHandler(w http.ResponseWriter, r *http.Request) {
	// todo: w.Write return the basic information of local node
	switch r.Method {
	case "GET":
		w.Write([]byte(r.Method))
	case "POST":
		w.Write([]byte(r.Method))
	default:
		w.Write([]byte(""))
	}
}

func (n *Node) runLocalServe() {
	http.Handle("/"+n.localNodeKey, websocket.Handler(n.wsHandler))
	http.HandleFunc("/node", n.nodeHandler)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", n.localPort), nil); err != nil {
		log.Error("Start local ws serve fail", err)
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
			if err := n.broadcastMsg(tpMsg); err != nil {
				log.Error(err)
				continue
			}

			log.Info("A new message", common.Hash2String(tpMsg.ID()), "just be created and broadcast")
		}
	}
}

func (n Node) broadcastMsg(msg *core.Message) error {
	for _, peer := range n.peers {
		if !peer.Connected() {
			continue
		}
		if err := peer.SendMsg(msg); err != nil {
			return err
		}
	}
	return nil
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
