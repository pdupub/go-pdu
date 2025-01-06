package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pdupub/go-pdu/internal/p2p"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pdu",
	Short: "PDU is a command line tool",
	Long:  `PDU is a command line tool with P2P functionality.`,
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(rpcCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the P2P node",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		node, err := p2p.NewNode(ctx)
		if err != nil {
			fmt.Printf("Failed to create node: %v\n", err)
			os.Exit(1)
		}

		defer func() {
			fmt.Println("Shutting down node...")
			if err := node.Close(); err != nil {
				fmt.Printf("Error closing node: %v\n", err)
			}
			fmt.Println("Node shutdown complete")
		}()

		fmt.Printf("P2P node started with ID: %s\n", node.Host.ID())
		// fmt.Printf("Listening on addresses:\n")
		// for _, addr := range node.Host.Addrs() {
		// 	fmt.Printf("  - %s/p2p/%s\n", addr, node.Host.ID())
		// }

		// 显示本地连接地址
		localAddr := node.GetLocalAddress()
		fmt.Printf("Local address for connections: %s\n", localAddr)
		// fmt.Println("\nUse this address in another terminal with:")
		// fmt.Printf("  pdu connect %s\n", localAddr)

		// 等待中断信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	},
}

var rpcCmd = &cobra.Command{
	Use:   "rpc [peerID] [msg]",
	Short: "Connect to RPC server",
	Long:  `Connect to a remote RPC server using the provided address.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.DialHTTP("http://127.0.0.1:8545")
		if err != nil {
			fmt.Printf("Failed to connect to RPC server: %v", err)
		}
		defer client.Close()

		// 2) 调用远程方法 "math_add"
		//    go-ethereum/rpc 调用形式为: client.Call(&result, "serviceName_methodName", param1, param2, ...)
		var result string
		err = client.Call(&result, "pdu_chat", args[1])
		if err != nil {
			fmt.Printf("RPC Chat call error: %v", err)
		}

		fmt.Printf("RPC chat result: %s\n", result)

		err = client.Call(&result, "pdu_message", args[0], args[1])
		if err != nil {
			fmt.Printf("RPC Message call error: %v", err)
		}

		fmt.Printf("RPC Message result: %s\n", result)

	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
