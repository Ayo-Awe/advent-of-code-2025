package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

func ParseInput(lines []string) map[string][]string {
	devices := make(map[string][]string, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ":")
		device := parts[0]
		outputs := strings.Fields(parts[1])
		devices[device] = outputs
	}
	return devices
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	devices := ParseInput(lines)

	fmt.Println("solution to part one: ", PartOne(devices))
	fmt.Println("solution to part two: ", PartTwo(devices))
}

func PartOne(devices map[string][]string) int {
	memo := make(map[string]int)
	return paths("you", "out", devices, memo)
}

func PartTwo(devices map[string][]string) int {
	// depending on your input fft might come before dac or vice versa

	soln := 1
	// svr -> fft -> dac -> out
	if count := paths("fft", "dac", devices, make(map[string]int)); count > 0 {
		soln *= count
		soln *= paths("svr", "fft", devices, make(map[string]int))
		soln *= paths("dac", "out", devices, make(map[string]int))
	} else {
		// svr -> dac -> fft -> out
		soln *= paths("svr", "dac", devices, make(map[string]int))
		soln *= paths("dac", "fft", devices, make(map[string]int))
		soln *= paths("fft", "out", devices, make(map[string]int))
	}

	return soln
}

func paths(src, dest string, devices map[string][]string, memo map[string]int) int {
	if count, exists := memo[src]; exists {
		return count
	}

	if src == dest {
		return 1
	}

	for _, output := range devices[src] {
		memo[src] += paths(output, dest, devices, memo)
	}

	return memo[src]
}
