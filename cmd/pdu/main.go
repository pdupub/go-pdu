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
	"os"

	"github.com/pdupub/go-pdu/node"
	"github.com/pdupub/go-pdu/udb"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "pdu",
		Short: "A decentralized P2P program",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("P2P node running")
			// 初始化并启动节点的代码
		},
	}

	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "Run test command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running test...")
			// test方法的代码
		},
	}

	var nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Run node command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("P2P node running")
			listenPort := 4001 // P2P 监听端口
			webPort := 8080    // Web 服务器端口
			rpcPort := 8545    // RPC 服务器端口
			node.Run(listenPort, webPort, rpcPort)
		},
	}

	var DBCmd = &cobra.Command{
		Use:   "db",
		Short: "Run db command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running db...")
			dbName := "udb.db"
			db, err := udb.InitDB(dbName)
			defer db.CloseDB()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// db方法的代码
		},
	}

	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(DBCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
