package main

import (
	"fmt"
	"os"

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
	// 这里可以添加命令行参数和子命令
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
