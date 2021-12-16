package main

import (
	"aoc-2021/solvers"
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

type Solver interface {
	Solve(input []string, part int) (int64, error)
}

var solutions []Solver = []Solver{
	solvers.Day1{},
	solvers.Day2{},
	solvers.Day3{},
	solvers.Day4{},
	solvers.Day5{},
	solvers.Day6{},
	solvers.Day7{},
	solvers.Day8{},
	solvers.Day9{},
	solvers.Day10{},
	solvers.Day11{},
	solvers.Day12{},
	solvers.Day13{},
	solvers.Day14{},
	solvers.Day15{},
	solvers.Day16{},
}

var inputFile string
var questionPart int
var question int
var runProfile bool

func main() {
	parseArgs()

	start := time.Now()
	lines, err := parseLines()
	if err != nil {
		panic(err)
	}

	if runProfile {
		profileMode(lines)
	} else {
		solutionMode(lines)
	}
	fmt.Println("including parse, took:", time.Since(start))
}

func solutionMode(lines []string) {
	start := time.Now()
	answer, err := solutions[question].Solve(lines, questionPart)
	if err != nil {
		panic(err)
	}
	fmt.Println("solution took", time.Since(start))
	fmt.Println("Answer to question part", questionPart, "was", answer)
}

func profileMode(lines []string) {
	runs := make([]time.Duration, 0, 1000)
	solver := solutions[question]
	for i := 0; i < 1000; i++ {
		start := time.Now()
		solver.Solve(lines, questionPart)
		runs = append(runs, time.Since(start))
	}

	var min, max, total time.Duration
	for i, runtime := range runs {
		if i == 0 {
			min = runtime
		}
		if runtime < min {
			min = runtime
		}
		if runtime > max {
			max = runtime
		}
		total += runtime
	}
	avg := time.Duration(total.Nanoseconds() / int64(len(runs)))
	fmt.Println("min:", min, "max:", max, "avg:", avg)
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
	flag.BoolVar(&runProfile, "prof", false, "Whether to run a profile. If enabled runs the solution 1000 times and grabs an average, min and max runtimes")

	flag.Parse()

	if questionPart < 1 || questionPart > 2 {
		panic("part must be either 1 or 2")
	}

	if question == 0 || question > len(solutions) {
		fmt.Println("question must be between 1 and", len(solutions))
	}
	question -= 1
}
