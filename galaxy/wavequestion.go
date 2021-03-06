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

import "github.com/pdupub/go-pdu/common"

// WaveQuestion implements the Wave interface and represents request info message.
type WaveQuestion struct {
	WaveID common.Hash `json:"waveID"`
	Cmd    string      `json:"cmd"`
	Args   [][]byte    `json:"args"`
}

// Command returns the protocol command string for the wave.
func (w *WaveQuestion) Command() string {
	return CmdQuestion
}
