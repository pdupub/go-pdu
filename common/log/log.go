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

const (
	// LvlError is log level error
	LvlError = iota
	// LvlWarn is log level warn
	LvlWarn
	// LvlInfo is log level info
	LvlInfo
	// LvlDebug is log level debug
	LvlDebug
	// LvlTrace is log level trace
	LvlTrace
)

func alignedString(lvl int) string {
	switch lvl {
	case LvlTrace:
		return "TRACE"
	case LvlDebug:
		return "DEBUG"
	case LvlInfo:
		return "INFO "
	case LvlWarn:
		return "WARN "
	case LvlError:
		return "ERROR"
	default:
		panic("bad level")
	}
}

func msgColor(lvl int) (color int) {
	switch lvl {
	case LvlError:
		color = 31
	case LvlWarn:
		color = 33
	case LvlInfo:
		color = 32
	case LvlDebug:
		color = 36
	case LvlTrace:
		color = 34
	default:
		panic("bad level")
	}
	return
}

// Info log in default color
func Info(v ...interface{}) {
	println(LvlInfo, v...)
}

// Trace log in green
func Trace(v ...interface{}) {
	println(LvlTrace, v...)
}

// Error log in red
func Error(v ...interface{}) {
	println(LvlError, v...)
}

// Warn log is yellow
func Warn(v ...interface{}) {
	println(LvlWarn, v...)
}

// Debug log is blue
func Debug(v ...interface{}) {
	println(LvlDebug, v...)
}

func println(lvl int, v ...interface{}) {
	fmt.Printf("%c[;;%dm", 0x1B, msgColor(lvl))
	v = append([]interface{}{alignedString(lvl)}, v...)
	log.Println(v...)
	fmt.Printf("%c[0m", 0x1B)
}
