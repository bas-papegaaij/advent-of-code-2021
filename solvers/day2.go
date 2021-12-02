package solvers

import (
	"fmt"
	"strconv"
	"strings"
)

type Day2 struct{}

func (d Day2) Solve(input []string, part int) (int64, error) {
	switch part {
	case 1:
		return d.part1(input)
	case 2:
		return d.part2(input)
	}

	return -1, invalidPartError(part)
}

func (d Day2) part1(input []string) (int64, error) {
	depth := 0
	pos := 0

	for _, command := range input {
		action, magnitude, err := d.parseAction(command)
		if err != nil {
			return -1, err
		}

		switch action {
		case "forward":
			pos += magnitude
		case "down":
			depth += magnitude
		case "up":
			depth -= magnitude
		default:
			fmt.Println("Got unexpected action", action, "ignoring")
		}
	}

	return int64(depth * pos), nil
}

func (d Day2) part2(input []string) (int64, error) {
	aim := 0
	pos := 0
	depth := 0

	for _, command := range input {
		action, magnitude, err := d.parseAction(command)
		if err != nil {
			return -1, err
		}

		switch action {
		case "down":
			aim += magnitude
		case "up":
			aim -= magnitude
		case "forward":
			pos += magnitude
			depth += aim * magnitude
		}
	}

	return int64(pos * depth), nil
}

func (Day2) parseAction(command string) (string, int, error) {
	parts := strings.Split(command, " ")
	if len(parts) != 2 {
		return "", -1, fmt.Errorf("invalid command: %s", command)
	}

	magnitude, err := strconv.Atoi(parts[1])
	return parts[0], magnitude, err
}
