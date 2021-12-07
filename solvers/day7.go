package solvers

import (
	"fmt"
	"sort"
	"strings"
)

type Day7 struct{}

func (d Day7) Solve(input []string, part int) (int64, error) {
	nums, err := inputsToInt(strings.Split(input[0], ","))
	if err != nil {
		return -1, err
	}

	// I've done a quick and dirty brute force solution to both parts
	// because I'm still not 100% sure why my solutions work (especially part 2)
	// or if they work for all inputs. So far it's consistently produced
	// identical results between the brute force method (definitely correct) and my solution
	switch part {
	case 1:

		// d.part1BruteForce(nums)
		return d.part1(nums)
	case 2:
		// d.part2BruteForce(nums)
		return d.part2(nums)
	}

	return -1, invalidPartError(part)
}

func (d Day7) part1(nums []int) (int64, error) {
	sort.Ints(nums)
	median := nums[len(nums)/2]
	total := 0
	for _, num := range nums {
		total += absInt(median - num)
	}

	return int64(total), nil
}

func (d Day7) part2(nums []int) (int64, error) {
	sum := sumIntSlice(nums)
	roundedMean := (sum + len(nums) - 1) / len(nums)

	results := []int{}
	for i := -1; i <= 0; i++ {
		total := 0
		for _, num := range nums {
			distance := absInt((roundedMean + i) - num)
			total += (distance * (distance + 1)) / 2
		}
		results = append(results, total)
	}

	sort.Ints(results)
	return int64(results[0]), nil
}

func (d Day7) part1BruteForce(nums []int) {
	results := make([]int, 0, len(nums))
	sort.Ints(nums)
	for i := 0; i <= nums[len(nums)-1]; i++ {
		total := 0
		for _, num := range nums {
			total += absInt(i - num)
		}
		results = append(results, total)
	}

	sort.Ints(results)
	fmt.Println(results[0])
}

func (d Day7) part2BruteForce(nums []int) {
	results := make([]int, 0, len(nums))
	sort.Ints(nums)
	for i := 0; i < nums[len(nums)-1]; i++ {
		total := 0
		for _, num := range nums {
			distance := absInt(i - num)
			total += (distance * (distance + 1)) / 2
		}
		results = append(results, total)
	}

	sort.Ints(results)
	fmt.Println(results[0])
}
