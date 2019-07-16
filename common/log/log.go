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

package log

import (
	"fmt"
	"log"
)

func Info(v ...interface{}) {
	log.Println(v...)
}

func Trace(v ...interface{}) {
	fmt.Printf("%c[;;%dm", 0x1B, 32)
	log.Println(v...)
	fmt.Printf("%c[0m", 0x1B)
}

func Error(v ...interface{}) {
	fmt.Printf("%c[;;%dm", 0x1B, 31)
	log.Println(v...)
	fmt.Printf("%c[0m", 0x1B)
}

func Split(info string) {

	Trace()
	Trace("#################################################################################################################")
	Trace("############################################## ", info, " #################################################")
	Trace("#################################################################################################################")
	Trace()
	Trace()
}
