package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ayo-awe/advent-of-code-2025/aoc"
)

func ParseInput(lines []string) ([][]int, []string, error) {
	operators := strings.Fields(lines[len(lines)-1])

	operands := make([][]int, len(lines)-1)
	for i, line := range lines[:len(lines)-1] {
		sNums := strings.Fields(line)

		nums := make([]int, len(sNums))
		for i, sNum := range sNums {
			num, err := strconv.Atoi(sNum)
			if err != nil {
				return nil, nil, err
			}

			nums[i] = num
		}

		operands[i] = nums
	}

	return operands, operators, nil
}

func main() {
	filename := flag.String("file", "input.txt", "input file name")
	flag.Parse()

	lines, err := aoc.ReadInputLineByLine(*filename)
	if err != nil {
		log.Fatal(err)
	}

	operands, operators, err := ParseInput(lines)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("solution to part one: ", PartOne(operands, operators))
	fmt.Println("solution to part two: ", PartTwo(lines))
}

func PartOne(operands [][]int, operators []string) int {
	var total int

	for i, operator := range operators {
		var res int
		if operator == "*" {
			res = 1
		}

		for _, ops := range operands {
			if operator == "*" {
				res *= ops[i]
			} else {
				res += ops[i]
			}
		}

		total += res
	}

	return total
}

func PartTwo(lines []string) int {
	operators := strings.Fields(lines[len(lines)-1])

	var opIdx int
	var nums []int
	var total int

	for col := range len(lines[0]) {
		var digits []int

		for row := range len(lines) - 1 {
			if digit := int(lines[row][col]) - '0'; digit >= 0 && digit <= 9 {
				digits = append(digits, digit)
			}
		}

		if len(digits) > 0 {
			nums = append(nums, digitsToNum(digits))
		} else {
			op := operators[opIdx]
			total += exec(op, nums)
			nums = nil
			opIdx++
		}
	}

	op := operators[opIdx]
	total += exec(op, nums)

	return total
}

func digitsToNum(digits []int) int {
	var num int
	for i := range len(digits) {
		num += pow10(len(digits)-1-i) * digits[i]
	}
	return num
}

func pow10(n int) int {
	pow := 1
	for range n {
		pow *= 10
	}
	return pow
}

func exec(op string, nums []int) int {
	switch op {
	case "+":
		var total int
		for i := range nums {
			total += nums[i]
		}
		return total
	case "*":
		total := 1
		for i := range nums {
			total *= nums[i]
		}
		return total
	default:
		panic("unknown operator")
	}
}
