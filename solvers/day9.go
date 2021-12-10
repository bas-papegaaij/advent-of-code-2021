package solvers

import (
	"sort"
	"strconv"
	"strings"
)

type Day9 struct{}

type depthPoint struct {
	value int
	basin *depthPoint
	pos   intVector
}
type depthMap struct {
	width  int
	height int
	values []depthPoint
}

func (d Day9) Solve(input []string, part int) (int64, error) {
	grid, err := d.parseHeightMap(input)
	if err != nil {
		return -1, err
	}

	if part == 0 || part > 2 {
		return -1, invalidPartError(part)
	}

	return d.solve(grid, part), nil
}

func (d Day9) solve(grid depthMap, part int) int64 {
	basins := make(map[*depthPoint]int)
	for y := 0; y < grid.height; y++ {
		for x := 0; x < grid.width; x++ {
			count, basin := grid.followToBasin(intVector{x, y})
			if basin != nil {
				basins[basin] += count
			}
		}
	}

	if part == 1 {
		sumBasins := 0
		for basin := range basins {
			sumBasins += basin.value + 1
		}
		return int64(sumBasins)
	}

	// Not the most efficient way to check for size, but the set is small enough where
	// it doesn't have a meaningful impact on the performance
	sizes := make([]int, len(basins))
	for _, size := range basins {
		sizes = append(sizes, size)
	}
	sort.Ints(sizes)
	back := len(sizes) - 1
	return int64(sizes[back] * sizes[back-1] * sizes[back-2])
}

func (d Day9) parseHeightMap(input []string) (depthMap, error) {
	width := len(input[0])
	size := width * len(input)
	grid := make([]depthPoint, 0, size)
	for y, line := range input {
		for x, p := range strings.Split(line, "") {
			val, err := strconv.Atoi(p)
			if err != nil {
				return depthMap{}, err
			}
			grid = append(grid, depthPoint{val, nil, intVector{x, y}})
		}
	}

	return depthMap{
		width,
		len(input),
		grid,
	}, nil
}

// assumes x and y are within the bounds of the grid
func (grid depthMap) isLowPoint(p *depthPoint) bool {
	// Technically we could get a tiny performance improvement by
	// doing the same check as getLowestAdjacentPoint but bailing early
	// if any points are lower. This feels a lot cleaner though
	lowest := grid.getLowestAdjacentPoint(p)
	if lowest == nil {
		return true
	}

	return lowest.value > p.value
}

func (grid depthMap) getPoint(x int, y int) *depthPoint {
	if x < 0 || x >= grid.width || y < 0 || y >= grid.height {
		return nil
	}

	return &grid.values[y*grid.width+x]
}

// returns how many previously uncounted points sit between this point
// and its nearest basin (including this point). And returns the basin connected to this point
func (grid depthMap) followToBasin(startingPoint intVector) (int, *depthPoint) {
	curPt := grid.getPoint(startingPoint.x, startingPoint.y)

	if curPt == nil || curPt.basin != nil || curPt.value == 9 {
		return 0, nil
	}

	if grid.isLowPoint(curPt) {
		curPt.basin = curPt
		return 1, curPt
	}

	nextPt := grid.getLowestAdjacentPoint(curPt)

	if nextPt.basin != nil {
		curPt.basin = nextPt.basin
		return 1, nextPt.basin
	}

	count, basin := grid.followToBasin(nextPt.pos)
	curPt.basin = basin
	return 1 + count, basin
}

func (grid depthMap) getLowestAdjacentPoint(current *depthPoint) *depthPoint {
	x, y := current.pos.x, current.pos.y
	adjacent := make([]*depthPoint, 4)
	adjacent[0] = grid.getPoint(x-1, y)
	adjacent[1] = grid.getPoint(x+1, y)
	adjacent[2] = grid.getPoint(x, y-1)
	adjacent[3] = grid.getPoint(x, y+1)

	lowest := 9
	var nextPt *depthPoint
	for _, pt := range adjacent {
		if pt == nil {
			continue
		}

		if pt.value < lowest {
			lowest = pt.value
			nextPt = pt
		}
	}

	return nextPt
}
