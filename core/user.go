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

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/pdupub/go-pdu/common"
	"github.com/pdupub/go-pdu/core/rule"
	"github.com/pdupub/go-pdu/crypto"
)

const (
	// DefaultDimensionNum is the default number of demension for calculating distance
	DefaultDimensionNum = 4
	// DefaultPerimeter is the default perimeter of universe for calculating distance
	DefaultPerimeter = 1e+10
	// MaxDimensionNum is the max number of demension
	MaxDimensionNum = 11
)

// User is the author of any msg in pdu
type User struct {
	Name       string   `json:"name"`
	BirthExtra string   `json:"extra"`
	Auth       *Auth    `json:"auth"`
	BirthMsg   *Message `json:"birthMsg"`
	LifeTime   uint64   `json:"lifeTime"`
}

// CreateRootUser try to create root user by public key
// Gender of the root user is depend on key,name and extra
func CreateRootUser(key crypto.PublicKey, name, extra string) *User {
	return &User{Name: name, BirthExtra: extra, Auth: &Auth{PublicKey: key}, BirthMsg: nil, LifeTime: rule.MaxLifeTime}
}

// CreateNewUser create new user by cosign message
// The msg must be signed by user in local user dag.
// Both parents must be in the local use dag.
// Both parents fit the nature rules.
// The Birth struct signed by both parents.
func CreateNewUser(universe *Universe, msg *Message) (*User, error) {
	if msg.Value.ContentType != TypeBirth {
		return nil, ErrContentTypeNotBirth
	}
	var contentBirth ContentBirth
	if err := json.Unmarshal(msg.Value.Content, &contentBirth); err != nil {
		return nil, err
	}
	newUser := contentBirth.User
	newUser.BirthMsg = msg
	// calculate the life time of new user
	p0 := universe.userD.GetVertex(contentBirth.Parents[0].UserID)
	if p0 == nil {
		return nil, ErrUserNotExist
	}
	maxParentLifeTime := p0.Value().(*User).LifeTime

	p1 := universe.userD.GetVertex(contentBirth.Parents[1].UserID)
	if p1 == nil {
		return nil, ErrUserNotExist
	}
	if maxParentLifeTime < p1.Value().(*User).LifeTime {
		maxParentLifeTime = p1.Value().(*User).LifeTime
	}
	if maxParentLifeTime == rule.MortalLifetime {
		newUser.LifeTime = rule.MortalLifetime
	} else {
		newUser.LifeTime = maxParentLifeTime / rule.LifetimeReduceRate
	}

	return &newUser, nil
}

// ID return the vertex.id, related to parents and value of the vertex
// ID cloud use as address of user account
func (u User) ID() common.Hash {
	hash := sha256.New()
	hash.Reset()

	auth, _ := json.Marshal(u.Auth)
	lifeTime := fmt.Sprintf("%v", u.LifeTime)
	var birthMsg string
	// todo : add init BirthMsg to rootUser
	// todo : so this condition can be deleted
	if u.BirthMsg != nil {
		birthMsg += fmt.Sprintf("%v", u.BirthMsg.SenderID)
		for _, v := range u.BirthMsg.Reference {
			birthMsg += fmt.Sprintf("%v%v", v.SenderID, v.MsgID)
		}
		birthMsg += fmt.Sprintf("%v%v%v", u.BirthMsg.Signature.Signature, u.BirthMsg.Signature.Source, u.BirthMsg.Signature.SigType)
		birthMsg += fmt.Sprintf("%v%v", u.BirthMsg.Value.Content, u.BirthMsg.Value.ContentType)
	}
	hash.Write(append(append(append(append([]byte(u.Name), u.BirthExtra...), auth...), birthMsg...), lifeTime...))
	return common.Bytes2Hash(hash.Sum(nil))
}

// Gender return the gender of user, true = male = end of ID is odd
func (u User) Gender() bool {
	hashID := u.ID()
	if uid := new(big.Int).SetBytes(hashID[:]); uid.Mod(uid, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
		return true
	}
	return false
}

// Value return the vertex.value
func (u User) Value() interface{} {
	return nil
}

// Distance return the distance from user to common.Hash base on setting perimeter and dimension
func (u User) Distance(location common.Hash, dimension int, perimeter *big.Int) (distance *big.Int, err error) {
	return u.distance(location, dimension, perimeter)
}

