// Copyright 2018 The PDU Authors
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

// gomobile bind -target=ios github.com/pdupub/go-pdu/cmd/pdum
// gomobile bind -target=android github.com/pdupub/go-pdu/cmd/pdum

// -target=adroid
// # runtime/cgo
// clang80: error: argument unused during compilation: '-stdlib=libc++' [-Werror,-Wunused-command-line-argument]

// vim /Users/tataufo/Library/Android/sdk/ndk-bundle/build/tools/make_standalone_toolchain.py
// :124
// #flags = '-target {} -stdlib=libc++'.format(target)
// flags = '-target {} '.format(target)

package pdum

func PAdd(a, b int) int {
	return a + b
}
