package day1

import (
	"fmt"
	"regexp"
	"strconv"
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

var strMapping map[string]int = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
	"1":     1,
	"2":     2,
	"3":     3,
	"4":     4,
	"5":     5,
	"6":     6,
	"7":     7,
	"8":     8,
	"9":     9,
}

func convertNumber(input string) (int, error) {
	result, ok := strMapping[input]
	if !ok {
		return 0, fmt.Errorf("unknown number [%s]", input)
	}

	return result, nil
}

func extractNumber(input string) int {
	index1 := len(input)
	dig1 := 0
	index2 := -1
	dig2 := 0
	for key, value := range strMapping {
		frIdx := strings.Index(input, key)
		if frIdx >= 0 && frIdx < index1 {
			index1 = frIdx
			dig1 = value
		}
		laIdx := strings.LastIndex(input, key)
		if laIdx >= 0 && laIdx+len(key) > index2 {
			index2 = laIdx + len(key)
			dig2 = value
		}
	}

	return 10*dig1 + dig2
}

func (s solution) Part1(input []string) (string, error) {
	re := regexp.MustCompile("[0-9]")

	total := 0

	for _, line := range input {
		matches := re.FindAllString(line, -1)

		if len(matches) > 0 {
			cand, _ := strconv.Atoi(matches[0] + matches[len(matches)-1])
			total = total + cand
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	total := 0

	for _, line := range input {
		cand := extractNumber(line)
		total += cand
	}

	return solver.Solved(total)
}
