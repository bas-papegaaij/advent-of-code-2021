package solvers

import (
	"fmt"
	"strings"
)

type Day12 struct{}

type cave struct {
	connectedCaves []*cave
	small          bool
	key            string
}

type caveStructure map[string]*cave
type caveList []*cave

func (d Day12) Solve(input []string, part int) (int64, error) {
	if part < 1 || part > 2 {
		return -1, invalidPartError(part)
	}
	structure := d.parseInput(input)
	allowDuplicate := part == 2
	return int64(structure.countPaths(structure["start"], nil, allowDuplicate)), nil
}

func (cs caveStructure) countPaths(start *cave, visited caveList, allowDuplicate bool) int {
	// we've found a path
	if start.key == "end" {
		return 1
	}

	if start.small {
		wasVisited := visited.contains(start.key)
		if wasVisited {
			if !allowDuplicate {
				return 0
			}
			allowDuplicate = false
		} else {
			visited = append(visited, start)
		}
	}

	total := 0
	for _, cave := range start.connectedCaves {
		total += cs.countPaths(cave, visited, allowDuplicate)
	}

	return total
}

func (d Day12) parseInput(input []string) caveStructure {
	structure := make(caveStructure)
	for _, line := range input {
		parts := strings.Split(line, "-")
		a, b := structure.getOrMakeCave(parts[0]), structure.getOrMakeCave(parts[1])
		// ignore the start so we don't try and re-traverse that
		if b.key != "start" {
			a.connectedCaves = append(a.connectedCaves, b)
		}
		if a.key != "start" {
			b.connectedCaves = append(b.connectedCaves, a)
		}
	}
	return structure
}

func (cs caveStructure) getOrMakeCave(key string) *cave {
	c, ok := cs[key]
	if !ok {
		c = &cave{
			small:          strings.ToLower(key) == key,
			connectedCaves: []*cave{},
			key:            key,
		}
		cs[key] = c
	}

	return c
}

func (cs caveStructure) print() {
	for k, v := range cs {
		fmt.Printf("%s: %+v\n", k, v)
	}
}

func (cl caveList) contains(key string) bool {
	for _, c := range cl {
		if c.key == key {
			return true
		}
	}
	return false
}
