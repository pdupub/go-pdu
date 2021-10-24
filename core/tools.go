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

// Combination return the cnt select n from m
func Combination(m, n int) (int, error) {
	if m < 1 || n < 1 || m < n {
		return 0, ErrCombinationParamsNotCorrect
	}
	val := 1
	for i := m; i > m-n; i-- {
		val *= i
	}
	for i := n; i > 1; i-- {
		val /= i
	}
	return val, nil
}

func CalcPopulation(glimit []*GenerationLimit, gNum []int) ([]int, error) {
	if len(gNum) == 0 {
		return nil, ErrInitialGenerationNumMissing
	}
	predictNum := []int{gNum[0]}

	for i := 0; i < len(glimit)-1; i++ {
		nextGNum := predictNum[i] * glimit[i].ChildrenMaxSize / glimit[i+1].ParentsMinSize
		if len(gNum) <= i+1 {
			predictNum = append(predictNum, nextGNum)
		} else if gNum[i+1] <= nextGNum {
			predictNum = append(predictNum, gNum[i+1])
		} else {
			return nil, ErrCalcGenerationNumberBeyondLimit
		}
	}
	return predictNum, nil
}
