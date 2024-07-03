// Copyright 2024 The PDU Authors
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
	"encoding/json"
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

func JsonToQuantum(body []byte) (*Quantum, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	quantum := &Quantum{Signature: Hex2Sig(jsonData["sig"].(string))}
	refs := []Sig{}
	for _, ref := range jsonData["refs"].([]interface{}) {
		refs = append(refs, Hex2Sig(ref.(string)))
	}
	cs := []*QContent{}
	for _, content := range jsonData["cs"].([]interface{}) {
		qc := &QContent{}
		if data, ok := content.(map[string]interface{})["data"]; ok {
			qc.Data = []byte(data.(string))
		}
		if fmt, ok := content.(map[string]interface{})["fmt"]; ok {
			qc.Format = fmt.(string)
		}

		if zipped, ok := content.(map[string]interface{})["zipped"]; ok {
			qc.Zipped = zipped.(bool)
		}
		cs = append(cs, qc)
	}
	quantum.Contents = cs
	quantum.Nonce = int(jsonData["nonce"].(float64))
	quantum.References = refs

	return quantum, nil
}
