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

type machine struct {
	lightMask int // bit representation of lights e.g [#..] = 100
	numLights int // number of lights
	buttons   [][]int
	joltage   []int
}

type buttonCombo struct {
	joltage []int
	presses int
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
	// we solve each machine as simulataneous equations
	for _, m := range machines {
		minSum += solve(m.joltage, patterns(m))
	}
	return minSum
}

func patterns(m machine) map[string][]buttonCombo {
	p := make(map[string][]buttonCombo)
	for i := range 1 << len(m.buttons) {
		joltage := make([]int, len(m.joltage))

		var presses int
		for buttonIdx := range len(m.buttons) {
			pressed := (1<<buttonIdx)&i != 0
			if pressed {
				presses++
				for _, light := range m.buttons[buttonIdx] {
					joltage[light]++
				}
			}
		}

		pattern := joltageParityKey(joltage)
		p[pattern] = append(p[pattern], buttonCombo{
			joltage: joltage,
			presses: presses,
		})
	}
	return p
}

func solve(joltage []int, patterns map[string][]buttonCombo) int {
	// base case: joltage is zeroed out e.g 0,0,0,0
	if slices.Max(joltage) == 0 {
		return 0
	}

	key := joltageParityKey(joltage)

	minPresses := 1_000_000
	for _, combo := range patterns[key] {
		newJoltage := slices.Clone(joltage)

		// subtract and halve
		for j := range combo.joltage {
			newJoltage[j] = (newJoltage[j] - combo.joltage[j]) / 2
		}

		// skip invalid joltage
		if slices.Min(newJoltage) < 0 {
			continue
		}

		presses := 2*solve(newJoltage, patterns) + combo.presses
		minPresses = min(minPresses, presses)
	}

	return minPresses
}

func joltageParityKey(joltage []int) string {
	var pattern string
	for _, v := range joltage {
		if v%2 == 0 {
			pattern += "e"
		} else {
			pattern += "o"
		}
	}
	return pattern
}
