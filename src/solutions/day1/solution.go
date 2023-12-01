package day1

import (
	"strings"

	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "1"
}

func extractNumber(input string, mapping map[string]int) int {
	index1 := len(input)
	dig1 := 0
	index2 := -1
	dig2 := 0
	for key, value := range mapping {
		frIdx := strings.Index(input, key)
		if frIdx >= 0 {
			index1 = min(frIdx, index1)
			dig1 = value
		}
		laIdx := strings.LastIndex(input, key)
		if laIdx >= 0 {
			index2 = max(laIdx+len(key), index2)
			dig2 = value
		}
	}

	return 10*dig1 + dig2
}

func (s solution) Part1(input []string) (string, error) {
	total := 0
	mapping := map[string]int{
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	}

	for _, line := range input {
		total += extractNumber(line, mapping)
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	total := 0
	mapping := map[string]int{
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9,
	}

	for _, line := range input {
		total += extractNumber(line, mapping)
	}

	return solver.Solved(total)
}
