package day9

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wthys/advent-of-code-2023/collections/set"
	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "9"
}

func (s solution) Part1(input []string) (string, error) {
	histories, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := Measurement(0)
	for _, history := range histories {
		next := history.Next()
		// fmt.Printf("%v => %v\n", history, next)
		total += next
	}
	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	histories, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := Measurement(0)
	for _, history := range histories {
		prev := history.Prev()
		// fmt.Printf("%v => %v\n", history, next)
		total += prev
	}
	return solver.Solved(total)
}

type (
	Measurement  int
	Measurements []Measurement
	Histories    []Measurements
)

func (ms Measurements) Diff() Measurements {
	if len(ms) < 2 {
		return Measurements{}
	}

	diffs := Measurements{}
	for idx, val := range ms[:len(ms)-1] {
		diffs = append(diffs, ms[idx+1]-val)
	}
	return diffs
}

func (ms Measurements) IsConstant() bool {
	if len(ms) <= 1 {
		return true
	}
	return set.New(ms...).Len() == 1
}

func (ms Measurements) Next() Measurement {
	diffs := ms.Diff()

	all0 := true
	for _, v := range diffs {
		if v != 0 {
			all0 = false
			break
		}
	}
	last := ms[len(ms)-1]
	if all0 {
		return last
	}
	return last + diffs.Next()
}

func (ms Measurements) Prev() Measurement {
	diffs := ms.Diff()

	all0 := true
	for _, v := range diffs {
		if v != 0 {
			all0 = false
			break
		}
	}
	first := ms[0]
	if all0 {
		return first
	}
	return first - diffs.Prev()
}

func ParseInput(input []string) (Histories, error) {
	reNum := regexp.MustCompile(`-?[0-9]+`)

	histories := Histories{}

	for _, line := range input {
		matches := reNum.FindAllString(line, -1)
		if len(matches) == 0 {
			continue
		}

		measurements := Measurements{}
		for _, m := range matches {
			num, _ := strconv.Atoi(m)
			measurements = append(measurements, Measurement(num))
		}

		histories = append(histories, measurements)
	}

	if len(histories) == 0 {
		return Histories{}, fmt.Errorf("no histories found")
	}

	return histories, nil
}
