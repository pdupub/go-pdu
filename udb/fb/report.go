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

package fb

import (
	"encoding/json"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
)

type Report struct {
	Content   *core.QContent `json:"c"`
	Signature core.Sig       `json:"sig"`
}

// Ecrecover recover
func (r *Report) Ecrecover() (identity.Address, error) {
	b, err := json.Marshal(r.Content)
	if err != nil {
		return identity.Address{}, err
	}
	return identity.Ecrecover(b, r.Signature)
}

func (r *Report) Sign(did *identity.DID) error {
	b, err := json.Marshal(r.Content)
	if err != nil {
		return err
	}
	sig, err := did.Sign(b)
	if err != nil {
		return err
	}
	r.Signature = sig
	return nil
}
