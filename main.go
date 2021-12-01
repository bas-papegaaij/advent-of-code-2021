package main

import (
	"flag"
	"fmt"
)

type AocSolver interface {
	Solve(input []string, part int) string
}

var solvers []AocSolver = []AocSolver{}

var inputFile string
var questionPart int
var question int

func main() {
	parseArgs()
	lines, err := parseLines()
	if err != nil {
		panic(err)
	}

	answer := solvers[question].Solve(lines, questionPart)
	fmt.Println("Answer to question part", questionPart, "was", answer)
}

func parseLines() ([]string, error) {
	return nil, nil
}

func parseArgs() {
	flag.StringVar(&inputFile, "input", "input.txt", "Location of the input file, relative to the working directory")
	flag.IntVar(&question, "question", 0, "Which question to answer")
	flag.IntVar(&questionPart, "part", 1, "The part of the question to answer - 1 or 2")

	if questionPart < 1 || questionPart > 2 {
		panic("part must be either 1 or 2")
	}

	if question == 0 || question > len(solvers) {
		fmt.Println("question must be between 1 and", len(solvers))
	}
	question -= 1
}
