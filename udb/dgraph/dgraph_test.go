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

package dgraph

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/pdupub/go-pdu/udb"
)

// url & api token
const (
	testUrl   = "https://blue-surf-570228.us-east-1.aws.cloud.dgraph.io/graphql"
	testToken = "OTNiNjMyMTM1ODU3ZWE2NWU3YWQyNmM2MzgwZmViODg="
)

func TestUDBQuantum(t *testing.T) {
	cdgd, err := New(testUrl, testToken)
	if err != nil {
		t.Error(err)
	}

	defer cdgd.Close()

	if err := cdgd.initSchema(); err != nil {
		t.Error(err)
	}

	if err := cdgd.dropData(); err != nil {
		t.Error(err)
	}

	sig1 := "0x001"
	sig2 := "0x002"
	sig3 := "0x003"
	sender1 := "0xa01"
	sender2 := "0xa02"

	// test1 : quantum with empty refs
	quantum := &udb.Quantum{
		Sig:  sig3,
		Type: 1,
		Refs: []*udb.Quantum{{Sig: sig1}, {Sig: sig2}},
		Contents: []*udb.Content{{
			Fmt:  1,
			Data: "Hello World!",
		}},
		Sender: &udb.Individual{Address: sender1},
	}

	if _, _, _, err := sgcheck(cdgd, t, quantum, sender1); err != nil {
		t.Error(err)
	}

	// test2 : quantum fill empty quantum
	quantum = &udb.Quantum{
		Sig:  sig1,
		Type: 1,
		Contents: []*udb.Content{{
			Fmt:  1,
			Data: "For the Light!",
		}},
		Sender: &udb.Individual{Address: sender2},
	}

	if _, _, _, err := sgcheck(cdgd, t, quantum, sender2); err != nil {
		t.Error(err)
	}
	// test3 : quantum fill empty quantum
	quantum = &udb.Quantum{
		Sig:  sig2,
		Type: 2,
		Refs: []*udb.Quantum{{Sig: sig1}},
		Contents: []*udb.Content{{
			Fmt:  3,
			Data: "{\"nickname\":[\"Hello\",\"World\"]}",
		}},
		Sender: &udb.Individual{Address: sender1},
	}

	if _, _, _, err := sgcheck(cdgd, t, quantum, sender1); err != nil {
		t.Error(err)
	}

}

func sgcheck(cdgd *CDGD, t *testing.T, quantum *udb.Quantum, sender string) (dbQ *udb.Quantum, uid string, sid string, err error) {

	if uid, sid, err = cdgd.NewQuantum(quantum); err != nil {
		return dbQ, uid, sid, err
	} else {
		res, _ := json.Marshal(quantum)
		t.Log("SetQuantum", "uid", uid, "sid", sid, "sender", sender, "quantum", string(res))
	}

	if dbQ, err = cdgd.GetQuantum(quantum.Sig); err != nil {
		return dbQ, uid, sid, err
	} else {
		res, _ := json.Marshal(dbQ)
		t.Log("GetQuantum", "uid", dbQ.UID, "sid", sid, "sender", dbQ.Sender.Address, "quantum", string(res))

	}

	if uid != dbQ.UID {
		return dbQ, uid, sid, errors.New("uid not match")
	}
	return dbQ, uid, sid, err
}
