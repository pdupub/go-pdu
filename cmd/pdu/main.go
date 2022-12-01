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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/node"
	"github.com/pdupub/go-pdu/params"
)

var (
	passwordPath string
	projectPath  string
	saltPath     string
	references   string
	message      string
	keyIndex     int
)

func main() {
	viper.New()
	rootCmd := &cobra.Command{
		Use:   "pdu",
		Short: "PDU command line interface (" + params.Version + ")",
		Long: `ParaDigi Universe
	A decentralized social networking service
	Website: https://pdu.pub`,
	}

	rootCmd.AddCommand(TestCmd())
	rootCmd.AddCommand(StartCmd())
	rootCmd.AddCommand(SendMsgCmd())
	rootCmd.AddCommand(CreateKeystoreCmd())
	rootCmd.PersistentFlags().StringVar(&projectPath, "projectPath", "./", "project root path")

	if err := rootCmd.Execute(); err != nil {
		return
	}
}

// TestCmd run some test functions, just for developers
func TestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test some methods",
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Println("testing")

			return nil
		},
	}
	return cmd
}

// StartCmd start run the node
func StartCmd() *cobra.Command {
	var firebaseKeyPath string
	var firebaseProjectID string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start run node",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, os.Kill)
			interval, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			if serv, err := node.New(interval, firebaseKeyPath, firebaseProjectID); err != nil {
				return err
			} else {
				serv.Run(c)
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&firebaseKeyPath, "fbKeyPath", "./udb/fb/test-firebase-adminsdk.json", "path of firebase json key")
	cmd.Flags().StringVar(&firebaseProjectID, "fbProjectID", "pdupub-a2bdd", "project ID")

	return cmd
}

// SendMsgCmd is send msg to node
func SendMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send sampel msg to node",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			if keyIndex < 0 || keyIndex > 99 {
				return errors.New("index out of range")
			}
			testKeyfile := projectPath + params.TestKeystore(keyIndex)

			did, _ := identity.New()
			if err := did.UnlockWallet(testKeyfile, params.TestPassword); err != nil {
				return err
			}

			fmt.Println("msg author\t", did.GetKey().Address.Hex())

			// TODO: write msg to firebase firestore as client

			return nil
		},
	}
	cmd.Flags().IntVar(&keyIndex, "key", 0, "index of key used")
	cmd.Flags().StringVar(&message, "msg", "Hello World!", "content of msg send")
	cmd.Flags().StringVar(&references, "refs", "", "references split by comma")
	return cmd
}

func CreateKeystoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ck [keystore_path]",
		Short: "Create keystores",
		Long:  "Create keystores by same pass & salt",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			pass, salt, err := getPassAndSalt()
			if err != nil {
				return err
			}
			dids, err := identity.CreateKeystoreArrayAndSaveLocal(args[0], pass, salt, 3)
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
