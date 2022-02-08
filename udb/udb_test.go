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

package udb

import "testing"

// url & api token
const (
	testUrl        = "https://blue-surf-570176.us-east-1.aws.cloud.dgraph.io/graphql"
	testToken      = "OWJmYWM4NmM1ZjJlMGEyYWMyNGQ4NzU4Mjk3ZTI5ZDU="
	testAddressHex = "0x123"
)

func TestUDB(t *testing.T) {
	udb, err := New(testUrl, testToken)
	if err != nil {
		t.Error(err)
	}

	defer udb.Close()

	if err := udb.initIndividual(); err != nil {
		t.Error(err)
	}

	if err := udb.addIndividual(testAddressHex); err != nil {
		t.Error(err)
	}

	if res, err := udb.queryIndividual(testAddressHex); err != nil {
		t.Error(err)
	} else {
		for _, item := range res {
			t.Log("Address", item.Address, "DType", item.DType)
			for _, attr := range item.Profile {
				t.Log("Key", attr.Name, "Value", attr.Value)
			}
		}
	}

	if err := udb.dropData(); err != nil {
		t.Error(err)
	}
}
