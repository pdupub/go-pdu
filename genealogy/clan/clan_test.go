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
	"github.com/TATAUFO/PDU/accounts"
	"github.com/TATAUFO/PDU/common"
	"github.com/ethereum/go-ethereum/crypto"
	mrand "math/rand"
	"testing"
	"time"
)

func TestUtil(t *testing.T) {
	seed := time.Now().UnixNano()
	mrand.Seed(seed)

	// generate two private key as Adam and Eve, first parents.
	var fad, mad common.Address
	fpk, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	fad.SetBytes(crypto.PubkeyToAddress(fpk.PublicKey).Bytes())
	mpk, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	mad.SetBytes(crypto.PubkeyToAddress(mpk.PublicKey).Bytes())
	fatherAccount := accounts.Account{fad, common.Signature{}, common.Signature{}, common.NatureTime{time.Now().UnixNano(), "a"}}
	motherAccount := accounts.Account{mad, common.Signature{}, common.Signature{}, common.NatureTime{time.Now().UnixNano(), "b"}}

	// create new genealogy
	clan, err := New(fatherAccount, motherAccount)
	if err != nil {
		t.Fatal(err)
	}

	// generate new private key
	var address common.Address
	cpk, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	address.SetBytes(crypto.PubkeyToAddress(cpk.PublicKey).Bytes())
	dob := common.NatureTime{time.Now().UnixNano(), "abc"}

	// build account
	account := accounts.Account{address, common.Signature{}, common.Signature{}, dob}
	accountHash := common.ToHash(account)
	// sign by father
	fSign, err := crypto.Sign(common.ToMD5(accountHash), fpk)
	if err != nil {
		t.Fatal(err)
	}

	// update account info
	account.FatherSign.SetBytes(fSign)
	accountHash = common.ToHash(account)
	// sign by mother
	mSign, err := crypto.Sign(common.ToMD5(accountHash), mpk)
	account.MotherSign.SetBytes(mSign)

	// add new account into genealogy
	err = clan.Add(account)
	if err != nil {
		t.Fatal(err)
	}
	// check the size of genealogy
	if clan.size != 3 {
		t.Fatal("clan size != 3")
	}
	// check the generation of genealogy
	if clan.generation != 1 {
		t.Fatal("clan generation != 1")
	}
}
