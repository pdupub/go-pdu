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
	"strings"

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/galaxy"
	"github.com/pdupub/go-pdu/peer"
)

const (
	maxQuestionPerWave = 30
)

func (n *Node) syncCreateUniverse() {
	log.Info("Start sync universe start", "create universe")
	for _, p := range n.peers {
		if !p.Connected() {
			continue
		}
		if err := p.SendQuestion(galaxy.CmdRoots); err != nil {
			log.Error(err)
			continue
		}

		w, err := galaxy.ReceiveWave(p.Conn)
		if err != nil {
			log.Error(err)
			continue
		}
		if w.Command() == galaxy.CmdRoots {
			mw := w.(*galaxy.WaveRoots)
			user0 := mw.Users[0]
			user1 := mw.Users[1]
			log.Info("user0", common.Hash2String(user0.ID()))
			log.Info("user1", common.Hash2String(user1.ID()))
			// update init step
			n.initStep = db.StepRootsSaved
			n.universe, err = core.NewUniverse(user0, user1)
			if err != nil {
				log.Error(err)
				continue
			}
			if err := db.SaveRootUsers(n.udb, mw.Users[:]); err != nil {
				log.Error(err)
				continue
			}
			break
		}
	}
}

func (n *Node) syncPeers() {
	log.Info("Start sync peers from othe node")
	for _, p := range n.peers {
		if !p.Connected() {
			continue
		}
		if err := p.SendQuestion(galaxy.CmdPeers); err != nil {
			log.Error(err)
			continue
		}

		w, err := galaxy.ReceiveWave(p.Conn)
		if err != nil {
			log.Error(err)
			continue
		}
		if w.Command() == galaxy.CmdPeers {
			mw := w.(*galaxy.WavePeers)
			for _, op := range mw.Peers {
				log.Debug("Peer address", op)
			}
			n.SetNodes(strings.Join(mw.Peers, ","))
		}
	}
}

func (n *Node) syncPingPong() {
	log.Info("Ping other peers")
	for _, p := range n.peers {
		if !p.Connected() {
			// todo : try to re dial peer
			continue
		}
		if err := p.SendPing(); err != nil {
			log.Error(err)
			continue
		}

		w, err := galaxy.ReceiveWave(p.Conn)
		if err != nil {
			log.Error(err)
			continue
		}
		if w.Command() == galaxy.CmdPong {
			log.Debug("Pong received success")
			continue
		}
		log.Error("Pong received fail")
	}
}

func (n *Node) syncMsgFromPeers() {
	log.Info("Start Sync message from peers")
	for _, p := range n.peers {
		if !p.Connected() {
			continue
		}
		// get current last message
		lastMsg, err := db.GetLastMsg(n.udb)
		var lastMsgID common.Hash
		if err != nil && err != db.ErrMessageNotFound {
			log.Error(err)
			return
		}
		if lastMsg != nil {
			lastMsgID = lastMsg.ID()
		}

		var msgs []*core.Message
		for i := 0; i < maxQuestionPerWave; i++ {
			resMsg, err := n.syncMsg(p, lastMsgID)
			if err != nil {
				log.Error(err)
				break
			}
			if len(resMsg) == 0 {
				break
			}

			lastMsg = resMsg[len(resMsg)-1]
			if lastMsgID == lastMsg.ID() {
				break
			}
			lastMsgID = lastMsg.ID()
			msgs = append(msgs, resMsg...)
		}
		if len(msgs) > 0 {
			log.Debug("Sync", len(msgs), "message from", common.Hash2String(msgs[0].ID()), "to", common.Hash2String(msgs[len(msgs)-1].ID()))
			for _, msg := range msgs {
				if err := n.saveMsg(msg); err != nil {
					log.Error(err)
					break
				}
			}
		}
	}
}

func (n *Node) syncMsg(p *peer.Peer, lastMsgID common.Hash) ([]*core.Message, error) {
	var msgs []*core.Message
	// send question
	if err := p.SendQuestion(galaxy.CmdMessages, lastMsgID); err != nil {
		return nil, err
	}

	// recevie message
	w, err := galaxy.ReceiveWave(p.Conn)
	if err != nil {
		return nil, err
	}

	// check msg cmd
	if w.Command() == galaxy.CmdMessages {
		mw := w.(*galaxy.WaveMessages)
		for _, mb := range mw.Msgs {
			var msg core.Message
			err := json.Unmarshal(mb, &msg)
			if err != nil {
				return msgs, err
			}
			msgs = append(msgs, &msg)
		}
	}
	return msgs, nil
}
