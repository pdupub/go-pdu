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

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/common/log"
	"github.com/pdupub/go-pdu/core/rule"
	"github.com/pdupub/go-pdu/crypto"
	"testing"
)

var (
	err                                   error
	Adam, Eve                             *User
	priKeyAdam, priKeyEve                 *crypto.PrivateKey
	universe                              *Universe
	firstMsgIDFromAdam, firstMsgIDFromEve common.Hash
	ref                                   MsgReference
	AdamPartMsgIDs                        []common.Hash
	universeEngine                        crypto.Engine
)

func TestNewUniverse(t *testing.T) {
	universeEngine, _ = SelectEngine(crypto.ETH)

	// Test 1: Create root users, Adam and Eve , create universe
	// The gender of user relate to public key (random), so createRootUser
	// will repeat until two root user with different gender be created.
	// The new universe will be created by those two root users.
	// There is no message(msg) right now, so now space-time in universe.

	Adam, Eve, priKeyAdam, priKeyEve, err = createAdamAndEve()
	if err != nil {
		t.Error("create root user fail", err)
	}
	universe, err = NewUniverse(Eve, Adam)
	if err != nil {
		t.Error("create msg dag fail", err)
	}
	if res := universe.GetUserInfo(Adam.ID(), Adam.ID()); res != nil {
		t.Error("should be nil")
	}
	if priKeyAdam.Source != priKeyEve.Source {
		t.Error("private key confuse")
	}

	// Test 2: Create txt msg, first msg with no reference
	// this msg is signed by Adam
	value := MsgValue{
		ContentType: TypeText,
		Content:     []byte("hello world!"),
	}
	msg, err := CreateMsg(Adam, &value, priKeyAdam)
	if err != nil {
		t.Error("create msg fail , err :", err)
	} else if msg.Value.ContentType != TypeText {
		t.Error("first msg type should be TypeText")
	}
	firstMsgIDFromAdam = msg.ID()

	// Test 3: Add msg from Adam into universe, this msg will
	// be used to create default space-time in universe.
	// The initialize process of universe have been finished.
	// New space-time can be created later, but the universe
	// will always be this one. (equal to genesis.json in ethereum config)

	if err = universe.AddMsg(msg); err != nil {
		t.Error("add msg fail , err :", err)
	}
	if newAdam := universe.GetUserByID(Adam.ID()); newAdam == nil {
		t.Error("get Adam from userDAG fail")
	}

	// Test 4: Verify msg
	// msg contain the Adam is and signature.
	// Adam's public key can be found from userDAG by Adam ID
	if err := verifyMsg(msg); err != nil {
		t.Error(err)
	}

	// Test 5: Create lots of txt msg with reference, add into universe,
	// Each valid message will add into messageDAG in universe, and the ID of msg
	// will be keep in timespaceDAG for check reference and check DOB
	// In this step, the count of msg is more than tp need for bod.
	value2 := MsgValue{
		ContentType: TypeText,
		Content:     []byte("hey u!"),
	}
	ref := MsgReference{SenderID: Adam.ID(), MsgID: msg.ID()}
	msg2, err := CreateMsg(Eve, &value2, priKeyEve, &ref)
	firstMsgIDFromEve = msg2.ID()
	if err != nil {
		t.Error("create msg fail", err)
	}
	// add msg2
	if err := universe.AddMsg(msg2); err != nil {
		t.Error("add msg2 fail", err)
	}
	// verify msg
	if err := verifyMsg(msg2); err != nil {
		t.Error("verify msg2 fail", err)
	}
	// loop to add msg dag
	if err := loopAddMsg(universe, Adam, priKeyAdam, msg.ID()); err != nil {
		t.Error(err)
	}

	if maxSeq := universe.GetMaxSeq(Adam.ID()); maxSeq != rule.ReproductionInterval+1 {
		t.Error("max seq for time proof :", maxSeq)
	}

}

