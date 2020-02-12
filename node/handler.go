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
	"math/big"

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/galaxy"
	"github.com/pdupub/go-pdu/peer"
	"golang.org/x/net/websocket"
)

func (n Node) handleMessages(ws *websocket.Conn, w galaxy.Wave) (*core.Message, error) {
	var msg core.Message
	wm := w.(*galaxy.WaveMessages)
	for _, wmsg := range wm.Msgs {
		if err := json.Unmarshal(wmsg, &msg); err != nil {
			return nil, err
		}
	}
	// save msg (universe & udb)
	if err := n.saveMsg(&msg); err != nil {
		return nil, err
	} else if err := n.broadcastMsg(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (n Node) handleQuestion(ws *websocket.Conn, w galaxy.Wave) error {
	wm := w.(*galaxy.WaveQuestion)
	//log.Debug("Received question", wm.Cmd)
	switch wm.Cmd {
	case galaxy.CmdRoots:
		user0, user1, err := db.GetRootUsers(n.udb)
		if err != nil {
			return err
		}
		p := peer.Peer{Conn: ws}
		if err = p.SendRoots(user0, user1); err != nil {
			return err
		}
	case galaxy.CmdMessages:
		var order, count *big.Int
		var err error
		var msgs []*core.Message
		msgID := common.Bytes2Hash(wm.Args[0])

		if msgID != common.Bytes2Hash([]byte{}) {
			order, count, err = db.GetOrderCntByMsg(n.udb, msgID)
			if err != nil {
				return err
			}
			order = order.Add(order, big.NewInt(1))
		} else {
			order = big.NewInt(0)
			count, err = db.GetMsgCount(n.udb)
			if err != nil {
				return err
			}
		}

		if order != nil && count != nil && count.Uint64()-order.Uint64() > peer.MaxMsgCountPerWave {
			//log.Debug("Send msg from order", order, "size", peer.MaxMsgCountPerWave)
			msgs = db.GetMsgByOrder(n.udb, order, peer.MaxMsgCountPerWave)
		}
		p := peer.Peer{Conn: ws}
		if err = p.SendMsgs(msgs); err != nil {
			return err
		}
	}
	return nil
}

func (n Node) wsHandler(ws *websocket.Conn) {
	for {
		w, err := galaxy.ReceiveWave(ws)
		if err != nil {
			log.Error(err)
			break
		}
		if w.Command() == galaxy.CmdMessages {
			if msg, err := n.handleMessages(ws, w); err != nil {
				//log.Error(err)
				log.Debug(err)
			} else {
				log.Info("Received message", common.Hash2String(msg.ID()))
			}
		} else if w.Command() == galaxy.CmdQuestion {
			if err := n.handleQuestion(ws, w); err != nil {
				log.Error(err)
			}
		}
	}
}
