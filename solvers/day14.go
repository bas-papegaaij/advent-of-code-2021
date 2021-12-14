package solvers

import (
	"math"
)

type Day14 struct{}

type polymerCount map[string]int64
type polymerInstructions map[string]string

type polymerizer struct {
	pairCounts    polymerCount
	elementCounts polymerCount
	instructions  polymerInstructions
}

func (d Day14) Solve(input []string, part int) (int64, error) {
	if part < 1 || part > 2 {
		return -1, invalidPartError(part)
	}

	poly := d.parseInput(input)

	steps := 10
	if part == 2 {
		steps = 40
	}

	// really hope we don't overflow...
	return d.solve(poly, steps), nil
}

func (d Day14) solve(poly polymerizer, steps int) int64 {
	for i := 0; i < steps; i++ {
		poly.updateCounts()
	}
	return poly.calculateScore()
}

func (d Day14) parseInput(input []string) polymerizer {
	template := input[0]
	instructions := make(map[string]string, len(input)-2)
	// We'll assume the inputs are valid to save having to write more parsing code
	for _, instruction := range input[2:] {
		instructions[string(instruction[0:2])] = string(instruction[6:])
	}

	pairCounts := make(polymerCount, len(instructions))
	elementCounts := make(polymerCount)
	elementCounts[string(template[0])]++
	for i := 0; i < len(template)-1; i++ {
		a, b := template[i], template[i+1]
		elementCounts[string(b)]++
		pairCounts[string(a)+string(b)]++
	}

	return polymerizer{
		pairCounts,
		elementCounts,
		instructions,
	}
}

func (p polymerizer) updateCounts() {
	tmp := make(polymerCount, len(p.pairCounts))
	for k, v := range p.pairCounts {
		tmp[k] = v
	}

	for pair, curVal := range tmp {
		if curVal == 0 {
			continue
		}
		replacement := p.instructions[pair]
		p.elementCounts[replacement] += curVal
		A, B := string(pair[0])+replacement, replacement+string(pair[1])
		p.pairCounts[A] += curVal
		p.pairCounts[B] += curVal
		p.pairCounts[pair] -= curVal
	}
}

func (p polymerizer) calculateScore() int64 {
	var min int64 = math.MaxInt64
	var max int64 = math.MinInt64
	for _, count := range p.elementCounts {
		if count < min {
			min = count
		}
		if count > max {
			max = count
		}
	}

	return max - min
}

func (p polymerizer) totalElements() int64 {
	var total int64
	for _, count := range p.elementCounts {
		total += count
	}
	return total
}
