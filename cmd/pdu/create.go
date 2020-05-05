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
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PDU Universe",
	RunE: func(_ *cobra.Command, args []string) error {
		udb, err := initNodeDir()
		if err != nil {
			return err
		}
		fmt.Println("Database initialized successfully", dataDir)

		if err := initUniverseAndSave(udb); err != nil {
			os.RemoveAll(dataDir)
			return err
		}
		if err := addUniverseSettings(udb); err != nil {
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

func addUniverseSettings(udb db.UDB) (err error) {
	var dimension, perimeter, redshift int64

	for {
		var dimensionInput string
		fmt.Printf("Universe Dimension (%d): ", core.DefaultDimensionNum)
		scanLine(&dimensionInput)
		if dimensionInput != "" {
			dimension, err = strconv.ParseInt(dimensionInput, 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if dimension <= 0 || dimension > core.MaxDimensionNum {
				fmt.Println(core.ErrDimensionNumberNotSuitable)
				continue
			}
		} else {
			dimension = core.DefaultDimensionNum
		}
		break
	}
	if err := udb.Set(db.BucketConfig, db.ConfigUniverseDimension, big.NewInt(dimension).Bytes()); err != nil {
		return err
	}

	for {
		var perimeterInput string
		fmt.Printf("Universe Perimeter (%e): ", core.DefaultPerimeter)
		scanLine(&perimeterInput)
		if perimeterInput != "" {
			perimeter, err = strconv.ParseInt(perimeterInput, 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if perimeter == 0 {
				fmt.Println(core.ErrPerimeterIsZero)
				continue
			}
			if perimeter < 0 {
				perimeter = -perimeter
			}
		} else {
			perimeter = core.DefaultPerimeter
		}
		break
	}
	if err := udb.Set(db.BucketConfig, db.ConfigUniversePerimeter, big.NewInt(perimeter).Bytes()); err != nil {
		return err
	}

	for {
		var redshiftInput string
		var defaultRedshift int64
		defaultRedshift = perimeter / 1e+4
		fmt.Printf("Universe Red-shift (%d): ", defaultRedshift)
		scanLine(&redshiftInput)
		if redshiftInput != "" {
			redshift, err = strconv.ParseInt(redshiftInput, 10, 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			redshift = defaultRedshift
		}
		break
	}
	if err := udb.Set(db.BucketConfig, db.ConfigUniverseRedshiftConstant, big.NewInt(redshift).Bytes()); err != nil {
		return err
	}

	return nil
}

func initUniverseAndSave(udb db.UDB) error {
	// create root users
	users, priKeys, err := createRootUsers()
	if err != nil {
		return err
	}

	if err := db.SaveRootUsers(udb, users); err != nil {
		return err
	}

	fmt.Println("Create root users successfully", users[0].Gender(), users[1].Gender())

	// create universe by root users
	universe, err := core.NewUniverse(users[0], users[1])
	if err != nil {
		return err
	}

	if universe.GetUserByID(users[0].ID()).ID() != users[0].ID() || universe.GetUserByID(users[1].ID()).ID() != users[1].ID() {
		return errors.New("root users miss match")
	}

	// create first msg
	msg, err := createFirstMsg(users, priKeys)
	if err != nil {
		return err
	}
	fmt.Println("First msg ID is ", common.Hash2String(msg.ID()))

	if err := universe.AddMsg(msg); err != nil {
		return err
	}

	if err := db.SaveMsg(udb, msg); err != nil {
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
			fmt.Println("ID", common.Hash2String(user.ID()), "name", user.Name, "extra", user.BirthExtra, "gender", user.Gender())
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

func init() {
	createCmd.PersistentFlags().StringVar(&dataDir, "datadir", "", fmt.Sprintf("(default $HOME/%s)", params.DefaultPath))
	rootCmd.AddCommand(createCmd)
}
