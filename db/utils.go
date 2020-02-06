// Copyright 2019 The PDU Authors
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

package db

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core"
)

var (
	// ErrMessageNotFound returns when the message not be found
	ErrMessageNotFound = errors.New("message can not be found")
)

func SaveRootUsers(udb UDB, users []*core.User) (err error) {
	// save root users
	var root0, root1 []byte
	if root0, err = json.Marshal(users[0]); err != nil {
		return err
	}
	if err = udb.Set(BucketConfig, ConfigRoot0, root0); err != nil {
		return err
	}
	if err = udb.Set(BucketUser, common.Hash2String(users[0].ID()), root0); err != nil {
		return err
	}

	if root1, err = json.Marshal(users[1]); err != nil {
		return err
	}

	if err = udb.Set(BucketConfig, ConfigRoot1, root1); err != nil {
		return err
	}

	if err = udb.Set(BucketUser, common.Hash2String(users[1].ID()), root1); err != nil {
		return err
	}

	if err := udb.Set(BucketConfig, ConfigCurrentStep, big.NewInt(StepRootsSaved).Bytes()); err != nil {
		return err
	}
	return nil
}

func GetRootUsers(udb UDB) (*core.User, *core.User, error) {
	var user0, user1 core.User
	var err error
	root0, err := udb.Get(BucketConfig, ConfigRoot0)
	if err != nil {
		return nil, nil, err
	}
	root1, err := udb.Get(BucketConfig, ConfigRoot1)
	if err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal(root0, &user0); err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal(root1, &user1); err != nil {
		return nil, nil, err
	}
	return &user0, &user1, nil
}

func SaveMsg(udb UDB, msg *core.Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	countBytes, err := udb.Get(BucketConfig, ConfigMsgCount)
	if err != nil {
		return err
	}
	count := new(big.Int).SetBytes(countBytes)
	err = udb.Set(BucketMsg, common.Hash2String(msg.ID()), msgBytes)
	if err != nil {
		return err
	}
	err = udb.Set(BucketMID, count.String(), common.Hash2Bytes(msg.ID()))
	if err != nil {
		return err
	}
	err = udb.Set(BucketMOD, common.Hash2String(msg.ID()), count.Bytes())
	if err != nil {
		return err
	}
	count = count.Add(count, big.NewInt(1))
	err = udb.Set(BucketConfig, ConfigMsgCount, count.Bytes())
	if err != nil {
		return err
	}

	err = udb.Set(BucketLastMID, common.Hash2String(msg.SenderID), common.Hash2Bytes(msg.ID()))
	if err != nil {
		return err
	}
	return nil
}

func GetLastMsg(udb UDB) (*core.Message, error) {
	var msg core.Message
	countBytes, err := udb.Get(BucketConfig, ConfigMsgCount)
	if err != nil {
		return nil, err
	}
	count := new(big.Int).SetBytes(countBytes)
	mid, err := udb.Get(BucketMID, count.Sub(count, big.NewInt(1)).String())
	if err != nil {
		return nil, err
	} else if mid == nil {
		return nil, ErrMessageNotFound
	}
	msgBytes, err := udb.Get(BucketMsg, common.Bytes2String(mid))
	if err != nil {
		return nil, err
	} else if msgBytes == nil {
		return nil, ErrMessageNotFound
	}
	err = json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func GetMsgByOrder(udb UDB, start *big.Int, size int) (msgs []*core.Message) {
	for ; size > 0; size-- {
		mid, err := udb.Get(BucketMID, start.String())
		if err != nil || mid == nil {
			continue
		}
		msgBytes, err := udb.Get(BucketMsg, common.Bytes2String(mid))
		if err != nil || msgBytes == nil {
			continue
		}
		var msg core.Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			continue
		}
		msgs = append(msgs, &msg)
		start = start.Add(start, big.NewInt(1))
	}
	return msgs
}

func GetOrderCntByMsg(udb UDB, mid common.Hash) (order *big.Int, count *big.Int, err error) {
	orderBytes, err := udb.Get(BucketMOD, common.Hash2String(mid))
	if err != nil {
		return nil, nil, err
	} else if orderBytes == nil {
		return nil, nil, ErrMessageNotFound
	}
	order = new(big.Int).SetBytes(orderBytes)

	count, err = GetMsgCount(udb)
	if err != nil {
		return nil, nil, err
	}
	return order, count, nil
}

func GetMsgCount(udb UDB) (count *big.Int, err error) {
	countBytes, err := udb.Get(BucketConfig, ConfigMsgCount)
	if err != nil {
		return nil, err
	}
	count = new(big.Int).SetBytes(countBytes)
	return count, nil
}

func GetLastMsgByUser(udb UDB, userID common.Hash) (*core.Message, error) {
	var msg core.Message
	lastMsgBytes, err := udb.Get(BucketLastMID, common.Hash2String(userID))
	if err != nil {
		return nil, err
	} else if lastMsgBytes == nil {
		return nil, ErrMessageNotFound
	}
	msgBytes, err := udb.Get(BucketMsg, common.Bytes2String(lastMsgBytes))
	if err != nil {
		return nil, err
	} else if msgBytes == nil {
		return nil, ErrMessageNotFound
	}
	err = json.Unmarshal(msgBytes, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
