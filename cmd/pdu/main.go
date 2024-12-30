package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pdupub/go-pdu/internal/p2p"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pdu",
	Short: "PDU is a command line tool",
	Long:  `PDU is a command line tool with various functionalities.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to PDU!")
	},
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
		defer node.Close()

		fmt.Printf("P2P node started with ID: %s\n", node.Host.ID())
		select {}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
