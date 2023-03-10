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
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pdupub/go-pdu/node"
)

// RunCmd start to run the node daemon
func RunCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "run [loop interval]",
		Short: "Start to run node daemon",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, os.Kill)
			interval, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			if serv, err := node.New(interval, firebaseKeyPath, firebaseProjectID); err != nil {
				return err
			} else {
				serv.Run(c)
			}

			return nil
		},
	}

	return cmd
}
