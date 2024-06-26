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

package main

import (
	"log"
	"os"

	"github.com/pdupub/go-pdu/node"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "pdu",
		Short: "A decentralized P2P program",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("P2P node running")
			listenPort := 4001
			webPort := 8546
			rpcPort := 8545
			dbName := "pdu.db"

			n, err := node.NewNode(listenPort, dbName)
			if err != nil {
				log.Fatalf("Failed to create node: %s", err)
				return
			}

			n.Run(webPort, rpcPort)
		},
	}

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Run test command",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Running test...")
			// test方法的代码
		},
	}

	rootCmd.AddCommand(testCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
