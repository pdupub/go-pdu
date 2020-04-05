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

	"github.com/pdupub/go-pdu/console"
	"github.com/spf13/cobra"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Console of pdu",
	RunE: func(_ *cobra.Command, args []string) error {
		var content string
		consoleInfoDisplay()
		cls, err := console.NewConsole()
		if err != nil {
			return err
		}
		for {
			fmt.Print("> ")
			scanLine(&content)
			if content == "quit" || content == "q" {
				cls.Close()
				break
			} else if content == "show" {
				cls.ExeShowCmd()
			} else {
				fmt.Println(content)
			}
		}

		return nil
	},
}

func consoleInfoDisplay() {
	fmt.Print(`
##############################################################
#                                                            #
#               Welcome to PDU console ;-)                   #
#                                                            #
##############################################################

`)
}

func init() {
	rootCmd.AddCommand(consoleCmd)
}
