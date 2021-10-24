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
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pdupub/go-pdu/contracts/poster"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

const (
	// PosterAddress is the address of
	PosterAddress = "0x0E90c6380a4b8696a65fF931AEad985C1C0eC76C"
	// RPC is the rpc of TTC Mainnet
	RPC = "http://35.229.26.170:8501"
)

type PosterRecord struct {
	Info     string
	Author   common.Address
	Deposit  *big.Int
	Start    *big.Int
	Update   *big.Int
	CoinsDay *big.Int
	Users    []common.Address
}

type Operator struct {
	key            *keystore.Key
	rpc            *rpc.Client
	chainID        *big.Int
	client         *ethclient.Client
	contractPoster *poster.Contract
}

func NewOperator() (*Operator, error) {
	operator := Operator{
		chainID: new(big.Int),
	}
	// dial rpc
	if client, err := rpc.Dial(RPC); err == nil {
		operator.rpc = client
		log.Println("Dial rpc success", "url", RPC)
	} else {
		return nil, err
	}
	// update chain id
	if err := operator.updateVersion(); err != nil {
		return nil, err
	}

	// init contract
	if err := operator.createContract(); err != nil {
		return nil, err
	}

	return &operator, nil
}

func (o *Operator) updateVersion() error {
	var response string
	if err := o.rpc.Call(&response, "net_version"); err != nil {
		return err

	}
	chainID, err := strconv.ParseUint(response, 10, 64)
	if err != nil {
		return err
	}
	o.chainID = new(big.Int).SetUint64(chainID)

	return nil
}

func (o *Operator) createContract() error {

	o.client = ethclient.NewClient(o.rpc)

	if contractPoster, err := poster.NewContract(common.HexToAddress(PosterAddress), o.client); err != nil {
		return err
	} else {
		o.contractPoster = contractPoster
	}
	return nil
}

func (o *Operator) SetKey(did *identity.DID) {
	o.key = did.GetKey()
}

func (o *Operator) GetNextRecordID() (*big.Int, error) {
	return o.contractPoster.NextRecordID(&bind.CallOpts{})
}

func (o *Operator) GetLastRecords(num int) (records []PosterRecord, err error) {
	nextRecordID, err := o.GetNextRecordID()
	if err != nil {
		return records, err
	}

	if nextRecordID.Cmp(big.NewInt(0)) <= 0 {
		return records, nil
	}

	for recordID := nextRecordID.Uint64() - 1; ; recordID-- {
		r := PosterRecord{}
		r.Info, r.Author, r.Deposit, r.Start, r.Update, r.CoinsDay, r.Users, err = o.contractPoster.GetRecord(&bind.CallOpts{}, new(big.Int).SetUint64(recordID))
		if err != nil {
			return records, err
		}
		records = append(records, r)
		if recordID == 0 || len(records) >= num {
			break
		}
	}
	return records, nil
}
