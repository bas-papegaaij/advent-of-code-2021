package solvers

import "strconv"

type Day3 struct{}

func (d Day3) Solve(input []string, part int) (int64, error) {
	switch part {
	case 1:
		return d.part1(input)
	case 2:
		return d.part2(input)
	}
	return -1, invalidPartError(part)
}

func (d Day3) part1(input []string) (int64, error) {
	gamma := 0
	epsilon := 0

	for i := 0; i < len(input[0]); i++ {
		off, on := d.getBitCount(input, i)
		if on >= off {
			gamma |= 1 << (len(input[0]) - 1 - i)
		} else {
			epsilon |= 1 << (len(input[0]) - 1 - i)
		}
	}

	return int64(gamma * epsilon), nil

	// First attempt of part 1 was below - possibly a bit over-engineered
	// bitCounts, err := d.getBitCounts(input)
	// if err != nil {
	// 	return -1, err
	// }

	// var gamma uint64
	// for i, count := range bitCounts {
	// 	if count >= len(input)/2 {
	// 		gamma |= 1 << uint64(len(input[0])-i-1)
	// 	}
	// }

	// // create a with the length of each bitfield
	// // we can XOR that together with the original to get the inverse, ignoring leading 0-bits
	// // Since the least common bit in each position is always guaranteed to be the opposite of the most common bit
	// // this will give us our epsilon
	// var mask uint64
	// mask -= 1
	// mask <<= len(input[0])
	// mask = ^mask
	// epsilon := gamma ^ mask
	// return int64(gamma * epsilon), nil
}

func (d Day3) part2(input []string) (int64, error) {
	o2GeneratorPredicate := func(offbits int, onBits int) byte {
		if onBits >= offbits {
			return '1'
		}
		return '0'
	}

	c02ScrubberPredicate := func(offBits int, onBits int) byte {
		if offBits <= onBits {
			return '0'
		}
		return '1'
	}

	o2value, err := d.getRating(input, o2GeneratorPredicate)
	if err != nil {
		return -1, err
	}
	scrubberValue, err := d.getRating(input, c02ScrubberPredicate)
	if err != nil {
		return -1, err
	}

	return o2value * scrubberValue, nil
}

type targetBitPredicate func(offBits int, onBits int) byte

func (d Day3) getRating(input []string, bitPredicate targetBitPredicate) (int64, error) {
	for bitPos := 0; len(input) > 1; bitPos++ {
		input = d.eliminateEntries(input, bitPos, bitPredicate)
	}

	return strconv.ParseInt(input[0], 2, 64)
}

func (d Day3) eliminateEntries(input []string, bitPos int, bitPredicate targetBitPredicate) []string {
	var remainingEntries []string
	offBits, onBits := d.getBitCount(input, bitPos)
	targetBit := bitPredicate(offBits, onBits)
	for _, str := range input {
		if str[bitPos] == targetBit {
			remainingEntries = append(remainingEntries, str)
		}
	}

	return remainingEntries
}

// func (d Day3) getBitCounts(input []string) ([]int, error) {
// 	// assuming all input bit strings are the same length
// 	bitCounts := make([]int, len(input[0]))
// 	for _, bitstr := range input {
// 		for i, bit := range bitstr {
// 			if bit == '1' {
// 				bitCounts[i]++
// 			}
// 		}
// 	}

// 	return bitCounts, nil
// }

func (d Day3) getBitCount(input []string, bitPosition int) (int, int) {
	onCount := 0
	offCount := 0
	for _, bitStr := range input {
		if bitStr[bitPosition] == '0' {
			offCount++
		} else {
			// We'll assume our input is all valid bit strings
			onCount++
		}
	}

	return offCount, onCount
}
