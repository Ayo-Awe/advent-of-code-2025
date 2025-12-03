package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

func ParseInput(input string) ([][2]int, error) {
	input = strings.TrimSpace(input)
	rangesStr := strings.Split(input, ",")

	ranges := make([][2]int, len(rangesStr))
	for i, rangeStr := range rangesStr {
		parts := strings.Split(rangeStr, "-")

		lower, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		upper, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		ranges[i] = [2]int{lower, upper}
	}
	return ranges, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	input, err := aoc.ReadInput(*filename)
	if err != nil {
		log.Fatal(err)
	}

	ranges, err := ParseInput(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("solution to part one: ", PartOne(ranges))
	fmt.Println("solution to part two: ", PartTwo(ranges))
}

func PartOne(ranges [][2]int) int {
	var total int

	for _, idRange := range ranges {
		lower, upper := idRange[0], idRange[1]
		for i := lower; i <= upper; i++ {
			iStr := fmt.Sprintf("%d", i)
			numDigits := len(iStr)
			if numDigits%2 == 1 {
				continue
			}

			if iStr[:numDigits/2] == iStr[numDigits/2:] {
				total += i
			}
		}
	}

	return total
}

func PartTwo(ranges [][2]int) int {
	var total int

	for _, idRange := range ranges {
		lower, upper := idRange[0], idRange[1]
		for i := lower; i <= upper; i++ {
			iStr := fmt.Sprintf("%d", i)
			if isRepeating(iStr) {
				total += i
			}
		}
	}

	return total
}

func isRepeating(s string) bool {
	if len(s) < 2 {
		return false
	}

	var ptr int
	seqLen := 1
	for i := 1; i < len(s); i++ {
		if s[ptr] == s[i] {
			ptr++
			continue
		}

		ptr = 0
		if s[ptr] == s[i] {
			ptr++
			seqLen = i
		} else {
			seqLen = 1 + i
		}
	}

	return ptr != 0 && len(s)%seqLen == 0
}
