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
	"os"
	"strconv"
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
		Short: "Operations on node",
	}
	cmd.Flags().StringVar(&firebaseKeyPath, "fbKeyPath", params.TestFirebaseAdminSDKPath, "path of firebase json key")
	cmd.Flags().StringVar(&firebaseProjectID, "fbProjectID", params.TestFirebaseProjectID, "project ID")

	cmd.AddCommand(DiagnosisCmd())
	cmd.AddCommand(BackupCmd())
	cmd.AddCommand(ExecuteCmd())
	cmd.AddCommand(UploadCmd())
	cmd.AddCommand(TruncateCmd())
	cmd.AddCommand(JudgeCmd())
	cmd.AddCommand(HideProcessedQuantumCmd())
	return cmd
}

func DiagnosisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnosis",
		Short: "diagnosis quantums on remote node",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
			if err != nil {
				return err
			}
			qs, err := fbu.GetQuantumsOnWait()
			if err != nil {
				return err
			}
			fmt.Println(len(qs))
			for i, quantum := range qs {
				sigHex := core.Sig2Hex(quantum.Signature)
				fmt.Println("(", i, ")", sigHex[:10], "...", sigHex[120:])
			}

			for {
				res := question("please select the number of sigHex to be diagnosis (q for quit)", false)
				if res == "q" || res == "Q" {
					break
				}
				num, err := strconv.Atoi(res)
				if err != nil {
					fmt.Println(err)
					continue
				} else if num >= len(qs) {
					fmt.Println("number not exist")
					continue
				}
				// start to diagnosis

			}
			return nil
		},
	}
	return cmd

}

// UploadCmd
func UploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload quantums to remote node",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {

			dataFilePath := args[0]
			data, err := os.ReadFile(dataFilePath)
			if err != nil {
				return err
			}

			var qs []*core.Quantum
			if err := json.Unmarshal(data, &qs); err != nil {
				return err
			}

			ctx := context.Background()
			fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
			if err != nil {
				return err
			}

			accept, wait, reject, err := fbu.ReceiveQuantums(qs)
			if err != nil {
				return err
			}
			fmt.Println("upload finished! node accept:", len(accept), "\twait:", len(wait), "\treject:", len(reject))

			return nil
		},
	}

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

			timestamp := time.Now().UnixMilli()
			if err = os.WriteFile(fmt.Sprintf("./backup_quantums_%d", timestamp), res, 0644); err != nil {
				return err
			}
			fmt.Println("quantums backup finished!")

			if inds, err := fbu.GetFBDataByTable("individual"); err != nil {
				return err
			} else {
				if err = os.WriteFile(fmt.Sprintf("./backup_individuals_%d", timestamp), inds, 0644); err != nil {
					return err
				}
				fmt.Println("individuals backup finished!")

			}

			if coms, err := fbu.GetFBDataByTable("community"); err != nil {
				return err
			} else {
				if err = os.WriteFile(fmt.Sprintf("./backup_communties_%d", timestamp), coms, 0644); err != nil {
					return err
				}
				fmt.Println("communities backup finished!")

			}

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

				removeLocalData := boolChoice("delete local data in " + configPath)
				if removeLocalData {
					if err := os.RemoveAll(configPath); err != nil {
						return err
					} else {
						fmt.Println(configPath + "have been deleted!")
					}
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

	if err := truncateTable(client, ctx, "community"); err != nil {
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
