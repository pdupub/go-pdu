// Copyright 2019 The PDU Authors
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
	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/db/bolt"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PDU Universe",
	RunE: func(_ *cobra.Command, args []string) error {
		udb, err := initDB(dataDir)
		if err != nil {
			return err
		}
		fmt.Println("Database initialized successfully", dataDir)

		priKey1, priKey2, err := unlockRootsKey()
		if err != nil {
			os.RemoveAll(dataDir)
			return err
		}
		fmt.Println("Unlock root key", priKey1, priKey2)
		fmt.Println("Create root users")

		fmt.Println("Create universe and space-time")

		if err := udb.Close(); err != nil {
			return err
		}
		fmt.Println("Database closed successfully")
		return nil
	},
}

func unlockRootsKey() (priKey1 *crypto.PrivateKey, priKey2 *crypto.PrivateKey, err error) {
	priKey1, err = unlockKey()
	if err != nil {
		return nil, nil, err
	}
	priKey2, err = unlockKey()
	return priKey1, priKey2, err
}

func unlockKey() (*crypto.PrivateKey, error) {
	var keyFile string
	fmt.Print("keyfile path: ")
	fmt.Scan(&keyFile)
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	fmt.Print("password: ")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}

	return core.DecryptKey(keyJson, string(passwd))
}

func initDB(dataDir string) (db.UDB, error) {
	if dataDir == "" {
		home, _ := homedir.Dir()
		dataDir = path.Join(home, params.DefaultPath)
	}
	err := os.Mkdir(dataDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	dbFilePath := path.Join(dataDir, "u.db")
	udb, err := bolt.NewDB(dbFilePath)
	if err != nil {
		return nil, err
	}

	bucketName := []byte("universe")
	if err := udb.CreateBucket(bucketName); err != nil {
		return nil, err
	}
	return udb, nil
}

func init() {
	rootCmd.AddCommand(createCmd)
}
