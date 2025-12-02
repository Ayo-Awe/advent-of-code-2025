#!/bin/bash

# Original script credits: 0xdf and sventec (GitHub)

# Usage:
# Set $AOC_SESSION to AOC session cookie: export AOC_SESSION=<cookie> or via a .env file
# Then: source setup.sh <day#>
#  e.g. source setup.sh 3
# Optionally, include the year for completing previous year's challenges:
#       source setup.sh <day#> <year>
#  e.g. source setup.sh 3 2022

day=$1
main_file_template=$(echo -n "package main

import (
    \"flag\"
    \"fmt\"
)

func main() {
	filename := flag.String(\"file\", \"input.txt\", \"input file name\")
	flag.Parse()

    fmt.Println(\"solution to part one: \", PartOne())
    fmt.Println(\"solution to part two: \", PartTwo())
}

func PartOne() int {
    return 0
}

func PartTwo() int {
    return 0
}")

if [ -z "$day" ]; then
    echo "Usage: $0 <day> [year]"
    exit 1
fi

if [[ ! "$day" =~ ^[0-9]+$ ]]; then
    echo "Invalid argument: $day is not a number"
    exit 1
fi

folder=$(printf "day_%02d" "$day")
main_file="$folder/main.go"
input_file="$folder/input.txt"

if [ ! -e "$folder" ]; then
    mkdir "$folder"
fi


if [ ! -e "$main_file" ]; then
    echo -n "$main_file_template" > "$main_file"
fi

# load .env if it exists
if [ -e ".env" ]; then
    set -a
    source .env
    set +a
fi

# fetch input for current day
if [ -z "$AOC_SESSION" ]; then
  echo "\$AOC_SESSION isn't set. Can't fetch your input. Quitting"
  return
fi

year=${2:-$(date +%Y)}
url_base="https://adventofcode.com/${year}/day/${day}"

curl -s "${url_base}/input" --cookie "session=${AOC_SESSION}" -o $input_file
