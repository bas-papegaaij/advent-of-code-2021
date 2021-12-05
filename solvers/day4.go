package solvers

import (
	"fmt"
	"sort"
	"strings"
)

type Day4 struct{}

func (d Day4) Solve(input []string, part int) (int64, error) {
	bingoInput, err := parseBingoInput(input)
	if err != nil {
		return -1, err
	}

	switch part {
	case 1:
		return d.part1(bingoInput)
	case 2:
		return d.part2(bingoInput)
	}
	return -1, invalidPartError(part)
}

func (d Day4) part1(input bingoInput) (int64, error) {
	for _, num := range input.drawNumbers {
		for _, board := range input.boards {
			idx := board.markNumber(num)
			if idx != -1 && board.hasBingo(idx) {
				return int64(num * board.sumUnmarkedSquares()), nil
			}
		}
	}
	return -1, fmt.Errorf("no boards had a bingo after all numbers were drawn")
}

func (d Day4) part2(input bingoInput) (int64, error) {
	remainingBoards := input.boards
	for _, num := range input.drawNumbers {
		removeIdxs := []int{}
		for boardIdx, board := range remainingBoards {
			markedIdx := board.markNumber(num)
			if markedIdx != -1 && board.hasBingo(markedIdx) {
				if len(remainingBoards) > 1 {
					removeIdxs = append(removeIdxs, boardIdx)
				} else {
					return int64(num * board.sumUnmarkedSquares()), nil
				}
			}
		}

		// Remove finished boards so we're not wasting time checking them for completion
		// Sort the indexes to remove in descending order so we can do a quick swap and pop
		// for each remove. This avoids having to copy loads of elements when culling the
		// boards that have already won
		sort.Slice(removeIdxs, func(a int, b int) bool {
			return b < a
		})

		for _, idx := range removeIdxs {
			remainingBoards[idx] = remainingBoards[len(remainingBoards)-1]
			remainingBoards = remainingBoards[:len(remainingBoards)-1]
		}
	}
	return -1, fmt.Errorf("no boards had a bingo at end of part 2")
}

type bingoInput struct {
	drawNumbers []int
	boards      []bingoBoard
}

func parseBingoInput(input []string) (bingoInput, error) {
	draw := strings.Split(input[0], ",")
	drawNumbers, err := inputsToInt(draw)
	if err != nil {
		return bingoInput{}, err
	}

	boards := []bingoBoard{}
	for _, boardInput := range splitByEmptyLines(input[2:]) {
		board, err := parseBoard(boardInput)
		if err != nil {
			return bingoInput{}, nil
		}
		boards = append(boards, board)
	}

	return bingoInput{
		drawNumbers: drawNumbers,
		boards:      boards,
	}, nil
}

func parseBoard(input []string) (bingoBoard, error) {
	numbers := []int{}
	var width int
	for _, row := range input {
		curRow := strings.Split(row, " ")
		cleanRow := []string{}
		for _, subStr := range curRow {
			if subStr != "" {
				cleanRow = append(cleanRow, subStr)
			}
		}

		nums, err := inputsToInt(cleanRow)
		if err != nil {
			return bingoBoard{}, err
		}
		width = len(nums)
		numbers = append(numbers, nums...)
	}

	return newBoard(numbers, width), nil
}

type bingoSquare struct {
	marked bool
	value  int
}

type bingoBoard struct {
	squares []bingoSquare

	width int
}

func newBoard(numbers []int, width int) bingoBoard {
	b := bingoBoard{
		width:   width,
		squares: make([]bingoSquare, len(numbers)),
	}

	for i, num := range numbers {
		b.squares[i] = bingoSquare{
			value:  num,
			marked: false,
		}
	}

	return b
}

// Mark the given number on the board.
// returns the number's index on the board if it was present
// -1 otherwise
func (b bingoBoard) markNumber(number int) int {
	for i, square := range b.squares {
		if square.value == number {
			b.squares[i].marked = true
			return i
		}
	}
	return -1
}

// Checks if the square at the given index is part of a bingo
func (b bingoBoard) hasBingo(index int) bool {
	row := index / b.width
	if b.isRowComplete(row) {
		return true
	}

	column := index % b.width
	return b.isColumnComplete(column)
}

func (b bingoBoard) isRowComplete(rowIndex int) bool {
	for i := 0; i < b.width; i++ {
		square := b.squares[rowIndex*b.width+i]
		if square.marked == false {
			return false
		}
	}

	return true
}

func (b bingoBoard) isColumnComplete(columnIndex int) bool {
	for i := 0; i < len(b.squares)/b.width; i++ {
		square := b.squares[i*b.width+columnIndex]
		if square.marked == false {
			return false
		}
	}

	return true
}

func (b bingoBoard) sumUnmarkedSquares() int {
	sum := 0
	for _, square := range b.squares {
		if square.marked == false {
			sum += square.value
		}
	}

	return sum
}
