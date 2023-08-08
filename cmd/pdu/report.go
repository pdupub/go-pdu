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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pdupub/go-pdu/core"
	"github.com/pdupub/go-pdu/identity"
	"github.com/pdupub/go-pdu/udb/fb"
	"github.com/spf13/cobra"
)

// ReportCmd is used to build report json
// report json is not part of pdu/core, is signed but not a quantum, this kind of msg
// is used for fit platform rule, such as AppStore ...

// the request is send by curl
// curl -X POST -H "Content-Type: application/json" -d '{"c":{"data": "0x123", "fmt": 6}, "sig":"..."}' <URL>

func ReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Create Report (For test)",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(_ *cobra.Command, args []string) (err error) {

			var currentDID *identity.DID

			// step 0: select the type of quantum and fill contents.
			r, err := initReport()
			if err != nil {
				return err
			}

			// step 1: unlock the test wallet
			if currentDID, err = unlockTestWallet(); err != nil {
				return err
			}

			// step 3: sign report.
			if err = r.Sign(currentDID); err != nil {
				return err
			}

			// display the information no matter which step got.
			displayReport(r)
			return nil
		},
	}
	return cmd
}
func initReport() (*fb.Report, error) {
	_, qType := multiChoice("please select the type of report", "Address", "Quantum")
	var content core.QContent
	switch qType {
	case 0:
		addrHex := question("please input the address hex", false)
		content = core.QContent{Data: []byte(addrHex), Format: core.QCFmtStringAddressHex}
	case 1:
		sigHex := question("please input the signature hex", false)
		content = core.QContent{Data: []byte(sigHex), Format: core.QCFmtStringSignatureHex}
	default:
		return nil, errors.New("unknown")
	}
	return &fb.Report{Content: &content}, nil
}

func displayReport(report *fb.Report) {
	fmt.Println("")
	fmt.Println("-------------------------------")
	fmt.Println()
	// display information here

	fmt.Println("sig \t", core.Sig2Hex(report.Signature))

	qBytes, err := json.Marshal(report)
	if err != nil {
		fmt.Println("display error ", err)
	} else {
		fmt.Println(string(qBytes))
	}

	fmt.Println()
	fmt.Println("-------------------------------")
	fmt.Println("")
}
