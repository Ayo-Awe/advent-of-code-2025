package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

func ParseInput(lines []string) [][]rune {
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		grid[i] = []rune(line)
	}
	return grid
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	grid := ParseInput(lines)

	fmt.Println("solution to part one: ", PartOne(grid))
	fmt.Println("solution to part two: ", PartTwo(grid))
}

func PartOne(grid [][]rune) int {
	var count int

	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == '@' && neighRolls(x, y, grid) < 4 {
				count++
			}
		}
	}

	return count
}

func PartTwo(grid [][]rune) int {
	var count int

	for {
		var removed [][2]int

		for y := range grid {
			for x := range grid[y] {
				if grid[y][x] == '@' && neighRolls(x, y, grid) < 4 {
					count++
					removed = append(removed, [2]int{x, y})
				}
			}
		}

		if len(removed) == 0 {
			break
		}

		for _, pos := range removed {
			x, y := pos[0], pos[1]
			grid[y][x] = 'x'
		}
	}

	return count
}

func neighRolls(x, y int, grid [][]rune) int {
	var count int

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			// home position
			if dx == 0 && dy == 0 {
				continue
			}

			// out of bounds
			nx, ny := x+dx, y+dy
			if nx < 0 || nx >= len(grid[0]) || ny < 0 || ny >= len(grid) {
				continue
			}

			if grid[ny][nx] == '@' {
				count++
			}
		}
	}

	return count
}
