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
	"github.com/pdupub/go-pdu/crypto/ethereum"
	"github.com/pdupub/go-pdu/crypto/pdu"
)

// Auth contain public key
type Auth struct {
	crypto.PublicKey
}

// UnmarshalJSON is used to unmarshal json
func (a *Auth) UnmarshalJSON(input []byte) error {
	aMap := make(map[string]interface{})
	err := json.Unmarshal(input, &aMap)
	if err != nil {
		return err
	}
	a.Source = aMap["source"].(string)
	a.SigType = aMap["sigType"].(string)
	switch a.Source {
	case pdu.SourceName:
		pk, err := pdu.UnmarshalJSON(input)
		if err != nil {
			return err
		}
		a.PublicKey = *pk
	case ethereum.SourceName:
		pk, err := ethereum.UnmarshalJSON(input)
		if err != nil {
			return err
		}
		a.PublicKey = *pk
	default:
		return crypto.ErrSourceNotMatch
	}
	return nil
}

// MarshalJSON marshal public key to json
func (a Auth) MarshalJSON() ([]byte, error) {
	switch a.Source {
	case pdu.SourceName:
		return pdu.MarshalJSON(a.PublicKey)
	case ethereum.SourceName:
		return ethereum.MarshalJSON(a.PublicKey)
	default:
		return nil, crypto.ErrSourceNotMatch
	}

}
