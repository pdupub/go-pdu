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

package msg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/pdupub/go-pdu/identity"
)

// SignedMsg is the struct of signed Message
type SignedMsg struct {
	Message
	Signature []byte `json:"signature"` // message signature
}

// Sign add signature to this msg
func (sm *SignedMsg) Sign(did *identity.DID) error {
	sig, err := sm.Message.Sign(did)
	if err != nil {
		return err
	}
	sm.Signature = sig
	return nil
}

// Verify msg by address
func (sm *SignedMsg) Verify(address common.Address) error {
	return sm.Message.Verify(sm.Signature, address)
}

// Ecrecover the author address from msg
func (sm *SignedMsg) Ecrecover() (common.Address, error) {
	return sm.Message.Ecrecover(sm.Signature)
}

// Post the signed message to target url
func (sm *SignedMsg) Post(url string) ([]byte, error) {

	smBytes, err := json.Marshal(sm)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: time.Second,
	}

	resp, err := client.Post(url,
		"application/json",
		bytes.NewReader(smBytes))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
