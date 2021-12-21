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
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/msg"
	"github.com/pdupub/go-pdu/p2p"
	"github.com/pdupub/go-pdu/params"
)

var (
	templatePath        string
	dbPath              string
	passwordPath        string
	projectPath         string
	saltPath            string
	references          string
	message             string
	imagePath           string
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
	rootCmd.AddCommand(AutoInitCmd())
	rootCmd.AddCommand(UploadCmd())
	rootCmd.AddCommand(SendMsgCmd())
	rootCmd.AddCommand(StartCmd())
	rootCmd.AddCommand(CreateKeystoreCmd())
	rootCmd.AddCommand(InitUniverseCmd())
	rootCmd.PersistentFlags().StringVar(&url, "url", "http://127.0.0.1", "target url")
	rootCmd.PersistentFlags().IntVar(&port, "port", params.DefaultPort, "port to start server or send msg")
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", params.DefaultDBPath, "path of database")
	rootCmd.PersistentFlags().StringVar(&projectPath, "projectPath", "./", "project root path")
	rootCmd.PersistentFlags().StringVar(&nodeList, "nodes", "", "node list")

	if err := rootCmd.Execute(); err != nil {
		return
	}
}

// StartCmd start run the node
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start run node",
		RunE: func(_ *cobra.Command, args []string) error {

			universe, err := core.NewUniverse()
			if err != nil {
				return err
			}

			entropy, err := core.NewEntropy()
			if err != nil {
				return err
			}
			universe.SetEntropy(entropy)

			society, err := core.NewSociety(core.GenesisRoots...)
			if err != nil {
				return err
			}
			universe.SetSociety(society)

			g := new(core.Genesis)
			g.SetUniverse(universe)
			p2p.New(templatePath, dbPath, g, port, ignoreUnknownSource, splitNodes())
			return nil
		},
	}
	cmd.Flags().BoolVar(&ignoreUnknownSource, "ignoreUS", true, "ignore unknown source")
	cmd.Flags().StringVar(&templatePath, "tpl", "p2p/public", "path of template")
	return cmd
}

func TestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test some tmp func",
		RunE: func(_ *cobra.Command, args []string) error {
			pMin := []int{0, 1, 1, 2, 2} // minimize parents cnt need to sign for creating ID on current generation
			cMax := []int{3, 1, 2, 2}    // maximize children cnt current generation can be envolve in creating
			gNum := []int{1, 3, 3}       // maximize ID cnt for current generation.

			// current_generation = max(parents_generation) + 1

			fmt.Println("pMin", pMin)
			fmt.Println("cMax", cMax)
			fmt.Println("gNum", gNum)

			gLimit := []*core.GenerationLimit{
				{ParentsMinSize: 0, ChildrenMaxSize: 3},
				{ParentsMinSize: 1, ChildrenMaxSize: 1},
				{ParentsMinSize: 1, ChildrenMaxSize: 2},
				{ParentsMinSize: 2, ChildrenMaxSize: 2},
				{ParentsMinSize: 2, ChildrenMaxSize: 2}}
			predictNum, err := core.CalcPopulation(gLimit, gNum)
			if err != nil {
				return err
			}
			fmt.Println("pNum", predictNum)
			return nil
		},
	}
	return cmd
}

// AutoInitCmd is tmp test for genesis
func AutoInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auto",
		Short: "Auto Initialize PDU",
		RunE: func(_ *cobra.Command, args []string) error {
			var err error
			// genesis generation struct
			ggs := make([]*core.GGInfo, 5)
			gNum := []int{1}
			for i, startPos := 0, 0; i < len(core.DefaultGLimit); i++ {
				ggs[i] = core.NewGGInfo()

				if i == 0 {
					ggs[i].CurN = gNum[0]
				} else {
					gNum, err := core.CalcPopulation(core.DefaultGLimit, gNum)
					if err != nil {
						return err
					}
					ggs[i].CurN = gNum[i]
				}

				ggs[i].Limit = core.DefaultGLimit[i]
				ggs[i].IDs = unlockTestDIDs(startPos, startPos+ggs[i].CurN)
				startPos += ggs[i].CurN
			}
			g, err := core.CreateGenesis(ggs)
			if err != nil {
				return err
			}
			p2p.New(templatePath, dbPath, g, port, true, splitNodes())
			return nil
		},
	}
	cmd.Flags().StringVar(&templatePath, "tpl", "p2p/public", "path of template")
	return cmd
}

