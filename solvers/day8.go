package solvers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Day8 struct {
}

type digitSignal struct {
	inputs  []int
	outputs []int
}

var digitsMap map[int]string = map[int]string{
	0b1110111: "0",
	0b0100100: "1",
	0b1011101: "2",
	0b1101101: "3",
	0b0101110: "4",
	0b1101011: "5",
	0b1111011: "6",
	0b0100101: "7",
	0b1111111: "8",
	0b1101111: "9",
}

func (d Day8) Solve(input []string, part int) (int64, error) {
	signals := d.parseInput(input)
	switch part {
	case 1:
		return d.part1(signals)
	case 2:
		return d.part2(signals)
	}

	return -1, invalidPartError(part)
}

func (d Day8) parseInput(input []string) []digitSignal {
	signals := make([]digitSignal, 0, len(input))
	for _, str := range input {
		split := strings.Split(str, "|")
		in := strings.TrimSpace(split[0])
		out := strings.TrimSpace(split[1])
		signals = append(signals, digitSignal{
			inputs:  d.stringSliceToBinary(strings.Split(in, " ")),
			outputs: d.stringSliceToBinary(strings.Split(out, " ")),
		})
	}
	return signals
}

func (d Day8) part1(input []digitSignal) (int64, error) {
	results := make([]int, 10)
	for _, signal := range input {
		for _, output := range signal.outputs {
			switch d.countSetBits(output) {
			case 2:
				results[1]++
			case 3:
				results[7]++
			case 4:
				results[4]++
			case 7:
				results[8]++
			}
		}
	}

	return int64(sumIntSlice(results)), nil
}

func (d Day8) part2(input []digitSignal) (int64, error) {
	sum := 0
	for _, entry := range input {
		signalsMap := d.getsignalMap(entry)
		sum += d.entryToNumber(signalsMap, entry)
	}
	return int64(sum), nil
}

func (d Day8) getsignalMap(signal digitSignal) map[int]int {
	// maps each position on the display to the associated one in our inputs for this entry
	// TODO: build reverse lookup so we can convert
	signalMap := make(map[int]int)
	sort.Slice(signal.inputs, func(i int, j int) bool { return d.countSetBits(signal.inputs[i]) < d.countSetBits(signal.inputs[j]) })
	numbers := make([]int, len(signal.inputs))

	// We know these numbers match up with the given signals for sure, since they are of a unique length
	numbers[1] = signal.inputs[0]
	numbers[4] = signal.inputs[2]
	numbers[7] = signal.inputs[1]
	numbers[8] = signal.inputs[9]

	// We can easily find the top 'a' signal, since it's the only signal not common between 1 and 7
	signalMap[1<<0] = d.findDiff(numbers[1], numbers[7])
	// if we compare number 4 with the top signal set with all 6-length numbers
	// we should get a diff of 1 only with the number 9, which will be "g"
	numbers[9], signalMap[1<<6] = d.getNumWithDiff(signal.inputs[6:9], numbers[4]|signalMap[1<<0], 1)
	// If we add "a" and "g" to 1, we should get a diff of only 1 with just the number 3
	// which will give us the "d" character
	numbers[3], signalMap[1<<3] = d.getNumWithDiff(signal.inputs[3:6], numbers[1]|signalMap[1<<0]|signalMap[1<<6], 1)
	// we can now find "b", if we add "d" to the number 1 and diff it with the number 4
	signalMap[1<<1] = d.findDiff(numbers[1]|signalMap[1<<3], numbers[4])
	// We can now construct a number that has a diff of 1 with 5, also identifying "f"
	numbers[5], signalMap[1<<5] = d.getNumWithDiff(signal.inputs[3:6], signalMap[1<<0]+signalMap[1<<1]+signalMap[1<<3]+signalMap[1<<6], 1)
	// now we can find the other character in the number 1
	signalMap[1<<2] = d.findDiff(numbers[1], signalMap[1<<5])
	// finally we can diff number 9 with number 8 to get "e"
	signalMap[1<<4] = d.findDiff(numbers[9], numbers[8])

	// need to reverse the map so we can actually look up our characters in it and get the resultant position
	result := make(map[int]int)
	for k, v := range signalMap {
		result[v] = k
	}
	return result
}

func (d Day8) getNumWithDiff(testNumbers []int, test int, expectedDiff int) (int, int) {
	for _, num := range testNumbers {
		diff := d.findDiff(num, test)
		if d.countSetBits(diff) == expectedDiff {
			return num, diff
		}
	}
	panic(fmt.Sprintf("no number in provided list has diff of length %d with %b", expectedDiff, test))
}

func (d Day8) findDiff(a int, b int) int {
	return a ^ b
}

func (d Day8) stringSliceToBinary(strings []string) []int {
	res := make([]int, 0, len(strings))
	for _, str := range strings {
		res = append(res, d.stringToBinary(str))
	}
	return res
}

// This just lets us create an easy lookup table
func (d Day8) stringToBinary(outputStr string) int {
	num := 0
	for _, r := range outputStr {
		num |= 1 << int(r-'a')
	}
	return num
}

func (d Day8) convert(signalsMap map[int]int, input int) int {
	res := 0
	for i := 0; i < 7; i++ {
		mask := 1 << i
		if input&mask != 0 {
			res |= signalsMap[mask]
		}
	}
	return res
}

func (d Day8) entryToNumber(signalsMap map[int]int, entry digitSignal) int {
	str := ""
	for _, output := range entry.outputs {
		converted := d.convert(signalsMap, output)
		digit, ok := digitsMap[converted]
		if !ok {
			panic("couldn't find digit in map")
		}

		str += digit
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}

	return num
}

func (d Day8) countSetBits(val int) int {
	count := 0
	for val > 0 {
		val = val & (val - 1)
		count++
	}
	return count
}
