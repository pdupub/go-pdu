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

package p2p

import (
	"testing"
)

// func TestNew(t *testing.T) {
// 	New(true, "./public", "testdata/test.db", nil)
// }

// curl --header "Content-Type: application/json" \
//   --request POST \
//   --data '{"address":"0xaEa768ddAd062bd341b6D03caeEfc371E675C1aE","content":"eyJhcnJheSI6WzEsMjAsM10sIm1hcCI6eyJhIjoyLCJiIjoxLCJiMjMiOjEwLCJmIjozLCJmZiI6MTAwMCwiaCI6MTB9LCJudW1iZXIiOjUsInN0cmluZyI6IkhlbGxvIFdvcmxkISJ9","nonce":14,"refs":["0x070d15041083041b48d0f2297357ce59ad18f6c608d70a1e6e04bcf494e366db","0x08fd3282eecbf25d31a9a5e51ed2d79a806f14281fbb583a5ee4024589b959d9"],"signature":"zLjcF4LV7hBjLklH7zKK8IBW7HZToAkMfPeHce8yfBMCeIsO9x6lBucb+RIx/Bhjd5eZVxQ0/nZgTfZtuMZQPwA="}' \
//   http://localhost:1323/peers

// curl --header "Content-Type: application/json"  \
//  	--request PUT \
//  	--data '{"address":"0xaEa768ddAd062bd341b6D03caeEfc371E675C1aE","content":"eyJhcnJheSI6WzEsMjAsM10sIm1hcCI6eyJhIjoyLCJiIjoxLCJiMjMiOjEwLCJmIjozLCJmZiI6MTAwMCwiaCI6MTB9LCJudW1iZXIiOjUsInN0cmluZyI6IkhlbGxvIFdvcmxkISJ9","nonce":14,"refs":["0x070d15041083041b48d0f2297357ce59ad18f6c608d70a1e6e04bcf494e366db","0x08fd3282eecbf25d31a9a5e51ed2d79a806f14281fbb583a5ee4024589b959d9"],"signature":"zLjcF4LV7hBjLklH7zKK8IBW7HZToAkMfPeHce8yfBMCeIsO9x6lBucb+RIx/Bhjd5eZVxQ0/nZgTfZtuMZQPwA="}' \
//    	http://localhost:1323/peers/1

func TestServer(t *testing.T) {
	t.Log("Test server")
}
