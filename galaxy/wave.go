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

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// WaveSize is the number of bytes in a wave
// WaveHeaderSize 24 bytes + waveBody
const WaveSize = 2048

// WaveHeaderSize is the number of bytes in a wave header
// command 12 bytes + length 4 bytes + checksum 4 bytes
const WaveHeaderSize = 24

// CommandSize is the fixed size of all commands
const CommandSize = 12

// Commands used in wave which describe the type of wave.
const (
	CmdQuestion = "question"
	CmdVersion  = "version"
	CmdRoots    = "roots"
	CmdMessages = "messages"
	CmdPing     = "ping"
	CmdPong     = "aler"
	CmdUser     = "user"
	CmdPeers    = "peers"
)

var (
	errWaveLengthTooLong = errors.New("wave length too long")
)

// Wave is an interface that describes a galaxy information.
type Wave interface {
	Command() string
}

func makeEmptyWave(command string) (Wave, error) {
	var wave Wave
	switch command {
	case CmdQuestion:
		wave = &WaveQuestion{}
	case CmdVersion:
		wave = &WaveVersion{}
	case CmdRoots:
		wave = &WaveRoots{}
	case CmdMessages:
		wave = &WaveMessages{}
	case CmdPing:
		wave = &WavePing{}
	case CmdPong:
		wave = &WavePong{}
	case CmdUser:
		wave = &WaveUser{}
	case CmdPeers:
		wave = &WavePeers{}
	default:
		return nil, fmt.Errorf("unhandled command [%s]", command)
	}
	return wave, nil
}

// SendWave send a wave message to w
func SendWave(w io.Writer, wave Wave) (int, error) {
	var magic, waveLen, checkSum [4]byte
	var command [CommandSize]byte
	cmd := wave.Command()
	if len(cmd) > CommandSize {
		return 0, fmt.Errorf("command [%s] is too long [max %v]", wave, CommandSize)
	}
	copy(command[:], []byte(cmd))
	copy(magic[:], []byte(""))
	copy(waveLen[:], []byte(""))
	copy(checkSum[:], []byte(""))

	waveHeader := bytes.NewBuffer(make([]byte, 0, WaveHeaderSize))
	waveHeader.Write((magic[:]))
	waveHeader.Write(command[:])
	waveHeader.Write(waveLen[:])
	waveHeader.Write(checkSum[:])
	waveBody, err := json.Marshal(wave)
	if err != nil {
		return 0, err
	}
	waveBytes := append(waveHeader.Bytes(), waveBody...)
	if len(waveBytes) > WaveSize {
		return 0, errWaveLengthTooLong
	}
	return w.Write(waveBytes)
}

// ReceiveWave receive a wave message from r
func ReceiveWave(r io.Reader) (Wave, error) {
	waveBytes := make([]byte, WaveSize)
	n, err := r.Read(waveBytes)
	if err != nil {
		return nil, err
	}
	waveHeader := waveBytes[:WaveHeaderSize]
	waveBody := waveBytes[WaveHeaderSize:n]

	// Strip trailing zeros from command string.
	command := string(bytes.TrimRight(waveHeader[4:CommandSize+4], string(0)))
	msg, err := makeEmptyWave(command)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(waveBody, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
