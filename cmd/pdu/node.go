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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pdupub/go-pdu/params"
)

// NodeCmd manually perform some actions on the node,
// like set address level or remove processed message from the node.
func NodeCmd() *cobra.Command {
	var firebaseKeyPath string
	var firebaseProjectID string

	cmd := &cobra.Command{
		Use:   "node",
		Short: "Perform some actions on the node",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Println("actions missing!")
			return nil
		},
	}
	cmd.Flags().StringVar(&firebaseKeyPath, "fbKeyPath", params.TestFirebaseAdminSDKPath, "path of firebase json key")
	cmd.Flags().StringVar(&firebaseProjectID, "fbProjectID", params.TestFirebaseProjectID, "project ID")

	cmd.AddCommand(NodeTestCmd())
	return cmd
}

// NodeTestCmd test connection by local firebase settings
func NodeTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test connection by local firebase settings",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {

			fmt.Println("test ...")

			return nil
		},
	}

	return cmd
}
