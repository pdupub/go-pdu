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

// QuantumVersion is version of Quantum, use for parse data
const QuantumVersion = 1

// ResImage is type of resource (QRes)
const ResImage = 1

const (
	// QuantumTypeInfo is information Quantum
	QuantumTypeInfo = iota
	// QuantumTypeBorn is create user Quantum
	QuantumTypeBorn
	// QuantumTypeProfile is update user profile
	QuantumTypeProfile
)

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
	URL      string `json:"url"`
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

// QProfile is Quantum profile struct
type QProfile struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	URL      string `json:"url"`
	Location string `json:"location"`
	Avatar   *QRes  `json:"avatar"`
	Extra    string `json:"extra"`
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
		Type:    pType,
		Version: QuantumVersion,
		Data:    data,
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
	pb := QProfile{
		Name:     name,
		Email:    email,
		Bio:      bio,
		URL:      url,
		Location: location,
		Avatar:   avatar,
		Extra:    extra,
	}

	return NewQuantum(QuantumTypeProfile, pb)
}

// GetProfile return profile information
func (p *Quantum) GetProfile() (*QProfile, error) {
	if p.Type != QuantumTypeProfile {
		return nil, ErrQuantumTypeNotCorrect
	}
	pp := new(QProfile)
	if err := json.Unmarshal(p.Data, pp); err != nil {
		return nil, err
	}
	return pp, nil
}

// NewBornQuantum create a create user Quantum
func NewBornQuantum(target common.Address) (*Quantum, error) {
	pb := PBorn{Addr: target}
	return NewQuantum(QuantumTypeBorn, pb)
}

// GetNewBorn get create individual address from Quantum
func (p *Quantum) GetNewBorn() (common.Address, error) {
	if p.Type != QuantumTypeBorn {
		return common.Address{}, ErrQuantumTypeNotCorrect
	}

	pb := new(PBorn)
	if err := json.Unmarshal(p.Data, pb); err != nil {
		return common.Address{}, err
	}
	return pb.Addr, nil
}

// ParentSign is used to sign a QuantumTypeBorn Quantum as parent
func (p *Quantum) ParentSign(did *identity.DID) error {
	if p.Type != QuantumTypeBorn {
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

// GetParents return parents of new create individual only if Quantum is QuantumTypeBorn
func (p *Quantum) GetParents() (parents []common.Address, err error) {
	if p.Type != QuantumTypeBorn {
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
