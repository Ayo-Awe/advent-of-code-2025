package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

var (
	shapesPattern  = regexp.MustCompile(`\d:\n(?:[#\.]+\n?)+`)
	regionsPattern = regexp.MustCompile(`\d+x\d+:(?: \d+)+`)
)

type region struct {
	dim        [2]int
	quantities []int
}

func ParseInput(input string) ([][]string, []region, error) {
	rawShapes := shapesPattern.FindAllString(input, -1)
	rawRegions := regionsPattern.FindAllString(input, -1)

	shapes := make([][]string, len(rawShapes))
	for i, rawShape := range rawShapes {
		rawShape = strings.Trim(rawShape, "\n")
		shapes[i] = strings.Split(rawShape, "\n")[1:]
	}

	regions := make([]region, len(rawRegions))
	for i, rawRegion := range rawRegions {
		parts := strings.Split(rawRegion, ":")
		rawDim := strings.Split(parts[0], "x")
		rawQty := strings.Fields(parts[1])

		w, err := strconv.Atoi(rawDim[0])
		if err != nil {
			return nil, nil, err
		}

		h, err := strconv.Atoi(rawDim[1])
		if err != nil {
			return nil, nil, err
		}

		quantities := make([]int, len(rawQty))
		for i := range rawQty {
			qty, err := strconv.Atoi(rawQty[i])
			if err != nil {
				return nil, nil, err
			}

			quantities[i] = qty
		}

		regions[i] = region{dim: [2]int{w, h}, quantities: quantities}
	}

	return shapes, regions, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	input, err := aoc.ReadInput(*filename)
	if err != nil {
		log.Fatal(err)
	}

	shapes, regions, err := ParseInput(input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("solution to part one: ", PartOne(shapes, regions))
}

func PartOne(shapes [][]string, regions []region) int {
	var count int
	for _, region := range regions {
		area := region.dim[0] * region.dim[1]

		totalQty := 0
		for _, qty := range region.quantities {
			totalQty += qty
		}

		reqArea := totalQty * 9
		if area >= reqArea {
			count++
		}
	}
	return count
}

