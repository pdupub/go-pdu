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
	// ErrPhotonTypeNotCorrect returns if photon type not correct
	ErrPhotonTypeNotCorrect = errors.New("photon type not correct")

	// ErrPhotonAlreadyExist returns if photon with same key already exist on entropy
	ErrPhotonAlreadyExist = errors.New("photon already exist")

	// ErrPhotonReferenceNotExist returns if photon reference not exist on entropy
	ErrPhotonReferenceNotExist = errors.New("photon reference not exist")

	// ErrPhotonNotExist returns if photon not exist on entropy
	ErrPhotonNotExist = errors.New("photon not exist")

	// ErrPhotonReferenceMissing returns when try to add new photon without references
	ErrPhotonReferenceMissing = errors.New("photon reference is missing")

	// ErrPhotonReferenceNotCorrect returns if refs[0] already have child from same author ...
	ErrPhotonReferenceNotCorrect = errors.New("photon reference is not correct")

	// ErrPhotonMissingReferenceByAuthor returns when photon missing reference by same author
	ErrPhotonMissingReferenceByAuthor = errors.New("photon missing reference by same author")

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
