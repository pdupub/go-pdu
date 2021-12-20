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

package main

import "C"

// other imports should be separate from the special Cgo import
import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/msg"
)

//export signMsg
func signMsg(privKeyHexC, messageC, refC *C.char) *C.char {
	privKey := C.GoString(privKeyHexC)
	message := C.GoString(messageC)
	ref := C.GoString(refC)

	did := new(identity.DID)
	did.LoadECDSA(privKey)
	bp, _ := core.NewQuantum(core.QuantumTypeInfo, []byte(message))

	content, _ := json.Marshal(bp)

	var refBytes []byte
	if len(ref) == 130 {
		refBytes, _ = hex.DecodeString(ref)
	} else {
		refBytes, _ = base64.StdEncoding.DecodeString(ref)
	}

	m := msg.New(content, refBytes)
	sm := msg.SignedMsg{Message: *m}
	sm.Sign(did)
	smBytes, _ := json.Marshal(sm)

	return C.CString(string(smBytes))
}

//export getAddress
func getAddress(privKeyHexC *C.char) *C.char {
	privKey := C.GoString(privKeyHexC)
	did := new(identity.DID)
	did.LoadECDSA(privKey)
	return C.CString(did.GetKey().Address.Hex())
}

//export generateKey
func generateKey() *C.char {
	privKey, err := identity.GeneratePrivateKey()
	if err != nil {
		return C.CString("")
	}
	return C.CString(hex.EncodeToString(privKey))
}

//export generateKeystore
func generateKeystore(privKeyHexC, passwordC *C.char) *C.char {
	privKeyBytes, err := hex.DecodeString(C.GoString(privKeyHexC))
	if err != nil {
		return C.CString("")
	}
	passwd := []byte(C.GoString(passwordC))

	fileJSON, err := identity.GenerateKeystore(privKeyBytes, passwd)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(fileJSON))
}

//export unlockKeystore
func unlockKeystore(fileJSONC, passwordC *C.char) *C.char {
	fileBytes := []byte(C.GoString(fileJSONC))

	privKey, err := identity.InspectKeystore(fileBytes, []byte(C.GoString(passwordC)))
	if err != nil {
		return C.CString("")
	}
	return C.CString(hex.EncodeToString(privKey))
}

//export ecrecover
func ecrecover(smC *C.char) *C.char {
	smBytes := []byte(C.GoString(smC))
	sm := new(msg.SignedMsg)
	if err := json.Unmarshal(smBytes, &sm); err != nil {
		return C.CString("")
	}
	author, err := sm.Ecrecover()
	if err != nil {
		return C.CString("")
	}
	return C.CString(author.Hex())
}

//export getParent
func getParent(bornAddrC, sigC *C.char) *C.char {
	bornAddr := C.GoString(bornAddrC)
	sig := C.GoString(sigC)
	born := common.HexToAddress(bornAddr)
	hash := crypto.Keccak256(born.Bytes())

	var sigBytes []byte
	if len(sig) == 130 {
		sigBytes, _ = hex.DecodeString(sig)
	} else {
		sigBytes, _ = base64.StdEncoding.DecodeString(sig)
	}

	pubkey, err := crypto.Ecrecover(hash, sigBytes)
	if err != nil {
		return C.CString("")
	}

	signer := common.Address{}
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	return C.CString(signer.Hex())
}

func main() {}
