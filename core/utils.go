// Copyright 2021 The PDU Authors
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

package core

import (
	"encoding/hex"
)

type Sig []byte

func (sig *Sig) toHex() string {
	return Sig2Hex(*sig)
}

func Sig2Hex(sig Sig) string {
	return "0x" + hex.EncodeToString(sig)
}

func Hex2Sig(str string) Sig {
	if str[:2] == "0x" || str[:2] == "0X" {
		str = str[2:]
	}
	h, _ := hex.DecodeString(str)
	return h
}
