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

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/spf13/cobra"
)

const (
	testUrl   = "https://blue-surf-550111.us-east-1.aws.cloud.dgraph.io/graphql"
	testToken = "ZGE5M2IwNWFmMmY5NTYwN2ExODg4ZGYxMWJiYmRkZDg="
)

func DgoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dgo",
		Short: "Test DGraph func",
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()

			// create conn
			conn, err := dgo.DialSlashEndpoint(testUrl, testToken)
			if err != nil {
				return err
			}
			defer conn.Close()

			// create client
			dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

			type Person struct {
				Uid   string   `json:"uid,omitempty"`
				Name  string   `json:"name,omitempty"`
				DType []string `json:"dgraph.type,omitempty"`
			}

			op := &api.Operation{}

			op.Schema = `
			name: string @index(exact) .
			type Person {
				name
			}
		`

			if err := dg.Alter(ctx, op); err != nil {
				return err
			}

			p := Person{
				Name:  "Alice",
				DType: []string{"Person"},
			}

			mu := &api.Mutation{
				CommitNow: true,
			}
			pb, err := json.Marshal(p)
			if err != nil {
				return err
			}

			mu.SetJson = pb
			_, err = dg.NewTxn().Mutate(ctx, mu)
			if err != nil {
				return err
			}

			variables := make(map[string]string)
			variables["$a"] = "Alice"
			q := `
				query Alice($a: string){
					me(func: eq(name, $a)) {
						name
						dgraph.type
					}
				}
			`

			resp, err := dg.NewTxn().QueryWithVars(ctx, q, variables)
			if err != nil {
				return err
			}

			type Root struct {
				Me []Person `json:"me"`
			}

			var r Root
			err = json.Unmarshal(resp.Json, &r)
			if err != nil {
				return err
			}

			fmt.Println(string(resp.Json))
			return nil
		},
	}
	return cmd
}
