package solvers

import (
	"fmt"
	"strings"
)

type Day11 struct{}

type octopus struct {
	power   int
	flashed bool
}

type octopusField struct {
	octopuses []octopus
	width     int
	height    int
}

func (d Day11) Solve(input []string, part int) (int64, error) {
	field, err := d.parseInput(input)
	if err != nil {
		return -1, err
	}

	switch part {
	case 1:
		return int64(field.countFlashes(100)), nil
	case 2:
		return int64(field.findSynchronisedFlash()), nil
	}

	return -1, invalidPartError(part)
}

func (d Day11) parseInput(input []string) (octopusField, error) {
	width, height := len(input[0]), len(input)
	octopuses := make([]octopus, 0, width*height)

	for _, line := range input {
		powerLevels, err := inputsToInt(strings.Split(line, ""))
		if err != nil {
			return octopusField{}, err
		}

		for _, power := range powerLevels {
			octopuses = append(octopuses, octopus{power, false})
		}
	}

	return octopusField{
		octopuses,
		width,
		height,
	}, nil
}

func (f *octopusField) countFlashes(simSteps int) int {
	total := 0
	for i := 0; i < simSteps; i++ {
		total += f.simulate()
	}
	return total
}

func (f *octopusField) findSynchronisedFlash() int {
	for step := 1; ; step++ {
		flashes := f.simulate()
		if flashes == len(f.octopuses) {
			return step
		}
	}
}

// runs the sim once, returning the number of flashes seen
func (f *octopusField) simulate() int {
	// increase power of each octopus
	for i := 0; i < len(f.octopuses); i++ {
		f.octopuses[i].power++
	}
	// Go through each octopus and do the flash
	for x := 0; x < f.width; x++ {
		for y := 0; y < f.height; y++ {
			f.doFlash(x, y)
		}
	}

	// reset the octopuses
	total := 0
	for i := 0; i < len(f.octopuses); i++ {
		if f.octopuses[i].reset() {
			total++
		}
	}

	return total
}

// Performs the "flash" for the octopus at the given position
// returns the number of flashes this caused including this one
func (f *octopusField) doFlash(x int, y int) {
	o := f.getOctopus(x, y)
	if o == nil || o.flashed || o.power <= 9 {
		return
	}

	o.flashed = true
	for xx := x - 1; xx <= x+1; xx++ {
		for yy := y - 1; yy <= y+1; yy++ {
			if xx == x && yy == y {
				continue
			}
			o2 := f.getOctopus(xx, yy)
			if o2 != nil {
				o2.power++
			}
			f.doFlash(xx, yy)
		}
	}
}

func (f *octopusField) getOctopus(x int, y int) *octopus {
	if x < 0 || x >= f.width || y < 0 || y >= f.height {
		return nil
	}

	return &f.octopuses[y*f.width+x]
}

func (o *octopus) reset() bool {
	if o.flashed {
		o.power = 0
		o.flashed = false
		return true
	}

	return false
}

func (f *octopusField) print() {
	for y := 0; y < f.height; y++ {
		fmt.Println(f.octopuses[y*f.width : y*f.width+f.width])
	}
}
