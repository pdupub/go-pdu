package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/pdupub/go-pdu/internal/p2p"
	"github.com/spf13/cobra"
)

var (
	rpcEnable bool // 是否开启 RPC 服务
	rpcPort   int  // 添加 RPC 端口变量

	dbPath string // 数据库文件地址

)

var rootCmd = &cobra.Command{
	Use:   "pdu",
	Short: "PDU is a command line tool",
	Long:  `PDU is a command line tool with P2P functionality.`,
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(rpcCmd)
	rootCmd.AddCommand(createKeyCmd)
	rootCmd.AddCommand(listKeysCmd)

	startCmd.Flags().BoolVar(&rpcEnable, "rpc", false, "Enable RPC ")
	startCmd.Flags().IntVarP(&rpcPort, "rpcport", "p", 8545, "RPC server port")
	startCmd.Flags().StringVar(&dbPath, "dbpath", "pdu.db", "Path of local database")
	rpcCmd.Flags().IntVarP(&rpcPort, "rpcport", "p", 8545, "RPC server port")

}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the P2P node",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		node, err := p2p.NewNode(ctx, dbPath)
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

		if rpcEnable {
			if err = node.StartRPC(rpcPort); err != nil {
				fmt.Println("RPC open fail")
			}
		}

		// 等待中断信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
	},
}

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Start interactive RPC command session",
	Long:  `Enter an interactive RPC command session with a remote RPC server.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := fmt.Sprintf("http://127.0.0.1:%d", rpcPort)

		client, err := rpc.DialHTTP(addr)
		if err != nil {
			log.Fatalf("Failed to connect to RPC server at %s: %v", addr, err)
		}
		defer client.Close()

		fmt.Println("Connected to RPC server.")
		fmt.Println("Type 'quit' or 'q' to exit.")
		fmt.Println("Type your RPC method and arguments in the format: method arg1 arg2 ...")

		// 使用 bufio.NewReader 读取用户输入
		reader := bufio.NewReader(os.Stdin)

		for {
			// 显示提示符
			fmt.Print("> ")

			// 读取整行输入
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				continue
			}

			// 去除行末的换行符
			input = strings.TrimSpace(input)

			// 检查退出条件
			if input == "quit" || input == "q" {
				fmt.Println("Exiting RPC session.")
				break
			}

			// 拆分用户输入为方法名和参数
			parts := strings.Fields(input)
			if len(parts) < 1 {
				fmt.Println("Invalid input. Please provide a method name and arguments.")
				continue
			}

			method := "pdu_" + parts[0]
			args := parts[1:]

			// 将 []string 转换为 []interface{}
			rpcArgs := make([]interface{}, len(args))
			for i, arg := range args {
				rpcArgs[i] = arg
			}

			// 调用 RPC 方法
			if parts[0] == "list" {

				var result []string
				err = client.Call(&result, method, rpcArgs...)

				if err != nil {
					fmt.Printf("RPC call error: %v\n", err)
				} else {
					fmt.Printf("RPC result: %s\n", result)
				}
			} else {
				var result string
				err = client.Call(&result, method, rpcArgs...)

				if err != nil {
					fmt.Printf("RPC call error: %v\n", err)
				} else {
					fmt.Printf("RPC result: %s\n", result)
				}
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
