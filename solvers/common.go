package solvers

import (
	"fmt"
	"strconv"
	"strings"
)

type intVector struct {
	x int
	y int
}

func invalidPartError(part int) error {
	return fmt.Errorf("no solution for question part %d", part)
}

func inputsToInt(input []string) ([]int, error) {
	res := make([]int, 0, len(input))
	for i, line := range input {
		num, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("unable to parse %s as int on line %d", line, i)
		}

		res = append(res, num)
	}

	return res, nil
}

func splitByEmptyLines(input []string) [][]string {
	res := make([][]string, 0)
	curBlock := []string{}
	for i := 0; i < len(input); i++ {
		curStr := input[i]
		if curStr == "" {
			res = append(res, curBlock)
			curBlock = []string{}
		} else {
			curBlock = append(curBlock, curStr)
		}
	}
	res = append(res, curBlock)
	return res
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

func sumIntSlice(vals []int) int {
	sum := 0
	for _, num := range vals {
		sum += num
	}
	return sum
}

func stringToIntvector(val string) (intVector, error) {
	split := strings.Split(val, ",")
	if len(split) != 2 {
		return intVector{}, fmt.Errorf("invalid intvector string %s had %d parts, expected 2", val, len(split))
	}
	x, err := strconv.Atoi(split[0])
	if err != nil {
		return intVector{}, fmt.Errorf("invalid intvector string had non-number x-part %s", split[0])
	}
	y, err := strconv.Atoi(split[1])
	if err != nil {
		return intVector{}, fmt.Errorf("invalid intvector string had non-number y-part %s", split[0])
	}

	return intVector{x, y}, nil
}
