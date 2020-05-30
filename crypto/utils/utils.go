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

package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/bitcoin"
	"github.com/pdupub/go-pdu/crypto/ethereum"
	"github.com/pdupub/go-pdu/crypto/pdu"
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

// DecryptKey decrypt private key from keyJSON file
func DecryptKey(keyJSON []byte, passwd string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	var engine crypto.Engine

	keyJM := make(map[string]interface{})
	if err := json.Unmarshal(keyJSON, &keyJM); err != nil {
		return nil, nil, err
	}
	source := keyJM["source"].(string)

	engine, err := SelectEngine(source)
	if err != nil {
		return nil, nil, err
	}

	return engine.DecryptKey(keyJSON, passwd)
}

// DisplayKey decrypt private key from keyJSON file
func DisplayKey(privKey *crypto.PrivateKey, pubKey *crypto.PublicKey) error {
	var engine crypto.Engine
	engine, err := SelectEngine(privKey.Source)
	if err != nil {
		return err
	}
	privM, pubM, err := engine.MappingKey(privKey, pubKey)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("Private Key:")
	if source, ok := privM["source"]; ok {
		fmt.Println("source\t", source)
	}
	if sigType, ok := privM["sigType"]; ok {
		fmt.Println("sigType\t", sigType)
	}
	if k, ok := privM["privKey"]; ok {
		fmt.Println("key\t", k)
	}
	fmt.Println()
	fmt.Println("Public Key:")
	if source, ok := pubM["source"]; ok {
		fmt.Println("source\t", source)
	}
	if sigType, ok := pubM["sigType"]; ok {
		fmt.Println("sigType\t", sigType)
	}
	if k, ok := pubM["pubKey"]; ok {
		fmt.Println("key\t", k)
	}

	return nil
}
