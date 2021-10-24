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

// PhotonVersion is version of photon, use for parse data
const PhotonVersion = 1

// ResImage is type of resource (PIRes)
const ResImage = 1

const (
	// PhotonTypeInfo is information photon
	PhotonTypeInfo = iota
	// PhotonTypeBorn is create user photon
	PhotonTypeBorn
	// PhotonTypeProfile is update user profile
	PhotonTypeProfile
)

// Photon is main struct
type Photon struct {
	Type    int    `json:"t"`
	Version int    `json:"v"`
	Data    []byte `json:"d"`
}

// PIRes is struct of resource, image, video, audio ...
type PIRes struct {
	Format   int    `json:"format"`
	Data     []byte `json:"data"`
	URL      string `json:"url"`
	Checksum []byte `json:"cs"` // sha256
}

// PInfo is photon information struct
type PInfo struct {
	Text      string   `json:"text"`
	Quote     []byte   `json:"quote"`
	Resources []*PIRes `json:"resources"`
}

// PBorn is photon born struct
type PBorn struct {
	Addr       common.Address `json:"addr"`
	Signatures [][]byte       `json:"sigs"`
}

// PProfile is photon profile struct
type PProfile struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	URL      string `json:"url"`
	Location string `json:"location"`
	Avatar   *PIRes `json:"avatar"`
	Extra    string `json:"extra"`
}

// NewPhoton is used to create Photon
func NewPhoton(pType int, sData interface{}) (*Photon, error) {
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

	photon := Photon{
		Type:    pType,
		Version: PhotonVersion,
		Data:    data,
	}
	return &photon, nil
}

// NewInfoPhoton is used to create PhotonTypeInfo Photon
func NewInfoPhoton(text string, quote []byte, res ...*PIRes) (*Photon, error) {
	pb := PInfo{
		Text:      text,
		Quote:     quote,
		Resources: res,
	}
	return NewPhoton(PhotonTypeInfo, pb)
}

// NewProfilePhoton is used to create PhotonTypeProfile Photon
func NewProfilePhoton(name, email, bio, url, location, extra string, avatar *PIRes) (*Photon, error) {
	pb := PProfile{
		Name:     name,
		Email:    email,
		Bio:      bio,
		URL:      url,
		Location: location,
		Avatar:   avatar,
		Extra:    extra,
	}

	return NewPhoton(PhotonTypeProfile, pb)
}

// GetProfile return profile information
func (p *Photon) GetProfile() (*PProfile, error) {
	if p.Type != PhotonTypeProfile {
		return nil, ErrPhotonTypeNotCorrect
	}
	pp := new(PProfile)
	if err := json.Unmarshal(p.Data, pp); err != nil {
		return nil, err
	}
	return pp, nil
}

// NewBornPhoton create a create user Photon
func NewBornPhoton(target common.Address) (*Photon, error) {
	pb := PBorn{Addr: target}
	return NewPhoton(PhotonTypeBorn, pb)
}

// GetNewBorn get create individual address from photon
func (p *Photon) GetNewBorn() (common.Address, error) {
	if p.Type != PhotonTypeBorn {
		return common.Address{}, ErrPhotonTypeNotCorrect
	}

	pb := new(PBorn)
	if err := json.Unmarshal(p.Data, pb); err != nil {
		return common.Address{}, err
	}
	return pb.Addr, nil
}

// ParentSign is used to sign a PhotonTypeBorn Photon as parent
func (p *Photon) ParentSign(did *identity.DID) error {
	if p.Type != PhotonTypeBorn {
		return ErrPhotonTypeNotCorrect
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

// GetParents return parents of new create individual only if Photon is PhotonTypeBorn
func (p *Photon) GetParents() (parents []common.Address, err error) {
	if p.Type != PhotonTypeBorn {
		return []common.Address{}, ErrPhotonTypeNotCorrect
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
