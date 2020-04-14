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

import "errors"

var (
	// ErrUserNotExist returns fail to find a user
	ErrUserNotExist = errors.New("user not exist")

	// ErrMsgAlreadyExist returns try to add a message which already exist in universe
	ErrMsgAlreadyExist = errors.New("msg already exist")

	// ErrMsgNotFound returns fail to find a message
	ErrMsgNotFound = errors.New("msg not found")

	// ErrTPAlreadyExist returns time proof message already exist
	ErrTPAlreadyExist = errors.New("time proof already exist")

	// ErrUserAlreadyExist returns try a add a user who is already exist in unverse
	ErrUserAlreadyExist = errors.New("user already exist")

	// ErrNotSupportYet returns not support temp
	ErrNotSupportYet = errors.New("not error, just not support yet")

	// ErrNewUserAddFail returns add new user fail for unknown reason
	ErrNewUserAddFail = errors.New("new user add fail")

	// ErrCreateSpaceTimeFail returns when create space time fail
	ErrCreateSpaceTimeFail = errors.New("create space time fail")

	// ErrAddUserToSpaceTimeFail returns when add user to space time fail
	ErrAddUserToSpaceTimeFail = errors.New("add user to space time fail")

	// ErrCreateRootUserFail returns when create root user fail
	ErrCreateRootUserFail = errors.New("create root user fail")

	// ErrContentTypeNotBirth returns when try to add user from a not birth message
	ErrContentTypeNotBirth = errors.New("content type is not TypeBirth")

	// ErrDimensionNumberNotSuitable returns if the dimension is zero or too large
	ErrDimensionNumberNotSuitable = errors.New("number of dimension is not suitable")

	// ErrPerimeterIsZero returns if perimeter is zero
	ErrPerimeterIsZero = errors.New("perimeter should not be zero")
)
