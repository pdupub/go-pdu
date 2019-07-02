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
			var privKeyAdam, privKeyEve []*ecdsa.PrivateKey
			for i := 0; i < retryCnt; i++ {
				if Adam == nil {
					privKeyAdam, Adam, _ = createRootUser(true)
				}
				if Eve == nil {
					privKeyEve, Eve, _ = createRootUser(false)
				}
				if Adam != nil && Eve != nil {

					log.Println("Adam ID :", crypto.Hash2String(Adam.ID()))
					log.Println("private key start ", "#########")
					for _, v := range privKeyAdam {
						log.Println(crypto.Bytes2String(v.D.Bytes()))
					}
					log.Println("private key end ", "#########")
					log.Println("Eve ID  :", crypto.Hash2String(Eve.ID()))
					log.Println("private key start :", "#########")
					for _, v := range privKeyEve {
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
			var privKeys []interface{}
			for _, k := range privKeyAdam {
				privKeys = append(privKeys, k)
			}
			privKey := crypto.PrivateKey{
				Source:  Adam.Auth().Source,
				SigType: Adam.Auth().SigType,
				PriKey:  privKeys,
			}
			value := core.MsgValue{
				Content: []byte("hello world!"),
			}
			msg, err := core.CreateMsg(Adam, &value, &privKey)
			if err != nil {
				log.Println("create msg fail , err :", err)
			} else {
				log.Println("first msg from Adam ", "sender", crypto.Hash2String(msg.SenderID))
				log.Println("first msg from Adam ", "value.content", string(msg.Value.Content))
				log.Println("first msg from Adam ", "reference", msg.Reference)
				log.Println("first msg from Adam ", "signature", msg.Signature)
			}

			// verify msg
			if msg.SenderID == Adam.ID() {
				// public key should be found in user_dag by user id.
				msg.Signature.PubKey = Adam.Auth().PubKey
				res, err := core.VerifyMsg(*msg)
				if err != nil {
					log.Println("verfiy fail, err :", err)
				} else {
					log.Println("verify result is: ", res)
				}
			}

			// new msg reference first msg
			// create msg
			var privKeys2 []interface{}
			for _, k := range privKeyEve {
				privKeys2 = append(privKeys2, k)
			}
			privKey2 := crypto.PrivateKey{
				Source:  Eve.Auth().Source,
				SigType: Eve.Auth().SigType,
				PriKey:  privKeys2,
			}
			value2 := core.MsgValue{
				Content: []byte("hey u!"),
			}
			ref := core.MsgReference{Adam, msg.ID()}
			msg2, err := core.CreateMsg(Eve, &value2, &privKey2, &ref)
			if err != nil {
				log.Println("create msg fail , err :", err)
			} else {
				log.Println("first msg from Adam ", "sender", crypto.Hash2String(msg2.SenderID))
				log.Println("first msg from Adam ", "value.content", string(msg2.Value.Content))
				log.Println("first msg from Adam ", "reference", msg2.Reference)
				log.Println("first msg from Adam ", "signature", msg2.Signature)
			}

			// verify msg
			if msg.SenderID == Eve.ID() {
				// public key should be found in user_dag by user id.
				msg2.Signature.PubKey = Eve.Auth().PubKey
				res, err := core.VerifyMsg(*msg2)
				if err != nil {
					log.Println("verfiy fail, err :", err)
				} else {
					log.Println("verify result is: ", res)
				}
			}

			return nil
		},
	}

	return cmd
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
