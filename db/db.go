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

package db

const (
	// BucketUser is used to save all users
	BucketUser = "user"

	// BucketMsg is used to save the msg (msg.ID/ msg)
	BucketMsg = "msg"

	// BucketMID is used to save msg.ID by order (order/msg.ID)
	BucketMID = "mid"

	// BucketMOD is used to save msg received sequence (msg.ID/order)
	BucketMOD = "mod"

	// BucketLastMID is used to save last msg.ID by user.ID
	BucketLastMID = "lmid"

	// BucketConfig is used to save config info when universe be created
	BucketConfig = "config"

	// BucketPeer is used to save the peer information
	BucketPeer = "peer"

	// ConfigRoot0 root user which gender is 0
	ConfigRoot0 = "root0"

	// ConfigRoot1 root user which gender is 1
	ConfigRoot1 = "root1"

	// ConfigMsgCount is the current message count in the universe
	ConfigMsgCount = "msg_count"

	// ConfigCurrentStep is the current step of initialize the universe
	// step 0 - create bucket
	// step 1 - roots saved
	ConfigCurrentStep = "current_step"

	// ConfigLocalNodeKey is the local node key
	ConfigLocalNodeKey = "local_node_key"
)

const (
	// StepInitDB is the step which all bucket in db have been created
	StepInitDB = iota

	// StepRootsSaved is the step which two roots have been saved into db
	StepRootsSaved
)

// Row is the key/value pair from db
type Row struct {
	K string
	V []byte
}

// UDB is a database interface for embed database, default db is bolt
type UDB interface {
	Close() error
	CreateBucket(string) error
	DeleteBucket(string) error
	Set(string, string, []byte) error
	Get(string, string) ([]byte, error)
	Find(string, string, ...int) ([]*Row, error)
}
