package day3

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/wthys/advent-of-code-2023/collections/set"
	"github.com/wthys/advent-of-code-2023/grid"
	"github.com/wthys/advent-of-code-2023/location"
	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "3"
}

func (s solution) Part1(input []string) (string, error) {
	partNumbers, parts := parseInput(input)

	total := 0
	for _, partNumber := range partNumbers {
		_, found := findPart(parts, partNumber)
		if found {
			total += partNumber.number
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	partNumbers, parts := parseInput(input)

	total := 0
	parts.Apply(func(loc location.Location, part rune) {
		if part != rune('*') {
			return
		}

		numbers := findPartNumbers(partNumbers, loc)
		if len(numbers) != 2 {
			return
		}

		ratio := 1
		for _, pn := range numbers {
			ratio *= pn.number
		}
		total += ratio
	})

	return solver.Solved(total)
}

type (
	PartNumber struct {
		number    int
		locations []location.Location
	}
)

func (pn PartNumber) Contains(loc location.Location) bool {
	for _, pnLoc := range pn.locations {
		if pnLoc == loc {
			return true
		}
	}
	return false
}

func (pn PartNumber) String() string {
	return fmt.Sprintf("PartNumber(#%v, loc=%v)", pn.number, pn.locations[0])
}

func findPartNumbers(partNumbers []PartNumber, loc location.Location) []PartNumber {
	indices := set.New[int]()
	for _, nbr := range loc.Neejbers() {
		for idx, pn := range partNumbers {
			if pn.Contains(nbr) {
				indices.Add(idx)
			}
		}
	}

	numbers := []PartNumber{}
	indices.Do(func(idx int) bool {
		numbers = append(numbers, partNumbers[idx])
		return true
	})
	return numbers
}

func findPart(parts *grid.Grid[rune], partNumber PartNumber) (location.Location, bool) {
	checked := map[location.Location]bool{}
	for _, loc := range partNumber.locations {
		for _, neejber := range loc.Neejbers() {
			_, ok := checked[neejber]
			if ok {
				continue
			}

			_, err := parts.Get(neejber)
			if err == nil {
				return neejber, true
			}
			checked[neejber] = true
		}
	}
	return location.Location{}, false
}

func parseInput(input []string) ([]PartNumber, *grid.Grid[rune]) {
	partNumberMatrix := map[location.Location]*PartNumber{}
	parts := grid.New[rune]()

	for y, line := range input {
		for x, char := range line {
			if char == rune('.') {
				continue
			}

			loc := location.New(x, y)
			value, err := strconv.Atoi(string(char))
			if err != nil {
				parts.Set(loc, char)
			} else {
				partNumber, ok := partNumberMatrix[location.New(x-1, y)]
				if !ok {
					partNumber = &PartNumber{0, []location.Location{}}
				}
				partNumber.number = partNumber.number*10 + value
				partNumber.locations = append(partNumber.locations, loc)
				partNumberMatrix[loc] = partNumber
			}
		}
	}

	partNumbers := []*PartNumber{}

	for _, partNumber := range partNumberMatrix {
		if slices.Contains(partNumbers, partNumber) {
			continue
		}
		partNumbers = append(partNumbers, partNumber)
	}

	byValue := []PartNumber{}
	for _, pn := range partNumbers {
		byValue = append(byValue, *pn)
	}

	return byValue, parts
}
