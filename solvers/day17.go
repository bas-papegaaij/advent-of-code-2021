package solvers

import (
	"strconv"
	"strings"
)

type Day17 struct{}

type projectile struct {
	velocity intVector
	position intVector
}

func (d Day17) Solve(input []string, part int) (int64, error) {
	min, max := d.parseInput(input[0])
	candidates := d.findProbeCandidates(min, max)

	if part == 1 {
		return int64(d.findHighestCandidate(candidates)), nil
	}

	if part == 2 {
		return int64(len(candidates)), nil
	}
	return -1, nil
}

func (d Day17) parseInput(input string) (intVector, intVector) {
	coords := input[13:]
	split := strings.Split(coords, ", ")
	xMinStr, xMaxStr := strings.Split(split[0][2:], "..")[0], strings.Split(split[0][2:], "..")[1]
	yMinStr, yMaxStr := strings.Split(split[1][2:], "..")[0], strings.Split(split[1][2:], "..")[1]

	xMin, _ := strconv.Atoi(xMinStr)
	xMax, _ := strconv.Atoi(xMaxStr)
	yMin, _ := strconv.Atoi(yMinStr)
	yMax, _ := strconv.Atoi(yMaxStr)

	return intVector{xMin, yMin}, intVector{xMax, yMax}
}

func (d Day17) findProbeCandidates(minTarget intVector, maxTarget intVector) []intVector {
	candidates := []intVector{}
	for initialY := -minTarget.y; initialY >= minTarget.y; initialY-- {
		for initialX := 1; initialX <= maxTarget.x; initialX++ {
			probe := projectile{
				velocity: intVector{initialX, initialY},
				position: intVector{0, 0},
			}

			// if initialY is positive, we can skip simulating the steps until it reaches 0 again
			// since it will do so in initialY*2 + 1 steps, at which point its velocity is -(initialY+1)
			if initialY > 0 {
				probe.velocity.y = -(initialY + 1)
				step := initialY*2 + 1
				// we can skip simulating x easily enough just by using triangle numbers
				offset := initialX - step
				if offset < 0 {
					offset = 0
				}
				probe.position.x = d.getMaxDistance(initialX) - d.getMaxDistance(offset)
				probe.velocity.x = initialX - step
				if probe.velocity.x < 0 {
					probe.velocity.x = 0
				}
			}

			done := false
			for !done {
				probe.step()
				if probe.position.x >= minTarget.x &&
					probe.position.x <= maxTarget.x &&
					probe.position.y >= minTarget.y &&
					probe.position.y <= maxTarget.y {
					candidates = append(candidates, intVector{initialX, initialY})
					break
				}

				if probe.position.x > maxTarget.x || probe.position.y < minTarget.y {
					done = true
				}
			}
		}
	}
	return candidates
}

func (d Day17) findHighestCandidate(candidates []intVector) int {
	max := 0
	for _, c := range candidates {
		if c.y < 0 {
			continue
		}

		cMax := d.getMaxDistance(c.y)
		if cMax > max {
			max = cMax
		}
	}
	return max
}

func (d Day17) getMaxDistance(initialV int) int {
	return (initialV * (initialV + 1)) / 2
}

func (p *projectile) step() {
	p.position.x += p.velocity.x
	p.position.y += p.velocity.y
	if p.velocity.x > 0 {
		p.velocity.x--
	} else if p.velocity.x < 0 {
		p.velocity.x++
	}

	p.velocity.y--
}
