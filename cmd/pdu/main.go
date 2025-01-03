package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
