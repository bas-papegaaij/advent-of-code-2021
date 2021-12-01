package solvers

import (
	"fmt"
	"strconv"
)

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
