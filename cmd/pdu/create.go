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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/db/bolt"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
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

		if err := initUniverseAndSave(udb); err != nil {
			os.RemoveAll(dataDir)
			return err
		}

		fmt.Println("Create universe and space-time successfully")

		if err := udb.Close(); err != nil {
			return err
		}
		fmt.Println("Database closed successfully")
		return nil
	},
}

func initUniverseAndSave(udb db.UDB) error {
	users, priKeys, err := createRootUsers()
	if err != nil {
		return err
	}

	if err := saveRootUsers(users, udb); err != nil {
		return err
	}

	fmt.Println("Create root users successfully", users[0].Gender(), users[1].Gender())
	universe, err := core.NewUniverse(users[0], users[1])
	if err != nil {
		return err
	}

	if universe.GetUserByID(users[0].ID()).ID() != users[0].ID() || universe.GetUserByID(users[1].ID()).ID() != users[1].ID() {
		return errors.New("root users miss match")
	}

	msg, err := createFirstMsg(users, priKeys)
	if err != nil {
		return err
	}
	fmt.Println("First msg ID is ", common.Hash2String(msg.ID()))

	if err := universe.AddMsg(msg); err != nil {
		return err
	}

	if err := saveMsg(msg, udb); err != nil {
		return err
	}
	return nil
}

func createFirstMsg(users []*core.User, priKeys []*crypto.PrivateKey) (*core.Message, error) {
	var userSelected int
	var user *core.User
	var priKey *crypto.PrivateKey
	var content string
	fmt.Println("Please select the user to create first message: [0/1] ")
	fmt.Println("user 0: ", common.Hash2String(users[0].ID()))
	fmt.Println("user 1: ", common.Hash2String(users[1].ID()))
	fmt.Scan(&userSelected)
	if userSelected == 0 {
		user = users[0]
		priKey = priKeys[0]
	} else if userSelected == 1 {
		user = users[1]
		priKey = priKeys[1]
	}

	fmt.Println("Please input the content for the first msg")
	scanLine(&content)

	value := core.MsgValue{
		ContentType: core.TypeText,
		Content:     []byte(content),
	}
	msg, err := core.CreateMsg(user, &value, priKey)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func scanLine(input *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*input = string(data)
}

func saveMsg(msg *core.Message, udb db.UDB) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	countBytes, err := udb.Get(db.BucketConfig, db.ConfigMsgCount)
	if err != nil {
		return err
	}
	count := new(big.Int).SetBytes(countBytes)
	err = udb.Set(db.BucketMsg, common.Hash2String(msg.ID()), msgBytes)
	if err != nil {
		return err
	}

	err = udb.Set(db.BucketMID, count.String(), []byte(common.Hash2String(msg.ID())))
	if err != nil {
		return err
	}
	count = count.Add(count, big.NewInt(1))
	err = udb.Set(db.BucketConfig, db.ConfigMsgCount, count.Bytes())
	if err != nil {
		return err
	}

	err = udb.Set(db.BucketLastMID, common.Hash2String(msg.SenderID), []byte(common.Hash2String(msg.ID())))
	if err != nil {
		return err
	}
	return nil
}

func saveRootUsers(users []*core.User, udb db.UDB) (err error) {
	// save root users
	var root0, root1 []byte
	if root0, err = json.Marshal(users[0]); err != nil {
		return err
	}
	if err = udb.Set(db.BucketConfig, db.ConfigRoot0, root0); err != nil {
		return err
	}
	if err = udb.Set(db.BucketUser, common.Hash2String(users[0].ID()), root0); err != nil {
		return err
	}

	if root1, err = json.Marshal(users[1]); err != nil {
		return err
	}

	if err = udb.Set(db.BucketConfig, db.ConfigRoot1, root1); err != nil {
		return err
	}

	if err = udb.Set(db.BucketUser, common.Hash2String(users[1].ID()), root1); err != nil {
		return err
	}

	return nil
}

func createRootUsers() (users []*core.User, priKeys []*crypto.PrivateKey, err error) {

	for i := 0; i < 2; i++ {
		priKey, pubKey, err := unlockKeyByCmd()
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("Unlock root key successfully")
		for {
			var rootName, rootExtra, isSave string
			fmt.Print("Name: ")
			scanLine(&rootName)
			fmt.Print("Extra: ")
			scanLine(&rootExtra)
			user := core.CreateRootUser(*pubKey, rootName, rootExtra)
			fmt.Println("ID", common.Hash2String(user.ID()), "name", user.Name, "extra", user.DOBExtra, "gender", user.Gender())
			fmt.Print("save new user (yes/no): ")
			fmt.Scan(&isSave)
			if strings.ToUpper(isSave) == "YES" || strings.ToUpper(isSave) == "Y" {
				users = append(users, user)
				priKeys = append(priKeys, priKey)
				break
			}
		}
	}
	return users, priKeys, err
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
	if err := udb.Set(db.BucketConfig, db.ConfigMsgCount, big.NewInt(0).Bytes()); err != nil {
		return nil, err
	}

	if err := udb.CreateBucket(db.BucketUser); err != nil {
		return nil, err
	}
	if err := udb.CreateBucket(db.BucketMsg); err != nil {
		return nil, err
	}
	if err := udb.CreateBucket(db.BucketMID); err != nil {
		return nil, err
	}
	if err := udb.CreateBucket(db.BucketLastMID); err != nil {
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
