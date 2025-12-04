package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

func ParseInput(lines []string) [][]int {
	banks := make([][]int, len(lines))
	for i := range lines {
		bank := make([]int, len(lines[i]))
		for j := range lines[i] {
			bank[j] = int(lines[i][j] - '0')
		}
		banks[i] = bank
	}
	return banks
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	banks := ParseInput(lines)
	fmt.Println("solution to part one: ", PartOne(banks))
	fmt.Println("solution to part two: ", PartTwo(banks))
}

func PartOne(banks [][]int) int {
	var total int
	for _, bank := range banks {
		total += joltage(bank, 2)
	}
	return total
}

func PartTwo(banks [][]int) int {
	var total int
	for _, bank := range banks {
		total += joltage(bank, 12)
	}
	return total
}

func joltage(bank []int, batteries int) int {
	var joltage int

	// start inclusive, end exclusive
	start := 0
	for i := batteries - 1; i >= 0; i-- {
		end := len(bank) - i

		// find index of the max element in the given partition
		maxIdx := start
		for idx := start; idx < end; idx++ {
			if bank[idx] > bank[maxIdx] {
				maxIdx = idx
			}
		}

		joltage += int(math.Pow10(i)) * bank[maxIdx]
		start = maxIdx + 1
	}

	return joltage
}

