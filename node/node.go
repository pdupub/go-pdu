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
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/peer"
	"golang.org/x/net/websocket"
)

const (
	displayInterval   = 1000
	maxLoadPeersCount = 1000
	checkPeerInterval = 10
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
	initStep             uint64
}

// New is used to create new node
func New(udb db.UDB) (node *Node, err error) {
	node = &Node{
		udb:        udb,
		tpInterval: uint64(1),
		localPort:  DefaultLocalPort,
		peers:      make(map[common.Hash]*peer.Peer),
	}

	if err := node.loadUniverse(); err != nil {
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

// AddPeer add peer to local node peers
func (n *Node) AddPeer(p *peer.Peer) error {
	peerBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = n.udb.Set(db.BucketPeer, common.Hash2String(p.ID()), peerBytes)
	if err != nil {
		return err
	}

	if po, ok := n.peers[p.ID()]; !ok || (po.Url() != p.Url() && p.NodeKey != n.localNodeKey) {
		n.peers[p.ID()] = p
	}
	return nil
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
		currentPeer, err := peer.New(ip, port, nodeKey)
		if err != nil {
			return err
		}
		userIDHash, err := common.String2Hash(userID)
		if err != nil {
			return err
		}
		currentPeer.SetUserID(userIDHash)
		err = n.AddPeer(currentPeer)
		if err != nil {
			return err
		}
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
			log.Error(err)
			continue
		}
		h, err := common.String2Hash(row.K)
		if err != nil {
			log.Error(err)
			continue
		}
		if newPeer.NodeKey != n.localNodeKey {
			n.peers[h] = &newPeer
			log.Info("Peers load", newPeer.Url(), "by", common.Hash2String(h))
		}
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
	// create universe if need
	if n.initStep < db.StepRootsSaved {
		for {
			if n.initStep >= db.StepRootsSaved {
				break
			}
			select {
			case <-time.After(time.Second):
				n.updatePeersStatus()
				n.syncCreateUniverse()
			case <-sig:
				log.Info("Stop server")
				close(wait)
				return

			}
		}
	}

	// run node
	for {
		select {
		case <-time.After(time.Second * time.Duration(checkPeerInterval)):
			log.Info("Update peers status")
			n.updatePeersStatus()
			n.syncPeers()
			n.syncPingPong()
			n.syncMsgFromPeers()
		case <-sig:
			log.Info("Stop server")
			close(wait)
			return
		}
	}
}

func (n *Node) updatePeersStatus() {
	for _, p := range n.peers {
		if !p.Connected() {
			if err := p.Dial(); err != nil {
				log.Error(err)
			}
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
			lastMsg, err := db.GetLastMsg(n.udb)
			if err != nil {
				log.Error(err)
				continue
			}
			refs = append(refs, &core.MsgReference{SenderID: lastMsg.SenderID, MsgID: lastMsg.ID()})
			// load last msg from unlock user if exist
			lastMsgByUser, err := db.GetLastMsgByUser(n.udb, n.tpUnlockedUser.ID())
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
	for _, p := range n.peers {
		if !p.Connected() {
			continue
		}
		if err := p.SendMsg(msg); err != nil {
			return err
		}
	}
	return nil
}

func (n Node) saveMsg(msg *core.Message) error {
	if err := n.universe.AddMsg(msg); err != nil {
		return err
	}
	if err := db.SaveMsg(n.udb, msg); err != nil {
		return err
	}
	return nil
}

func (n *Node) loadUniverse() (err error) {
	stepBytes, err := n.udb.Get(db.BucketConfig, db.ConfigCurrentStep)
	if err != nil {
		return err
	}
	currentStep := new(big.Int).SetBytes(stepBytes).Uint64()
	var user0, user1 *core.User
	if currentStep < db.StepRootsSaved {
		return nil
	}
	user0, user1, err = db.GetRootUsers(n.udb)
	if err != nil {
		return err
	}
	// update init step
	n.initStep = db.StepRootsSaved
	log.Info("root0", common.Hash2String(user0.ID()))
	log.Info("root1", common.Hash2String(user1.ID()))
	n.universe, err = core.NewUniverse(user0, user1)
	if err != nil {
		return err
	}
	msgCount, err := db.GetMsgCount(n.udb)
	if err != nil {
		return err
	}
	for i := uint64(0); i < msgCount.Uint64(); i++ {
		// todo : replace by db.GetMsgByOrder()
		mid, err := n.udb.Get(db.BucketMID, new(big.Int).SetUint64(i).String())
		if err != nil {
			return err
		}
		msgBytes, err := n.udb.Get(db.BucketMsg, common.Bytes2String(mid))
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
