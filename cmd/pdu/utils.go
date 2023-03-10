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
)

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
