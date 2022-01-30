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

package udb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/pdupub/go-pdu/newv/core"
	"google.golang.org/grpc"
)

type UDB struct {
	ctx  context.Context
	conn *grpc.ClientConn
	dg   *dgo.Dgraph
}

// struct / table
type Person struct {
	Uid   string   `json:"uid,omitempty"`
	Name  string   `json:"name,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type DIndividual struct {
	core.Individual
	DType []string `json:"dgraph.type,omitempty"`
}

func (udb *UDB) initIndividual() error {
	// empty operation
	op := &api.Operation{}
	// Person, has only 1 field name (string)
	op.Schema = `
			name: string @index(exact) .
			type Person {
				name
			}
		`
	// update schema = add new schema
	return udb.dg.Alter(udb.ctx, op)
}

func (udb *UDB) addIndividual() error {
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
	_, err = udb.dg.NewTxn().Mutate(udb.ctx, mu)
	if err != nil {
		return err
	}
	return nil
}

func (udb *UDB) queryIndividual() error {
	// query from database
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

	resp, err := udb.dg.NewTxn().QueryWithVars(udb.ctx, q, variables)
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
}

func (udb *UDB) dropData() error {
	return udb.dg.Alter(udb.ctx, &api.Operation{DropOp: api.Operation_DATA})
}

func New(url, token string) (*UDB, error) {
	ctx := context.Background()

	conn, err := dgo.DialSlashEndpoint(testUrl, testToken)
	if err != nil {
		return nil, err
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	udb := UDB{
		ctx:  ctx,
		conn: conn,
		dg:   dg,
	}
	return &udb, nil
}

func (udb *UDB) Close() error {
	return udb.conn.Close()
}
