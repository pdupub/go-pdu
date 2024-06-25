// Copyright 2024 The PDU Authors
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

package identity

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/google/uuid"
)

type Address common.Address

func (addr Address) Hex() string {
	return common.Address(addr).Hex()
}

type DID struct {
	key *keystore.Key
}

func New() (*DID, error) {
	return &DID{}, nil
}

func BytesToAddress(bytes []byte) Address {
	return Address(common.BytesToAddress(bytes))
}

// UnlockWallet used to unlock wallet
func (d *DID) UnlockWallet(keyFilePath, password string) error {
	keyJSON, err := os.ReadFile(keyFilePath)
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return err
	}
	d.key = key
	return nil
}

func (d *DID) GetKey() *keystore.Key {
	return d.key
}

func (d *DID) GetAddress() Address {
	return Address(d.key.Address)
}

func (d *DID) Inspect(showPrivate bool) (addr string, pubKey string, privKey string, err error) {
	addr = d.key.Address.Hex()
	pubKey = hex.EncodeToString(crypto.FromECDSAPub(&d.key.PrivateKey.PublicKey))
	if showPrivate {
		privKey = hex.EncodeToString(crypto.FromECDSA(d.key.PrivateKey))
	}
	return addr, pubKey, privKey, nil
}

func (d *DID) LoadECDSA(privateKeyHex string) error {
	key, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return err
	}
	privateKey, err := crypto.ToECDSA(key)
	if err != nil {
		return err
	}

	currentAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	UUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	d.key = &keystore.Key{
		Id:         UUID,
		Address:    currentAddress,
		PrivateKey: privateKey,
	}
	return nil
}

func (d *DID) Sign(b []byte) ([]byte, error) {
	if d.key == nil {
		return nil, fmt.Errorf("wallet is not unlocked")
	}
	hash := crypto.Keccak256Hash(b)
	return crypto.Sign(hash.Bytes(), d.GetKey().PrivateKey)
}
