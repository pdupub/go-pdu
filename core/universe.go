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
	"github.com/pdupub/go-pdu/identity"
)

// Universe is an interface that describes an PDU Universe. A type that implements Universe
// contains all quantums which be received and accepted and has ability to process the quantums.
// The owner of node have right to decide accept or reject any quantums as they wish, so any
// universe may same or not with other's, usually one universe only contains part of whole
// exist quantums, which is not self-conflict.
type Universe interface {
	// ReceiveQuantum will proccess quantums, verify the signature of each quantum. reject the quantums
	// from address in blacklist or conflict with quantums already exist. Both accept and wait quantums
	// will be saved. wait quantums usually cause by ref[0] is missing, so can not be broadcast and implement.
	ReceiveQuantum(originQuantums []*Quantum) (accept []Sig, wait []Sig, reject []Sig, err error)

	// ProcessSingleQuantum verify the signature, decide whether to accept or not, process the quantum by QType
	// return err if quantum not accept. casuse of verif-fail, signer in blacklist, conflict or any reason.
	ProcessSingleQuantum(sig Sig) error

	// ProcessQuantum have same function of ReceiveQuantum, but process quantums already in universe, which is
	// quantums from return "wait" of ReceiveQuantum.
	ProcessQuantum(skip, limit int) (accept []Sig, wait []Sig, reject []Sig, err error)

	// JudgeIndividual update the attitude towards Individual and how to process quantums from this signer.
	JudgeIndividual(address identity.Address, level int, judgment string, evidence ...[]Sig) error

	// JudgeCommunity update the attitude towards Community, decide if filter the individual in this community or not
	JudgeCommunity(sig Sig, level int, statement string) error

	// QueryQuantum query quantums from whole accepted quantums if address is nil, not filter by type if qType is 0
	QueryQuantum(address identity.Address, qType int, skip int, limit int, desc bool) ([]*Quantum, error)

	// QueryIndividual query Individual from whole universe if community sig is nil.
	QueryIndividual(sig Sig, skip int, limit int, desc bool) ([]*Individual, error)

	// GetCommunity return community by signature of community create signature.
	GetCommunity(sig Sig) (*Community, error)

	// GetIndividual return individual by address.
	GetIndividual(address identity.Address) (*Individual, error)

	// GetQuantum return quantum by signature of quantum.
	GetQuantum(sig Sig) (*Quantum, error)
}
