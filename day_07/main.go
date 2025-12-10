package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

const (
	X, Y = 0, 1
)

func ParseInput(lines []string) ([][]rune, [2]int) {
	var start [2]int
	grid := make([][]rune, len(lines))

	for y := range lines {
		grid[y] = []rune(lines[y])
		for x := range lines[y] {
			if grid[y][x] == 'S' {
				start = [2]int{x, y}
			}
		}
	}

	return grid, start
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	grid, start := ParseInput(lines)

	fmt.Println("solution to part one: ", PartOne(grid, start))
	fmt.Println("solution to part two: ", PartTwo(grid, start))
}

func PartOne(grid [][]rune, start [2]int) int {
	seen := map[[2]int]bool{}
	queue := [][2]int{start}
	splits := map[[2]int]struct{}{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if seen[node] {
			continue
		}

		x, y := node[X], node[Y]
		for ny := y; ny < len(grid); ny++ {
			seen[node] = true

			if grid[ny][x] != '^' {
				continue
			}

			// split node
			xl, xr := x-1, x+1
			queue = append(queue, [2]int{xl, ny})
			queue = append(queue, [2]int{xr, ny})

			splits[[2]int{x, ny}] = struct{}{}
			break
		}

	}

	return len(splits)
}

func PartTwo(grid [][]rune, start [2]int) int {
	memo := make(map[[2]int]int)
	return timelines(grid, start, memo)
}

func timelines(grid [][]rune, start [2]int, memo map[[2]int]int) int {
	x, y := start[X], start[Y]

	for ny := y; ny < len(grid); ny++ {
		if grid[ny][x] == '^' {
			// we've seen this split before
			if branches, seen := memo[[2]int{x, ny}]; seen {
				return branches
			}

			// l and r are always within bounds based on the puzzle input
			lsplit := [2]int{x - 1, ny}
			rsplit := [2]int{x + 1, ny}

			ltimelines := timelines(grid, lsplit, memo)
			rtimelines := timelines(grid, rsplit, memo)

			memo[[2]int{x, ny}] = ltimelines + rtimelines
			return ltimelines + rtimelines
		}
	}

	return 1
}

