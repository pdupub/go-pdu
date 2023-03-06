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

package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pdupub/go-pdu/identity"
)

var (
	passwordPath string
	saltPath     string
	didCnt       int
)

// KeyCmd is used to create or unlock private key.
func KeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key [keystore_path]",
		Short: "Operations on keystores",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			pass, salt, err := getPassAndSalt()
			if err != nil {
				return err
			}
			dids, err := identity.CreateKeystoreArrayAndSaveLocal(args[0], pass, salt, didCnt)
			if err != nil {
				return err
			}
			for _, did := range dids {
				fmt.Println(did.GetKey().Address.Hex())
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&passwordPath, "pass", "", "path of password file")
	cmd.Flags().StringVar(&saltPath, "salt", "", "path of salt file")
	cmd.Flags().IntVar(&didCnt, "cnt", 3, "count of keystores")
	return cmd
}

func getPassAndSalt() (pass []byte, salt []byte, err error) {

	pass, err = ioutil.ReadFile(passwordPath)
	pass = []byte(strings.Replace(string(pass), "\n", "", -1))
	if err != nil {
		return
	}

	salt, err = ioutil.ReadFile(saltPath)
	salt = []byte(strings.Replace(string(salt), "\n", "", -1))
	if err != nil {
		return
	}

	return
}
