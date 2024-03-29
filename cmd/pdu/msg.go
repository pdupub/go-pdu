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
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/params"
	"github.com/pdupub/go-pdu/udb/fb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MsgCmd is used to create, add references, sign and upload Message
func MsgCmd() *cobra.Command {
	var importQuantumJSON string

	cmd := &cobra.Command{
		Use:   "msg",
		Short: "Create and Broadcast Message (For test your own node)",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			var q *core.Quantum
			nextStep := false
			addrHex := ""
			if importQuantumJSON != "" {
				importQuantumJSON = strings.TrimSpace(importQuantumJSON)
				tmpQ := core.Quantum{}
				err = json.Unmarshal([]byte(importQuantumJSON), &tmpQ)
				if err != nil {
					return err
				}
				q = &tmpQ
				addr, err := q.Ecrecover()
				if err != nil {
					return err
				}
				addrHex = addr.Hex()
				nextStep = boolChoice("upload to Firebase?")
			} else {
				var currentDID *identity.DID

				// step 0: select the type of quantum and fill contents.
				q, err = initQuantum()
				if err != nil {
					return err
				}
				nextStep = boolChoice("unlock the wallet?")

				// step 1: unlock the test wallet
				if nextStep {
					if currentDID, err = unlockTestWallet(); err != nil {
						return err
					}
					addrHex = currentDID.GetAddress().Hex()
					nextStep = boolChoice("continue to add references?")
				}

				// step 2: add the references list
				if nextStep {
					if q, err = addRefs(q, currentDID); err != nil {
						return err
					}
					nextStep = boolChoice("sign the quantum now?")
				}

				// step 3: sign quantum.
				if nextStep {
					if err = q.Sign(currentDID); err != nil {
						return err
					}
					nextStep = boolChoice("upload to Firebase?")
				}
			}
			// stpe 4: upload to firebase.
			// In the reproduction environment, this step should
			// be performed on the mobile phone
			if nextStep {
				ctx := context.Background()
				fbu, err := fb.NewFBUniverse(ctx, firebaseKeyPath, firebaseProjectID)
				if err != nil {
					return err
				}

				_, _, reject, err := fbu.ReceiveQuantums([]*core.Quantum{q})
				if err != nil {
					return err
				}

				if len(reject) == 0 {
					jsonData, _ := json.Marshal(q)
					record := make(map[string]interface{})
					record["sig"] = core.Sig2Hex(q.Signature)
					record["json"] = jsonData
					addRecord(addrHex, record)
				}
			}
			// display the information no matter which step got.
			display(q)
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&importQuantumJSON, "quantumJson", "", "JSON of signed quantum, wrapped with '''")
	return cmd
}

func initQuantum() (*core.Quantum, error) {
	_, qType := multiChoice("please select the type of message", "Information", "Integration", "Speciation", "Identification", "Termination")

	switch qType {
	case core.QuantumTypeInformation:
		var qcs []*core.QContent
		text := question("please input the text content", true)
		qcs = append(qcs, core.CreateTextContent(text))
		imgUrl := question("please input the resource url (return if not)", false)
		if len(imgUrl) != 0 {
			qc, _ := core.NewContent(core.QCFmtStringURL, []byte(imgUrl))
			qcs = append(qcs, qc)
		}
		return core.CreateInfoQuantum(qcs, core.FirstQuantumReference)
	case core.QuantumTypeIntegration:
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
		return core.CreateIntegrationQuantum(profiles, core.FirstQuantumReference)
	case core.QuantumTypeSpeciation:
		note := question("note of species", true)
		minCosignCnt, err := strconv.Atoi(question("minimum count of cosign", false))
		if err != nil {
			return nil, err
		}
		maxIdentifyCnt, err := strconv.Atoi(question("maximum count of Identify", false))
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
		return core.CreateSpeciesQuantum(note, minCosignCnt, maxIdentifyCnt, initAddrsHex, core.FirstQuantumReference)
	case core.QuantumTypeIdentification:
		speciesHex := question("please input the target species hex address", false)

		var addrsHex []string
		for i := 0; i < 7; i++ {
			addr := question("please input address to be identified (return if not)", false)
			if len(addr) == 0 {
				break
			}
			addrsHex = append(addrsHex, addr)
		}
		return core.CreateIdentificationQuantum(core.Hex2Sig(speciesHex), addrsHex, core.FirstQuantumReference)
	case core.QuantumTypeTermination:
		return core.CreateTerminationQuantum(core.FirstQuantumReference)
	}
	return nil, nil
}

