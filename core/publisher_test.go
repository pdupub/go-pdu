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
	"testing"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

func TestNewPublisher(t *testing.T) {
	// Create a new DID for the publisher
	did, err := identity.New()
	if err != nil {
		t.Errorf("error creating new DID: %v", err)
	}

	// Unlock the DID
	err = did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)
	if err != nil {
		t.Errorf("error unlocking wallet: %v", err)
	}

	address := did.GetAddress()
	publisher := NewPublisher(address)

	if publisher.GetAddress() != address {
		t.Errorf("expected address %s, got %s", address, publisher.GetAddress())
	}

	if publisher.Attitude.Level != AttitudeAccept {
		t.Errorf("expected attitude level %d, got %d", AttitudeAccept, publisher.Attitude.Level)
	}
}

func TestPublisherUpsertProfile(t *testing.T) {
	// Create a new DID for the publisher
	did, err := identity.New()
	if err != nil {
		t.Errorf("error creating new DID: %v", err)
	}

	// Unlock the DID
	err = did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)
	if err != nil {
		t.Errorf("error unlocking wallet: %v", err)
	}

	address := did.GetAddress()
	publisher := NewPublisher(address)

	content := &QContent{Data: []byte("test data"), Format: "txt", Zipped: false}
	err = publisher.UpsertProfile("key1", content)
	if err != nil {
		t.Errorf("error upserting profile: %v", err)
	}

	if publisher.Profile["key1"] != content {
		t.Errorf("expected profile content %v, got %v", content, publisher.Profile["key1"])
	}

	err = publisher.UpsertProfile("key1", nil)
	if err != nil {
		t.Errorf("error upserting profile: %v", err)
	}

	if publisher.Profile["key1"] != nil {
		t.Errorf("expected profile content to be nil, got %v", publisher.Profile["key1"])
	}
}

func TestPublisherUpdateAttitude(t *testing.T) {
	// Create a new DID for the publisher
	did, err := identity.New()
	if err != nil {
		t.Errorf("error creating new DID: %v", err)
	}

	// Unlock the DID
	err = did.UnlockWallet("../"+params.TestKeystore(0), params.TestPassword)
	if err != nil {
		t.Errorf("error unlocking wallet: %v", err)
	}

	address := did.GetAddress()
	publisher := NewPublisher(address)

	newAttitude := &Attitude{
		Level:    AttitudeBroadcast,
		Judgment: "trusted publisher",
	}

	err = publisher.UpdateAttitude(newAttitude)
	if err != nil {
		t.Errorf("error updating attitude: %v", err)
	}

	if publisher.Attitude.Level != AttitudeBroadcast {
		t.Errorf("expected attitude level %d, got %d", AttitudeBroadcast, publisher.Attitude.Level)
	}

	if publisher.Attitude.Judgment != "trusted publisher" {
		t.Errorf("expected judgment %s, got %s", "trusted publisher", publisher.Attitude.Judgment)
	}
}
