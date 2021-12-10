package solvers

import (
	"fmt"
	"strconv"
	"strings"
)

type Day5 struct{}

type lineSegment struct {
	start intVector
	end   intVector
}

type ventMap struct {
	values []int
	width  int
}

func (d Day5) Solve(input []string, part int) (int64, error) {
	segments, err := d.inputToLineSegments(input)
	if err != nil {
		return -1, err
	}
	switch part {
	case 1:
		segments = d.getStraightLines(segments)
	case 2:
		break
	default:
		return -1, invalidPartError(part)
	}

	vents := d.makeVentMap(segments)
	vents.markAllSpots(segments)

	return int64(vents.countDangerousCoords()), nil
}

func (d Day5) getStraightLines(segments []lineSegment) []lineSegment {
	res := make([]lineSegment, 0)
	for _, segment := range segments {
		// If either x stays the same, or y stays the same, we have a straight line
		if segment.start.x == segment.end.x || segment.start.y == segment.end.y {
			res = append(res, segment)
		}
	}
	return res
}

func (d Day5) getMaxExtents(segments []lineSegment) intVector {
	res := intVector{}
	for _, segment := range segments {
		if segment.start.x > res.x {
			res.x = segment.start.x
		}
		if segment.end.x > res.x {
			res.x = segment.end.x
		}

		if segment.start.y > res.y {
			res.y = segment.start.y
		}
		if segment.end.y > res.y {
			res.y = segment.end.y
		}
	}

	return res
}

func (d Day5) inputToLineSegments(input []string) ([]lineSegment, error) {
	segments := make([]lineSegment, 0, len(input))
	for _, line := range input {
		segment, err := d.parseLineSegment(line)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}
	return segments, nil
}
func (d Day5) parseLineSegment(input string) (lineSegment, error) {
	ends := strings.Split(input, "->")
	start, err := d.stringToIntVector(ends[0])
	if err != nil {
		return lineSegment{}, err
	}

	end, err := d.stringToIntVector(ends[1])
	if err != nil {
		return lineSegment{}, err
	}

	return lineSegment{
		start,
		end,
	}, nil
}

func (d Day5) stringToIntVector(input string) (intVector, error) {
	nums := strings.Split(strings.TrimSpace(input), ",")
	if len(nums) != 2 {
		return intVector{}, fmt.Errorf("invalid intvector input %s, expected 2 numbers, got %d", input, len(input))
	}

	x, err := strconv.Atoi(nums[0])
	if err != nil {
		return intVector{}, err
	}

	y, err := strconv.Atoi(nums[1])
	if err != nil {
		return intVector{}, err
	}

	return intVector{
		x,
		y,
	}, nil
}

func (d Day5) makeVentMap(segments []lineSegment) ventMap {
	maxExtents := d.getMaxExtents(segments)
	field := make([]int, (maxExtents.x+1)*(maxExtents.y+1))
	return ventMap{
		values: field,
		width:  maxExtents.x,
	}
}

func (v ventMap) markSpot(x int, y int) {
	v.values[y*v.width+x]++
}

func (v ventMap) markSpotsInLine(segment lineSegment) {
	x := segment.start.x
	xIncrement := segment.end.x - segment.start.x
	if xIncrement != 0 {
		xIncrement /= absInt(xIncrement)
	}

	y := segment.start.y
	yIncrement := segment.end.y - segment.start.y
	if yIncrement != 0 {
		yIncrement /= absInt(yIncrement)
	}

	v.markSpot(x, y)
	for x != segment.end.x || y != segment.end.y {
		x += xIncrement
		y += yIncrement
		v.markSpot(x, y)
	}
}

func (v ventMap) markAllSpots(segments []lineSegment) {
	for _, segment := range segments {
		v.markSpotsInLine(segment)
	}
}

func (v ventMap) countDangerousCoords() int {
	total := 0
	for _, value := range v.values {
		if value > 1 {
			total++
		}
	}

	return total
}
