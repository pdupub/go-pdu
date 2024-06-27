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
	var peerPort int
	var webPort int
	var rpcPort int
	var dbName string

	var rootCmd = &cobra.Command{
		Use:   "pdu",
		Short: "A decentralized P2P program",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("P2P node running")

			n, err := node.NewNode(peerPort, dbName)
			if err != nil {
				log.Fatalf("Failed to create node: %s", err)
				return
			}

			n.Run(webPort, rpcPort)
		},
	}

	rootCmd.Flags().IntVarP(&peerPort, "peerPort", "p", 4001, "Port to listen for P2P connections")
	rootCmd.Flags().IntVarP(&webPort, "webPort", "w", 8546, "Port for the web server")
	rootCmd.Flags().IntVarP(&rpcPort, "rpcPort", "r", 8545, "Port for the RPC server")
	rootCmd.Flags().StringVarP(&dbName, "dbName", "d", "pdu.db", "Database name")

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
