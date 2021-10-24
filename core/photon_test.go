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
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestPhoton(t *testing.T) {

	msg := "hello world!"
	photon, err := NewPhoton(PhotonTypeInfo, []byte(msg))
	if err != nil {
		t.Error(err)
	}
	if string(photon.Data) != msg {
		t.Error("photon data not correct")
	}
	if photon.Type != PhotonTypeInfo {
		t.Error("photon type not correct")
	}
	if photon.Version != PhotonVersion {
		t.Error("photon version not correct")
	}

	b, err := json.Marshal(photon)
	if err != nil {
		t.Error(err)
	}

	pt := new(Photon)
	err = json.Unmarshal(b, pt)
	if err != nil {
		t.Error(err)
	}

	if string(pt.Data) != string(photon.Data) {
		t.Error("data not match")
	}
}

func TestBornPhoton(t *testing.T) {
	addr := common.HexToAddress("0xDa6bdC0Cd00fbaB9B33D1B4370fb32B8f6331376")
	photon, err := NewBornPhoton(addr)
	if err != nil {
		t.Error(err)
	}

	did0, _ := identity.New()
	if err := did0.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword); err != nil {
		t.Error(err)
	}
	t.Log(did0.GetKey().Address.Hex())
	if err := photon.ParentSign(did0); err != nil {
		t.Error(err)
	}

	did1, _ := identity.New()
	if err := did1.UnlockWallet("../"+params.TestKeystore(1), params.TestPassword); err != nil {
		t.Error(err)
	}
	t.Log(did1.GetKey().Address.Hex())

	if err := photon.ParentSign(did1); err != nil {
		t.Error(err)
	}

	if parents, err := photon.GetParents(); err != nil {
		t.Error(err)
	} else if len(parents) != 2 {
		t.Error(errors.New("parents number not correct"))
	}

	b, err := json.Marshal(photon)
	if err != nil {
		t.Error(err)
	}

	pt := new(Photon)
	err = json.Unmarshal(b, pt)
	if err != nil {
		t.Error(err)
	}

	if pt.Type != PhotonTypeBorn {
		t.Error("photon type not correct")
	}
	if pt.Version != PhotonVersion {
		t.Error("photon version not correct")
	}

	pb := new(PBorn)
	if err := json.Unmarshal(pt.Data, pb); err != nil {
		t.Error(err)
	}
	if pb.Addr != addr {
		t.Error("target address not match")
	}
}

func TestProfilePhoton(t *testing.T) {
	photon, err := NewProfilePhoton("PDU", "hi@pdu.pub", "hello world!", "https://pdu.pub", "Earth", "", nil)
	if err != nil {
		t.Error(err)
	}

	b, err := json.Marshal(photon)
	if err != nil {
		t.Error(err)
	}

	pt := new(Photon)
	err = json.Unmarshal(b, pt)
	if err != nil {
		t.Error(err)
	}

	if pt.Type != PhotonTypeProfile {
		t.Error("photon type not correct")
	}
	if pt.Version != PhotonVersion {
		t.Error("photon version not correct")
	}

	pb := new(PProfile)
	if err := json.Unmarshal(pt.Data, pb); err != nil {
		t.Error(err)
	}
	if pb.Name != "PDU" {
		t.Error("name not match")
	}
}
