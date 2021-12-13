package solvers

import (
	"fmt"
	"strconv"
	"strings"
)

type Day13 struct{}

type foldDirection int

const (
	xFold foldDirection = 0
	yFold               = iota
)

type origamiInstruction struct {
	axis  foldDirection
	value int
}

type origamiInput struct {
	points       []intVector
	instructions []origamiInstruction
}

func (d Day13) Solve(input []string, part int) (int64, error) {
	o, err := d.parseInput(input)
	if err != nil {
		return -1, err
	}

	if part < 1 || part > 2 {
		return -1, invalidPartError(part)
	}

	totalFolds := 1
	if part == 2 {
		totalFolds = len(o.instructions)
	}

	for i := 0; i < totalFolds; i++ {
		o.fold(i)
	}

	if part == 2 {
		o.print()
	}

	return int64(len(o.points)), nil
}

func (d Day13) parseInput(input []string) (origamiInput, error) {
	points := []intVector{}
	curLine, line := 0, input[0]
	for ; line != ""; curLine++ {
		vec, err := stringToIntvector(line)
		if err != nil {
			return origamiInput{}, err
		}
		points = append(points, vec)
		line = input[curLine+1]
	}
	curLine++

	instructions := []origamiInstruction{}
	for _, line := range input[curLine:] {
		axisPart, value := strings.Split(line, "=")[0], strings.Split(line, "=")[1]
		var axis foldDirection
		if axisPart[len(axisPart)-1] == 'x' {
			axis = xFold
		} else if axisPart[len(axisPart)-1] == 'y' {
			axis = yFold
		} else {
			return origamiInput{}, fmt.Errorf("invalid fold direction %s", string(axisPart[len(axisPart)-1]))
		}

		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return origamiInput{}, err
		}
		instructions = append(instructions, origamiInstruction{
			axis,
			parsedValue,
		})
	}

	return origamiInput{
		points,
		instructions,
	}, nil
}

func (o *origamiInput) fold(instruction int) {
	inst := o.instructions[instruction]
	for i := 0; i < len(o.points); i++ {
		if inst.axis == xFold && o.points[i].x > inst.value {
			diff := o.points[i].x - inst.value
			o.points[i].x = inst.value - diff
		} else if inst.axis == yFold && o.points[i].y > inst.value {
			diff := o.points[i].y - inst.value
			o.points[i].y = inst.value - diff
		}
	}
	o.removeDuplicatePoints()
}

func (o *origamiInput) removeDuplicatePoints() {
	for i := 0; i < len(o.points); i++ {
		for j := i + 1; j < len(o.points); j++ {
			if o.points[i] == o.points[j] {
				o.points[i] = o.points[len(o.points)-1]
				o.points = o.points[:len(o.points)-1]
				i--
				break
			}
		}
	}
}

func (o *origamiInput) print() {
	// find the width of our grid
	width, height := 0, 0
	for _, p := range o.points {
		fmt.Println(p)
		if p.x > width {
			width = p.x
		}
		if p.y > height {
			height = p.y
		}
	}
	width++
	height++

	lines := make([][]rune, height)
	for i := 0; i < height; i++ {
		line := make([]rune, width)
		for j := 0; j < width; j++ {
			line[j] = '.'
		}
		lines[i] = line
	}

	fmt.Println("height", height)
	for _, p := range o.points {
		lines[p.y][p.x] = '#'
	}

	for i := 0; i < height; i++ {
		fmt.Println(string(lines[i]))
	}
}
