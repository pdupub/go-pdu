// Copyright 2019 The TTC Authors
// This file is part of the TTC library.
//
// The TTC library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The TTC library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the TTC library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/bitcoin"
	"github.com/pdupub/go-pdu/crypto/ethereum"
	"github.com/pdupub/go-pdu/crypto/pdu"
)

func SelectEngine(source string) (crypto.Engine, error) {
	var engine crypto.Engine
	switch source {
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
