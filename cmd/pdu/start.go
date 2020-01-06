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
	"github.com/mitchellh/go-homedir"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start to run PDU Universe",
	//SilenceUsage:  true,
	//SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		err := initLoad()
		if err != nil {
			return err
		}

		log.Info("Starting p2p node")
		log.Info("CONFIG_NAME", viper.GetString("CONFIG_NAME"))
		return nil
	},
}

// initLoad reads in config file and ENV variables if set.
func initLoad() error {
	if dataDir != "" {
		// Use config file from the flag.
		viper.SetConfigFile(path.Join(dataDir, params.DefaultConfigFile))
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		// Search config in $HOME/.pdu directory with name "config" (without extension).
		dataDir = path.Join(home, params.DefaultPath)
		viper.AddConfigPath(dataDir)
		viper.SetConfigName(params.DefaultConfigFile)
	}
	viper.SetConfigType("yml")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func init() {
	startCmd.PersistentFlags().StringVar(&dataDir, "datadir", "", fmt.Sprintf("(default $HOME/%s)", params.DefaultPath))
	rootCmd.AddCommand(startCmd)
}
