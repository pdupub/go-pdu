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
	"crypto/ecdsa"
	"encoding/json"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/ethereum"
	"testing"
)

func TestAuth_MarshalJSON(t *testing.T) {

	_, pubKey, err := ethereum.GenKey(crypto.Signature2PublicKey)
	if err != nil {
		t.Errorf("pdu genereate key fail, err: %s", err)
	}
	auth := Auth{
		*pubKey,
	}

	authBytes, err := json.Marshal(auth)
	if err != nil {
		t.Errorf("auth marshal fail, err: %s", err)
	}

	var targetAuth Auth
	err = json.Unmarshal(authBytes, &targetAuth)
	if err != nil {
		t.Errorf("auth unmarshal fail, err: %s", err)
	}

	if pubKey.Source != targetAuth.PublicKey.Source || pubKey.SigType != targetAuth.PublicKey.SigType {
		t.Errorf("pubkey info mismatch")
	}

	pubKey1 := pubKey.PubKey.(ecdsa.PublicKey)
	pubKey2 := targetAuth.PubKey.(ecdsa.PublicKey)
	if pubKey1.X.Cmp(pubKey2.X) != 0 || pubKey1.Y.Cmp(pubKey2.Y) != 0 {
		t.Errorf("public key mismatch")
	}

}
