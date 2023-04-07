// Copyright 2021 The PDU Authors
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
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/pdupub/go-pdu/node"
)

// RunCmd start to run the node daemon
func RunCmd() *cobra.Command {
	var interval int64
	var echoMode bool
	var port int64

	cmd := &cobra.Command{
		Use:   "run [loop interval]",
		Short: "Run node daemon",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {

			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, os.Kill)

			if serv, err := node.New(interval, firebaseKeyPath, firebaseProjectID); err != nil {
				return err
			} else {
				if echoMode {
					serv.RunEcho(port, c)
				} else {
					serv.Run(c)
				}
			}

			return nil
		},
	}

	cmd.PersistentFlags().Int64Var(&interval, "interval", 5, "time interval between consecutive processing on node")
	cmd.PersistentFlags().BoolVar(&echoMode, "echo", false, "only exe on received request")
	cmd.PersistentFlags().Int64Var(&port, "port", 8123, "http server started on")

	return cmd
}
