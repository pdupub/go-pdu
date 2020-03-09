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
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/mitchellh/go-homedir"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/utils"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/db/bolt"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/viper"
)

// initNodeDir initialize node dir and config file and db, and open db
func initNodeDir() (db.UDB, error) {
	if err := initDir(); err != nil {
		return nil, err
	}

	if err := initConfig(); err != nil {
		return nil, err
	}

	udb, err := initDB()
	if err != nil {
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
	if err := udb.CreateBucket(db.BucketMOD); err != nil {
		return nil, err
	}
	if err := udb.CreateBucket(db.BucketLastMID); err != nil {
		return nil, err
	}
	if err := udb.CreateBucket(db.BucketPeer); err != nil {
		return nil, err
	}
	if err := udb.Set(db.BucketConfig, db.ConfigCurrentStep, big.NewInt(db.StepInitDB).Bytes()); err != nil {
		return nil, err
	}
	return udb, nil
}

func unlockKeyByCmd() (*crypto.PrivateKey, *crypto.PublicKey, error) {
	var keyFile string
	fmt.Print("KeyFile path: ")
	fmt.Scan(&keyFile)
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}
	fmt.Print("Password: ")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, nil, err
	}

	return utils.DecryptKey(keyJson, string(passwd))
}

func unlockKeyByFile(keyFile, passFile string) (*crypto.PrivateKey, *crypto.PublicKey, error) {
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	passwd, err := ioutil.ReadFile(passFile)
	if err != nil {
		return nil, nil, err
	}

	return utils.DecryptKey(keyJson, strings.TrimSpace(string(passwd)))
}

func scanLine(input *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*input = string(data)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