// StandardDistance return the distance from user to common.Hash base on default setting perimeter and dimension
func (u User) StandardDistance(location common.Hash) (distance *big.Int, err error) {
	return u.distance(location, DefaultDimensionNum, big.NewInt(DefaultPerimeter))
}

func (u User) distance(location common.Hash, dimension int, perimeter *big.Int) (distance *big.Int, err error) {
	coordinates, err := u.GetCoordinates(location, dimension, perimeter)
	if err != nil {
		return nil, err
	}
	squareSum := big.NewInt(0)
	for _, coordinate := range coordinates {
		squareSum = squareSum.Add(squareSum, new(big.Int).Mul(coordinate, coordinate))
	}
	return squareSum.Sqrt(squareSum), nil
}

// GetCoordinates returns the coordinates of location, use u.ID as origin point
func (u User) GetCoordinates(location common.Hash, dimension int, perimeter *big.Int) ([]*big.Int, error) {
	if dimension <= 0 || dimension > MaxDimensionNum {
		return nil, ErrDimensionNumberNotSuitable
	}
	if perimeter.Cmp(big.NewInt(0)) == 0 {
		return nil, ErrPerimeterIsZero
	}
	perimeter.Abs(perimeter)
	coordinates := make([]*big.Int, dimension)
	dimensionSize := common.HashLength / dimension
	selfLocation := u.ID()
	maxDistance := new(big.Int).Rsh(perimeter, 1)

	for i := 0; i < dimension; i++ {
		if i == dimension-1 {
			coordinates[i] = new(big.Int).Sub(new(big.Int).SetBytes(location[i*dimensionSize:]), new(big.Int).SetBytes(selfLocation[i*dimensionSize:]))
		} else {
			coordinates[i] = new(big.Int).Sub(new(big.Int).SetBytes(location[i*dimensionSize:(i+1)*dimensionSize]), new(big.Int).SetBytes(selfLocation[i*dimensionSize:(i+1)*dimensionSize]))
		}
		coordinates[i].Add(coordinates[i], maxDistance)
		coordinates[i].Mod(coordinates[i], perimeter)
		coordinates[i].Sub(coordinates[i], maxDistance)

	}
	return coordinates, nil
}

// ParentsID return the ID of user parents,
// res[0] should be the female parent (id end by even)
// res[1] should be the male parent (id end by odd)
func (u User) ParentsID() [2]common.Hash {
	var parentsID [2]common.Hash
	if u.BirthMsg != nil {
		// get parents from birthMsg
		var contentBirth ContentBirth
		if err := json.Unmarshal(u.BirthMsg.Value.Content, &contentBirth); err != nil {
			return parentsID
		}
		parentsID[0] = contentBirth.Parents[0].UserID
		parentsID[1] = contentBirth.Parents[1].UserID
	}
	return parentsID
}

// UnmarshalJSON is used to unmarshal json
func (u *User) UnmarshalJSON(input []byte) error {
	userMap := make(map[string]interface{})
	err := json.Unmarshal(input, &userMap)
	if err != nil {
		return err
	}
	u.Name = userMap["name"].(string)
	u.BirthExtra = userMap["birthExtra"].(string)
	u.LifeTime, err = strconv.ParseUint(userMap["lifeTime"].(string), 0, 64)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(userMap["birthMsg"].(string)), &u.BirthMsg)
	json.Unmarshal([]byte(userMap["auth"].(string)), &u.Auth)

	return nil
}

// MarshalJSON marshal user to json
func (u *User) MarshalJSON() ([]byte, error) {
	userMap := make(map[string]interface{})
	userMap["name"] = u.Name
	userMap["birthExtra"] = u.BirthExtra
	userMap["lifeTime"] = fmt.Sprintf("%v", u.LifeTime)

	auth, err := json.Marshal(&u.Auth)
	if err != nil {
		return []byte{}, err
	}
	userMap["auth"] = string(auth)
	birthMsg, err := json.Marshal(&u.BirthMsg)
	if err != nil {
		return []byte{}, err
	}
	userMap["birthMsg"] = string(birthMsg)

	return json.Marshal(userMap)
}
