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

package rule

const (
	// MortalLifetime is valid life time for mortal account
	// It means any msg from this account, if need to proof the msg be sent during the account life time.
	// The reference of that msg must contain at least one msg from the time proof account, which used
	// in DOBMsg when create this account. And the reference must choose from the msg from used as reference
	// dob message to next 2^16 messages from the same time proof account.
	MortalLifetime uint64 = 1 << 12 // 16

	// ReproductionInterval is valid time interval for reproduction
	// The interval is 1/4 of life time, so any account can participate in reproduction for three times.
	// Each reproduction need two account to cosign, so the max reproduce rate for mortal is 1.5
	ReproductionInterval = MortalLifetime >> 2

	// MaxLifeTime is valid life time for 0 generation account (root accounts Adam & Eve)
	MaxLifeTime = MortalLifetime << 16

	// LifetimeReduceRate is life time of account will reduce except the mortal.
	// life_time_of_child = max(life_time_of_parent) / LifetimeReduceRate
	// if life_time_of_child < MortalLifetime then life_time_of_child = MortalLifetime
	LifetimeReduceRate = 2
)
