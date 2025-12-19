package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

const (
	X, Y = 0, 1
)

func ParseInput(lines []string) ([][2]int, error) {
	corners := make([][2]int, len(lines))
	for i := range lines {
		parts := strings.Split(lines[i], ",")

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		y, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		corners[i] = [2]int{x, y}
	}
	return corners, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	corners, err := ParseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("solution to part one: ", PartOne(corners))
	fmt.Println("solution to part two: ", PartTwo(corners))
}

func PartOne(corners [][2]int) int {
	var maxArea int
	for i := range corners {
		for j := i + 1; j < len(corners); j++ {
			a := corners[i]
			b := corners[j]
			dx := (max(a[X], b[X]) - min(a[X], b[X]) + 1)
			dy := (max(a[Y], b[Y]) - min(a[Y], b[Y]) + 1)
			maxArea = max(maxArea, dx*dy)
		}
	}
	return maxArea
}

func PartTwo(redTiles [][2]int) int {
	xset := make(map[int]struct{})
	yset := make(map[int]struct{})

	// slice of unique x and y coordinates
	xs := make([]int, 0, len(redTiles))
	ys := make([]int, 0, len(redTiles))

	for i := range redTiles {
		x, y := redTiles[i][X], redTiles[i][Y]

		if _, exists := xset[x]; !exists {
			xs = append(xs, x)
			xset[x] = struct{}{}
		}

		if _, exists := yset[y]; !exists {
			ys = append(ys, y)
			yset[y] = struct{}{}
		}
	}

	// sort both sets
	sort.Ints(xs)
	sort.Ints(ys)

	// lookup to translate  real coordiantes to compressed coordinates
	xlookup := make(map[int]int)
	ylookup := make(map[int]int)

	for i := range len(ys) {
		ylookup[ys[i]] = i
	}

	for i := range len(xs) {
		xlookup[xs[i]] = i
	}

	grid := make([][]rune, len(ys))
	for i := range grid {
		grid[i] = []rune(strings.Repeat(".", len(xs)))
	}

	// build compressed grid
	for i := range redTiles {
		a, b := redTiles[i], redTiles[(i+1)%len(redTiles)]

		// compressed coordinates of tile a and b
		cax, cay := xlookup[a[X]], ylookup[a[Y]]
		cbx, cby := xlookup[b[X]], ylookup[b[Y]]

		// mark current tile (a) as red
		grid[cay][cax] = '#'

		// mark all tiles between a and b as green
		// NOTE: adjacent tiles in the input are either on the same row or column
		if cay-cby == 0 {
			// tile a & b are on the same row
			for x := min(cax, cbx) + 1; x < max(cax, cbx); x++ {
				grid[cay][x] = 'X'
			}
		} else {
			// tile a & b are on the same column
			for y := min(cay, cby) + 1; y < max(cay, cby); y++ {
				grid[y][cax] = 'X'
			}
		}
	}

	// keep a set of all points outside the polygon
	outside := make(map[[2]int]struct{})
	queue := [][2]int{{-1, -1}} // we add an imaginary padding around the grid

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		x, y := curr[X], curr[Y]

		// bounds check factoring in the imaginary padding
		if x < -1 || x > len(grid[0]) || y < -1 || y > len(grid) {
			continue
		}

		// for non-padding tiles, skip red/green tiles
		if x >= 0 && x < len(grid[0]) && y >= 0 && y < len(grid) {
			if grid[y][x] == 'X' || grid[y][x] == '#' {
				continue
			}
		}

		if _, seen := outside[curr]; seen {
			continue
		}

		outside[curr] = struct{}{}

		cards := [][2]int{
			{-1, 0},
			{1, 0},
			{0, -1},
			{0, 1},
		}

		for _, card := range cards {
			nx, ny := x+card[X], y+card[Y]
			queue = append(queue, [2]int{nx, ny})
		}
	}

	var maxArea int
	for i := range redTiles {
		for j := i + 1; j < len(redTiles); j++ {
			a := redTiles[i]
			b := redTiles[j]
			dx := (max(a[X], b[X]) - min(a[X], b[X]) + 1)
			dy := (max(a[Y], b[Y]) - min(a[Y], b[Y]) + 1)

			// verify that all grid points are within the polygon
			cax, cay := xlookup[a[X]], ylookup[a[Y]]
			cbx, cby := xlookup[b[X]], ylookup[b[Y]]

			if !isAreaWithinPolygon([2]int{cax, cay}, [2]int{cbx, cby}, outside) {
				continue
			}

			maxArea = max(maxArea, dx*dy)
		}
	}

	return maxArea
}

func isAreaWithinPolygon(a, b [2]int, outside map[[2]int]struct{}) bool {
	cax, cay := a[X], a[Y]
	cbx, cby := b[X], b[Y]

	for y := min(cay, cby); y <= max(cay, cby); y++ {
		for x := min(cax, cbx); x <= max(cax, cbx); x++ {
			if _, exists := outside[[2]int{x, y}]; exists {
				return false
			}
		}
	}
	return true
}
