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
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pdupub/go-pdu/identity"
)

const (
	// QResTypeImageB is type of resource is bytes data
	QResTypeImageB = iota
	// QResTypeImageU is type of resource is url address
	QResTypeImageU
)

const (
	// QuantumTypeInfo is information Quantum
	QuantumTypeInfo = iota
	// QuantumTypeAgree is
	QuantumTypeAgree
	// QuantumTypeDisagree is
	QuantumTypeDisagree
	// QuantumTypeCreate is create user Quantum
	QuantumTypeCreate
	// QuantumTypeProfile is update user profile
	QuantumTypeProfile
)

// type Quantum struct {
// 	id       []byte      // fill by signature when this quantum be received and validate
// 	Type     int         `json:"t"` // QuantumTypeInfo ...
// 	Contents []*QContent `json:"cs"`
// }

// type QContent struct {
// 	Format int    `json:"f"` // 0 < StringInfo, StringUrl < 32 < Bitmap, mp4 < 128 < address ,sig (as quote) ...
// 	Data   []byte `json:"d"`
// }

// Quantum is main struct
type Quantum struct {
	Type    int    `json:"t"`
	Version int    `json:"v"`
	Data    []byte `json:"d"`
}

// QRes is struct of resource, image, video, audio ...
type QRes struct {
	Format   int    `json:"format"`
	Data     []byte `json:"data"`
	Checksum []byte `json:"cs"` // sha256
}

// QData is Quantum information struct
type QData struct {
	Text      string  `json:"text"`
	Quote     []byte  `json:"quote"`
	Resources []*QRes `json:"resources"`
}

// PBorn is Quantum born struct
type PBorn struct {
	Addr       common.Address `json:"addr"`
	Signatures [][]byte       `json:"sigs"`
}

// NewQuantum is used to create Quantum
func NewQuantum(pType int, sData interface{}) (*Quantum, error) {
	var data []byte
	var err error
	switch sData.(type) {
	case []byte:
		data = sData.([]byte)
	default:
		data, err = json.Marshal(sData)
		if err != nil {
			return nil, err
		}
	}

	quantum := Quantum{
		Type: pType,
		Data: data,
	}
	return &quantum, nil
}

// NewInfoQuantum is used to create QuantumTypeInfo Quantum
func NewInfoQuantum(text string, quote []byte, res ...*QRes) (*Quantum, error) {
	pb := QData{
		Text:      text,
		Quote:     quote,
		Resources: res,
	}
	return NewQuantum(QuantumTypeInfo, pb)
}

// NewProfileQuantum is used to create QuantumTypeProfile Quantum
func NewProfileQuantum(name, email, bio, url, location, extra string, avatar *QRes) (*Quantum, error) {
	qp := map[string]*QData{
		"name":     {Text: name},
		"email":    {Text: email},
		"bio":      {Text: bio},
		"url":      {Text: url},
		"location": {Text: location},
		"extra":    {Text: extra},
		"avatar":   {Resources: []*QRes{avatar}}}
	return NewQuantum(QuantumTypeProfile, qp)
}

// GetProfile return profile information
func (p *Quantum) GetProfile() (map[string]*QData, error) {
	if p.Type != QuantumTypeProfile {
		return nil, ErrQuantumTypeNotCorrect
	}

	var qp map[string]*QData
	if err := json.Unmarshal(p.Data, &qp); err != nil {
		return nil, err
	}

	return qp, nil
}

// NewBornQuantum create a create user Quantum
func NewBornQuantum(target common.Address) (*Quantum, error) {
	pb := PBorn{Addr: target}
	return NewQuantum(QuantumTypeCreate, pb)
}

// GetNewBorn get create individual address from Quantum
func (p *Quantum) GetNewBorn() (common.Address, error) {
	if p.Type != QuantumTypeCreate {
		return common.Address{}, ErrQuantumTypeNotCorrect
	}

	pb := new(PBorn)
	if err := json.Unmarshal(p.Data, pb); err != nil {
		return common.Address{}, err
	}
	return pb.Addr, nil
}

// ParentSign is used to sign a QuantumTypeCreate Quantum as parent
func (p *Quantum) ParentSign(did *identity.DID) error {
	if p.Type != QuantumTypeCreate {
		return ErrQuantumTypeNotCorrect
	}

	pb := new(PBorn)
	if err := json.Unmarshal(p.Data, pb); err != nil {
		return err
	}

	hash := crypto.Keccak256(pb.Addr.Bytes())
	sig, err := crypto.Sign(hash, did.GetKey().PrivateKey)

	if err != nil {
		return err
	}
	pb.Signatures = append(pb.Signatures, sig)

	data, err := json.Marshal(pb)
	if err != nil {
		return err
	}
	p.Data = data
	return nil
}

// GetParents return parents of new create individual only if Quantum is QuantumTypeCreate
func (p *Quantum) GetParents() (parents []common.Address, err error) {
	if p.Type != QuantumTypeCreate {
		return []common.Address{}, ErrQuantumTypeNotCorrect
	}

	pb := new(PBorn)
	if err := json.Unmarshal(p.Data, pb); err != nil {
		return []common.Address{}, err
	}

	hash := crypto.Keccak256(pb.Addr.Bytes())
	for _, signature := range pb.Signatures {
		pubkey, err := crypto.Ecrecover(hash, signature)
		if err != nil {
			return []common.Address{}, err
		}
		signer := common.Address{}
		copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
		parents = append(parents, signer)
	}

	return parents, nil
}
