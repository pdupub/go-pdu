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
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/db"
	"github.com/pdupub/go-pdu/db/bolt"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		var user0, user1 core.User
		root0, err := udb.Get(db.BucketConfig, db.ConfigRoot0)
		if err != nil {
			return err
		}
		root1, err := udb.Get(db.BucketConfig, db.ConfigRoot1)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(root0, &user0); err != nil {
			return err
		} else {
			log.Info(common.Hash2String(user0.ID()), user0.Gender())
		}
		if err := json.Unmarshal(root1, &user1); err != nil {
			return err
		} else {
			log.Info(common.Hash2String(user1.ID()), user1.Gender())
		}

		universe, err := core.NewUniverse(&user0, &user1)
		if err != nil {
			return err
		}

		log.Info("Create universe and space-time successfully", universe.GetSpaceTimeIDs())
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
	rootCmd.AddCommand(startCmd)
}
