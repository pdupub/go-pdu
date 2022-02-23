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

	"github.com/pdupub/go-pdu/params"
	"github.com/pdupub/go-pdu/udb"
)

func TestUDBQuantum(t *testing.T) {
	cdgd, err := New(params.TestDGraphURL, params.TestDGraphToken)
	if err != nil {
		t.Error(err)
	}

	defer cdgd.Close()

	if err := cdgd.SetSchema(); err != nil {
		t.Error(err)
	}

	if err := cdgd.DropData(); err != nil {
		t.Error(err)
	}

	sig1 := "0x001" // used by test 2
	sig2 := "0x002" // used by test 3
	sig3 := "0x003" // used by test 1
	sig4 := "0x004" // used by test 5
	sender1 := "0xa01"
	sender2 := "0xa02"
	sender3 := "0xa03"

	// test1 : quantum with empty refs
	t.Log("----------------------------------------------------------------------------------")
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
	t.Log("----------------------------------------------------------------------------------")
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
	t.Log("----------------------------------------------------------------------------------")
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

	_, _, sid, err := sgcheck(cdgd, t, quantum, sender1)
	if err != nil {
		t.Error(err)
	}

	dbi, err := cdgd.GetIndividual(sender1)
	if err != nil {
		t.Error(err)
	}

	t.Log(dbi)
	individual := &udb.Individual{
		UID: sid,
	}
	t.Log(individual.Address)
	dbi, err = cdgd.GetIndividual(sender1)
	if err != nil {
		t.Error(err)
	}

	t.Log(dbi)

	// test 4 query quantum by user
	t.Log("----------------------------------------------------------------------------------")
	quantums, err := cdgd.QueryQuantum(sender1, 0, 1, 10, false)
	if err != nil {
		t.Error(err)
	}
	for _, q := range quantums {
		res, _ := json.Marshal(q)
		t.Log("quantum", string(res))
	}

	// test 5 add & get community
	// direct save as community, not by quantum
	t.Log("----------------------------------------------------------------------------------")
	quantum = &udb.Quantum{
		Sig:  sig4,
		Type: 3,
		Refs: []*udb.Quantum{{Sig: sig2}, {Sig: sig3}},
		Contents: []*udb.Content{{
			Fmt:  1,
			Data: "First Community", // note
		}, {
			Fmt:  33,
			Data: "", // base community
		}, {
			Fmt:  4,
			Data: "1", // minCosignCnt
		}, {
			Fmt:  4,
			Data: "2", // maxInviteCnt
		}, {
			Fmt:  6,
			Data: sender2, // initMembers
		}, {
			Fmt:  6,
			Data: sender3, //  initMembers
		},
		},
		Sender: &udb.Individual{Address: sender1},
	}

	_, _, _, err = sgcheck(cdgd, t, quantum, sender1)
	if err != nil {
		t.Error(err)
	}

	community := &udb.Community{
		MaxInviteCnt: 2,
		MinCosignCnt: 1,
		Define:       &udb.Quantum{Sig: sig4},
		InitMembers:  []*udb.Individual{{Address: sender2}, {Address: sender3}},
	}
	cid, err := cdgd.NewCommunity(community)
	if err != nil {
		t.Error(err)
	}

	t.Log("cid", cid)
	community2, err := cdgd.GetCommunity(sig4)
	if err != nil {
		t.Error(err)
	}
	c2, _ := json.Marshal(community2)

	t.Log(string(c2))
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
