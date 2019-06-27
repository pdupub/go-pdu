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
	"encoding/json"
	"fmt"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/pdu"
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
				var pubKey interface{}
				err := json.Unmarshal([]byte(aMap["pubKey"].(string)), &pubKey)
				if err != nil {
					return err
				} else {
					fmt.Println("s2pk", pubKey)
					a.PubKey = pubKey
				}
			} else if a.SigType == pdu.MultipleSignatures {
				pubKey := make([]interface{}, 0)
				err := json.Unmarshal([]byte(aMap["pubKey"].(string)), &pubKey)
				if err != nil {
					return err
				} else {
					fmt.Println("ms:", pubKey)
					a.PubKey = pubKey
				}
			} else {

			}
		}

	}
	return nil
}

func (a *Auth) MarshalJSON() ([]byte, error) {
	aMap := make(map[string]interface{})
	aMap["source"] = a.Source
	aMap["sigType"] = a.SigType
	if pubKey, err := json.Marshal(a.PubKey); err != nil {
		return []byte{}, err
	} else {
		aMap["pubKey"] = string(pubKey)
	}
	return json.Marshal(aMap)
}