func TestUniverse_AddBODMsg(t *testing.T) {
	// Test 6: Create dob msg(create new user msg), and verify
	// new msg reference first & second msg
	valueDob := MsgValue{
		ContentType: TypeDOB,
	}
	_, pubKeyA2, err := universeEngine.GenKey(crypto.MultipleSignatures, 5)
	if err != nil {
		t.Error("generate public key fail", err)
	}

	auth := Auth{PublicKey: *pubKeyA2}
	content, err := CreateDOBMsgContent("A2", "1234", &auth)
	if err != nil {
		t.Error("create bod content fail", err)
	}
	content.SignByParent(Adam, *priKeyAdam)
	content.SignByParent(Eve, *priKeyEve)

	valueDob.Content, err = json.Marshal(content)
	if err != nil {
		t.Error("content marshal fail", err)
	}

	msgDob, err := CreateMsg(Eve, &valueDob, priKeyEve, &ref, &MsgReference{SenderID: Eve.ID(), MsgID: firstMsgIDFromEve})
	if err != nil {
		t.Error("create msg fail", err)
	}
	verifyMsg(msgDob)

	// Test 7: Marshal & unmarshal JSON for msg
	msgBytes, err := json.Marshal(msgDob)
	//log.Info(common.Bytes2String(msgBytes))
	var msgDob2 Message
	if err != nil {
		t.Error("marshal fail", err)
	} else {
		json.Unmarshal(msgBytes, &msgDob2)
		verifyMsg(&msgDob2)
		if msgBytes2, err := json.Marshal(msgDob2); err != nil || common.Bytes2String(msgBytes) != common.Bytes2String(msgBytes2) {
			t.Error("marshal or unmarshal fail", err)
		}
	}

	// verify the signature in the content of DOBMsg
	var dobContent DOBMsgContent
	json.Unmarshal(msgDob2.Value.Content, &dobContent)
	jsonBytes, _ := json.Marshal(dobContent.User)
	sigAdam := crypto.Signature{Signature: dobContent.Parents[1].Signature,
		PublicKey: universe.GetUserByID(dobContent.Parents[1].UserID).Auth.PublicKey}
	sigEve := crypto.Signature{Signature: dobContent.Parents[0].Signature,
		PublicKey: universe.GetUserByID(dobContent.Parents[0].UserID).Auth.PublicKey}

	if res, err := universeEngine.Verify(jsonBytes, &sigAdam); err != nil || res == false {
		t.Error("verify Adam fail", err)
	}
	if res, err := universeEngine.Verify(jsonBytes, &sigEve); err != nil || res == false {
		t.Error("verify Eve fail", err)
	}

	// Test 8: Create new User from dob message
	// user create from msg3 and msg4 should be same user
	if err := universe.AddMsg(msgDob); err != nil {
		t.Error("add msg3 fail", err)
	}
	if err := universe.AddMsg(&msgDob2); err != errMsgAlreadyExist {
		t.Errorf("add msg4 fail, err should be %s, but now err : %s", errMsgAlreadyExist, err)
	}

	if maxSeq := universe.GetMaxSeq(Eve.ID()); maxSeq != 0 {
		t.Error("max seq for Eve time proof, should be 0 :", maxSeq)
	}

}

func TestUniverse_AddSpaceTime(t *testing.T) {
	// Test 9: Create new space-time by msg from Eve
	ref = MsgReference{SenderID: Eve.ID(), MsgID: firstMsgIDFromEve}
	var msgNewST *Message
	for i := uint64(0); i < rule.ReproductionInterval; i++ {
		v := MsgValue{
			ContentType: TypeText,
			Content:     []byte(fmt.Sprintf("msg:%d", i)),
		}
		refT := MsgReference{SenderID: Adam.ID(), MsgID: AdamPartMsgIDs[i]}
		msgT, err := CreateMsg(Eve, &v, priKeyEve, &ref, &refT)
		if err != nil {
			t.Error("loop :", i, " err:", err)
		}
		if i == 50 {
			msgNewST = msgT
		}
		err = universe.AddMsg(msgT)
		if err != nil {
			log.Error("loop :", i, " err:", err)
		}
		ref = MsgReference{SenderID: Eve.ID(), MsgID: msgT.ID()}
		verifyMsg(msgT)
	}

	if err = universe.AddSpaceTime(msgNewST, msgNewST.Reference[0]); err != errCreateSpaceTimeFail {
		t.Error("should be err : ", errCreateSpaceTimeFail)
	}
	if err = universe.AddSpaceTime(msgNewST, msgNewST.Reference[1]); err != nil {
		t.Error("add space time fail, err :", err)
	} else if maxSeq := universe.GetMaxSeq(Eve.ID()); maxSeq == 0 {

		t.Error("max seq for Eve time proof, should be larger than 0 :", maxSeq)
	}

}

