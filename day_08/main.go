package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

const (
	X = iota
	Y
	Z
)

func ParseInput(lines []string) ([][3]int, error) {
	jboxes := make([][3]int, len(lines))
	for i := range lines {
		dims := strings.Split(lines[i], ",")

		var jbox [3]int
		for j := range dims {
			dim, err := strconv.Atoi(dims[j])
			if err != nil {
				return nil, err
			}
			jbox[j] = dim
		}

		jboxes[i] = jbox
	}

	return jboxes, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	jboxes, err := ParseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	// we could optimize this and preallocate nCr size
	pairs := make([][2]int, 0, len(jboxes))
	for i := 0; i < len(jboxes)-1; i++ {
		for j := i + 1; j < len(jboxes); j++ {
			pairs = append(pairs, [2]int{i, j})
		}
	}

	// sort pairs by distance ascending
	slices.SortFunc(pairs, func(a, b [2]int) int {
		distA := dist(jboxes[a[0]], jboxes[a[1]])
		distB := dist(jboxes[b[0]], jboxes[b[1]])

		if distA > distB {
			return 1
		} else if distA < distB {
			return -1
		}

		return 0
	})

	fmt.Println("solution to part one: ", PartOne(jboxes, pairs))
	fmt.Println("solution to part two: ", PartTwo(jboxes, pairs))
}

func PartOne(jboxes [][3]int, pairs [][2]int) int {
	// we represent the sets using trees
	// we use a slice to represent the tree
	// such that the parent of node "i" is given by parents[i]
	parents := make([]int, len(jboxes))
	for i := range parents {
		parents[i] = i
	}

	n := 1000
	for i := range n {
		pair := pairs[i]
		merge(parents, pair[0], pair[1])
	}

	sizes := make([]int, len(jboxes))
	for i := range len(jboxes) {
		r := root(parents, i)
		sizes[r]++
	}

	slices.SortFunc(sizes, func(a, b int) int { return b - a })

	// return sum of the top-3 circuits
	return sizes[0] * sizes[1] * sizes[2]
}

func PartTwo(jboxes [][3]int, pairs [][2]int) int {
	parents := make([]int, len(jboxes))
	for i := range parents {
		parents[i] = i
	}

	circuits := len(jboxes)
	for i := range pairs {
		pair := pairs[i]

		if merge(parents, pair[0], pair[1]) {
			circuits--
		}

		if circuits == 1 {
			return jboxes[pair[0]][X] * jboxes[pair[1]][X]
		}
	}

	fmt.Println(circuits)

	return -1
}

func dist(a, b [3]int) float64 {
	dx := a[X] - b[X]
	dy := a[Y] - b[Y]
	dz := a[Z] - b[Z]

	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

func root(parents []int, node int) int {
	parent := parents[node]

	if parent == node {
		return node
	}

	// optimization: make node a direct descendant of the root
	parents[node] = root(parents, parent)
	return parents[node]
}

// returns true if the jboxes didn't belong in the same circuit
// i.e there were merged into a single circuit
func merge(parents []int, a, b int) bool {
	ra := root(parents, a)
	rb := root(parents, b)
	parents[ra] = rb
	return ra != rb
}
