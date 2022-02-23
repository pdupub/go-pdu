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

package dgraph

// Schema defines data struct by DQL, if current server accept unchecked refs, then
// create new quantum if sig in refs is not exist in system. or reject the quantum if not.
const Schema = `
# -- common --
pdu.type: string @index(hash) .
# -- community --
community.note: uid .	# uid of content, first content of quantum which define this community
community.base: uid .	# uid of community, base.define is the base quantum of define quantum
community.maxInviteCnt: int .
community.minCosignCnt: int .
community.define: uid @reverse . # uid of quantum, quantum which define current community
community.initMembers: [uid] .   # uid of individual, only init members, not creator, not invited individuals
# -- content --
content.data: string .
content.fmt: int @index(int) .
# -- individual --
individual.address: string @index(hash) .
individual.communities: [uid] @reverse . # uid of community, which current individual is member
# -- quantum --
quantum.contents: [uid] .  # uid of contents
quantum.refs: [uid] @reverse .  # uid of quantums
quantum.sender: uid @reverse .  # uid of individual
quantum.sig: string @index(hash) .
quantum.type: int @index(int) .
quantum.timestamp: int @index(int) .
# -- types --
type community {
	community.note
	community.base
	community.maxInviteCnt
	community.minCosignCnt
	community.define
	community.initMembers
}
type content {
	content.data
	content.fmt
}
type individual {
	individual.address
	individual.communities
}
type quantum {
	quantum.contents
	quantum.refs
	quantum.sender
	quantum.sig
	quantum.type
	quantum.timestamp
}`
