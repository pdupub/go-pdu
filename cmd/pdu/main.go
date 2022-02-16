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

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
)

var (
	templatePath        string
	passwordPath        string
	projectPath         string
	saltPath            string
	references          string
	message             string
	resByData           bool
	keyIndex            int
	port                int
	url                 string
	nodeList            string
	profileName         string
	profileEmail        string
	profileUrl          string
	ignoreUnknownSource bool
)

func main() {
	viper.New()
	rootCmd := &cobra.Command{
		Use:   "pdu",
		Short: "PDU command line interface (" + params.Version + ")",
		Long: `Parallel Digital Universe
	A decentralized identity-based social network
	Website: https://pdu.pub`,
	}

	rootCmd.AddCommand(TestCmd())
	rootCmd.AddCommand(StartCmd())
	rootCmd.AddCommand(SendMsgCmd())
	rootCmd.AddCommand(CreateKeystoreCmd())
	rootCmd.PersistentFlags().StringVar(&url, "url", "http://127.0.0.1", "target url")
	rootCmd.PersistentFlags().IntVar(&port, "port", params.DefaultPort, "port to start server or send msg")
	rootCmd.PersistentFlags().StringVar(&projectPath, "projectPath", "./", "project root path")

	if err := rootCmd.Execute(); err != nil {
		return
	}
}

// TestCmd run some test functions, just for developers
func TestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test some methods",
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Println("testing")
			return nil
		},
	}
	return cmd
}

// StartCmd start run the node
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start run node",
		RunE: func(_ *cobra.Command, args []string) error {

			// universe, err := core.NewUniverse()
			// if err != nil {
			// 	return err
			// }

			// entropy, err := core.NewEntropy()
			// if err != nil {
			// 	return err
			// }
			// universe.SetEntropy(entropy)

			// society, err := core.NewSociety(core.GenesisRoots...)
			// if err != nil {
			// 	return err
			// }
			// universe.SetSociety(society)

			// g := new(core.Genesis)
			// g.SetUniverse(universe)
			// p2p.New(templatePath, dbPath, g, port, ignoreUnknownSource, splitNodes())
			return nil
		},
	}
	cmd.Flags().BoolVar(&ignoreUnknownSource, "ignoreUS", true, "ignore unknown source")
	cmd.Flags().StringVar(&templatePath, "tpl", "p2p/public", "path of template")
	cmd.Flags().StringVar(&nodeList, "nodes", "", "node list")
	return cmd
}

// SendMsgCmd is send msg to node
func SendMsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send hello msg to node",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			if keyIndex < 0 || keyIndex > 99 {
				return errors.New("index out of range")
			}
			testKeyfile := projectPath + params.TestKeystore(keyIndex)

			did, _ := identity.New()
			if err := did.UnlockWallet(testKeyfile, params.TestPassword); err != nil {
				return err
			}

			fmt.Println("msg author\t", did.GetKey().Address.Hex())

			// ress := []*core.QRes{}

			// bp := new(core.Quantum)

			// bp, err = core.NewInfoQuantum(message, nil, ress...)
			// if err != nil {
			// 	return err
			// }

			// content, err := json.Marshal(bp)
			// if err != nil {
			// 	return err
			// }

			// var refs [][]byte
			// if references != "" {
			// 	for _, r := range strings.Split(references, ",") {
			// 		if len(r) == 130 {
			// 			refs = append(refs, common.Hex2Bytes(r))
			// 		} else {
			// 			ref, err := base64.StdEncoding.DecodeString(r)
			// 			if err != nil {
			// 				return err
			// 			}

			// 			refs = append(refs, ref)
			// 		}
			// 	}
			// } else {
			// 	// request latest msg from same author
			// 	resp, err := http.Get(fmt.Sprintf("%s:%d/info/latest/%s", url, port, did.GetKey().Address.Hex()))
			// 	if err != nil {
			// 		return err
			// 	}

			// 	defer resp.Body.Close()
			// 	body, err := ioutil.ReadAll(resp.Body)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	resMsg := new(msg.SignedMsg)
			// 	if err := json.Unmarshal(body, resMsg); err != nil {
			// 		return err
			// 	}
			// 	if resMsg.Signature != nil {
			// 		refs = append(refs, resMsg.Signature)
			// 	}
			// }
			// m := msg.New(content, refs...)
			// sm := msg.SignedMsg{Message: *m}
			// if err := sm.Sign(did); err != nil {
			// 	return err
			// }
			// fmt.Println("signature\t", common.Bytes2Hex(sm.Signature))
			// resp, err := sm.Post(fmt.Sprintf("%s:%d", url, port))
			// if err != nil {
			// 	return err
			// }

			// fmt.Println("whole resp\t", string(resp))
			return nil
		},
	}
	cmd.Flags().IntVar(&keyIndex, "key", 0, "index of key used")
	cmd.Flags().StringVar(&message, "msg", "Hello World!", "content of msg send")
	cmd.Flags().StringVar(&references, "refs", "", "references split by comma")
	cmd.Flags().StringVar(&profileName, "pname", "PDU-)", "profile name")
	cmd.Flags().StringVar(&profileEmail, "pemail", "hi@pdu.pub", "profile email")
	cmd.Flags().StringVar(&profileUrl, "purl", "https://pdu.pub", "profile url")
	cmd.Flags().BoolVar(&resByData, "rbd", false, "build resource by image data")
	return cmd
}

func CreateKeystoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ck [keystore_path]",
		Short: "Create keystores",
		Long:  "Create keystores by same pass & salt",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			pass, salt, err := getPassAndSalt()
			if err != nil {
				return err
			}
			dids, err := identity.CreateKeystoreArrayAndSaveLocal(args[0], pass, salt, 3)
			if err != nil {
				return err
			}
			for _, did := range dids {
				fmt.Println(did.GetKey().Address.Hex())
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&passwordPath, "pass", "", "path of password file")
	cmd.Flags().StringVar(&saltPath, "salt", "", "path of salt file")
	return cmd
}

func splitNodes() []string {
	nodes := []string{}
	for _, n := range strings.Split(nodeList, ",") {
		if len(n) == 0 {
			continue
		}
		if !strings.HasPrefix(n, "https://") && !strings.HasPrefix(n, "http://") {
			continue
		}
		nodes = append(nodes, n)
	}
	return nodes
}

func unlockTestDIDs(startPos, endPos int) []*identity.DID {
	var dids []*identity.DID
	for i := startPos; i < endPos; i++ {
		did := new(identity.DID)
		did.UnlockWallet(projectPath+params.TestKeystore(i), params.TestPassword)
		dids = append(dids, did)
	}
	return dids
}

func getPassAndSalt() (pass []byte, salt []byte, err error) {
	if len(passwordPath) > 0 {
		pass, err = ioutil.ReadFile(passwordPath)
		pass = []byte(strings.Replace(fmt.Sprintf("%s", pass), "\n", "", -1))
	} else {
		fmt.Printf("keystore pass: ")
		pass, err = gopass.GetPasswd()
	}
	if err != nil {
		return
	}

	if len(saltPath) > 0 {
		salt, err = ioutil.ReadFile(saltPath)
		salt = []byte(strings.Replace(fmt.Sprintf("%s", salt), "\n", "", -1))
	} else {
		fmt.Printf("keystore salt: ")
		salt, err = gopass.GetPasswd()
	}
	if err != nil {
		return
	}

	return
}
