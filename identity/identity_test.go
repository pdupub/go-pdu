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

package identity

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pdupub/go-pdu/params"
)

type TestInfo struct {
	Array  []uint64          `json:"array"`
	Map    map[string]uint64 `json:"map"`
	Number uint64            `json:"number"`
	String string            `json:"string"`
}

func TestSignMsg(t *testing.T) {
	did, _ := New()
	privateKeyHex := "689ac13dc3f424c8d5a6ef07a2e443311fc40ae4c370dac127bf5c1267e1ac98"
	t.Log("Private key :", privateKeyHex)
	if err := did.LoadECDSA(privateKeyHex); err != nil {
		t.Error(err)
	}
	t.Log("Public key :", did.key.PrivateKey.PublicKey)

	message := "Hello"
	t.Log("Message :", message)
	sig, err := did.Sign([]byte(message))
	if err != nil {
		t.Error(err)
	}
	t.Log("address", did.key.Address.Hex())
	t.Log("signature", common.Bytes2Hex(sig))
}

func TestNew(t *testing.T) {
	did, _ := New()
	did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)

	// build content
	rand.Seed(time.Now().UnixNano())
	testInfo := TestInfo{Number: uint64(rand.Intn(100)), String: "Hello World!", Map: make(map[string]uint64), Array: []uint64{1, 20, 3}}
	testInfo.Map["b"] = 1
	testInfo.Map["a"] = 2
	testInfo.Map["f"] = 3
	testInfo.Map["h"] = 10
	testInfo.Map["ff"] = 1000
	testInfo.Map["b23"] = 10

	contentBytes, err := json.Marshal(testInfo)
	if err != nil {
		t.Error(err)
	}

	hash := crypto.Keccak256(contentBytes)
	sig, err := crypto.Sign(hash, did.key.PrivateKey)
	if err != nil {
		t.Error(err)
	}

	// verify
	pubkey, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		t.Error(err)
	}

	signer := common.Address{}
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	if did.GetKey().Address.Hex() != signer.Hex() {
		t.Error("verify signer fail")
	}

	addr, pub, private, err := did.Inspect(true)
	if err != nil {
		t.Error(err)
	}
	t.Log(addr)
	t.Log(pub)
	t.Log(private)

	did2 := new(DID)
	did2.LoadECDSA(private)
	addr, pub, private, err = did2.Inspect(true)
	if err != nil {
		t.Error(err)
	}
	t.Log(addr)
	t.Log(pub)
	t.Log(private)

}
