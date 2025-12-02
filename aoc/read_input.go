package aoc

import (
	"bufio"
	"os"
)

func ReadInputLineByLine(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, err
		}

		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func ReadInput(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
