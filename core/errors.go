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

package core

import "errors"

var (
	// ErrQuantumTypeNotCorrect returns if quantum type not correct
	ErrQuantumTypeNotCorrect = errors.New("quantum type not correct")

	// ErrQuantumAlreadyExist returns if quantum with same key already exist on entropy
	ErrQuantumAlreadyExist = errors.New("quantum already exist")

	// ErrQuantumReferenceNotExist returns if quantum reference not exist on entropy
	ErrQuantumReferenceNotExist = errors.New("quantum reference not exist")

	// ErrQuantumNotExist returns if quantum not exist on entropy
	ErrQuantumNotExist = errors.New("quantum not exist")

	// ErrQuantumReferenceMissing returns when try to add new quantum without references
	ErrQuantumReferenceMissing = errors.New("quantum reference is missing")

	// ErrQuantumReferenceNotCorrect returns if refs[0] already have child from same author ...
	ErrQuantumReferenceNotCorrect = errors.New("quantum reference is not correct")

	// ErrQuantumMissingReferenceByAuthor returns when quantum missing reference by same author
	ErrQuantumMissingReferenceByAuthor = errors.New("quantum missing reference by same author")

	// ErrSocietyIDConflict returns if address in the value not equal to key
	ErrSocietyIDConflict = errors.New("Address conflict")

	// ErrCombinationParamsNotCorrect return if m < 1 or n < 1 or n > m
	ErrCombinationParamsNotCorrect = errors.New("combination params m, n not correct")

	// ErrInitialGenerationNumMissing return if try to calc generation population without setting initial generation number
	ErrInitialGenerationNumMissing = errors.New("initial generation number missing")

	// ErrCalcGenerationNumberBeyondLimit return if generation number is larger than the theoretical maximize value
	ErrCalcGenerationNumberBeyondLimit = errors.New("calculate generation number beyond limit")

	// ErrAddIndividualWithoutEnoughParents return if try to create new ID without enough parents sign
	ErrAddIndividualWithoutEnoughParents = errors.New("add individual without enough parents")

	// ErrAddIndividualBeyondChildrenMaxLimit return if parents sign the create msg have create children beyond limit
	ErrAddIndividualBeyondChildrenMaxLimit = errors.New("add individual beyond children max limit")

	// ErrIndividualNotExistInSociety return if individual not exist in society
	ErrIndividualNotExistInSociety = errors.New("individual not exist in society")
)
