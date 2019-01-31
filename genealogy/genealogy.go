// Copyright 2018 The PDU Authors
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

package genealogy

import (
	"github.com/pdupub/PDU/accounts"
	"github.com/pdupub/PDU/common"
)

// Individual is account link to parents, except two in genesis config
type Individual struct {
	Account    accounts.Account `json:"account"`    // Account
	Father     *Individual      `json:"father"`     // Father
	Mother     *Individual      `json:"mother"`     // Mother
	Generation uint64           `json:"generation"` // Generation
	Gender     bool             `json:"gender"`     // Gender
	Credit     uint64           `json:"credit"`     // Credit

}

// Genealogy represents topology of all known accounts
type Genealogy interface {
	// Add returns error if add account to genealogy fail
	Add(account accounts.Account) error
	// Get returns the Individual by address
	Get(address common.Address) (Individual, error)
	// Size returns the size of all account
	Size() uint64
	// Generation returns the max generation
	Generation() uint64
}
