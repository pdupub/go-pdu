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

package clan

import (
	"errors"
	"github.com/TATAUFO/PDU/accounts"
	"github.com/TATAUFO/PDU/common"
	"github.com/TATAUFO/PDU/genealogy"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	defaultInitIndividualCredit = 1000 // Number of credit when individual add to clan
)

// Various error messages to mark accounts invalid
var (
	//ErrInvalidOriginParent is returned when Origin Parent is invalid
	ErrInvalidOriginParent = errors.New("invalid parent")
)

// Clan is DAG topology of genealogy
type Clan struct {
	dag        map[common.Address]*genealogy.Individual
	size       uint64
	generation uint64
}

// New create new clan
func New(father, mother accounts.Account) (*Clan, error) {
	if len(father.Address) != common.AddressLength || len(mother.Address) != common.AddressLength {
		return nil, ErrInvalidOriginParent
	}
	Adam := genealogy.Individual{
		father,
		nil,
		nil,
		0,
		true,
		defaultInitIndividualCredit}
	Eve := genealogy.Individual{
		mother,
		nil,
		nil,
		0,
		false,
		defaultInitIndividualCredit}
	clan := Clan{
		dag:        make(map[common.Address]*genealogy.Individual),
		size:       2,
		generation: 0,
	}
	clan.dag[father.Address] = &Adam
	clan.dag[mother.Address] = &Eve
	return &clan, nil
}

// Add account into genealogy, return error if parents is missing
func (c *Clan) Add(account accounts.Account) error {
	if _, ok := c.dag[account.Address]; ok {
		return nil
	}
	fatherAddress, err := ecrecover(accounts.Account{Address: account.Address, DOB: account.DOB}, account.FatherSign[:])
	if err != nil {
		return err
	}
	if _, ok := c.dag[fatherAddress]; !ok {
		return ErrInvalidOriginParent
	}
	motherAddress, err := ecrecover(accounts.Account{Address: account.Address, DOB: account.DOB, FatherSign: account.FatherSign}, account.MotherSign[:])
	if err != nil {
		return err
	}
	if _, ok := c.dag[motherAddress]; !ok {
		return ErrInvalidOriginParent
	}
	// The generation of account is larger generation of parents plus one
	parentGeneration := c.dag[fatherAddress].Generation
	if parentGeneration < c.dag[motherAddress].Generation {
		parentGeneration = c.dag[motherAddress].Generation
	}

	c.dag[account.Address] = &genealogy.Individual{
		account,
		c.dag[fatherAddress],
		c.dag[motherAddress],
		parentGeneration + 1,
		gender(account),
		defaultInitIndividualCredit}
	c.size += 1
	if parentGeneration+1 > c.generation {
		c.generation = parentGeneration + 1
	}
	return nil
}

func gender(account accounts.Account) bool {
	accountHash := common.ToHash(account)
	return new(big.Int).Mod(new(big.Int).SetBytes(accountHash), big.NewInt(2)).Uint64() == uint64(1)

}

func ecrecover(account accounts.Account, sig []byte) (common.Address, error) {
	accountHash := common.ToHash(account)
	sigPUblicKey, err := crypto.Ecrecover(common.ToMD5(accountHash), sig)
	if err != nil {
		return common.Address{}, err
	}
	publicKey, err := crypto.UnmarshalPubkey(sigPUblicKey)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(crypto.PubkeyToAddress(*publicKey).Bytes()), nil
}
