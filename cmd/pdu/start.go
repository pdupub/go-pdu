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
	"github.com/mitchellh/go-homedir"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/db/bolt"
	"github.com/pdupub/go-pdu/node"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start to run PDU Universe",
	RunE: func(_ *cobra.Command, args []string) error {
		err := initConfigLoad()
		if err != nil {
			return err
		}

		log.Info("Starting p2p node")
		log.Info("CONFIG_NAME", viper.GetString("CONFIG_NAME"))

		udb, err := initDBLoad()
		if err != nil {
			return err
		}

		pn, err := node.New(udb)
		if err != nil {
			return err
		}

		// for all node mode need to unlock account
		var unlockedUser core.User
		var unlockedPrivateKey *crypto.PrivateKey
		if nodeTPEnable {
			var unlockedPublicKey *crypto.PublicKey
			unlockedPrivateKey, unlockedPublicKey, err = unlockKeyByFile(unlockKeyFile, unlockPassFile)
			if err != nil {
				return err
			}

			if len(unlockUserIDPrefix) < 5 {
				return errors.New("user ID not have enough prefix")
			}

			rows, err := udb.Find(db.BucketUser, unlockUserIDPrefix, 2)
			if err != nil {
				return err
			}
			if len(rows) == 0 {
				return errors.New("user ID can not be found")
			}
			if len(rows) > 1 {
				return errors.New("user ID which has this prefix are not unique")
			}

			// check public key match
			unlockUserBytes, err := udb.Get(db.BucketUser, rows[0].K)
			if err != nil {
				return err
			}
			err = json.Unmarshal(unlockUserBytes, &unlockedUser)
			if err != nil {
				return err
			}

			p1, err := json.Marshal(unlockedUser.Auth.PubKey)
			if err != nil {
				return err
			}
			p2, err := json.Marshal(unlockedPublicKey.PubKey)
			if err != nil {
				return err
			}
			if unlockedUser.Auth.Source != unlockedPublicKey.Source ||
				unlockedUser.Auth.SigType != unlockedPublicKey.SigType ||
				common.Bytes2String(p1) != common.Bytes2String(p2) {
				return errors.New("public key not match")
			}

			log.Info("Account unlocked success", rows[0].K)
		}

		if nodeTPEnable {
			if err := pn.EnableTP(&unlockedUser, unlockedPrivateKey, nodeTPInterval); err != nil {
				return err
			}
		}

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill)
		pn.SetLocalPort(localPort)
		if nodeAddressList != "" {
			pn.SetNodes(nodeAddressList)
		}
		pn.Run(c)

		return nil
	},
}

// initConfigLoad reads in config file and ENV variables if set.
func initConfigLoad() error {
	if dataDir == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		dataDir = path.Join(home, params.DefaultPath)

	}

	viper.SetConfigFile(path.Join(dataDir, params.DefaultConfigFile))
	viper.SetConfigType("yml")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func initDBLoad() (db.UDB, error) {
	dbFilePath := path.Join(dataDir, "u.db")
	udb, err := bolt.NewDB(dbFilePath)
	if err != nil {
		return nil, err
	}
	return udb, nil
}

func init() {
	startCmd.PersistentFlags().StringVar(&dataDir, "datadir", "", fmt.Sprintf("(default $HOME/%s)", params.DefaultPath))
	startCmd.PersistentFlags().StringVar(&nodeAddressList, "nodes", "", "pdu nodes list, split by comma [userid@ip:port/nodeKey]")
	startCmd.PersistentFlags().Uint64Var(&localPort, "localPort", node.DefaultLocalPort, "local port")

	// time proof
	startCmd.PersistentFlags().BoolVar(&nodeTPEnable, "tp", false, "time proof enable")
	startCmd.PersistentFlags().Uint64Var(&nodeTPInterval, "tpInterval", node.DefaultTimeProofInterval, "time proof interval")

	// unlock account
	startCmd.PersistentFlags().StringVar(&unlockUserIDPrefix, "user", "", "user ID prefix")
	startCmd.PersistentFlags().StringVar(&unlockKeyFile, "key", "", "key file")
	startCmd.PersistentFlags().StringVar(&unlockPassFile, "pass", "", "pass file")

	rootCmd.AddCommand(startCmd)
}
