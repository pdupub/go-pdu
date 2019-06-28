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

package user

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/pdu"
	"math/big"
)

type Auth struct {
	crypto.PublicKey
}

func (a *Auth) UnmarshalJSON(input []byte) error {
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return err
	} else {
		a.Source = aMap["source"].(string)
		a.SigType = aMap["sigType"].(string)
		if a.Source == pdu.SourceName {
			if a.SigType == pdu.Signature2PublicKey {
				pubKey := new(ecdsa.PublicKey)
				pubKey.Curve = elliptic.P256()
				pubKey.X, pubKey.Y = big.NewInt(0), big.NewInt(0)
				pk := aMap["pubKey"].([]interface{})
				pubKey.X.UnmarshalText([]byte(pk[0].(string)))
				pubKey.Y.UnmarshalText([]byte(pk[1].(string)))
				a.PubKey = *pubKey
			} else if a.SigType == pdu.MultipleSignatures {
				pk := aMap["pubKey"].([]interface{})
				var pubKeys []ecdsa.PublicKey
				for i := 0; i < len(pk)/2; i++ {
					pubKey := new(ecdsa.PublicKey)
					pubKey.Curve = elliptic.P256()
					pubKey.X, pubKey.Y = big.NewInt(0), big.NewInt(0)
					pubKey.X.UnmarshalText([]byte(pk[i*2].(string)))
					pubKey.Y.UnmarshalText([]byte(pk[i*2+1].(string)))
					pubKeys = append(pubKeys, *pubKey)
				}
				a.PubKey = pubKeys
			} else {
				// todo : err
			}
		}

	}
	return nil
}

func (a Auth) MarshalJSON() ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if a.Source == pdu.SourceName {
		if a.SigType == pdu.Signature2PublicKey {
			pk := a.PubKey.(ecdsa.PublicKey)
			pubKey := make([]string, 2)
			pubKey[0] = pk.X.String()
			pubKey[1] = pk.Y.String()
			aMap["pubKey"] = pubKey
		} else if a.SigType == pdu.MultipleSignatures {
			pks := a.PubKey.([]interface{})
			pubKey := make([]string, len(pks)*2)
			for i, v := range pks {
				pk := v.(ecdsa.PublicKey)
				pubKey[i*2] = pk.X.String()
				pubKey[i*2+1] = pk.Y.String()
			}
			aMap["pubKey"] = pubKey
		}
	}
	return json.Marshal(aMap)
}