func addRefs(quantum *core.Quantum, did *identity.DID) (*core.Quantum, error) {
	var refs []core.Sig
	selfRef := ""

	fbu, err := fb.NewFBUniverse(context.Background(), firebaseKeyPath, firebaseProjectID)
	if err != nil {
		fmt.Println(err)
		ignoreRemoteErr := boolChoice("ignore the error above ?")
		if !ignoreRemoteErr {
			return nil, err
		}
	} else {
		currentIndividual, err := fbu.GetIndividual(did.GetAddress())
		if err != nil {
			fmt.Println(err)
			ignoreRemoteErr := boolChoice("ignore the error above ?")
			if !ignoreRemoteErr {
				return nil, err
			}
		} else {
			selfRef = core.Sig2Hex(currentIndividual.LastSig)
		}
	}

	//
	useLocalLastSig := false
	records, err := loadRecords(did.GetAddress().Hex(), 5)
	if err == nil && len(records) > 0 {
		// remote individual record not avaliable
		if len(selfRef) == 0 {
			useLocalLastSig = true
		}
		// remote record is not latest
		for i, record := range records {
			if i > 0 && record["sig"].(string) == selfRef {
				useLocalLastSig = true
			}
		}
	}

	if len(selfRef) > 0 && !useLocalLastSig {
		// append remote one to local
		record := make(map[string]interface{})
		record["sig"] = selfRef
		record["json"] = ""
		addRecord(did.GetAddress().Hex(), record)
	}

	if useLocalLastSig {
		selfRef = records[0]["sig"].(string)
	}

	if len(selfRef) == 0 {
		selfRef = question("please input the self-reference(return if not)", false)
	} else {
		fmt.Println("self reference is ", selfRef[:16]+"..."+selfRef[110:])
		if boolChoice("use another self-reference") {
			selfRef = question("please input new self-reference", false)
		}
	}

	if len(selfRef) != 0 {
		refs = append(refs, core.Hex2Sig(selfRef))
	} else {
		refs = append(refs, core.FirstQuantumReference)
	}

	if len(records) > 1 {
		fmt.Println("some reference can be used:")
		for i := len(records) - 1; i > 0; i-- {
			fmt.Println("-\t", records[i]["sig"])
		}
	}

	for {
		ref := question("please input other references (return if not)", false)
		if len(ref) == 0 {
			break
		}
		refs = append(refs, core.Hex2Sig(ref))
	}

	quantum.References = refs
	return quantum, nil
}

func unlockTestWallet() (*identity.DID, error) {
	var keyfilePath, password string
	for {
		useTestKey := boolChoice("do you want to use test keys?")
		if useTestKey {
			keyIndex, err := strconv.Atoi(question("please select the index from test key group (0~99)", false))
			if err != nil {
				continue
			}
			if keyIndex < 0 || keyIndex > 99 {
				continue
			}
			keyfilePath = projectPath + params.TestKeystore(keyIndex)
			password = params.TestPassword
		} else {
			keyfilePath = question("please input full path of your keystore", false)
			password = question("please input the password for your keystore", false)
		}
		did, _ := identity.New()
		if err := did.UnlockWallet(keyfilePath, password); err != nil {
			return nil, err
		}
		fmt.Println("msg author is\t", did.GetKey().Address.Hex())

		return did, nil
	}
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

func initRecords(addrHex string) (int, error) {

	viper.New()
	viper.SetConfigName("record_" + addrHex) // name of config file (without extension)
	viper.SetConfigType("json")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configPath)          // call multiple times to add many search paths

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Set("sseq", 0) // sequence for address-self
			viper.SafeWriteConfig()
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
	viper.Set("sseq", sseq+1)
	viper.WriteConfig()

	return nil
}
