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

package udb

// Schema defines data struct by DQL, if current server accept unchecked refs, then
// create new quantum if sig in refs is not exist in system. or reject the quantum if not.
const Schema = `
community.base: uid .
community.invitations: [uid] .
community.maxInviteCnt: int .
community.members: [uid] .
community.minCosignCnt: int .
community.rule: uid .
content.data: string .
content.fmt: int @index(int) .
individual.address: string @index(hash) .
individual.communities: [uid] .
individual.quantums: [uid] .
quantum.contents: [uid] .
quantum.refs: [uid] .
quantum.sender: uid .
quantum.sig: string @index(hash) .
quantum.type: int @index(int) .
type community {
	community.base
	community.invitations
	community.maxInviteCnt
	community.members
	community.minCosignCnt
	community.rule
}
type content {
	content.data
	content.fmt
}
type individual {
	individual.address
	individual.communities
	individual.quantums
}
type quantum {
	quantum.contents
	quantum.refs
	quantum.sender
	quantum.sig
	quantum.type
}`
