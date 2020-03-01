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
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/galaxy"
)

func (n *Node) askPeers(pid common.Hash) error {
	p := n.peers[pid]
	localPeerBytes, err := json.Marshal(n.localPeer())
	if err != nil {
		return err
	}
	waveID := common.CreateHash()
	if err := p.SendQuestion(waveID, galaxy.CmdPeers, localPeerBytes); err != nil {
		return err
	}
	if err := n.recordQuestion(pid, waveID); err != nil {
		return err
	}
	return nil
}

func (n *Node) askPing(pid common.Hash) error {
	p := n.peers[pid]
	// ping each of peer
	waveID := common.CreateHash()
	if err := p.SendPing(waveID); err != nil {
		n.removePeer(pid)
		return err
	}
	if err := n.recordPing(pid, waveID); err != nil {
		return err
	}
	return nil
}

func (n *Node) askRoots(pid common.Hash) error {
	p := n.peers[pid]
	waveID := common.CreateHash()
	if err := p.SendQuestion(waveID, galaxy.CmdRoots); err != nil {
		return err
	} else {
		if err := n.recordQuestion(pid, waveID); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) askMsg(pid common.Hash) error {
	p := n.peers[pid]
	// get current last message
	lastMsg, err := db.GetLastMsg(n.udb)
	var lastMsgID common.Hash
	if err != nil && err != db.ErrMessageNotFound {
		return err
	}
	if lastMsg != nil {
		lastMsgID = lastMsg.ID()
	}

	if lastMsg != nil && lastMsgID == n.lastSyncMsg {
		return errNoNewMsgSync
	}

	waveID := common.CreateHash()
	if err := p.SendQuestion(waveID, galaxy.CmdMessages, lastMsgID); err != nil {
		return err
	}
	if err := n.recordQuestion(pid, waveID); err != nil {
		return err
	}
	return nil
}

func (n *Node) reAskMsg(waveID common.Hash) error {
	if r, ok := n.questionRecord[waveID]; ok {
		if left, ok := n.peerSyncCnt[r.pid]; ok && left > 0 {
			n.peerSyncCnt[r.pid] = left - 1
			if err := n.askMsg(r.pid); err != nil {
				n.peerSyncCnt[r.pid] = 0
				return err
			}
		}
	}
	return nil
}

func (n *Node) delRecord(waveID common.Hash, cmd string) {
	if cmd == galaxy.CmdPong {
		delete(n.pingpongRecord, waveID)
	} else {
		delete(n.questionRecord, waveID)
	}
}

func (n *Node) recordQuestion(peerID, waveID common.Hash) error {
	if _, ok := n.questionRecord[waveID]; !ok {
		n.questionRecord[waveID] = &Record{pid: peerID, delay: 0}
	} else {
		return errDuplicateWaveID
	}
	return nil
}

func (n *Node) recordPing(peerID, waveID common.Hash) error {
	if _, ok := n.pingpongRecord[waveID]; !ok {
		n.pingpongRecord[waveID] = &Record{pid: peerID, delay: 0}
	} else {
		return errDuplicateWaveID
	}
	return nil
}
