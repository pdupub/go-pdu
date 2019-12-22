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
	"errors"
	"fmt"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/bitcoin"
	"github.com/pdupub/go-pdu/crypto/ethereum"
	"github.com/pdupub/go-pdu/crypto/pdu"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/gopass"
)

const (
	operGenerate = "generate"
	operInspect  = "inspect"
)

var (
	errUnknownOperation = errors.New("unknown operation")
	errPasswordNotMatch = errors.New("password not match")
	errKeyFileMissing   = errors.New("key file missing")
)

var crypt, output, keyFile, sigType string
var msCount int

// accountCmd represents the create command
var accountCmd = &cobra.Command{
	Use:   "account [generate/inspect]",
	Short: "Account generate or inspect",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {

		var engine crypto.Engine
		switch strings.ToUpper(crypt) {
		case crypto.PDU:
			engine = pdu.New()
		case crypto.ETH:
			engine = ethereum.New()
		case crypto.BTC:
			engine = bitcoin.New()
		default:
			return crypto.ErrSourceNotMatch
		}

		switch strings.ToLower(args[0]) {
		case operGenerate:
			fmt.Printf("password: ")
			passwd, err := gopass.GetPasswd()
			if err != nil {
				return err
			}
			fmt.Printf("repeat password: ")
			passwd2, err := gopass.GetPasswd()
			if err != nil {
				return err
			}
			if string(passwd) != string(passwd2) {
				return errPasswordNotMatch
			}
			if _, err := os.Stat(keyFile); err == nil {
				return fmt.Errorf("keyfile already exists at %s", keyFile)
			} else if !os.IsNotExist(err) {
				return err
			}

			if strings.ToUpper(sigType) != crypto.Signature2PublicKey && strings.ToUpper(sigType) != crypto.MultipleSignatures {
				return crypto.ErrSigTypeNotSupport
			}
			privateKey, _, err := engine.GenKey(strings.ToUpper(sigType), msCount)
			if err != nil {
				return err
			}
			keyJson, err := engine.EncryptKey(privateKey, string(passwd))
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(output), 0700); err != nil {
				return fmt.Errorf("could not create directory %s", filepath.Dir(keyFile))
			}
			if err := ioutil.WriteFile(output, keyJson, 0600); err != nil {
				return fmt.Errorf("failed to write keyfile to %s: %v", keyFile, err)
			}
			fmt.Println(output, "is created success.")

		case operInspect:
			if keyFile == "" {
				return errKeyFileMissing
			}
			keyjson, err := ioutil.ReadFile(keyFile)
			if err != nil {
				return err
			}
			fmt.Printf("password: ")
			passwd, err := gopass.GetPasswd()
			if err != nil {
				return err
			}
			pk, err := engine.DecryptKey(keyjson, string(passwd))
			if err != nil {
				return err
			}
			fmt.Println(pk.PriKey)

		default:
			return errUnknownOperation
		}

		return nil
	},
}

func init() {

	accountCmd.PersistentFlags().StringVar(&sigType, "sigType", crypto.Signature2PublicKey, "sig type (S2PK/MS)")
	accountCmd.PersistentFlags().IntVar(&msCount, "msCount", 1, "count number of MS")
	accountCmd.PersistentFlags().StringVar(&keyFile, "key", "", "key file")
	accountCmd.PersistentFlags().StringVar(&crypt, "crypt", "PDU", "type of crypt (default is PDU)")
	accountCmd.PersistentFlags().StringVarP(&output, "output", "o", "key.json", "output file")
	rootCmd.AddCommand(accountCmd)
}
