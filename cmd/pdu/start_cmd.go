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
	"github.com/pdupub/go-pdu/common/log"
	"github.com/spf13/cobra"
)

// StartCmd start the pdu server.
// At first step of our road map, this server is only for test.
// At later steps, the server will have different functions, like time proof server, broadcast server etc.
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [port]",
		Short: "Start PDU Server on ",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			port := args[0]
			log.Info("listen on", port)
			return nil
		},
	}

	return cmd
}
