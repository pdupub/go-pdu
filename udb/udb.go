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

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
)

type UDB struct {
	ctx  context.Context
	conn *grpc.ClientConn
	dg   *dgo.Dgraph
}

// struct / table
type Individual struct {
	Address string   `json:"address,omitempty"`
	DType   []string `json:"dgraph.type,omitempty"`
}

func (udb *UDB) initIndividual() error {
	// empty operation
	op := &api.Operation{}
	// Individual, has only 1 field address (string)
	op.Schema = `
			address: string @index(exact) .
			type Individual {
				address
			}
		`
	// update schema = add new schema
	return udb.dg.Alter(udb.ctx, op)
}

func (udb *UDB) addIndividual(address string) error {
	p := Individual{
		Address: address,
		DType:   []string{"Individual"},
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

func (udb *UDB) queryIndividual(address string) ([]Individual, error) {
	// query from database
	variables := make(map[string]string)
	variables["$a"] = address
	q := `
					query QueryIndividual($a: string){
						queryRes(func: eq(address, $a)) {
							address
							dgraph.type
						}
					}
				`

	resp, err := udb.dg.NewTxn().QueryWithVars(udb.ctx, q, variables)
	if err != nil {
		return nil, err
	}

	type Root struct {
		Me []Individual `json:"queryRes"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(resp.Json))
	return r.Me, nil
}

func (udb *UDB) dropData() error {
	return udb.dg.Alter(udb.ctx, &api.Operation{DropOp: api.Operation_DATA})
}

func New(url, token string) (*UDB, error) {
	ctx := context.Background()

	conn, err := dgo.DialSlashEndpoint(url, token)
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
