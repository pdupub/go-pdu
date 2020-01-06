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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/db/bolt"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PDU Universe",
	RunE: func(_ *cobra.Command, args []string) error {

		if err := initDir(); err != nil {
			return err
		}

		if err := initConfig(); err != nil {
			return err
		}

		udb, err := initDB()
		if err != nil {
			return err
		}
		fmt.Println("Database initialized successfully", dataDir)

		pubKeys, err := unlockRootsKeys(2)
		if err != nil {
			os.RemoveAll(dataDir)
			return err
		}
		fmt.Println("Unlock root key successfully")

		users, err := createRootUsers(pubKeys)
		if err != nil {
			os.RemoveAll(dataDir)
			return err
		}

		if err := saveRootUsers(users, udb); err != nil {
			os.RemoveAll(dataDir)
			return err
		}

		fmt.Println("Create root users successfully", users[0].Gender(), users[1].Gender())
		universe, err := core.NewUniverse(users[0], users[1])
		if err != nil {
			os.RemoveAll(dataDir)
			return err
		}

		if universe.GetUserByID(users[0].ID()).ID() != users[0].ID() || universe.GetUserByID(users[1].ID()).ID() != users[1].ID() {
			os.RemoveAll(dataDir)
			return errors.New("root users miss match")
		}
		fmt.Println("Create universe and space-time successfully")

		if err := udb.Close(); err != nil {
			return err
		}
		fmt.Println("Database closed successfully")
		return nil
	},
}

func saveRootUsers(users []*core.User, udb db.UDB) error {
	// save root users
	if root0, err := json.Marshal(users[0]); err != nil {
		return err
	} else if err := udb.Set(db.BucketConfig, db.ConfigRoot0, root0); err != nil {
		return err
	}

	if root1, err := json.Marshal(users[1]); err != nil {
		return err
	} else if err := udb.Set(db.BucketConfig, db.ConfigRoot1, root1); err != nil {
		return err
	}
	return nil
}

func createRootUsers(pubKeys []*crypto.PublicKey) (users []*core.User, err error) {
	for _, pubKey := range pubKeys {
		for {
			var rootName, rootExtra, isSave string
			fmt.Print("name: ")
			fmt.Scan(&rootName)
			fmt.Print("extra: ")
			fmt.Scan(&rootExtra)
			user := core.CreateRootUser(*pubKey, rootName, rootExtra)
			fmt.Println("ID", common.Hash2String(user.ID()), "name", user.Name, "extra", user.DOBExtra, "gender", user.Gender())
			fmt.Print("save new user (yes/no): ")
			fmt.Scan(&isSave)
			if strings.ToUpper(isSave) == "YES" || strings.ToUpper(isSave) == "Y" {
				users = append(users, user)
				break
			}
		}
	}
	return users, err
}

func unlockRootsKeys(cnt int) (pubKeys []*crypto.PublicKey, err error) {
	for i := 0; i < cnt; i++ {
		_, pubKey, err := unlockKey()
		if err != nil {
			return pubKeys, err
		}
		pubKeys = append(pubKeys, pubKey)
	}
	return pubKeys, err
}

func unlockKey() (*crypto.PrivateKey, *crypto.PublicKey, error) {
	var keyFile string
	fmt.Print("keyfile path: ")
	fmt.Scan(&keyFile)
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	fmt.Print("password: ")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, nil, err
	}

	return core.DecryptKey(keyJson, string(passwd))
}

func initDB() (db.UDB, error) {
	dbFilePath := path.Join(dataDir, "u.db")
	udb, err := bolt.NewDB(dbFilePath)
	if err != nil {
		return nil, err
	}

	if err := udb.CreateBucket(db.BucketConfig); err != nil {
		return nil, err
	}
	if err := udb.CreateBucket(db.BucketUser); err != nil {
		return nil, err
	}
	if err := udb.CreateBucket(db.BucketMsg); err != nil {
		return nil, err
	}
	return udb, nil
}

func initConfig() error {
	viper.SetConfigType(params.DefaultConfigType)
	viper.Set("CONFIG_NAME", "PDU")
	return viper.WriteConfigAs(path.Join(dataDir, params.DefaultConfigFile))
}

func initDir() error {
	if dataDir == "" {
		home, _ := homedir.Dir()
		dataDir = path.Join(home, params.DefaultPath)
	}
	err := os.Mkdir(dataDir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	createCmd.PersistentFlags().StringVar(&dataDir, "datadir", "", fmt.Sprintf("(default $HOME/%s)", params.DefaultPath))
	rootCmd.AddCommand(createCmd)
}
