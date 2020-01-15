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

package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"io/ioutil"
	"strings"
)

func unlockKeyByCmd() (*crypto.PrivateKey, *crypto.PublicKey, error) {
	var keyFile string
	fmt.Print("KeyFile path: ")
	fmt.Scan(&keyFile)
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	fmt.Print("Password: ")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, nil, err
	}

	return core.DecryptKey(keyJson, string(passwd))
}

func unlockKeyByFile(keyFile, passFile string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	passwd, err := ioutil.ReadFile(passFile)
	if err != nil {
		return nil, nil, err
	}

	return core.DecryptKey(keyJson, strings.TrimSpace(string(passwd)))
}
