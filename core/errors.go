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
	ErrUserNotExist           = errors.New("user not exist")
	ErrMsgAlreadyExist        = errors.New("msg already exist")
	ErrMsgNotFound            = errors.New("msg not found")
	ErrTPAlreadyExist         = errors.New("time proof already exist")
	ErrUserAlreadyExist       = errors.New("user already exist")
	ErrNotSupportYet          = errors.New("not error, just not support yet")
	ErrNewUserAddFail         = errors.New("new user add fail")
	ErrCreateSpaceTimeFail    = errors.New("create space time fail")
	ErrAddUserToSpaceTimeFail = errors.New("add user to space time fail")
	ErrCreateRootUserFail     = errors.New("create root user fail")
	ErrContentTypeNotDOB      = errors.New("content type is not TypeDOB")
)
