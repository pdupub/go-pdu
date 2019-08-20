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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/core/rule"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/pdu"
	"github.com/spf13/cobra"
)

var (
	errUserNotExist       = errors.New("user not exist in this system")
	errCreateRootUserFail = errors.New("create root user fail")
)

// TestCmd do the test related to step of create the pdu universe and
// split the universe for different space-time
func TestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test on pdu",
		RunE: func(_ *cobra.Command, args []string) error {

			// Test 1: Create root users, Adam and Eve , create universe
			// The gender of user relate to public key (random), so createRootUser
			// will repeat until two root user with different gender be created.
			// The new universe will be created by those two root users.
			// There is no message(msg) right now, so now space-time in universe.
			Adam, Eve, privKeyAdam, privKeyEve, err := createAdamAndEve()
			universe, err := core.NewUniverse(Eve, Adam)
			if err != nil {
				log.Error("create msg dag fail, err:", err)
			}
			if res := universe.GetUserInfo(Adam.ID(), Adam.ID()); res != nil {
				log.Error("should be nil")
			}
			log.Split("Test 1 finish")

			// Test 2: Create txt msg, first msg with no reference
			// this msg is signed by Adam
			value := core.MsgValue{
				ContentType: core.TypeText,
				Content:     []byte("hello world!"),
			}
			msg, err := core.CreateMsg(Adam, &value, privKeyAdam)
			if err != nil {
				log.Error("create msg fail , err :", err)
			} else {
				log.Info("first msg from Adam ", "sender", common.Hash2String(msg.SenderID))
				if msg.Value.ContentType == core.TypeText {
					log.Info("first msg from Adam ", "value.content", string(msg.Value.Content))
				}
				log.Info("first msg from Adam ", "reference", msg.Reference)
				//log.Info("first msg from Adam ", "signature", msg.Signature)
			}
			log.Split("Test 2 finish")

			// Test 3: Add msg from Adam into universe, this msg will
			// be used to create default space-time in universe.
			// The initialize process of universe have been finished.
			// New space-time can be created later, but the universe
			// will always be this one. (equal to genesis.json in ethereum config)
			err = universe.AddMsg(msg)
			if err != nil {
				log.Error("add msg fail , err :", err)
			} else {
				log.Info("msg dag add msg", common.Hash2String(universe.GetMsgByID(msg.ID()).ID()))
			}
			if newAdam := universe.GetUserByID(Adam.ID()); newAdam != nil {
				log.Info("get Adam from userDAG :", common.Hash2String(newAdam.ID()))
			}
			displayAllSpaceTimeUserState(universe)
			log.Split("Test 3 finish")

			// Test 4: Verify msg
			// msg contain the Adam is and signature.
			// Adam's public key can be found from userDAG by Adam ID
			verifyMsg(universe, msg, true)
			log.Split("Test 4 finish")

			// Test 5: Create lots of txt msg with reference, add into universe,
			// Each valid message will add into messageDAG in universe, and the ID of msg
			// will be keep in timespaceDAG for check reference and check DOB
			// In this step, the count of msg is more than tp need for bod.
			value2 := core.MsgValue{
				ContentType: core.TypeText,
				Content:     []byte("hey u!"),
			}
			ref := core.MsgReference{SenderID: Adam.ID(), MsgID: msg.ID()}
			msg2, err := core.CreateMsg(Eve, &value2, privKeyEve, &ref)
			if err != nil {
				log.Error("create msg fail , err :", err)
			} else {
				log.Info("first msg from Eve ", "sender", common.Hash2String(msg2.SenderID))
				if msg2.Value.ContentType == core.TypeText {
					log.Info("first msg from Eve ", "value.content", string(msg2.Value.Content))
				}
				log.Info("first msg from Eve ", "reference", msg2.Reference)
			}

			// add msg2
			if err := universe.AddMsg(msg2); err != nil {
				log.Error("add msg2 fail, err:", err)
			} else {
				log.Info("msg dag add msg2", common.Hash2String(universe.GetMsgByID(msg2.ID()).ID()))
			}

			// verify msg
			verifyMsg(universe, msg2, true)

			// loop to add msg dag
			ref = core.MsgReference{SenderID: Adam.ID(), MsgID: msg.ID()}
			var AdamPartMsgIDs []common.Hash
			for i := uint64(0); i < rule.REPRODUCTION_INTERVAL; i++ {
				v := core.MsgValue{
					ContentType: core.TypeText,
					Content:     []byte(fmt.Sprintf("msg:%d", i)),
				}
				msgT, err := core.CreateMsg(Adam, &v, privKeyAdam, &ref)
				if err != nil {
					log.Error("loop :", i, " err:", err)
				}
				err = universe.AddMsg(msgT)
				if err != nil {
					log.Error("loop :", i, " err:", err)
				}
				ref = core.MsgReference{SenderID: Adam.ID(), MsgID: msgT.ID()}
				AdamPartMsgIDs = append(AdamPartMsgIDs, msgT.ID())
				if i%(rule.REPRODUCTION_INTERVAL>>3) == 0 {
					log.Info("add ", i, "msgs")
				}
				verifyMsg(universe, msgT, false)
			}

			maxSeq := universe.GetMaxSeq(Adam.ID())
			log.Info("max seq for time proof :", maxSeq)
			displayAllSpaceTimeUserState(universe)
			log.Split("Test 5 finish")
			// Test 6: create dob msg, and verify
			// new msg reference first & second msg
			valueDob := core.MsgValue{
				ContentType: core.TypeDOB,
			}
			_, pubKeyA2, err := pdu.GenKey(pdu.MultipleSignatures, 5)

			auth := core.Auth{PublicKey: *pubKeyA2}
			content, err := core.CreateDOBMsgContent("A2", "1234", &auth)
			if err != nil {
				log.Error("create bod content fail, err:", err)
			}
			content.SignByParent(Adam, *privKeyAdam)
			content.SignByParent(Eve, *privKeyEve)

			valueDob.Content, err = json.Marshal(content)
			if err != nil {
				log.Error("content marshal fail , err:", err)
			}

			ref2 := core.MsgReference{SenderID: Eve.ID(), MsgID: msg2.ID()}
			msgDob, err := core.CreateMsg(Eve, &valueDob, privKeyEve, &ref, &ref2)
			if err != nil {
				log.Error("create msg fail , err :", err)
			} else {
				log.Info("first dob msg ", "sender", common.Hash2String(msgDob.SenderID))
				if msgDob.Value.ContentType == core.TypeText {
					log.Info("first dob msg ", "value.content", string(msgDob.Value.Content))
				} else if msgDob.Value.ContentType == core.TypeDOB {
					log.Info("first dob msg ", "bod.content", string(msgDob.Value.Content)[:60]+"...")
				}
				//log.Info("first dob msg ", "reference", msgDob.Reference)
				//log.Info("first dob msg ", "signature", msgDob.Signature)
			}

			verifyMsg(universe, msgDob, true)
			log.Split("Test 6 finish")

			// Test 7: json marshal & unmarshal for msg
			msgBytes, err := json.Marshal(msgDob)
			//log.Info(common.Bytes2String(msgBytes))

			var msgDob2 core.Message
			if err != nil {
				log.Error("marshal fail err :", err)
			} else {
				err = json.Unmarshal(msgBytes, &msgDob2)
				if err != nil {
					log.Error("unmarshal fail err:", err)
				}
				verifyMsg(universe, &msgDob2, true)
				msgBytes, err = json.Marshal(msgDob2)
				if err != nil {
					log.Error("marshal fail err:", err)
				}
				//log.Info(common.Bytes2String(msgBytes))
			}

			// verify the signature in the content of DOBMsg
			if msgDob2.Value.ContentType == core.TypeDOB {
				var dobContent core.DOBMsgContent
				err = json.Unmarshal(msgDob2.Value.Content, &dobContent)
				if err != nil {
					log.Error("dob message can not be unmarshl, err:", err)
				}

				jsonBytes, err := json.Marshal(dobContent.User)
				if err != nil {
					log.Error("user to json fail , err:", err)
				}

				sigAdam := crypto.Signature{Signature: dobContent.Parents[1].Signature,
					PublicKey: universe.GetUserByID(dobContent.Parents[1].UserID).Auth.PublicKey}
				sigEve := crypto.Signature{Signature: dobContent.Parents[0].Signature,
					PublicKey: universe.GetUserByID(dobContent.Parents[0].UserID).Auth.PublicKey}

				if res, err := pdu.Verify(jsonBytes, sigAdam); err != nil || res == false {
					log.Error("verify Adam fail, err", err)
				} else {
					log.Info("verify Adam true")
				}

				if res, err := pdu.Verify(jsonBytes, sigEve); err != nil || res == false {
					log.Info("verify Eve fail, err", err)
				} else {
					log.Info("verify Eve true")
				}
			} else {
				log.Error("should be dob msg")
			}
			log.Split("Test 7 finish ")
			// Test 8: create new User from dob message
			// user create from msg3 and msg4 should be same user

			if err := universe.AddMsg(msgDob); err != nil {
				log.Error("add msg3 fail , err", err)
			} else {
				log.Info("add msg3 success")
			}
			if err := universe.AddMsg(&msgDob2); err != core.ErrMsgAlreadyExist {
				log.Error("add msg4 fail, err should be %s, but now err : %s", core.ErrMsgAlreadyExist, err)
			}

			maxSeq = universe.GetMaxSeq(Eve.ID())
			log.Info("max seq for Eve time proof, should be 0 :", maxSeq)
			displayAllSpaceTimeUserState(universe)

			log.Split("Test 8 finish")

			ref = core.MsgReference{SenderID: Eve.ID(), MsgID: msg2.ID()}
			var msgNewST *core.Message
			for i := uint64(0); i < rule.REPRODUCTION_INTERVAL; i++ {
				v := core.MsgValue{
					ContentType: core.TypeText,
					Content:     []byte(fmt.Sprintf("msg:%d", i)),
				}
				refT := core.MsgReference{SenderID: Adam.ID(), MsgID: AdamPartMsgIDs[i]}
				msgT, err := core.CreateMsg(Eve, &v, privKeyEve, &ref, &refT)
				if err != nil {
					log.Error("loop :", i, " err:", err)
				}
				if i == 50 {
					msgNewST = msgT
				}
				err = universe.AddMsg(msgT)
				if err != nil {
					log.Error("loop :", i, " err:", err)
				}
				ref = core.MsgReference{SenderID: Eve.ID(), MsgID: msgT.ID()}
				verifyMsg(universe, msgT, false)
			}

			if err = universe.AddSpaceTime(msgNewST, msgNewST.Reference[0]); err != core.ErrCreateSpaceTimeFail {
				log.Error("should be err : ", core.ErrCreateSpaceTimeFail)
			}
			if err = universe.AddSpaceTime(msgNewST, msgNewST.Reference[1]); err != nil {
				log.Error("add space time fail, err :", err)
			} else {
				maxSeq = universe.GetMaxSeq(Eve.ID())
				log.Info("max seq for Eve time proof, should be larger than 0 :", maxSeq)
			}
			log.Split("Test 9 finish")

			for i := uint64(0); i < rule.REPRODUCTION_INTERVAL; i++ {
				v := core.MsgValue{
					ContentType: core.TypeText,
					Content:     []byte(fmt.Sprintf("msg 2:%d", i)),
				}
				msgT, err := core.CreateMsg(Eve, &v, privKeyEve, &ref)
				if err != nil {
					log.Error("loop :", i, " err:", err)
				}
				err = universe.AddMsg(msgT)
				if err != nil {
					log.Error("loop :", i, " err:", err)
				}
				ref = core.MsgReference{SenderID: Eve.ID(), MsgID: msgT.ID()}
				verifyMsg(universe, msgT, false)
			}
			maxSeq = universe.GetMaxSeq(Eve.ID())
			log.Info("max seq for Eve time proof, should be larger than 0 :", maxSeq)
			log.Split("Test 10 Finish")

			valueDob = core.MsgValue{
				ContentType: core.TypeDOB,
			}
			_, pubKeyA3, err := pdu.GenKey(pdu.MultipleSignatures, 3)

			auth = core.Auth{PublicKey: *pubKeyA3}
			content, err = core.CreateDOBMsgContent("A3", "789", &auth)
			if err != nil {
				log.Error("create bod content fail, err:", err)
			}
			content.SignByParent(Adam, *privKeyAdam)
			content.SignByParent(Eve, *privKeyEve)

			valueDob.Content, err = json.Marshal(content)
			if err != nil {
				log.Error("content marshal fail , err:", err)
			}
			refAdam := core.MsgReference{SenderID: Adam.ID(), MsgID: AdamPartMsgIDs[len(AdamPartMsgIDs)-1]}
			msgDob, err = core.CreateMsg(Eve, &valueDob, privKeyEve, &ref, &refAdam)
			if err != nil {
				log.Error("create msg fail , err :", err)
			} else {
				log.Info("dob msg ", "sender", common.Hash2String(msgDob.SenderID))
				log.Info("dob msg ", "bod.content", string(msgDob.Value.Content)[:60]+"...")

			}
			if err := universe.AddMsg(msgDob); err != nil {
				log.Error("add user dob msg fail, err:", err)
			}

			displayAllSpaceTimeUserState(universe)
			log.Split("Test 11 Finish")
			return nil
		},
	}

	return cmd
}

