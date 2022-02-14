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
	"github.com/pdupub/go-pdu/newv/core"
	"google.golang.org/grpc"
)

// Value of DType is same with the name of type, so can be used to expand, like expand(DTypeQuantum)
const (
	DTypeQuantum    = "quantum"
	DTypeContent    = "content"
	DTypeIndividual = "individual"
	DTypeCommunity  = "community"
)

type UDB struct {
	ctx  context.Context
	conn *grpc.ClientConn
	dg   *dgo.Dgraph
}
type QueryResultRoot struct {
	QuantumRes    []Quantum    `json:"quantum"`
	ContentRes    []Content    `json:"content"`
	IndividualRes []Individual `json:"individual"`
	CommunityRes  []Community  `json:"community"`
}

type Quantum struct {
	UID      string      `json:"uid,omitempty"`
	Sig      string      `json:"quantum.sig,omitempty"`
	Type     int         `json:"quantum.type,omitempty"`
	Refs     []*Quantum  `json:"quantum.refs,omitempty"`
	Contents []*Content  `json:"quantum.contents,omitempty"`
	Sender   *Individual `json:"quantum.sender,omitempty"`
	DType    []string    `json:"dgraph.type,omitempty"` // Q
}

type Content struct {
	UID   string   `json:"uid,omitempty"`
	Fmt   int      `json:"content.fmt,omitempty"`
	Data  string   `json:"content.data,omitempty"`
	DType []string `json:"dgraph.type,omitempty"` // C
}

type Individual struct {
	UID         string       `json:"uid,omitempty"`
	Address     string       `json:"individual.address,omitempty"`
	Communities []*Community `json:"individual.communities,omitempty"`
	Quantums    []*Quantum   `json:"individual.quantums,omitempty"`
	DType       []string     `json:"dgraph.type,omitempty"` // I
}

type Community struct {
	UID          string        `json:"uid,omitempty"`
	Base         *Quantum      `json:"community.base,omitempty"`
	Invitations  []*Quantum    `json:"community.invitations,omitempty"`
	MaxInviteCnt int           `json:"community.maxInviteCnt,omitempty"`
	MinCosignCnt int           `json:"community.minCosignCnt,omitempty"`
	Members      []*Individual `json:"community.members,omitempty"`
	Rule         *Quantum      `json:"community.rule,omitempty"`
	DType        []string      `json:"dgraph.type,omitempty"` // Community
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

func (udb *UDB) initSchema() error {
	// empty operation
	op := &api.Operation{}
	// Individual, has only 1 field address (string)
	op.Schema = Schema
	// update schema = add new schema
	return udb.dg.Alter(udb.ctx, op)
}

func (udb *UDB) dropData() error {
	return udb.dg.Alter(udb.ctx, &api.Operation{DropOp: api.Operation_DATA})
}

// SetQuantum works like set operation in DQL (upsert), but by the signature not uid.
func (udb *UDB) SetQuantum(quantum *core.Quantum, address string) (uid string, err error) {
	// try to find quantum with same sig in db
	_, lastQuantum, err := udb.GetQuantum(quantum.Signature)
	if err != nil {
		return uid, err
	}

	// if exist
	if lastQuantum != nil {
		uid = lastQuantum.UID
		// return if exist, not empty
		// exist quantum can not be change by any condition
		if lastQuantum.Type != 0 {
			return lastQuantum.UID, nil
		}
	}

	refs := []*Quantum{}
	for _, v := range quantum.References {
		_, dbQ, err := udb.GetQuantum(v)
		if err != nil {
			return uid, err
		}
		if dbQ != nil {
			// each predicate in struct must be omitempty to avoid remove value
			refs = append(refs, &Quantum{UID: dbQ.UID})
		} else {
			refs = append(refs, udb.buildEmptyQuantum(v))
		}
	}
	contents := []*Content{}
	for _, v := range quantum.Contents {
		contents = append(contents, &Content{Fmt: v.Format, Data: string(v.Data), DType: []string{DTypeContent}})
	}
	sender, err := udb.GetIndividual(address)
	if err != nil {
		return uid, err
	}
	if sender == nil {
		sender = udb.buildEmptyIndividual(address)
	}
	p := Quantum{
		UID:      uid,
		DType:    []string{DTypeQuantum},
		Sig:      string(quantum.Signature),
		Type:     quantum.Type,
		Refs:     refs,
		Contents: contents,
		Sender:   sender,
	}

	if uid == "" {
		p.UID = "_:QuantumUID"
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(p)
	if err != nil {
		return uid, err
	}

	mu.SetJson = pb
	resp, err := udb.dg.NewTxn().Mutate(udb.ctx, mu)
	if err != nil {
		return uid, err
	}

	if uid == "" {
		uid = resp.Uids["QuantumUID"]
	}

	return uid, nil
}

func (udb *UDB) buildEmptyQuantum(sig core.Sig) *Quantum {
	return &Quantum{Sig: string(sig),
		DType: []string{DTypeQuantum}}
}

// GetQuantum query the Quantum by Signature of Quantum
func (udb *UDB) GetQuantum(sig core.Sig) (quantum *core.Quantum, dbq *Quantum, err error) {
	// query quantum by signature
	params := make(map[string]string)
	params["$sig"] = string(sig)
	params["$type"] = DTypeQuantum
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
	resp, err := udb.dg.NewTxn().QueryWithVars(udb.ctx, q, params)
	if err != nil {
		return nil, nil, err
	}

	var r QueryResultRoot
	err = json.Unmarshal(resp.Json, &r)
	if err != nil || len(r.QuantumRes) != 1 {
		return nil, nil, err
	}

	dbq = &r.QuantumRes[0]
	return udb.parseQuantum(dbq), &r.QuantumRes[0], nil
}

func (udb *UDB) parseQuantum(quantum *Quantum) *core.Quantum {
	if quantum == nil {
		return nil
	}
	refs := []core.Sig{}
	for _, v := range quantum.Refs {
		refs = append(refs, core.Sig(v.Sig))
	}
	contents := []*core.QContent{}
	for _, v := range quantum.Contents {
		contents = append(contents, &core.QContent{Format: v.Fmt, Data: []byte(v.Data)})
	}
	return &core.Quantum{
		Signature: []byte(quantum.Sig),
		UnsignedQuantum: core.UnsignedQuantum{
			Type:       quantum.Type,
			Contents:   contents,
			References: refs,
		},
	}
}

func (udb *UDB) buildEmptyIndividual(address string) *Individual {
	return &Individual{Address: address,
		DType: []string{DTypeIndividual}}
}

// GetIndividual query the Individual by Address Hex
func (udb *UDB) GetIndividual(address string) (dbi *Individual, err error) {
	// query quantum by signature
	params := make(map[string]string)
	params["$addr"] = address
	params["$type"] = DTypeIndividual
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
	resp, err := udb.dg.NewTxn().QueryWithVars(udb.ctx, q, params)
	if err != nil {
		return nil, err
	}

	var r QueryResultRoot
	err = json.Unmarshal(resp.Json, &r)
	if err != nil || len(r.IndividualRes) != 1 {
		return nil, err
	}

	return &r.IndividualRes[0], nil
}
