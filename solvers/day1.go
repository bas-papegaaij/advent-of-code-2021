package solvers

type Day1 struct{}

func (d Day1) Solve(input []string, part int) (int64, error) {
	switch part {
	case 1:
		return d.solve(input, 1)
	case 2:
		return d.solve(input, 3)
	}

	return -1, invalidPartError(part)
}

func (d Day1) solve(input []string, windowSize int) (int64, error) {
	nums, err := inputsToInt(input)
	if err != nil {
		return -1, err
	}

	increaseCount := 0
	lastWin := sumIntSlice(nums[:windowSize])
	for i := 1; i < len(nums)-(windowSize-1); i++ {
		curWin := sumIntSlice(nums[i : i+windowSize])
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
