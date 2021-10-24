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

package contracts

import (
	"math/big"
	"testing"
)

func TestNew(t *testing.T) {
	op, err := NewOperator()
	if err != nil {
		t.Error(err)
	}
	if op.chainID.Cmp(big.NewInt(8848)) != 0 {
		t.Error("chain ID is not 8848")
	}

	nextRecordID, err := op.GetNextRecordID()
	if err != nil {
		t.Error(err)
	}
	t.Log("nextRecordID", nextRecordID)

	records, err := op.GetLastRecords(5)
	if err != nil {
		t.Error(err)
	}
	for _, r := range records {
		t.Log(r)
	}
}
