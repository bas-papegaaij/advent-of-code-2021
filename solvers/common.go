package solvers

import (
	"fmt"
	"strconv"
)

// defined for convenience here so we can
type genericSolver interface {
	part1(input []string) (int64, error)
	part2(input []string) (int64, error)
}

// convenient where we don't need to do additional input pre-processing
// e.g. each line from the input is sufficient for the solution functions
func genericSolve(solver genericSolver, input []string, part int) (int64, error) {
	switch part {
	case 1:
		return solver.part1(input)
	case 2:
		return solver.part2(input)
	}

	return -1, fmt.Errorf("no solution for question part %d", part)
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
