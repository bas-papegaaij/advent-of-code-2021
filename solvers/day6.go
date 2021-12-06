package solvers

import (
	"strings"
)

type Day6 struct{}

func (d Day6) Solve(input []string, part int) (int64, error) {
	counts := make([]int, 9)
	numStrings := strings.Split(input[0], ",")
	startingCount, err := inputsToInt(numStrings)
	if err != nil {
		return -1, err
	}

	for _, counter := range startingCount {
		counts[counter]++
	}

	switch part {
	case 1:
		return d.solve(counts, 80)
	case 2:
		return d.solve(counts, 256)
	default:
		return -1, invalidPartError(part)
	}
}

func (d Day6) solve(counters []int, simDays int) (int64, error) {
	// let's start with a naive - brute-ish force approach
	// Should still effectively only be O(n) wrt the number of days
	// we need to simulate

	for i := 0; i < simDays; i++ {
		newFish := counters[0]

		for j := 0; j < len(counters)-1; j++ {
			// decrease each counter by 1 - i.e. all 6es are now 5s etc
			counters[j] = counters[j+1]
		}

		counters[len(counters)-1] = newFish
		counters[6] += newFish
	}
	return int64(sumIntSlice(counters)), nil
}
