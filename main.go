package main

import (
	"aoc-2021/solvers"
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Solver interface {
	Solve(input []string, part int) (int64, error)
}

var solutions []Solver = []Solver{
	solvers.Day1{},
	solvers.Day2{},
}

var inputFile string
var questionPart int
var question int

func main() {
	parseArgs()

	lines, err := parseLines()
	if err != nil {
		panic(err)
	}

	answer, err := solutions[question].Solve(lines, questionPart)
	if err != nil {
		panic(err)
	}

	fmt.Println("Answer to question part", questionPart, "was", answer)
}

func parseLines() ([]string, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func parseArgs() {
	flag.StringVar(&inputFile, "input", "input.txt", "Location of the input file, relative to the working directory")
	flag.IntVar(&question, "question", 0, "Which question to answer")
	flag.IntVar(&questionPart, "part", 1, "The part of the question to answer - 1 or 2")

	flag.Parse()

	if questionPart < 1 || questionPart > 2 {
		panic("part must be either 1 or 2")
	}

	if question == 0 || question > len(solutions) {
		fmt.Println("question must be between 1 and", len(solutions))
	}
	question -= 1
}
