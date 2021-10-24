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

/* build for mac
export CGO_ENABLED=1
export GOARCH=amd64
go build -buildmode=c-archive -o ../../client/PDUM/PDUM/pdu.a ./

"pdu.a" is missing one or more architectures required by this target: arm64
Adding "arm64" to Project -> Build Settings -> Excluded Architecture fixed the issue
*/

/* build for ios
export CGO_ENABLED=1
export GOOS=darwin
export GOARCH=arm64
export SDK=iphoneos
export CC=/usr/local/go/misc/ios/clangwrap.sh
### export CGO_CFLAGS="-fembed-bitcode"
go build -buildmode=c-archive -tags ios -o pdu.a ./

Build Settings -> Build Options -> Enable Bitcode : No
*/
