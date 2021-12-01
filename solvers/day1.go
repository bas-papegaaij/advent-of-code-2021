package solvers

type Day1 struct{}

func (d Day1) Solve(input []string, part int) (int64, error) {
	return genericSolve(d, input, part)
}

func (d Day1) part1(input []string) (int64, error) {
	nums, err := inputsToInt(input)
	if err != nil {
		return -1, err
	}

	increaseCount := 0
	for i, num := range nums[1:] {
		if num > nums[i] {
			increaseCount += 1
		}
	}

	return int64(increaseCount), nil
}

func (d Day1) part2(input []string) (int64, error) {
	nums, err := inputsToInt(input)
	if err != nil {
		return -1, err
	}

	increaseCount := 0
	lastWin := sumIntSlice(nums[:3])
	for i := 1; i < len(nums)-2; i++ {
		curWin := sumIntSlice(nums[i : i+3])
		if curWin > lastWin {
			increaseCount++
		}

		lastWin = curWin
	}

	return int64(increaseCount), nil
}

func sumIntSlice(vals []int) int {
	sum := 0
	for _, num := range vals {
		sum += num
	}
	return sum
}
