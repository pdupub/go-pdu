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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var (
	references string
	message    string
	keyIndex   int
)

// MsgCmd is used to create, add references, sign and upload Message
func MsgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "msg",
		Short: "Operations on message",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			var currentDID *identity.DID
			// step 1: select the type of quantum and fill contents.
			q, err := initQuantum()
			if err != nil {
				return err
			}
			nextStep := boolChoice("continue to add references?")
			// step 2: add the references list
			if nextStep {
				if q, err = addRefs(q); err != nil {
					return err
				}
				nextStep = boolChoice("sign the quantum now?")
			}
			// step 3: unlock the private key and sign.
			if nextStep {
				if q, currentDID, err = signQuantum(q); err != nil {
					return err
				}
				nextStep = boolChoice("upload to Firebase?")
			}
			// stpe 4: upload to firebase.
			if nextStep {
				if json, err := upload2FireBase(q); err != nil {
					return err
				} else {
					// save to local
					record := make(map[string]interface{})
					record["sig"] = core.Sig2Hex(q.Signature)
					record["json"] = json
					addRecord(currentDID.GetAddress().Hex(), record)
				}
			}
			// display the information no matter which step got.
			display(q)
			return nil
		},
	}
	cmd.Flags().IntVar(&keyIndex, "key", 0, "index of key used")
	cmd.Flags().StringVar(&message, "msg", "Hello World!", "content of msg send")
	cmd.Flags().StringVar(&references, "refs", "", "references split by comma")
	return cmd
}

func initQuantum() (*core.Quantum, error) {
	_, qType := multiChoice("please select the type of message", "Information", "Profile", "Community Define", "Invitation", "End Account")

	switch qType {
	case core.QuantumTypeInfo:
		var qcs []*core.QContent
		text := question("please input the text content", true)
		qcs = append(qcs, core.CreateTextContent(text))
		imgUrl := question("please input the resource url (return if not)", false)
		if len(imgUrl) != 0 {
			qc, _ := core.NewContent(core.QCFmtStringURL, []byte(imgUrl))
			qcs = append(qcs, qc)
		}
		return core.CreateInfoQuantum(qcs, core.FirstQuantumReference)
	case core.QuantumTypeProfile:
		profiles := make(map[string]interface{})
		for i := 0; i < 6; i++ {
			k := question("please input the attribute name", false)
			if len(k) == 0 {
				break
			}
			v := question("please input the attribute value", false)
			if len(v) == 0 {
				break
			}
			profiles[k] = v
		}
		return core.CreateProfileQuantum(profiles, core.FirstQuantumReference)
	case core.QuantumTypeCommunity:
		note := question("note of community", true)
		minCosignCnt, err := strconv.Atoi(question("minimum count of cosign", false))
		if err != nil {
			return nil, err
		}
		maxInviteCnt, err := strconv.Atoi(question("maximum count of Invite", false))
		if err != nil {
			return nil, err
		}
		var initAddrsHex []string
		for i := 0; i < 7; i++ {
			addr := question("please input initialized address (return if not)", false)
			if len(addr) == 0 {
				break
			}
			initAddrsHex = append(initAddrsHex, addr)
		}
		return core.CreateCommunityQuantum(note, minCosignCnt, maxInviteCnt, initAddrsHex, core.FirstQuantumReference)
	case core.QuantumTypeInvitation:
		communityHex := question("please input the target community hex address", false)

		var addrsHex []string
		for i := 0; i < 7; i++ {
			addr := question("please input address to be invited (return if not)", false)
			if len(addr) == 0 {
				break
			}
			addrsHex = append(addrsHex, addr)
		}
		return core.CreateInvitationQuantum(core.Hex2Sig(communityHex), addrsHex, core.FirstQuantumReference)
	case core.QuantumTypeEnd:
		return core.CreateEndQuantum(core.FirstQuantumReference)
	}
	return nil, nil
}

func addRefs(quantum *core.Quantum) (*core.Quantum, error) {
	var refs []core.Sig
	selfRef := question("please input the self-reference(return if not)", false)
	if len(selfRef) != 0 {
		refs = append(refs, core.Hex2Sig(selfRef))
	} else {
		refs = append(refs, core.FirstQuantumReference)
	}
	for {
		ref := question("please input reference (return if not)", false)
		if len(ref) == 0 {
			break
		}
		refs = append(refs, core.Hex2Sig(ref))
	}

	quantum.References = refs
	return quantum, nil
}