// UploadCmd is used to upload file to node
func UploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload file to node",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			filename := args[0]

			// unlock did
			if keyIndex < 0 || keyIndex > 99 {
				return errors.New("index out of range")
			}
			testKeyfile := projectPath + params.TestKeystore(keyIndex)
			did, _ := identity.New()
			if err := did.UnlockWallet(testKeyfile, params.TestPassword); err != nil {
				return err
			}

			uploadResp, err := uploadFile(did, filename)
			if err != nil {
				return nil
			}
			fmt.Printf("%s:%d/%s\n", url, port, uploadResp.Path)
			fmt.Println(uploadResp.Hash)

			// curl -F "hash=b026324c6904b2a9cb4b88d6d61c81d1" -F "signature=nono" -F "file=@/Users/tataufo/tmp/pass" http://127.0.0.1:1323/upload
			return nil
		},
	}
	cmd.Flags().IntVar(&keyIndex, "key", 0, "index of key used")
	return cmd
}

func uploadFile(did *identity.DID, filename string) (*p2p.UploadResponse, error) {
	// create hash
	_, hashHex, err := p2p.HashFile(filename)
	if err != nil {
		return nil, err
	}

	// sign
	m := msg.New([]byte(hashHex))
	sm := msg.SignedMsg{Message: *m}
	if err := sm.Sign(did); err != nil {
		return nil, err
	}

	// build url
	urlUpload := fmt.Sprintf("%s:%d/upload", url, port)

	extraParams := map[string]string{
		"hash":      hashHex,
		"author":    did.GetKey().Address.Hex(),
		"signature": common.Bytes2Hex(sm.Signature),
	}

	request, err := newfileUploadRequest(urlUpload, extraParams, filename)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	uploadResp := new(p2p.UploadResponse)
	json.Unmarshal(body, uploadResp)

	// curl -F "hash=b026324c6904b2a9cb4b88d6d61c81d1" -F "signature=nono" -F "file=@/Users/tataufo/tmp/pass" http://127.0.0.1:1323/upload
	return uploadResp, nil
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
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

			images := strings.Split(imagePath, ",")
			ress := []*core.QRes{}
			for _, image := range images {
				if image != "" {
					var data, cs []byte
					res := new(core.QRes)
					if strings.HasPrefix(image, "http") {
						_, hashStr, err := p2p.HashFile(image)
						if err != nil {
							return err
						}
						cs, err = hex.DecodeString(hashStr)
						if err != nil {
							return err
						}
						res.Format = core.QResTypeImageU
						res.Data = []byte(image)
						res.Checksum = cs
					} else if resByData {
						data, _, err = p2p.HashFile(image)
						if err != nil {
							return err
						}
						// _, URL = path.Split(image)
						res.Format = core.QResTypeImageB
						res.Data = data
						res.Checksum = cs
					} else {
						resp, err := uploadFile(did, image)
						if err != nil {
							return err
						}
						URL := fmt.Sprintf("%s:%d/%s", url, port, resp.Path)
						cs, err = hex.DecodeString(resp.Hash)
						if err != nil {
							return err
						}
						res.Checksum = cs
						res.Format = core.QResTypeImageU
						res.Data = []byte(URL)
					}
					ress = append(ress, res)
				}
			}

			bp := new(core.Quantum)
			if len(args) > 0 {
				if args[0] == "ids" {
					addrs := params.TestAddrs()
					bp, err = core.NewBornQuantum(addrs[len(addrs)-1])
					if err != nil {
						return err
					}
					if err = bp.ParentSign(did); err != nil {
						return err
					}
					fmt.Println("req addr\t", addrs[len(addrs)-1].Hex())
					parents, err := bp.GetParents()
					if err != nil {
						return err
					}
					fmt.Println("req parent\t", parents[0].Hex())
				} else if args[0] == "pf" {
					var avatorRes *core.QRes = nil
					if len(ress) > 0 {
						avatorRes = ress[0]
					}
					bp, err = core.NewProfileQuantum(profileName, profileEmail, "hello world", profileUrl, "Earth", "", avatorRes)
					if err != nil {
						return err
					}
				}
			} else {
				bp, err = core.NewInfoQuantum(message, nil, ress...)
				if err != nil {
					return err
				}
			}

			content, err := json.Marshal(bp)
			if err != nil {
				return err
			}

			var refs [][]byte
			if references != "" {
				for _, r := range strings.Split(references, ",") {
					if len(r) == 130 {
						refs = append(refs, common.Hex2Bytes(r))
					} else {
						ref, err := base64.StdEncoding.DecodeString(r)
						if err != nil {
							return err
						}

						refs = append(refs, ref)
					}
				}
			} else {
				// request latest msg from same author
				resp, err := http.Get(fmt.Sprintf("%s:%d/info/latest/%s", url, port, did.GetKey().Address.Hex()))
				if err != nil {
					return err
				}

				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				resMsg := new(msg.SignedMsg)
				if err := json.Unmarshal(body, resMsg); err != nil {
					return err
				}
				if resMsg.Signature != nil {
					refs = append(refs, resMsg.Signature)
				}
			}
			m := msg.New(content, refs...)
			sm := msg.SignedMsg{Message: *m}
			if err := sm.Sign(did); err != nil {
				return err
			}
			fmt.Println("signature\t", common.Bytes2Hex(sm.Signature))
			resp, err := sm.Post(fmt.Sprintf("%s:%d", url, port))
			if err != nil {
				return err
			}

			fmt.Println("whole resp\t", string(resp))
			return nil
		},
	}
	cmd.Flags().IntVar(&keyIndex, "key", 0, "index of key used")
	cmd.Flags().StringVar(&message, "msg", "Hello World!", "content of msg send")
	cmd.Flags().StringVar(&imagePath, "image", "", "path of image")
	cmd.Flags().StringVar(&references, "refs", "", "references split by comma")
	cmd.Flags().StringVar(&profileName, "pname", "PDU-)", "profile name")
	cmd.Flags().StringVar(&profileEmail, "pemail", "hi@pdu.pub", "profile email")
	cmd.Flags().StringVar(&profileUrl, "purl", "https://pdu.pub", "profile url")
	cmd.Flags().BoolVar(&resByData, "rbd", false, "build resource by image data")
	return cmd
}

func InitUniverseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [file_path]",
		Short: "Initialize PDU",
		Long:  "Initialize PDU DAG & Keystores",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var n []byte
			var err error
			fmt.Printf("generation number: ")
			fmt.Scanln(&n)
			num, err := strconv.Atoi(string(n))
			if err != nil {
				return err
			}
			// genesis generation struct
			ggs := make([]*core.GGInfo, num)
			for i := 0; i < num; i++ {
				ggs[i] = core.NewGGInfo()
			}

			for i := 0; i < num-1; i++ {

				fmt.Println("generation ", i, " | use current generation as parents ")

				if i == 0 {
					fmt.Printf("init root number : ")
					fmt.Scanln(&n)
					ggs[i].CurN, err = strconv.Atoi(string(n))
					if err != nil {
						return err
					}
					ggs[0].Limit.ParentsMinSize = 0
				} else {
					fmt.Printf("parents number : ")
					fmt.Scanln(&n)
					ggs[i].Limit.ParentsMinSize, err = strconv.Atoi(string(n))
					if err != nil {
						return err
					}
				}

				fmt.Printf("children number: ")
				fmt.Scanln(&n)
				ggs[i].Limit.ChildrenMaxSize, err = strconv.Atoi(string(n))
				if err != nil {
					return err
				}
				if i > 0 {
					ggs[i].CurN = ggs[i-1].CurN * ggs[i-1].Limit.ChildrenMaxSize / ggs[i].Limit.ParentsMinSize
				}
			}

			idsSum := 0
			for _, gg := range ggs {
				idsSum += gg.CurN
			}

			pass, salt, err := getPassAndSalt()
			if err != nil {
				return err
			}
			keystorePath := path.Join(args[0], "keystore")
			dids, err := identity.CreateKeystoreArrayAndSaveLocal(keystorePath, pass, salt, idsSum)
			if err != nil {
				return err
			}

			lastN := 0
			for _, ggInfo := range ggs {
				ggInfo.IDs = dids[lastN : lastN+ggInfo.CurN]
				lastN = lastN + ggInfo.CurN
			}
			g, err := core.CreateGenesis(ggs)
			if err != nil {
				return err
			}

			p2p.New(templatePath, dbPath, g, port, true, splitNodes())
			return nil
		},
	}
	cmd.Flags().StringVar(&templatePath, "tpl", "p2p/public", "path of template")
	cmd.Flags().StringVar(&passwordPath, "pass", "", "path of password file")
	cmd.Flags().StringVar(&saltPath, "salt", "", "path of salt file")
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
