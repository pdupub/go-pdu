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

// Universe is interface contain all quantums which be received, selected and accepted by yourself.
// Your universe may same or not with other's, usually your universe only contains part of whole
// exist quantums (not conflict). By methods in Universe, communities be created by quantum and individuals
// be invited into community can be found. Universse also have some aggregate infomation on quantums.
type Universe interface {
	// ReceiveQuantum just receive origin quantums, not verify signature
	ReceiveQuantum(originQuantums []*Quantum) (receive []Sig, err error)

	// ProcessSingleQuantum verify the signature, decide whether to accept or not, process the quantum by QType
	// return err if quantum not accept. casuse of verif-fail, signer be punished, conflict or any reason from U)
	ProcessSingleQuantum(sig Sig) error

	// ProcessQuantum do RecvQuantum with more efficient way, err will not be return if quantum not be accepted.
	// signature of quantums in accept or rejected is processed.
	ProcessQuantum(skip, limit int) (accept []Sig, wait []Sig, reject []Sig, err error)

	// JudgeIndividual update judgement info of Individual and process all/part of quantums from this signer if necessary.
	JudgeIndividual(address identity.Address, level int, judgment string, evidence ...[]Sig) error

	// JudgeCommunity update the attitude of Community, filter the individual in this community or not
	JudgeCommunity(sig Sig, level int, statement string) error

	// QueryQuantum query quantums from whole accepted quantums if address is nil, not filter by type if qType is 0
	QueryQuantum(address identity.Address, qType int, skip int, limit int, desc bool) ([]*Quantum, error)

	// QueryIndividual query Individual from whole universe if community sig is nil.
	QueryIndividual(sig Sig, skip int, limit int, desc bool) ([]*Individual, error)

	// GetCommunity return nil if not exist
	GetCommunity(sig Sig) *Community

	// GetIndividual return nil if not exist
	GetIndividual(address identity.Address) *Individual

	// GetQuantum return nil if not exist
	GetQuantum(sig Sig) *Quantum
}
