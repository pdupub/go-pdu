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
	"fmt"
	"golang.org/x/net/websocket"
)

// Peer contain the info of websocket connection
type Peer struct {
	IP      string `json:"ip"`
	Port    uint64 `json:"port"`
	NodeKey string `json:"nodeKey"`
	conn    *websocket.Conn
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
	p.conn = conn
	return nil
}

// Close the ws connection,
func (p *Peer) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// Url show the Peer ws url address
func (p Peer) Url() string {
	return fmt.Sprintf("ws://%s:%d/%s", p.IP, p.Port, p.NodeKey)
}

// origin used when peer dial
func (p Peer) origin() string {
	return fmt.Sprintf("http://%s:%d/", p.IP, p.Port)
}