func displayAllSpaceTimeUserState(universe *core.Universe) {
	for _, stID := range universe.GetSpaceTimeIDs() {
		for _, id := range universe.GetUserIDs(stID) {
			if uInfo := universe.GetUserInfo(id, stID); uInfo != nil {
				log.Info("ST:", universe.GetUserByID(stID).Name, uInfo.String(), "userID:", common.Hash2String(id)[:5])
			} else {
				log.Error("can not find user info")
			}
		}
	}
}

func verifyMsg(universe *core.Universe, msg *core.Message, show bool) {
	// verify msg
	sender := universe.GetUserByID(msg.SenderID)
	if sender != nil {
		msg.Signature.PubKey = sender.Auth.PubKey
		res, err := core.VerifyMsg(*msg)
		if err != nil {
			log.Error("verfiy fail, err :", err)
		} else if show {
			log.Info("verify result is: ", res)
		}
	} else {
		log.Error("verify fail, err:", errUserNotExist)
	}

}

func createAdamAndEve() (*core.User, *core.User, *crypto.PrivateKey, *crypto.PrivateKey, error) {
	retryCnt := 100
	var err error
	var Adam, Eve *core.User
	var privKeyAdam, privKeyEve *crypto.PrivateKey
	for i := 0; i < retryCnt; i++ {
		if Adam == nil {
			privKeyAdam, Adam, _ = createRootUser(true)
		}
		if Eve == nil {
			privKeyEve, Eve, _ = createRootUser(false)
		}
		if Adam != nil && Eve != nil {
			log.Info("Adam ID :", common.Hash2String(Adam.ID()))
			log.Info("Eve ID  :", common.Hash2String(Eve.ID()))
			break
		}
	}
	return Adam, Eve, privKeyAdam, privKeyEve, err
}

func createRootUser(male bool) (*crypto.PrivateKey, *core.User, error) {
	keyCnt := 7
	if !male {
		keyCnt = 3
	}
	privKey, pubKey, err := pdu.GenKey(pdu.MultipleSignatures, keyCnt)
	if err != nil {
		return nil, nil, err
	}

	users, err := core.CreateRootUsers(*pubKey)
	if err != nil {
		return nil, nil, err
	}

	if male && users[1] != nil {
		return privKey, users[1], nil
	} else if !male && users[0] != nil {
		return privKey, users[0], nil
	}
	return nil, nil, errCreateRootUserFail
}
