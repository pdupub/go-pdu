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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func main() {

	viper.New()
	rootCmd := &cobra.Command{
		Use:   "pdu",
		Short: "PDU command line interface",
	}
	rootCmd.AddCommand(InitializeCmd())

	rootCmd.Execute()
}

func InitializeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [generation num]",
		Short: "Initialize the n generations of PDU",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			genNum := args[0]
			log.Println("generations", genNum)

			return nil
		},
	}

	return cmd
}
