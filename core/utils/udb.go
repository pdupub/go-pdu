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

package utils

import (
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/udb"
)

// FromUDBIndividual from *udb.Individual to *core.Individual and sender uid
func FromUDBIndividual(dbIndividual *udb.Individual) (individual *core.Individual, senderDBID string) {
	if dbIndividual == nil {
		return individual, senderDBID
	}
	profile := make(map[string]*core.QContent)
	for _, v := range dbIndividual.Quantums {
		if v.Type == core.QuantumTypeProfile {
			for i := 0; i < len(v.Contents); i += 2 {
				c, _ := core.NewContent(v.Contents[i+1].Fmt, []byte(v.Contents[i+1].Data))
				profile[string(v.Contents[i].Data)] = c
			}
		}
	}

	// sender uid from db
	senderDBID = dbIndividual.UID

	individual = &core.Individual{
		Address: identity.HexToAddress(dbIndividual.Address),
		Profile: profile,
	}

	return individual, senderDBID
}

func ToUDBIndividual(individual *core.Individual, senderDBID string) *udb.Individual {
	dbIndividual := udb.Individual{
		UID:     senderDBID,
		Address: individual.Address.Hex(),
		DType:   []string{udb.DTypeIndividual},
	}

	return &dbIndividual
}

// FromUDBQuantum cover from *udb.Quantum to *core.Quantum and quantum uid and sender uid
func FromUDBQuantum(dbQuantum *udb.Quantum) (quantum *core.Quantum, quantumDBID string, senderDBID string) {
	if dbQuantum == nil {
		return quantum, quantumDBID, senderDBID
	}
	refs := []core.Sig{}
	for _, v := range dbQuantum.Refs {
		refs = append(refs, core.Sig(v.Sig))
	}
	contents := []*core.QContent{}
	for _, v := range dbQuantum.Contents {
		contents = append(contents, &core.QContent{Format: v.Fmt, Data: []byte(v.Data)})
	}

	// quantum uid from db
	quantumDBID = dbQuantum.UID

	// sender uid from db
	senderDBID = dbQuantum.Sender.UID

	quantum = &core.Quantum{
		Signature: []byte(dbQuantum.Sig),
		UnsignedQuantum: core.UnsignedQuantum{
			Type:       dbQuantum.Type,
			Contents:   contents,
			References: refs,
		},
	}

	return quantum, quantumDBID, senderDBID
}

func ToUDBQuantum(quantum *core.Quantum, quantumDBID string, senderDBID string) *udb.Quantum {
	refs := []*udb.Quantum{}
	for _, v := range quantum.References {
		refs = append(refs, &udb.Quantum{Sig: string(v)})
	}

	contents := []*udb.Content{}
	for _, v := range quantum.Contents {
		contents = append(contents, &udb.Content{Fmt: v.Format, Data: string(v.Data), DType: []string{udb.DTypeContent}})
	}

	sender := &udb.Individual{UID: senderDBID}
	dbQuantum := udb.Quantum{
		UID:      quantumDBID,
		Sig:      string(quantum.Signature),
		Type:     quantum.Type,
		Refs:     refs,
		Contents: contents,
		Sender:   sender,
		DType:    []string{udb.DTypeQuantum},
	}
	return &dbQuantum
}
