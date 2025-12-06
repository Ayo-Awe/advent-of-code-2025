package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

func ParseInput(input string) ([][2]int, []int, error) {
	parts := strings.Split(input, "\n\n")
	sRanges := strings.Split(strings.TrimSuffix(parts[0], "\n"), "\n")
	sIngredients := strings.Split(strings.TrimSuffix(parts[1], "\n"), "\n")

	ranges := make([][2]int, len(sRanges))
	for i, sRange := range sRanges {
		parts := strings.Split(sRange, "-")

		lower, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, nil, err
		}

		upper, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil, err
		}

		ranges[i] = [2]int{lower, upper}
	}

	ingredients := make([]int, len(sIngredients))
	for i, sIngredient := range sIngredients {
		ingredient, err := strconv.Atoi(sIngredient)
		if err != nil {
			return nil, nil, err
		}

		ingredients[i] = ingredient
	}

	return ranges, ingredients, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	input, err := aoc.ReadInput(*filename)
	if err != nil {
		log.Fatal(err)
	}

	ranges, ingredients, err := ParseInput(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("solution to part one: ", PartOne(ranges, ingredients))
	fmt.Println("solution to part two: ", PartTwo(ranges))
}

func PartOne(ranges [][2]int, ingredients []int) int {
	// brute force
	var fresh int

	for _, ingr := range ingredients {
		for _, ingRange := range ranges {
			lower, upper := ingRange[0], ingRange[1]
			// check if ingredient lies within this range
			if ingr >= lower && ingr <= upper {
				fresh++
				break
			}
		}
	}

	return fresh
}

func PartTwo(ranges [][2]int) int {
	ranges = slices.Clone(ranges)

	for i := range ranges {
		for j := i + 1; j < len(ranges); j++ {
			if overlap(ranges[i], ranges[j]) {
				ranges[j] = merge(ranges[i], ranges[j])

				// empty out ranges[i] as it's now ineffective
				// since it has been merged with ranges[j]
				ranges[i] = [2]int{1, 0} // value := 0 - 1 + 1 = 0
				break
			}
		}
	}

	var total int
	for i := range len(ranges) {
		total += ranges[i][1] - ranges[i][0] + 1
	}

	return total
}

func overlap(a, b [2]int) bool {
	start, end := 0, 1
	noOverlap := (a[start] < b[start] && a[end] < b[start]) ||
		(a[start] > b[end] && a[end] > b[end])
	return !noOverlap
}

func merge(a, b [2]int) [2]int {
	return [2]int{min(a[0], b[0]), max(a[1], b[1])}
}
