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

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
)

//export signMsgSample
func signMsgSample(privKeyHexC, messageC, refC *C.char) *C.char {
	// load did
	privKey := C.GoString(privKeyHexC)
	did := new(identity.DID)
	did.LoadECDSA(privKey)

	// reference
	ref := C.GoString(refC)
	var refBytes []byte
	if len(ref) == 130 {
		refBytes, _ = hex.DecodeString(ref)
	} else {
		refBytes, _ = base64.StdEncoding.DecodeString(ref)
	}

	// quantum
	message := C.GoString(messageC)
	content, _ := core.NewContent(core.QCFmtStringTEXT, []byte(message))
	q, _ := core.NewQuantum(core.QuantumTypeInformation, []*core.QContent{content}, refBytes)

	// add signature
	q.Sign(did)
	qBytes, _ := json.Marshal(q)

	return C.CString(string(qBytes))
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
func ecrecover(qC *C.char) *C.char {
	qBytes := []byte(C.GoString(qC))
	q := new(core.Quantum)
	if err := json.Unmarshal(qBytes, &q); err != nil {
		return C.CString("")
	}
	author, err := q.Ecrecover()
	if err != nil {
		return C.CString("")
	}
	return C.CString(author.Hex())
}

func main() {}
