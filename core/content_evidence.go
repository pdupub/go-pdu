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

package core

const (
	// LeakedPrivateKey is the evidence can prove a private key has been leaked.
	// Usually contain 1 anonymous msg(TypeText), the content is private key
	// and 1 msg sign by this private key.
	LeakedPrivateKey = iota

	// ExcessiveBirth is the evidence can prove a user build msg try to create next
	// new user against the nature rule.
	ExcessiveBirth
)

// ContentEvidence is the evidence of user illegal behavior,
// Recevier can punish the user by there own will.
// SenderID can be nil, for anonymous evidence
type ContentEvidence struct {
	EvidenceType int
	Msgs         []*Message
}

// CreateContentEvidence create the evidence msg content contain proof of user illeagal behavior
func CreateContentEvidence(evidenceType int, msgs []*Message) *ContentEvidence {
	return &ContentEvidence{EvidenceType: evidenceType, Msgs: msgs}
}
