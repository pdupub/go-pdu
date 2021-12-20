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
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/msg"
)

type GGInfo struct {
	Limit *GenerationLimit
	CurN  int
	IDs   []*identity.DID
}

func NewGGInfo() *GGInfo {
	return &GGInfo{Limit: &GenerationLimit{}}
}

type Genesis struct {
	universe *Universe
	msgs     []*msg.SignedMsg
}

func (g *Genesis) SetUniverse(universe *Universe) {
	g.universe = universe
}

func (g *Genesis) GetUniverse() *Universe {
	return g.universe
}

func (g *Genesis) GetMsgs() []*msg.SignedMsg {
	return g.msgs
}

func newGenesis() (*Genesis, error) {
	uni, err := NewUniverse()
	if err != nil {
		return nil, err
	}
	genesis := &Genesis{universe: uni}

	return genesis, nil
}

func CreateGenesis(ggs []*GGInfo) (*Genesis, error) {
	genesis, err := newGenesis()
	if err != nil {
		return nil, err
	}

	entropy, err := NewEntropy()
	if err != nil {
		return nil, err
	}
	genesis.universe.SetEntropy(entropy)

	var society *Society
	var refs [][]byte
	auSig := make(map[common.Address][][]byte)
	for i := range ggs {
		fmt.Println("")
		fmt.Println("----------- generation", i+1, "---------------")
		if i == 0 {
			society, err = NewSociety(GenesisRoots...)
			if err != nil {
				return nil, err
			}
			for _, id := range ggs[i].IDs {
				fmt.Println("v:", id.GetKey().Address.Hex())
			}
			genesis.universe.SetSociety(society)
		} else {
			ps := combine(ggs[i-1].IDs, ggs[i-1].Limit.ChildrenMaxSize, ggs[i].Limit.ParentsMinSize)
			for j, id := range ggs[i].IDs {
				// build Quantum
				bp, err := NewBornQuantum(id.GetKey().Address)
				if err != nil {
					return nil, err
				}
				fmt.Println("v:", id.GetKey().Address.Hex())
				for _, did := range ps[j] {
					if err = bp.ParentSign(did); err != nil {
						return nil, err
					}
					fmt.Println("---p:", did.GetKey().Address.Hex())
				}
				// build msg
				m := new(msg.Message)
				content, err := json.Marshal(bp)
				if err != nil {
					return nil, err
				}
				if len(refs) > 1 {
					m = msg.New(content, auSig[ps[j][0].GetKey().Address][len(auSig[ps[j][0].GetKey().Address])-1], refs[len(refs)-1])
				} else if len(refs) == 1 {
					m = msg.New(content, auSig[ps[j][0].GetKey().Address][len(auSig[ps[j][0].GetKey().Address])-1])
				} else {
					m = msg.New(content)
				}
				// sign msg
				sm := msg.SignedMsg{Message: *m}
				if err := sm.Sign(ps[j][0]); err != nil {
					return nil, err
				}
				fmt.Println("---s:", common.Bytes2Hex(sm.Signature))
				// msg create by id
				auSig[ps[j][0].GetKey().Address] = append(auSig[ps[j][0].GetKey().Address], sm.Signature)
				// msg create id
				auSig[id.GetKey().Address] = append(auSig[id.GetKey().Address], sm.Signature)
				// msg
				refs = append(refs, sm.Signature)

				if _, err := genesis.universe.ReceiveMsg(ps[j][0].GetKey().Address, sm.Signature, m.Content, m.References...); err != nil {
					return nil, err
				}

				genesis.msgs = append(genesis.msgs, &sm)
			}
		}
	}
	fmt.Println("----------------------------------------")

	return genesis, nil
}

func LoadGenesis() (*Genesis, error) {
	return newGenesis()
}

func transIDsToAddrs(IDs []*identity.DID) []common.Address {
	var addrs []common.Address
	for _, id := range IDs {
		addrs = append(addrs, id.GetKey().Address)
	}
	return addrs
}

func combine(dids []*identity.DID, cn, pn int) [][]*identity.DID {

	mCnt := make(map[common.Address]int)
	mDID := make(map[common.Address]*identity.DID)
	for _, did := range dids {
		mCnt[did.GetKey().Address] = cn
		mDID[did.GetKey().Address] = did
	}

	var cb [][]*identity.DID
	var updateStep int
	var addrRow []common.Address

MainLoop:
	for {
		for _, did := range dids {
			k := did.GetKey().Address

			if mCnt[k] > 0 {
				addrRow = append(addrRow, k)
				updateStep = 0
				mCnt[k] -= 1

				if len(addrRow) == pn {
					var didRow []*identity.DID
					for _, addr := range addrRow {
						didRow = append(didRow, mDID[addr])
					}
					cb = append(cb, didRow)
					addrRow = []common.Address{}
				}
			} else {
				updateStep += 1
			}
			if updateStep > len(mCnt)+1 {
				break MainLoop
			}
		}
	}

	return cb
}
