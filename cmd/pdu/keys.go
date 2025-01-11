package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
)

var createKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create ETH keystore file",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 指定 keystore 文件保存路径
		keystoreDir := "./keystore"
		if _, err := os.Stat(keystoreDir); os.IsNotExist(err) {
			err := os.Mkdir(keystoreDir, 0700)
			if err != nil {
				log.Fatalf("Failed to create keystore directory: %v", err)
			}
		}

		// 2. 从命令行输入密码
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter password for new account: ")
		password, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read password: %v", err)
		}
		password = strings.TrimSpace(password)

		if len(password) == 0 {
			log.Fatal("Password cannot be empty")
		}

		// 3. 创建一个新的 keystore 管理器
		ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

		// 4. 生成一个新的以太坊账户
		account, err := ks.NewAccount(password)
		if err != nil {
			log.Fatalf("Failed to create new account: %v", err)
		}

		// 5. 重命名 keystore 文件为地址名
		newFileName := account.Address.Hex()
		newFilePath := fmt.Sprintf("%s/%s", keystoreDir, newFileName)
		if err := os.Rename(account.URL.Path, newFilePath); err != nil {
			log.Fatalf("Failed to rename keystore file: %v", err)
		}

		// 6. 输出生成的账户地址和 keystore 文件位置
		fmt.Printf("New account created: %s\n", account.Address.Hex())
		fmt.Printf("Keystore file saved as: %s\n", newFileName)
	},
}
