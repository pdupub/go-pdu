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
	"os"

	"github.com/spf13/cobra"

	"github.com/pdupub/go-pdu/params"
)

var (
	projectPath string
	configPath  string

	firebaseKeyPath   string
	firebaseProjectID string
)

func main() {
	if err := initConfig(); err != nil {
		return
	}
	rootCmd := &cobra.Command{
		Use:   "pdu",
		Short: "PDU command line interface (" + params.Version + ")",
		Long: `ParaDigi Universe
	A decentralized social networking service
	Website: https://pdu.pub`,
	}

	rootCmd.AddCommand(RunCmd())
	rootCmd.AddCommand(MsgCmd())
	rootCmd.AddCommand(KeyCmd())
	rootCmd.AddCommand(NodeCmd())

	rootCmd.PersistentFlags().StringVar(&projectPath, "projectPath", "./", "project root path")
	rootCmd.PersistentFlags().StringVar(&firebaseKeyPath, "fbKeyPath", params.TestFirebaseAdminSDKPath, "path of firebase json key")
	rootCmd.PersistentFlags().StringVar(&firebaseProjectID, "fbProjectID", params.TestFirebaseProjectID, "project ID")

	rootCmd.Version = params.Version
	if err := rootCmd.Execute(); err != nil {
		return
	}
}

func initConfig() error {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath = home + "/.pdu/"
	_, err = os.Stat(configPath)
	if err != nil && !os.IsExist(err) {
		if err := os.Mkdir(configPath, os.ModePerm); err != nil { // perm 0666
			fmt.Println("create config fail", err)
			return err
		}
	}
	return nil
}
