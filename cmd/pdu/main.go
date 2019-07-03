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
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/crypto"
	"github.com/pdupub/go-pdu/crypto/pdu"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	DefaultNodeHome = os.ExpandEnv("$HOME/.pdu")
)

func main() {

	viper.New()
	rootCmd := &cobra.Command{
		Use:   "pdu",
		Short: "PDU command line interface",
	}
	rootCmd.AddCommand(InitializeCmd())
	rootCmd.AddCommand(StartCmd())

	rootCmd.Execute()
}

func InitializeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [generation num]",
		Short: "Initialize the root user",
		RunE: func(_ *cobra.Command, args []string) error {
			// create root users
			retryCnt := 100
			var Adam, Eve *core.User
			var privKeyAdamGroup, privKeyEveGroup []*ecdsa.PrivateKey
			for i := 0; i < retryCnt; i++ {
				if Adam == nil {
					privKeyAdamGroup, Adam, _ = createRootUser(true)
				}
				if Eve == nil {
					privKeyEveGroup, Eve, _ = createRootUser(false)
				}
				if Adam != nil && Eve != nil {

					log.Println("Adam ID :", crypto.Hash2String(Adam.ID()))
					log.Println("private key start ", "#########")
					for _, v := range privKeyAdamGroup {
						log.Println(crypto.Bytes2String(v.D.Bytes()))
					}
					log.Println("private key end ", "#########")
					log.Println("Eve ID  :", crypto.Hash2String(Eve.ID()))
					log.Println("private key start :", "#########")
					for _, v := range privKeyEveGroup {
						log.Println(crypto.Bytes2String(v.D.Bytes()))
					}
					log.Println("private key end :", "#########")

					break
				}
			}

			// add root user into dag
			userDAG, err := core.NewUserDag(Eve, Adam)
			if err != nil {
				log.Println("create new user dag fail, err :", err)
			} else {
				log.Println("user dag :", userDAG)
			}

			// try get Adam
			newAdam := userDAG.GetUserByID(Adam.ID())
			if newAdam != nil {
				log.Println("get Adam from userDAG :", crypto.Hash2String(newAdam.ID()))
			}

			// create msg
			privKeyAdam := buildPrivateKey(privKeyAdamGroup)
			value := core.MsgValue{
				ContentType: core.TypeText,
				Content:     []byte("hello world!"),
			}
			msg, err := core.CreateMsg(Adam, &value, &privKeyAdam)
			if err != nil {
				log.Println("create msg fail , err :", err)
			} else {
				log.Println("first msg from Adam ", "sender", crypto.Hash2String(msg.SenderID))
				if msg.Value.ContentType == core.TypeText {
					log.Println("first msg from Adam ", "value.content", string(msg.Value.Content))
				}
				log.Println("first msg from Adam ", "reference", msg.Reference)
				log.Println("first msg from Adam ", "signature", msg.Signature)
			}

			// verify msg
			verifyMsg(userDAG, msg)

			// new msg reference first msg
			// create msg
			privKeyEve := buildPrivateKey(privKeyEveGroup)
			value2 := core.MsgValue{
				ContentType: core.TypeText,
				Content:     []byte("hey u!"),
			}
			ref := core.MsgReference{Sender: Adam, MsgID: msg.ID()}
			msg2, err := core.CreateMsg(Eve, &value2, &privKeyEve, &ref)
			if err != nil {
				log.Println("create msg fail , err :", err)
			} else {
				log.Println("first msg from Eve ", "sender", crypto.Hash2String(msg2.SenderID))
				if msg2.Value.ContentType == core.TypeText {
					log.Println("first msg from Eve ", "value.content", string(msg2.Value.Content))
				}
				log.Println("first msg from Eve ", "reference", msg2.Reference)
				log.Println("first msg from Eve ", "signature", msg2.Signature)
			}

			// verify msg
			verifyMsg(userDAG, msg2)

			// new msg reference first & second msg
			// create bod msg
			value3 := core.MsgValue{
				ContentType: core.TypeDOB,
			}
			// test, same with adam
			auth := core.Auth{PublicKey: Adam.Auth.PublicKey}
			content, err := core.CreateDOBMsgContent("A2", nil, &auth)
			if err != nil {
				log.Println("create bod content fail, err:", err)
			}
			content.SignByParent(privKeyAdam, true)
			content.SignByParent(privKeyEve, false)
			value3.Content, err = json.Marshal(content)
			log.Println()
			if err != nil {
				log.Println("content marshal fail , err:", err)
			}

			ref2 := core.MsgReference{Sender: Eve, MsgID: msg2.ID()}
			msg3, err := core.CreateMsg(Eve, &value3, &privKeyEve, &ref, &ref2)
			if err != nil {
				log.Println("create msg fail , err :", err)
			} else {
				log.Println("first dob msg ", "sender", crypto.Hash2String(msg3.SenderID))
				if msg3.Value.ContentType == core.TypeText {
					log.Println("first dob msg ", "value.content", string(msg3.Value.Content))
				} else if msg3.Value.ContentType == core.TypeDOB {
					log.Println("first dob msg ", "bod.content", string(msg3.Value.Content))
				}
				log.Println("first dob msg ", "reference", msg3.Reference)
				log.Println("first dob msg ", "signature", msg3.Signature)
			}

			verifyMsg(userDAG, msg3)

			msgBytes, err := json.Marshal(msg3)
			log.Println(crypto.Bytes2String(msgBytes))

			if err != nil {
				log.Println("marshal fail err :", err)
			} else {
				var msg4 core.Message
				err = json.Unmarshal(msgBytes, &msg4)
				if err != nil {
					log.Println("unmarshal fail err:", err)
				}
				msgBytes, err = json.Marshal(msg4)
				log.Println(crypto.Bytes2String(msgBytes))
				verifyMsg(userDAG, &msg4)
			}

			return nil
		},
	}

	return cmd
}

