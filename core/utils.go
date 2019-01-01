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

package core

import (
	"encoding/json"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/bitcoin"
	"github.com/pdupub/go-pdu/crypto/ethereum"
	"github.com/pdupub/go-pdu/crypto/pdu"
	"strings"
)

// SelectEngine return a new engine by source type
func SelectEngine(source string) (crypto.Engine, error) {
	var engine crypto.Engine
	switch strings.ToUpper(source) {
	case crypto.BTC:
		engine = bitcoin.New()
	case crypto.PDU:
		engine = pdu.New()
	case crypto.ETH:
		engine = ethereum.New()
	default:
		return nil, crypto.ErrSourceNotMatch
	}
	return engine, nil
}

// DecryptKey decrypt private key from keyJson file
func DecryptKey(keyjson []byte, passwd string) (*crypto.PrivateKey, error) {
	var engine crypto.Engine

	keyJM := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &keyJM); err != nil {
		return nil, err
	}
	source := keyJM["source"].(string)

	engine, err := SelectEngine(source)
	if err != nil {
		return nil, err
	}

	pk, err := engine.DecryptKey(keyjson, passwd)
	if err != nil {
		return nil, err
	}
	return pk, nil
}
