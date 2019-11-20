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
	"github.com/pdupub/go-pdu/common/log"
	"github.com/spf13/cobra"
)

var (
	errUserNotExist       = errors.New("user not exist in this system")
	errCreateRootUserFail = errors.New("create root user fail")
)

// TestCmd do the test related to step of create the pdu universe and
// split the universe for different space-time
func TestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test on pdu",
		RunE: func(_ *cobra.Command, args []string) error {
			log.Info("test")
			return nil
		},
	}

	return cmd
}
