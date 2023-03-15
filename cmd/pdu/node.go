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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
	"github.com/pdupub/go-pdu/udb/fb"
)

// NodeCmd manually perform some actions on the node,
// like set address level or remove processed message from the node.
func NodeCmd() *cobra.Command {
	var firebaseKeyPath string
	var firebaseProjectID string

	cmd := &cobra.Command{
		Use:   "node",
		Short: "Perform some actions on the node",
	}
	cmd.Flags().StringVar(&firebaseKeyPath, "fbKeyPath", params.TestFirebaseAdminSDKPath, "path of firebase json key")
	cmd.Flags().StringVar(&firebaseProjectID, "fbProjectID", params.TestFirebaseProjectID, "project ID")

	cmd.AddCommand(BackupCmd())
	cmd.AddCommand(ExecuteCmd())
	cmd.AddCommand(TruncateCmd())
	cmd.AddCommand(JudgeCmd())
	cmd.AddCommand(HideProcessedQuantumCmd())
	return cmd
}

// BackupCmd do process quantum once on node
func BackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup processed quantums to local",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
			if err != nil {
				return err
			}
			limit := 10
			skip := 0
			var backup []*core.Quantum
			for {
				records, err := fbu.GetQuantums(limit, skip, true)
				if err != nil {
					return err
				}
				skip += len(records)
				if len(records) == 0 {
					break
				}
				fmt.Println("already download", skip, "quantums")
				backup = append(backup, records...)
			}

			res, err := json.Marshal(backup)
			if err != nil {
				return err
			}
			fmt.Println()
			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

func HideProcessedQuantumCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hide",
		Short: "Hide processed Quantum in node",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {

			ctx := context.Background()
			fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
			if err != nil {
				return err
			}

			sigHex := question("please provide the signature which you want to hide", false)
			if err = fbu.HideProcessedQuantum(core.Hex2Sig(sigHex)); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

// JudgeCmd used to judge Individual and Community on your own node.
func JudgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "judge",
		Short: "Judge Individual & Community on your own node",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {

			ctx := context.Background()
			fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
			if err != nil {
				return err
			}

			_, target := multiChoice("please select the target", "Individual", "Community")
			if target == 0 {
				addrHex := question("please provide the address which you want to reset the level", false)
				// TODO: addrHex should check

				_, level := multiChoice("please select the new status of this address", "reject on refs", "reject", "ignore content", "accept", "broadcast")
				if level < 2 {
					level -= 2
				} else {
					level -= 1
				}

				if err = fbu.JudgeIndividual(identity.HexToAddress(addrHex), level, ""); err != nil {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}

// TruncateCmd will clear up all data on firebase collections
func TruncateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "truncate",
		Short: "Clear up all data on firebase collections",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			confirm := boolChoice("ARE YOU SURE YOU WANT TO DELETE ALL DATA ON Firebase!!!")
			if confirm {
				delStr := question("Please type input \"Delete\"", false)
				if delStr == "Delete" {
					fmt.Println("start to delete data")
					truncate()
					fmt.Println("deleting finished!")
				}
			}

			return nil
		},
	}

	return cmd
}

// ExecuteCmd do process quantum once on node
func ExecuteCmd() *cobra.Command {
	var limit, skip int
	cmd := &cobra.Command{
		Use:   "exe",
		Short: "Do process quantum once on node",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
			if err != nil {
				return err
			}
			fbu.ProcessQuantums(limit, skip)
			return nil
		},
	}
	cmd.PersistentFlags().IntVar(&limit, "limit", 20, "process quantums number limit")
	cmd.PersistentFlags().IntVar(&skip, "skip", 0, "process quantums need to skip")

	return cmd
}

func truncate() error {
	ctx := context.Background()
	opt := option.WithCredentialsFile(projectPath + firebaseKeyPath)
	config := &firebase.Config{ProjectID: firebaseProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := truncateTable(client, ctx, "communtiy"); err != nil {
		return err
	}
	if err := truncateTable(client, ctx, "individual"); err != nil {
		return err
	}
	if err := truncateTable(client, ctx, "quantum"); err != nil {
		return err
	}

	// initialize data in universe
	configCollection := client.Collection("universe")
	configDocRef := configCollection.Doc("status")
	configMap := make(map[string]interface{})
	configMap["lastSequence"] = 0
	configMap["lastSigHex"] = ""
	configMap["updateTime"] = time.Now().UnixMilli()
	configDocRef.Set(ctx, configMap, firestore.Merge([]string{"lastSequence"}, []string{"lastSigHex"}, []string{"updateTime"}))

	return nil
}

func truncateTable(client *firestore.Client, ctx context.Context, collectionName string) error {
	currentCol := client.Collection(collectionName)
	docRefs, err := currentCol.DocumentRefs(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, docRef := range docRefs {
		docRef.Delete(ctx) // ignore err here
	}
	fmt.Println("collection", collectionName, "have been truncated")

	return nil
}
