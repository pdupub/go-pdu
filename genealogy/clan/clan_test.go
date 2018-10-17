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

	cpk, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	var address, fad, mad common.Address
	address.SetBytes(crypto.PubkeyToAddress(cpk.PublicKey).Bytes())

	dob := common.NatureTime{123, "abc"}
	account := accounts.Account{address, common.Signature{}, common.Signature{}, dob}
	accountHash := common.ToHash(account)

	fpk, err := crypto.GenerateKey()
	fad.SetBytes(crypto.PubkeyToAddress(fpk.PublicKey).Bytes())
	if err != nil {
		t.Fatal(err)
	}
	fSign, err := crypto.Sign(common.ToMD5(accountHash), fpk)
	if err != nil {
		t.Fatal(err)
	}
	account.FatherSign.SetBytes(fSign)

	mpk, err := crypto.GenerateKey()
	mad.SetBytes(crypto.PubkeyToAddress(mpk.PublicKey).Bytes())
	if err != nil {
		t.Fatal(err)
	}
	accountHash = common.ToHash(account)
	mSign, err := crypto.Sign(common.ToMD5(accountHash), mpk)
	account.MotherSign.SetBytes(mSign)

	fatherAccount := accounts.Account{fad, common.Signature{}, common.Signature{}, common.NatureTime{123, "a"}}
	motherAccount := accounts.Account{mad, common.Signature{}, common.Signature{}, common.NatureTime{234, "b"}}
	clan, err := New(fatherAccount, motherAccount)
	if err != nil {
		t.Fatal(err)
	}

	err = clan.Add(account)
	if err != nil {
		t.Fatal(err)
	}
}
