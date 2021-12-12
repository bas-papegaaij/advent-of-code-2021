package solvers

import (
	"sort"
	"strings"
)

type Day10 struct{}

var closingChars map[string]struct{} = map[string]struct{}{
	")": {},
	"]": {},
	"}": {},
	">": {},
}

var chunkMap map[string]string = map[string]string{
	"(": ")",
	"[": "]",
	"{": "}",
	"<": ">",
}

var corruptScoreLookup map[string]int = map[string]int{
	")": 3,
	"]": 57,
	"}": 1197,
	">": 25137,
}

var completionScoreLookup map[string]int = map[string]int{
	")": 1,
	"]": 2,
	"}": 3,
	">": 4,
}

type charStack struct {
	chars string
}

type instructionState struct {
	stack *charStack
	input string
}

func (d Day10) Solve(input []string, part int) (int64, error) {
	corruptScore, states := d.discardCorruptLines(input)

	switch part {
	case 1:
		return int64(corruptScore), nil
	case 2:
		return int64(d.calculateCompletionScore(states)), nil
	}

	return -1, invalidPartError(part)
}

func (d Day10) calculateCompletionScore(states []instructionState) int {
	scores := make([]int, 0, len(states))
	for _, state := range states {
		scores = append(scores, state.calculateCompletionScore())
	}
	sort.Ints(scores)
	return scores[len(scores)/2]
}

func (s instructionState) calculateCompletionScore() int {
	score := 0
	c := s.stack.pop()
	completionString := ""
	for c != "" {
		completion := chunkMap[c]
		completionString += completion
		score = score*5 + completionScoreLookup[completion]
		c = s.stack.pop()
	}
	return score
}

// assumes ordering is not important
func (d Day10) discardCorruptLines(input []string) (int, []instructionState) {
	states := make([]instructionState, 0, len(input))
	corruptScore := 0
	for i := 0; i < len(input); i++ {
		line := input[i]
		corrupt, state := d.findCorruptChar(line)
		if corrupt != "" {
			corruptScore += corruptScoreLookup[corrupt]
		} else {
			states = append(states, state)
		}
	}
	return corruptScore, states
}

func (d Day10) findCorruptChar(line string) (string, instructionState) {
	stack := charStack{}
	for _, c := range strings.Split(line, "") {
		_, isClosing := closingChars[c]
		if isClosing {
			back := stack.peek()
			if back != "" && chunkMap[back] != c {
				return c, instructionState{&stack, line}
			}
			if back != "" {
				stack.pop()
			}
		} else {
			stack.push(c)
		}
	}

	return "", instructionState{&stack, line}
}

func (c *charStack) peek() string {
	if len(c.chars) == 0 {
		return ""
	}
	return string(c.chars[len(c.chars)-1])
}
func (c *charStack) pop() string {
	last := c.peek()
	if last != "" {
		c.chars = c.chars[:len(c.chars)-1]
	}
	return last
}

func (c *charStack) push(val string) {
	c.chars += val
}
