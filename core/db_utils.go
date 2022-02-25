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

package core

import (
	"encoding/hex"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/udb"
)

func Sig2Hex(sig Sig) string {
	return hex.EncodeToString(sig)
}

func Hex2Sig(str string) Sig {
	h, _ := hex.DecodeString(str)
	return h
}

// FromUDBIndividual from *udb.Individual to *Individual and sender uid
func FromUDBIndividual(dbIndividual *udb.Individual) (individual *Individual, senderDBID string) {
	if dbIndividual == nil {
		return individual, senderDBID
	}

	// sender uid from db
	senderDBID = dbIndividual.UID
	individual = NewIndividual(identity.HexToAddress(dbIndividual.Address))

	// update communities
	communities := []*Community{}
	for _, v := range dbIndividual.Communities {
		communities = append(communities, &Community{Define: Hex2Sig(v.Base.Define.Sig)})
	}
	individual.Communities = communities

	// update attitude
	evidence := []*Quantum{}
	if dbIndividual.Attitude != nil {
		for _, v := range dbIndividual.Attitude.Evidence {
			evidence = append(evidence, &Quantum{Signature: Hex2Sig(v.Sig)})
		}
		individual.Attitude = &Attitude{
			Level:    1,
			Judgment: "",
			Evidence: evidence,
		}
	}
	return individual, senderDBID
}

// ToUDBIndividual should not be used for any reason
// func ToUDBIndividual(individual *Individual, senderDBID string) *udb.Individual {
// 	// TODO: add communities & quantums
// 	dbIndividual := udb.Individual{
// 		UID:     senderDBID,
// 		Address: individual.Address.Hex(),
// 		DType:   udb.DTypeIndividual,
// 	}

// 	return &dbIndividual
// }

// FromUDBQuantum cover from *udb.Quantum to *Quantum and quantum uid and sender uid
func FromUDBQuantum(dbQuantum *udb.Quantum) (quantum *Quantum, quantumDBID string, senderDBID string) {
	if dbQuantum == nil {
		return quantum, quantumDBID, senderDBID
	}
	refs := []Sig{}
	for _, v := range dbQuantum.Refs {
		refs = append(refs, Hex2Sig(v.Sig))
	}
	contents := []*QContent{}
	for _, v := range dbQuantum.Contents {
		contents = append(contents, &QContent{Format: v.Fmt, Data: []byte(v.Data)})
	}

	// quantum uid from db
	quantumDBID = dbQuantum.UID

	// sender uid from db
	senderDBID = dbQuantum.Sender.UID

	quantum = &Quantum{
		Signature: Hex2Sig(dbQuantum.Sig),
		UnsignedQuantum: UnsignedQuantum{
			Type:       dbQuantum.Type,
			Contents:   contents,
			References: refs,
		},
	}

	return quantum, quantumDBID, senderDBID
}

// ToUDBQuantum
func ToUDBQuantum(quantum *Quantum, quantumDBID string, senderDBID string) *udb.Quantum {
	refs := []*udb.Quantum{}
	for _, v := range quantum.References {
		refs = append(refs, &udb.Quantum{Sig: Sig2Hex(v)})
	}

	contents := []*udb.Content{}
	for _, v := range quantum.Contents {
		contents = append(contents, &udb.Content{Fmt: v.Format, Data: string(v.Data), DType: udb.DTypeContent})
	}

	addr, _ := quantum.Ecrecover()
	sender := &udb.Individual{UID: senderDBID, Address: addr.Hex()}

	dbQuantum := udb.Quantum{
		UID:      quantumDBID,
		Sig:      Sig2Hex(quantum.Signature),
		Type:     quantum.Type,
		Refs:     refs,
		Contents: contents,
		Sender:   sender,
		DType:    udb.DTypeQuantum,
	}
	return &dbQuantum
}

// FromUDBCommunity
func FromUDBCommunity(dbCommunity *udb.Community) (community *Community, communityDBID string) {
	initMembers := []identity.Address{}
	for _, v := range dbCommunity.InitMembers {
		initMembers = append(initMembers, identity.HexToAddress(v.Address))
	}

	communityDBID = dbCommunity.UID
	community = &Community{
		Note:          &QContent{Format: dbCommunity.Note.Fmt, Data: []byte(dbCommunity.Note.Data)},
		Define:        Hex2Sig(dbCommunity.Define.Sig),
		Creator:       identity.HexToAddress(dbCommunity.Define.Sender.Address),
		BaseCommunity: Hex2Sig(dbCommunity.Base.Define.Sig),
		MinCosignCnt:  dbCommunity.MinCosignCnt,
		MaxInviteCnt:  dbCommunity.MaxInviteCnt,
		InitMembers:   initMembers,
	}
	return community, communityDBID
}

func ToUDBCommunity(community *Community, communityDBID string) *udb.Community {
	initMembers := []*udb.Individual{}
	for _, v := range community.InitMembers {
		initMembers = append(initMembers, &udb.Individual{Address: v.Hex()})
	}

	// ignore content here
	dbCommunity := udb.Community{
		UID:          communityDBID,
		Base:         &udb.Community{Define: &udb.Quantum{Sig: Sig2Hex(community.BaseCommunity)}},
		MaxInviteCnt: community.MaxInviteCnt,
		MinCosignCnt: community.MinCosignCnt,
		Define:       &udb.Quantum{Sig: Sig2Hex(community.Define)},
		InitMembers:  initMembers,
		DType:        udb.DTypeCommunity,
	}
	return &dbCommunity
}