func signQuantum(quantum *core.Quantum) (*core.Quantum, *identity.DID, error) {

	for {
		keyIndex, err := strconv.Atoi(question("please select the index from test key group (0~99)", false))
		if err != nil {
			continue
		}
		if keyIndex < 0 || keyIndex > 99 {
			continue
		}
		testKeyfile := projectPath + params.TestKeystore(keyIndex)
		did, _ := identity.New()

		if err := did.UnlockWallet(testKeyfile, params.TestPassword); err != nil {
			return nil, nil, err
		}
		fmt.Println("msg author is\t", did.GetKey().Address.Hex())

		if err = quantum.Sign(did); err != nil {
			return nil, nil, err
		}
		return quantum, did, nil
	}
}

func upload2FireBase(quantum *core.Quantum) ([]byte, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(projectPath + params.TestFirebaseAdminSDKPath)
	config := &firebase.Config{ProjectID: params.TestFirebaseProjectID}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	collection := client.Collection("quantum")

	docID := core.Sig2Hex(quantum.Signature)
	docRef := collection.Doc(docID)

	dMap := make(map[string]interface{})
	qBytes, err := json.Marshal(quantum)
	if err != nil {
		return nil, err
	}

	dMap["recv"] = qBytes
	if _, err = docRef.Set(ctx, dMap); err != nil {
		return nil, err
	}

	return qBytes, nil
}

func display(quantum *core.Quantum) {
	fmt.Println("")
	fmt.Println("-------------------------------")
	fmt.Println()
	// display information here

	fmt.Println("sig \t", core.Sig2Hex(quantum.Signature))

	qBytes, err := json.Marshal(quantum)
	if err != nil {
		fmt.Println("display error ", err)
	} else {
		fmt.Println(string(qBytes))
	}

	fmt.Println()
	fmt.Println("-------------------------------")
	fmt.Println("")
}

func boolChoice(tip string) bool {
	for {
		answer := strings.ToLower(question(tip+" (y/n) ", false))
		if answer == "yes" || answer == "y" {
			return true
		} else if answer == "no" || answer == "n" {
			return false
		}
	}
}

func multiChoice(tip string, targets ...string) (string, int) {
	if len(targets) == 0 {
		return "", -1
	}
	newTip := tip + " :"
	for i := 1; i < len(targets)+1; i++ {
		newTip += " " + strconv.Itoa(i) + ")" + targets[i-1]
	}
	newTip += "\t"
	for {
		answer := question(newTip, false)
		intVar, err := strconv.Atoi(answer)
		if err == nil && 0 < intVar && intVar <= len(targets) {
			return targets[intVar-1], intVar - 1
		}
	}
}

func question(tip string, isMultiple bool) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(tip)
	inputText := ""
	for {
		scanner.Scan()
		text := scanner.Text()
		if len(text) != 0 {
			inputText += text
			inputText += "\n"
		} else {
			break
		}
		if !isMultiple {
			break
		}
	}
	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}

	return strings.TrimSuffix(inputText, "\n")
}

func initRecords(addrHex string) (int, error) {

	viper.New()
	viper.SetConfigName("record_" + addrHex) // name of config file (without extension)
	viper.SetConfigType("json")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configPath)          // call multiple times to add many search paths

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Set("sseq", 0) // sequence for address-self
			viper.WriteConfig()
			return 0, nil
		} else {
			return 0, err
		}
	}

	return viper.GetInt("sseq"), nil
}

func loadRecords(addrHex string, limit int) ([]map[string]interface{}, error) {

	var records []map[string]interface{}
	sseq, err := initRecords(addrHex)
	if err != nil {
		return records, err
	}
	if sseq == 0 {
		return records, nil
	}

	for i := sseq; i > 0 && i > sseq-limit; i-- {
		record := viper.GetStringMap(strconv.Itoa(i))
		if record == nil {
			break
		}
		records = append(records, record)
	}
	return records, nil
}

func addRecord(addrHex string, record map[string]interface{}) error {
	sseq, err := initRecords(addrHex)
	if err != nil {
		return err
	}

	viper.Set(strconv.Itoa(sseq+1), record)
	viper.SafeWriteConfig()

	viper.Set("sseq", sseq+1)
	viper.WriteConfig()

	return nil
}
