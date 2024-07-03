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
	"fmt"
	"log"
	"os"

	"github.com/pdupub/go-pdu/node"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
)

func main() {
	var peerPort int
	var webPort int
	var rpcPort int
	var dbName string
	var nodeKey string
	var testMode bool

	var rootCmd = &cobra.Command{
		Use:   "pdu",
		Short: "PDU command line interface (" + params.Version + ")",
		Long: `ParaDigi Universe
	A Peer-to-Peer Social Network Service
	Website: https://pdu.pub`,

		Run: func(cmd *cobra.Command, args []string) {
			if testMode {
				peerPort = 4002
				webPort = 8556
				rpcPort = 8555
				dbName = "pdu_test.db"
				nodeKey = "node_test.key"

				// Check for any specified flags in test mode and warn the user
				if cmd.Flags().Changed("peerPort") || cmd.Flags().Changed("webPort") || cmd.Flags().Changed("rpcPort") || cmd.Flags().Changed("dbName") || cmd.Flags().Changed("nodeKey") {
					fmt.Println("Warning: Test mode will ignore the specified parameters")
				}
			}

			log.Println("P2P node running")

			n, err := node.NewNode(peerPort, nodeKey, dbName)
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
	rootCmd.Flags().StringVarP(&nodeKey, "nodeKey", "n", "node.key", "Node key")
	rootCmd.Flags().BoolVarP(&testMode, "test", "t", false, "Run in test mode")

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
