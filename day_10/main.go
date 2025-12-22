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

type machine struct {
	lightMask int // bit representation of lights e.g [#..] = 100
	numLights int // number of lights
	buttons   [][]int
	joltage   []int
}

func (m machine) buttonmMasks() []int {
	masks := make([]int, len(m.buttons))
	for i := range m.buttons {
		masks[i] = buttonMask(m.buttons[i], m.numLights)
	}
	return masks
}

func ParseInput(lines []string) ([]machine, error) {
	machines := make([]machine, len(lines))
	for i, line := range lines {
		parts := strings.Fields(line)

		rawLights := parts[0]
		rawJoltage := parts[len(parts)-1]
		rawButtons := parts[1 : len(parts)-1]

		lightMask := parseLights(rawLights)

		joltage, err := parseJoltage(rawJoltage)
		if err != nil {
			return nil, err
		}

		buttons, err := parseButtons(rawButtons)
		if err != nil {
			return nil, err
		}

		machines[i] = machine{
			joltage:   joltage,
			lightMask: lightMask,
			buttons:   buttons,
			numLights: len(rawLights) - 2, // -2 to ignore the square brackets
		}
	}
	return machines, nil
}

func parseLights(raw string) int {
	// trim off square brackets
	trimmed := strings.Trim(strings.Trim(raw, "["), "]")

	//.##. = 0110
	var lights int
	for i, v := range trimmed {
		// state is either 0 or 1 i.e on/off
		var state int
		if v == '#' {
			state = 1
		}
		lights += state << (len(trimmed) - 1 - i)
	}

	return lights
}

func parseButtons(rawButtons []string) ([][]int, error) {
	buttons := make([][]int, len(rawButtons))
	for i, rawButton := range rawButtons {
		trimmed := strings.Trim(strings.Trim(rawButton, "("), ")")

		var button []int
		for _, lightStr := range strings.Split(trimmed, ",") {
			light, err := strconv.Atoi(lightStr)
			if err != nil {
				return nil, err
			}
			button = append(button, light)
		}

		buttons[i] = button
	}
	return buttons, nil
}

func parseJoltage(raw string) ([]int, error) {
	trimmed := strings.Trim(strings.Trim(raw, "{"), "}")
	parts := strings.Split(trimmed, ",")

	joltage := make([]int, len(parts))
	for i, joltStr := range parts {
		jolt, err := strconv.Atoi(joltStr)
		if err != nil {
			return nil, err
		}
		joltage[i] = jolt
	}

	return joltage, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	machines, err := ParseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("solution to part one: ", PartOne(machines))
	fmt.Println("solution to part two: ", PartTwo(machines))
}

func PartOne(machines []machine) int {
	var minSum int
	for _, m := range machines {
		presses := minPresses(m)
		minSum += presses
	}
	return minSum
}

func minPresses(m machine) int {
	queue := [][3]int{}

	// add initial nodes to queue
	buttonMasks := m.buttonmMasks()
	for _, b := range buttonMasks {
		queue = append(queue, [3]int{0, b, 0})
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		lights, button, presses := curr[0], curr[1], curr[2]

		// we similate presses by applying the button as a bitmask
		lights = lights ^ button
		presses++

		// we've found the minimum presses
		if lights == m.lightMask {
			return presses
		}

		// we want to simulate new presses
		for _, b := range buttonMasks {
			queue = append(queue, [3]int{lights, b, presses})
		}
	}

	// should be impossible to hit
	return 0
}

func buttonMask(button []int, numLights int) int {
	var mask int
	for _, light := range button {
		mask = mask | 1<<(numLights-light-1)
	}
	return mask
}

func PartTwo(machines []machine) int {
	var minSum int
	for _, m := range machines {
		presses := minPressesJoltage(m)
		minSum += presses
	}
	return minSum
}

func minPressesJoltage(m machine) int {
	var presses int
	joltage := slices.Clone(m.joltage)

	// sort buttons by size in desc order
	slices.SortFunc(m.buttons, func(a, b []int) int { return len(b) - len(a) })

	excludedButtons := map[int]bool{}
	for slices.Max(joltage) > 0 {

		target := math.MaxInt
		idx := 0
		for i, v := range joltage {
			if v != 0 && v < target {
				target = v
				idx = i
			}
		}

		var found bool
		for i, button := range m.buttons {
			if excludedButtons[i] {
				continue
			}

			if slices.Contains(button, idx) {
				if !found {
					// apply button target times
					for _, level := range button {
						joltage[level] -= target
					}
					found = true
				}

				// exclude all buttons that contain the target index
				excludedButtons[i] = true
			}
		}

		presses += target
	}

	if slices.Min(joltage) != 0 {
		fmt.Println("oops", joltage)
	}

	return presses
}