func TestUniverse_AddUserOnSpaceTime(t *testing.T) {
	// Test 10: Create lots of msg, add them into universe, user Eve
	// is valid in both space-time, so msg from Eve is valid in both space-time.
	for i := uint64(0); i < rule.ReproductionInterval; i++ {
		v := MsgValue{
			ContentType: TypeText,
			Content:     []byte(fmt.Sprintf("msg 2:%d", i)),
		}
		msgT, err := CreateMsg(Eve, &v, priKeyEve, &ref)
		if err != nil {
			log.Error("loop :", i, " err:", err)
		}
		err = universe.AddMsg(msgT)
		if err != nil {
			log.Error("loop :", i, " err:", err)
		}
		ref = MsgReference{SenderID: Eve.ID(), MsgID: msgT.ID()}
		verifyMsg(msgT)
	}
	if maxSeq := universe.GetMaxSeq(Eve.ID()); maxSeq == 0 {
		t.Error("max seq for Eve time proof, should be larger than 0 :", maxSeq)
	}

	// Test 11: Create new user, new user is valid in both
	// space-time with different life length left.
	valueDob := MsgValue{
		ContentType: TypeDOB,
	}
	_, pubKeyA3, err := universeEngine.GenKey(crypto.MultipleSignatures, 3)
	if err != nil {
		t.Error("generate public key fail", err)
	}

	auth := Auth{PublicKey: *pubKeyA3}
	content, err := CreateDOBMsgContent("A3", "789", &auth)
	if err != nil {
		t.Error("create bod content fail, err:", err)
	}
	content.SignByParent(Adam, *priKeyAdam)
	content.SignByParent(Eve, *priKeyEve)

	valueDob.Content, err = json.Marshal(content)
	if err != nil {
		t.Error("content marshal fail , err:", err)
	}
	refAdam := MsgReference{SenderID: Adam.ID(), MsgID: AdamPartMsgIDs[len(AdamPartMsgIDs)-1]}
	if msgDob, err := CreateMsg(Eve, &valueDob, priKeyEve, &ref, &refAdam); err != nil {
		t.Error("create msg fail , err :", err)
	} else if err := universe.AddMsg(msgDob); err != nil {
		log.Error("add user dob msg fail, err:", err)
	}

}

func loopAddMsg(universe *Universe, user *User, priKey *crypto.PrivateKey, lastMsgID common.Hash) error {
	ref = MsgReference{SenderID: user.ID(), MsgID: lastMsgID}
	// loop to add msg dag

	for i := uint64(0); i < rule.ReproductionInterval; i++ {
		v := MsgValue{
			ContentType: TypeText,
			Content:     []byte(fmt.Sprintf("msg:%d", i)),
		}
		msgT, err := CreateMsg(user, &v, priKey, &ref)
		if err != nil {
			return err
		}
		err = universe.AddMsg(msgT)
		if err != nil {
			return err
		}
		ref = MsgReference{SenderID: user.ID(), MsgID: msgT.ID()}
		AdamPartMsgIDs = append(AdamPartMsgIDs, msgT.ID())
		if err := verifyMsg(msgT); err != nil {
			return err
		}
	}
	return nil
}

func verifyMsg(msg *Message) error {
	// verify msg
	sender := universe.GetUserByID(msg.SenderID)
	if sender == nil {
		return errUserNotExist
	}
	msg.Signature.PubKey = sender.Auth.PubKey

	if verify, err := VerifyMsg(*msg); err != nil {
		return err
	} else if !verify {
		return errors.New("verify result is fail ")
	}
	return nil
}

func createAdamAndEve() (*User, *User, *crypto.PrivateKey, *crypto.PrivateKey, error) {
	retryCnt := 100
	var err error
	var Adam, Eve *User
	var privKeyAdam, privKeyEve *crypto.PrivateKey
	for i := 0; i < retryCnt; i++ {
		if Adam == nil {
			privKeyAdam, Adam, _ = createRootUser(true)
		}
		if Eve == nil {
			privKeyEve, Eve, _ = createRootUser(false)
		}
		if Adam != nil && Eve != nil {
			break
		}
	}
	return Adam, Eve, privKeyAdam, privKeyEve, err
}

func createRootUser(male bool) (*crypto.PrivateKey, *User, error) {
	keyCnt := 7
	if !male {
		keyCnt = 3
	}
	privKey, pubKey, err := universeEngine.GenKey(crypto.MultipleSignatures, keyCnt)
	if err != nil {
		return nil, nil, err
	}
	user := CreateRootUser(*pubKey, "name", "extra")
	if male == user.Gender() {
		return privKey, user, nil
	}
	return nil, nil, errCreateRootUserFail
}

func displayAllSpaceTimeUserState() {
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
