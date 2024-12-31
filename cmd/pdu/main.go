package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	rootCmd.AddCommand(connectCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the P2P node",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		node, err := p2p.NewNode(ctx, "pdu", "1.0.0")
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
var connectCmd = &cobra.Command{
	Use:   "connect [address]",
	Short: "Connect to a peer using its address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// 创建节点
		node, err := p2p.NewNode(ctx, "pdu", "1.0.0")
		if err != nil {
			fmt.Printf("Failed to create node: %v\n", err)
			os.Exit(1)
		}
		defer node.Close()

		// 连接到指定地址
		address := args[0]
		if err := node.ConnectToPeer(address); err != nil {
			fmt.Printf("Failed to connect: %v\n", err)
			os.Exit(1)
		}

		// 创建一个channel用于处理中断信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// 创建一个channel用于通知主循环退出
		done := make(chan bool)

		// 在goroutine中处理信号
		go func() {
			<-sigChan
			fmt.Println("Disconnecting...")
			done <- true
		}()

		// 交互式命令行
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Connection established. Type 'q' or 'quit' to exit.")

		for {
			select {
			case <-done:
				return
			default:
				fmt.Print("> ")
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					continue
				}

				// 去除输入末尾的换行符
				input = strings.TrimSpace(input)

				// 检查是否退出
				if input == "q" || input == "quit" {
					return
				}

				// 回显输入的内容
				// fmt.Printf("You typed: %s\n", input)
			}
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
