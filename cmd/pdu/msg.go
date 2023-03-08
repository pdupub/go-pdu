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
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pdupub/go-pdu/core"
	"github.com/spf13/cobra"
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
				if q, err = signQuantum(q); err != nil {
					return err
				}
				nextStep = boolChoice("upload to Firebase?")
			}
			// stpe 4: upload to firebase.
			if nextStep {
				if err = upload2FireBase(q); err != nil {
					return err
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

	return nil, nil
}

func addRefs(quantum *core.Quantum) (*core.Quantum, error) {

	return nil, nil
}

func signQuantum(quantum *core.Quantum) (*core.Quantum, error) {

	// if keyIndex < 0 || keyIndex > 99 {
	// 	return errors.New("index out of range")
	// }
	// testKeyfile := projectPath + params.TestKeystore(keyIndex)

	// did, _ := identity.New()
	// if err := did.UnlockWallet(testKeyfile, params.TestPassword); err != nil {
	// 	return err
	// }

	// fmt.Println("msg author\t", did.GetKey().Address.Hex())
	return nil, nil
}

func upload2FireBase(quantum *core.Quantum) error {
	return nil
}

func display(quantum *core.Quantum) {
	fmt.Println("")
	fmt.Println("-------------------------------")
	// display information here

	fmt.Println("-------------------------------")
	fmt.Println("")
}

func boolChoice(tip string) bool {
	for {
		answer := strings.ToLower(question(tip, false))
		if answer == "yes" || answer == "y" {
			return true
		} else if answer == "no" || answer == "n" {
			return false
		}
	}
}

func multiChoice(tip string, targets ...string) string {
	if len(targets) == 0 {
		return ""
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
			return targets[intVar-1]
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
