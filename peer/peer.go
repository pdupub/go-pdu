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

package peer

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/galaxy"
	"golang.org/x/net/websocket"
)

var (
	errPeerNotReachable = errors.New("this peer not reachable right now")
	errArgsNotSupport   = errors.New("arguments not support")
	errMsgsNeedSplit    = errors.New("messages need split into waves")
)

const (
	// MaxMsgCountPerWave is the max number of msg per wave
	MaxMsgCountPerWave = 2
)

// Peer contain the info of websocket connection
type Peer struct {
	IP      string `json:"ip"`
	Port    uint64 `json:"port"`
	NodeKey string `json:"nodeKey"`
	Conn    *websocket.Conn
}

// New create new Peer
func New(ip string, port uint64, nodeKey string) (*Peer, error) {
	return &Peer{IP: ip, Port: port, NodeKey: nodeKey}, nil
}

// Dial build ws connection
func (p *Peer) Dial() error {
	conn, err := websocket.Dial(p.Url(), "", p.origin())
	if err != nil {
		return err
	}
	p.Conn = conn
	return nil
}

// Close the ws connection,
func (p *Peer) Close() error {
	if p.Conn != nil {
		return p.Conn.Close()
	}
	return nil
}

// Url show the Peer ws url address
func (p Peer) Url() string {
	return fmt.Sprintf("ws://%s:%d/%s", p.IP, p.Port, p.NodeKey)
}

// Connected return true if this peer is connected right now
func (p *Peer) Connected() bool {
	if p.Conn != nil {
		return true
	}
	return false
}

// SendQuestion is used to send question to peer
func (p *Peer) SendQuestion(cmd string, args ...interface{}) error {
	if !p.Connected() {
		return errPeerNotReachable
	}

	newArgs, err := p.buildArgs(args...)
	if err != nil {
		return err
	}
	wave := &galaxy.WaveQuestion{
		Cmd:  cmd,
		Args: newArgs,
	}
	_, err = galaxy.SendWave(p.Conn, wave)
	if err != nil {
		return err
	}
	return nil
}

func (p Peer) buildArgs(args ...interface{}) (result [][]byte, err error) {
	for _, arg := range args {
		var item []byte
		switch arg.(type) {
		case uint64:
			item = new(big.Int).SetUint64(arg.(uint64)).Bytes()
		case string:
			item = []byte(arg.(string))
		case *big.Int:
			item = arg.(*big.Int).Bytes()
		case []byte:
			item = arg.([]byte)
		case common.Hash:
			item = common.Hash2Bytes(arg.(common.Hash))
		default:
			return nil, errArgsNotSupport
		}
		result = append(result, item)
	}
	return result, nil
}

// SendMsg is used to send msg to peer
func (p *Peer) SendMsg(msg *core.Message) error {
	return p.SendMsgs([]*core.Message{msg})
}

// SendMsgs is used to send mulitiple msgs
func (p *Peer) SendMsgs(msgs []*core.Message) error {
	if len(msgs) > MaxMsgCountPerWave {
		msgs = msgs[:MaxMsgCountPerWave]
	}
	if !p.Connected() {
		return errPeerNotReachable
	}
	var msgsB [][]byte
	for _, msg := range msgs {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		msgsB = append(msgsB, msgBytes)
	}
	wave := &galaxy.WaveMessages{
		Msgs: msgsB,
	}
	if _, err := galaxy.SendWave(p.Conn, wave); err != nil {
		return err
	}
	return nil
}

// SendRoots is used to send 2 roots to peer
func (p *Peer) SendRoots(user0, user1 *core.User) error {
	if !p.Connected() {
		return errPeerNotReachable
	}
	var users [2]*core.User
	users[0] = user0
	users[1] = user1
	wave := &galaxy.WaveRoots{
		Users: users,
	}

	_, err := galaxy.SendWave(p.Conn, wave)
	if err != nil {
		return err
	}
	return nil
}

// origin used when peer dial
func (p Peer) origin() string {
	return fmt.Sprintf("http://%s:%d/", p.IP, p.Port)
}
