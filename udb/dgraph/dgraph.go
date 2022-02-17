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

package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/pdupub/go-pdu/udb"
	"google.golang.org/grpc"
)

// CDGD is root struct for operate Cloud DGraph Database
type CDGD struct {
	ctx  context.Context
	conn *grpc.ClientConn
	dg   *dgo.Dgraph
}

// QueryResultRoot is used for UnMarshal query response
type QueryResultRoot struct {
	QuantumRes    []udb.Quantum    `json:"quantum"`
	ContentRes    []udb.Content    `json:"content"`
	IndividualRes []udb.Individual `json:"individual"`
	CommunityRes  []udb.Community  `json:"community"`
}

func New(url, token string) (*CDGD, error) {
	ctx := context.Background()

	conn, err := dgo.DialSlashEndpoint(url, token)
	if err != nil {
		return nil, err
	}

	dg := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	udb := CDGD{
		ctx:  ctx,
		conn: conn,
		dg:   dg,
	}
	return &udb, nil
}

func (cdgd *CDGD) Close() error {
	return cdgd.conn.Close()
}

func (cdgd *CDGD) initSchema() error {
	// empty operation
	op := &api.Operation{}
	// Individual, has only 1 field address (string)
	op.Schema = Schema
	// update schema = add new schema
	return cdgd.dg.Alter(cdgd.ctx, op)
}

func (cdgd *CDGD) dropData() error {
	return cdgd.dg.Alter(cdgd.ctx, &api.Operation{DropOp: api.Operation_DATA})
}

// NewQuantum works like set operation in DQL (upsert), but by the signature not uid.
func (cdgd *CDGD) NewQuantum(quantum *udb.Quantum) (uid string, sid string, err error) {
	// try to find quantum with same sig in db
	lastQuantum, err := cdgd.GetQuantum(quantum.Sig)
	if err != nil {
		return uid, sid, err
	}

	// if exist
	if lastQuantum != nil {
		uid = lastQuantum.UID
		// return if exist, not empty
		// exist quantum can not be change by any condition
		if lastQuantum.Type != 0 {
			sid = lastQuantum.Sender.UID
			return uid, sid, nil
		}
	}

	refs := []*udb.Quantum{}
	for _, v := range quantum.Refs {
		dbQ, err := cdgd.GetQuantum(v.Sig)
		if err != nil {
			return uid, sid, err
		}
		if dbQ != nil {
			// each predicate in struct must be omitempty to avoid remove value
			refs = append(refs, &udb.Quantum{UID: dbQ.UID})
		} else {
			refs = append(refs, cdgd.buildEmptyQuantum(v.Sig))
		}
	}
	contents := []*udb.Content{}
	for _, v := range quantum.Contents {
		contents = append(contents, &udb.Content{Fmt: v.Fmt, Data: v.Data, DType: []string{udb.DTypeContent}})
	}
	sender, err := cdgd.GetIndividual(quantum.Sender.Address)
	if err != nil {
		return uid, sid, err
	}
	if sender == nil {
		sender = cdgd.buildEmptyIndividual(quantum.Sender.Address)
	} else {
		sid = sender.UID
	}
	p := udb.Quantum{
		UID:      uid,
		DType:    []string{udb.DTypeQuantum},
		Sig:      quantum.Sig,
		Type:     quantum.Type,
		Refs:     refs,
		Contents: contents,
		Sender:   sender,
	}

	if uid == "" {
		p.UID = "_:QuantumUID"
	}

	if sid == "" {
		p.Sender.UID = "_:IndividualUID"
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(p)
	if err != nil {
		return uid, sid, err
	}

	mu.SetJson = pb
	resp, err := cdgd.dg.NewTxn().Mutate(cdgd.ctx, mu)
	if err != nil {
		return uid, sid, err
	}

	if uid == "" {
		uid = resp.Uids["QuantumUID"]
	}
	if sid == "" {
		sid = resp.Uids["IndividualUID"]
	}

	return uid, sid, nil
}

func (cdgd *CDGD) buildEmptyQuantum(sig string) *udb.Quantum {
	return &udb.Quantum{Sig: sig,
		DType: []string{udb.DTypeQuantum}}
}

// GetQuantum query the Quantum by Signature of Quantum
func (cdgd *CDGD) GetQuantum(sig string) (dbq *udb.Quantum, err error) {
	// query quantum by signature
	params := make(map[string]string)
	params["$sig"] = sig
	params["$type"] = udb.DTypeQuantum
	// DQL
	q := `
			query QueryQuantum($sig: string, $type: string){
				quantum(func: eq(quantum.sig, $sig)) @filter(eq(dgraph.type, $type)){
					uid
					quantum.contents {
						uid
						expand(content)
						dgraph.type
					}
					quantum.refs {
						uid
						expand(quantum)
						dgraph.type
					}
					quantum.sender{
						uid
						expand(individual)
						dgraph.type
					}
					quantum.sig
					quantum.type
					dgraph.type
				}
			}
		`
	// send query request
	resp, err := cdgd.dg.NewTxn().QueryWithVars(cdgd.ctx, q, params)
	if err != nil {
		return nil, err
	}

	var r QueryResultRoot
	err = json.Unmarshal(resp.Json, &r)
	if err != nil || len(r.QuantumRes) != 1 {
		return nil, err
	}

	dbq = &r.QuantumRes[0]
	return dbq, nil
}

func (cdgd *CDGD) buildEmptyIndividual(address string) *udb.Individual {
	return &udb.Individual{Address: address,
		DType: []string{udb.DTypeIndividual}}
}

func (cdgd *CDGD) NewIndividual(address string) (uid string, err error) {
	lastIndividual, err := cdgd.GetIndividual(address)
	if err != nil {
		return uid, err
	}

	// if exist
	if lastIndividual != nil {
		uid = lastIndividual.UID
		return uid, nil
	}

	p := cdgd.buildEmptyIndividual(address)
	p.UID = "_:IndividualUID"

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(p)
	if err != nil {
		return uid, err
	}

	mu.SetJson = pb
	resp, err := cdgd.dg.NewTxn().Mutate(cdgd.ctx, mu)
	if err != nil {
		return uid, err
	}
	uid = resp.Uids["IndividualUID"]

	return uid, nil
}

// GetIndividual query the Individual by Address Hex
func (cdgd *CDGD) GetIndividual(address string) (dbi *udb.Individual, err error) {
	// query quantum by signature
	params := make(map[string]string)
	params["$addr"] = address
	params["$type"] = udb.DTypeIndividual
	// DQL
	q := `
			query QueryIndividual($addr: string, $type: string){
				individual(func: eq(individual.address, $addr)) @filter(eq(dgraph.type, $type)){
					uid
					individual.address
					individual.communities {
						uid
						expand(community)
						dgraph.type
					}
					individual.quantums {
						uid
						expand(quantums)
						dgraph.type
					}
					dgraph.type
				}
			}
		`
	// send query request
	resp, err := cdgd.dg.NewTxn().QueryWithVars(cdgd.ctx, q, params)
	if err != nil {
		return dbi, err
	}

	var r QueryResultRoot
	err = json.Unmarshal(resp.Json, &r)
	if err != nil || len(r.IndividualRes) != 1 {
		return dbi, err
	}
	dbi = &r.IndividualRes[0]
	return dbi, nil
}