func buildPrivateKey(privKeyGroup []*ecdsa.PrivateKey) crypto.PrivateKey {
	var privKeys []interface{}
	for _, k := range privKeyGroup {
		privKeys = append(privKeys, k)
	}
	return crypto.PrivateKey{
		Source:  pdu.SourceName,
		SigType: pdu.MultipleSignatures,
		PriKey:  privKeys,
	}
}

func verifyMsg(userDAG *core.UserDAG, msg *core.Message) {

	// verify msg
	sender := userDAG.GetUserByID(msg.SenderID)
	if sender != nil {
		msg.Signature.PubKey = sender.Auth.PubKey
		res, err := core.VerifyMsg(*msg)
		if err != nil {
			log.Println("verfiy fail, err :", err)
		} else {
			log.Println("verify result is: ", res)
		}
	} else {
		log.Println("verify fail, err:", errors.New("user not exist in system"))
	}

}

func createRootUser(male bool) ([]*ecdsa.PrivateKey, *core.User, error) {
	keyCnt := 7
	if !male {
		keyCnt = 3
	}

	var privateKeyPool []*ecdsa.PrivateKey
	for i := 0; i < keyCnt; i++ {
		pk, err := pdu.GenerateKey()
		if err != nil {
			i--
			continue
		} else {
			privateKeyPool = append(privateKeyPool, pk)
		}
	}

	var pubKeys []interface{}
	for _, pk := range privateKeyPool {
		pubKeys = append(pubKeys, pk.PublicKey)
	}

	users, err := core.CreateRootUsers(crypto.PublicKey{Source: pdu.SourceName, SigType: pdu.MultipleSignatures, PubKey: pubKeys})
	if err != nil {
		return privateKeyPool, nil, err
	}

	if male && users[1] != nil {
		return privateKeyPool, users[1], nil
	} else if !male && users[0] != nil {
		return privateKeyPool, users[0], nil
	}
	return privateKeyPool, nil, errors.New("create root user fail")
}

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [port]",
		Short: "Start PDU Server on ",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			port := args[0]
			log.Println("listen on", port)

			return nil
		},
	}

	return cmd
}
