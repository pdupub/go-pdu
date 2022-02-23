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
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/pdupub/go-pdu/udb"
	"google.golang.org/grpc"
)

var (
	ErrQuantumMissing    = errors.New("quantum missing")
	ErrIndividualMissing = errors.New("individual missing")
)

// CDGD is root struct for operate Cloud DGraph Database
type CDGD struct {
	ctx  context.Context
	conn *grpc.ClientConn
	dg   *dgo.Dgraph
}

// QueryResultRoot is used for UnMarshal query response
type QueryResultRoot struct {
	QuantumRes    []*udb.Quantum    `json:"quantum"`
	ContentRes    []*udb.Content    `json:"content"`
	IndividualRes []*udb.Individual `json:"individual"`
	CommunityRes  []*udb.Community  `json:"community"`
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

func (cdgd *CDGD) SetSchema() error {
	// empty operation
	op := &api.Operation{}
	// Individual, has only 1 field address (string)
	op.Schema = Schema
	// update schema = add new schema
	return cdgd.dg.Alter(cdgd.ctx, op)
}

func (cdgd *CDGD) DropData() error {
	return cdgd.dg.Alter(cdgd.ctx, &api.Operation{DropOp: api.Operation_DATA})
}

// NewQuantum works like set operation in DQL (upsert), but by the signature not uid.
func (cdgd *CDGD) NewQuantum(quantum *udb.Quantum) (qid string, sid string, err error) {
	// try to find quantum with same sig in db
	lastQuantum, err := cdgd.GetQuantum(quantum.Sig)
	if err != nil {
		return qid, sid, err
	}

	// if exist
	if lastQuantum != nil {
		qid = lastQuantum.UID
		// return if exist, not empty
		// exist quantum can not be change by any condition
		if lastQuantum.Sender != nil {
			sid = lastQuantum.Sender.UID
			return qid, sid, nil
		}
	}

	refs := []*udb.Quantum{}
	for _, v := range quantum.Refs {
		dbQ, err := cdgd.GetQuantum(v.Sig)
		if err != nil {
			return qid, sid, err
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
		contents = append(contents, &udb.Content{Fmt: v.Fmt, Data: v.Data, DType: udb.DTypeContent})
	}
	sender, err := cdgd.GetIndividual(quantum.Sender.Address)
	if err != nil {
		return qid, sid, err
	}
	if sender == nil {
		sender = cdgd.buildEmptyIndividual(quantum.Sender.Address)
	} else {
		sid = sender.UID
	}
	p := udb.Quantum{
		UID:       qid,
		DType:     udb.DTypeQuantum,
		Sig:       quantum.Sig,
		Type:      quantum.Type,
		Refs:      refs,
		Contents:  contents,
		Sender:    sender,
		Timestamp: int(time.Now().UnixNano()),
	}

	if qid == "" {
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
		return qid, sid, err
	}

	mu.SetJson = pb
	resp, err := cdgd.dg.NewTxn().Mutate(cdgd.ctx, mu)
	if err != nil {
		return qid, sid, err
	}

	if qid == "" {
		qid = resp.Uids["QuantumUID"]
	}
	if sid == "" {
		sid = resp.Uids["IndividualUID"]
	}

	return qid, sid, nil
}

func (cdgd *CDGD) buildEmptyQuantum(sig string) *udb.Quantum {
	return &udb.Quantum{Sig: sig,
		DType: udb.DTypeQuantum}
}

// QueryQuantum is query quantums by params
// if address == "" then ignore the sender
// if qType == 0    then ignore the type of quantum
// if pageIndex < 1 then pageIndex = 1
// if pageSize < 0  then ignore the page size
func (cdgd *CDGD) QueryQuantum(address string, qType int, pageIndex int, pageSize int, desc bool) ([]*udb.Quantum, error) {

	if pageIndex < 1 {
		pageIndex = 1
	}
	params := make(map[string]string)
	params["$address"] = address
	params["$dtype"] = udb.DTypeIndividual
	params["$first"] = strconv.Itoa(pageSize)
	params["$offset"] = strconv.Itoa((pageIndex - 1) * pageSize)

	// DQL
	q := `
			query QueryQuantum($address: string, $dtype: string, $first: int, $offset: int){
				res(func: eq( individual.address, $address)) @filter(eq(pdu.type, $dtype)){
					uid
					quantum: ~quantum.sender FILTER_QTYPE (first: $first, offset: $offset ORDER) {
						uid
						expand(quantum){
							uid
							expand(quantum, individual, content)
							pdu.type
						}
						pdu.type
					}
				}
			}
		`
	//
	if qType <= 0 {
		q = strings.Replace(q, "FILTER_QTYPE", "", -1)
	} else {
		q = strings.Replace(q, "FILTER_QTYPE", "@filter(eq(quantum.type, "+strconv.Itoa(qType)+"))", -1)
	}

	if desc {
		q = strings.Replace(q, "ORDER", ", orderdesc: quantum.timestamp", -1)
	} else {
		q = strings.Replace(q, "ORDER", ", orderasc: quantum.timestamp", -1)
	}

	// send query request
	resp, err := cdgd.dg.NewTxn().QueryWithVars(cdgd.ctx, q, params)
	if err != nil {
		return nil, err
	}

	type RespRes struct {
		Result []*QueryResultRoot `json:"res"`
	}
	var r RespRes
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	return r.Result[0].QuantumRes, nil
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
				quantum(func: eq(quantum.sig, $sig)) @filter(eq(pdu.type, $type)){
					uid
					expand(quantum){
						uid
						expand(quantum, individual, content)
						pdu.type
					}
					pdu.type
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

	dbq = r.QuantumRes[0]
	return dbq, nil
}

func (cdgd *CDGD) buildEmptyIndividual(address string) *udb.Individual {
	return &udb.Individual{Address: address,
		DType: udb.DTypeIndividual}
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
				individual(func: eq(individual.address, $addr)) @filter(eq(pdu.type, $type)){
					uid
					expand(individual) {
						uid
					  	expand(community)
					  	pdu.type
					}
					pdu.type
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
	dbi = r.IndividualRes[0]
	return dbi, nil
}

// NewCommunity
func (cdgd *CDGD) NewCommunity(community *udb.Community) (cid string, err error) {
	// check quantum uid must exist
	if community.Define == nil {
		return cid, ErrQuantumMissing
	}

	// get define by community define sig
	if community.Define.Sig == "" {
		return cid, ErrQuantumMissing
	}
	quantum, err := cdgd.GetQuantum(community.Define.Sig)
	if err != nil {
		return cid, err
	}
	if quantum.Sender == nil {
		return cid, ErrIndividualMissing
	}

	community.Define = quantum

	// check countents length
	if len(quantum.Contents) < 4 {
		return cid, ErrQuantumMissing
	}

	// use quantum's first content uid as uid of note
	community.Note = &udb.Content{UID: quantum.Contents[0].UID}

	// set base by uid
	// if uid of base community exist , it must alse be exist in database
	if community.Base != nil && community.Base.Define != nil {
		baseCommunity, err := cdgd.GetCommunity(community.Base.Define.Sig)
		if err != nil {
			return cid, err
		}
		community.Base = &udb.Community{UID: baseCommunity.UID}
	}

	// set initMembers by uid
	for _, v := range community.InitMembers {
		if v.Address == "" {
			return cid, ErrIndividualMissing
		}
		individual, err := cdgd.GetIndividual(v.Address)
		if err != nil {
			return cid, err
		}

		if individual == nil {
			community.InitMembers = append(community.InitMembers, cdgd.buildEmptyIndividual(v.Address))
		} else {
			community.InitMembers = append(community.InitMembers, individual)
		}
	}

	community.DType = udb.DTypeCommunity
	community.UID = "_:CommunityUID"

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(community)
	if err != nil {
		return cid, err
	}

	mu.SetJson = pb
	resp, err := cdgd.dg.NewTxn().Mutate(cdgd.ctx, mu)
	if err != nil {
		return cid, err
	}

	cid = resp.Uids["CommunityUID"]

	// creator join community
	cdgd.JoinCommunity(community.Define.Sig, community.Define.Sender.Address)

	// init members join community
	for _, v := range community.InitMembers {
		cdgd.JoinCommunity(community.Define.Sig, v.Address)
	}
	return cid, err
}

func (cdgd *CDGD) JoinCommunity(defineSig string, address string) (cid string, sid string, err error) {

	// update individual.community of creator & initMembers
	community, _ := cdgd.GetCommunity(defineSig)
	cid = community.UID
	individual, _ := cdgd.GetIndividual(address)
	sid = individual.UID
	individual.Communities = append(individual.Communities, community)
	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(individual)
	if err != nil {
		return cid, sid, err
	}

	mu.SetJson = pb
	_, err = cdgd.dg.NewTxn().Mutate(cdgd.ctx, mu)
	if err != nil {
		return cid, sid, err
	}

	return cid, sid, nil
}

func (cdgd *CDGD) GetCommunity(sig string) (dbc *udb.Community, err error) {

	params := make(map[string]string)
	params["$sig"] = sig
	params["$dtype"] = udb.DTypeQuantum

	// DQL
	q := `
			query QueryCommunity($sig: string, $dtype: string){
				res(func: eq(quantum.sig, $sig)) @filter(eq(pdu.type, $dtype)){
					uid
					community: ~community.define {
						uid
						expand(community){
							uid
							expand(community, content, quantum, individual)
							pdu.type
						}
						pdu.type
					}
				}
			}
		`

	// send query request
	resp, err := cdgd.dg.NewTxn().QueryWithVars(cdgd.ctx, q, params)
	if err != nil {
		return nil, err
	}

	type RespRes struct {
		Result []*QueryResultRoot `json:"res"`
	}
	var r RespRes
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	if len(r.Result) > 0 && len(r.Result[0].CommunityRes) > 0 {
		dbc = r.Result[0].CommunityRes[0]
	}

	return dbc, nil
}
