package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

const (
	Dir  = 0
	Dist = 1

	L = -1
	R = 1
)

func ParseInput(lines []string) ([][2]int, error) {
	rotations := make([][2]int, len(lines))
	for i, line := range lines {
		var dir int
		if line[0] == 'L' {
			dir = L
		} else {
			dir = R
		}

		dist, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to int at line %d: %w", line[1:], i+1, err)
		}

		rotations[i] = [2]int{dir, dist}
	}
	return rotations, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()
	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	rotations, err := ParseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("solution to part one: ", PartOne(rotations))
	fmt.Println("solution to part two: ", PartTwo(rotations))
}

func PartOne(rotations [][2]int) int {
	dial := 50

	var count int
	for _, rot := range rotations {
		dir, dist := rot[Dir], rot[Dist]
		dial = (dial + dir*dist) % 100

		if dial == 0 {
			count++
		}
	}

	return count
}

func PartTwo(rotations [][2]int) int {
	dial := 50

	var count int
	for _, rot := range rotations {
		dir, dist := rot[Dir], rot[Dist]

		// calculate the distance from the current dial position to zero
		var dist2Zero int
		if dial == 0 {
			dist2Zero = 100
		} else if dir == L {
			dist2Zero = dial
		} else {
			dist2Zero = (100 - dial)
		}

		// rationale: move the dial to zero and count how many times we can go 360 round the dial
		// increment count with 1 (reaching the dial initially) + number of 360s (reaching the dial subsequently)
		if dist >= dist2Zero {
			count += 1 + (dist-dist2Zero)/100
		}

		dial = mod100(dial + dir*dist)
	}

	return count
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// returns all mod values as +ve integers
func mod100(val int) int {
	return (val + (abs(val)/100+1)*100) % 100
}
