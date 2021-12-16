package solvers

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Day15 struct{}

type chitonNode struct {
	distance int
	risk     int
	visited  bool
	index    int
	h        int
	f        int
}

type chitonGrid struct {
	width int
	nodes []*chitonNode
}

type chitonPQueue struct {
	elems []*chitonNode
}

func (d Day15) Solve(input []string, part int) (int64, error) {
	grid, err := d.parseInput(input)
	if err != nil {
		return -1, err
	}

	switch part {
	case 1:
		grid.calculateRisks()
		return int64(grid.nodes[len(grid.nodes)-1].distance), nil
	case 2:
		grid.expandMap(5)
		grid.calculateRisks()
		// grid.printPath()
		return int64(grid.nodes[len(grid.nodes)-1].distance), nil
	}

	return -1, invalidPartError(part)
}

func (d Day15) parseInput(input []string) (*chitonGrid, error) {
	width, height := len(input[0]), len(input)
	nodes := make([]*chitonNode, 0, width*height)
	x, y := 0, 0
	for _, line := range input {
		vals, err := inputsToInt(strings.Split(line, ""))
		if err != nil {
			return nil, err
		}

		for _, val := range vals {
			nodes = append(nodes, &chitonNode{
				distance: math.MaxInt,
				risk:     val,
				index:    len(nodes),
				// h:        ((width - 1) - x) + ((height - 1) - y),
				h: 0,
				f: math.MaxInt,
			})
			x++
		}
		x = 0
		y++
	}

	return &chitonGrid{
		width: width,
		nodes: nodes,
	}, nil
}

// This calculates the length of the shortest path to every node from the starting node
// This operation updates the nodes in-place
func (g *chitonGrid) calculateRisks() {
	steps := 0
	defer func() {
		fmt.Println("total steps", steps)
	}()
	g.nodes[0].distance = 0
	fmt.Println("node count:", len(g.nodes))
	pqueue := chitonPQueue{}
	pqueue.put(g.nodes[0])
	var cur *chitonNode
	for pqueue.length() > 0 {
		steps++
		cur = pqueue.get()
		// fmt.Println("selected x", cur.index%g.width, "y", cur.index/g.width)
		// cur.visited = true
		for _, connected := range g.getNeighbours(cur.index) {
			// if connected.visited {
			// 	continue
			// }

			if connected == g.nodes[len(g.nodes)-1] {
				connected.distance = cur.distance + connected.risk
				return
			}

			// calculate my heuristic as the shortest "straight-line" distance to the end node
			newVal := cur.distance + connected.risk
			if newVal < connected.distance {
				connected.distance = newVal
				connected.f = connected.distance + connected.h
				pqueue.put(connected)
			}
		}
	}
}

func (g *chitonGrid) expandMap(factor int) {
	originalWidth := g.width
	originalHeight := len(g.nodes) / g.width

	newNodes := make([]*chitonNode, len(g.nodes)*factor*factor)

	for column := 0; column < factor; column++ {
		for row := 0; row < factor; row++ {
			// copy the original nodes, adding row + column to each number in each square
			for i := 0; i < len(g.nodes); i++ {
				origX := i % originalWidth
				origY := i / originalWidth
				newX := origX + column*originalWidth
				newY := origY + row*originalHeight
				newIdx := newY*originalWidth*factor + newX
				risk := g.nodes[i].risk + row + column
				for risk > 9 {
					risk -= 9
				}
				newNodes[newIdx] = &chitonNode{
					distance: math.MaxInt,
					risk:     risk,
					index:    newIdx,
					f:        math.MaxInt,
					// h:        (((originalWidth * factor) - 1) - newX) + (((originalHeight * factor) - 1) - newY),
					h: 0,
				}
			}
		}
	}

	g.nodes = newNodes
	g.width = originalWidth * factor
}

func (g chitonGrid) getNeighbours(index int) []*chitonNode {
	row := index / g.width
	neighbours := []*chitonNode{}

	left, right, up, down := index-1, index+1, index-g.width, index+g.width

	if left >= 0 && left/g.width == row {
		neighbours = append(neighbours, g.nodes[left])
	}
	if right < len(g.nodes) && right/g.width == row {
		neighbours = append(neighbours, g.nodes[right])
	}
	if up >= 0 {
		neighbours = append(neighbours, g.nodes[up])
	}
	if down < len(g.nodes) {
		neighbours = append(neighbours, g.nodes[down])
	}

	return neighbours
}

func (q *chitonPQueue) put(n *chitonNode) {
	q.elems = append(q.elems, n)
	sort.Slice(q.elems, func(a int, b int) bool {
		return q.elems[a].f < q.elems[b].f
	})

	var builder strings.Builder
	for _, elem := range q.elems {
		builder.WriteString(fmt.Sprintf("%d->", elem.f))
	}
	fmt.Println(builder.String())
}

func (q *chitonPQueue) get() *chitonNode {
	first := q.elems[0]
	q.elems = q.elems[1:]
	return first
}

func (q *chitonPQueue) length() int {
	return len(q.elems)
}

func (g chitonGrid) printBigMap() {
	origWidth := g.width / 5
	origHeight := (len(g.nodes) / g.width) / 5
	var curLine strings.Builder
	for i := 0; i < len(g.nodes); i++ {
		curLine.WriteString(strconv.Itoa(g.nodes[i].risk))
		if (i+1)%origWidth == 0 {
			curLine.WriteString("|")
		}

		if (i+1)/g.width > i/g.width {
			fmt.Println(curLine.String())
			curLine = strings.Builder{}
			if ((i+1)/g.width)%origHeight == 0 {
				for i := 0; i < g.width+5; i++ {
					curLine.WriteString("-")
				}
				fmt.Println(curLine.String())
				curLine = strings.Builder{}
			}
		}
	}
}

func (g *chitonGrid) printPath() {
	var builder strings.Builder
	cur := g.nodes[len(g.nodes)-1]
	totalRisk := 0
	steps := 0
	for cur != g.nodes[0] {
		steps++
		totalRisk += cur.risk
		neighbours := g.getNeighbours(cur.index)
		minDistance := cur.distance
		var next *chitonNode
		for _, n := range neighbours {
			if n.distance < minDistance {
				minDistance = n.distance
				next = n
			}
		}
		builder.WriteString(fmt.Sprintf("%d -> ", next.risk))
		cur = next
	}
	fmt.Println(builder.String())
	fmt.Println(totalRisk)
}
