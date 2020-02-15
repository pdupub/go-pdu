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

package galaxy

// WavePeers implements the Wave interface and represent node peers.
type WavePeers struct {
	Peers []string `json:"peers"` //"927171CCD704F7A5D481C04B6F16D0E73F823C69EEC9D49866258E658366D136@127.0.0.1:8341/9F1A60DF0D424EF0E80BBF0F3F90D7F2"
}

// Command returns the protocol command string for the wave.
func (w *WavePeers) Command() string {
	return CmdPeers
}
