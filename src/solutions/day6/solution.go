package day6

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wthys/advent-of-code-2023/solver"
	"github.com/wthys/advent-of-code-2023/util/interval"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "6"
}

func (s solution) Part1(input []string) (string, error) {
	races, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 1
	for _, race := range races {
		ivl := race.Options()
		total = total * ivl.Len()
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	race, err := parseInput2(input)
	if err != nil {
		return solver.Error(err)
	}

	ivl := race.Options()
	return solver.Solved(ivl.Len())
}

type (
	Race struct {
		time   int
		record int
	}
)

func (r Race) String() string {
	return fmt.Sprintf("Race %vms, %vmm", r.time, r.record)
}

func (r Race) race(t int) int {
	return max(0, r.time-t) * t
}

func (race Race) Options() interval.Interval {
	m := (race.time + 1) / 2

	lolo := 0
	lohi := m
	for lohi-lolo > 1 {
		lomid := (lolo + lohi) / 2
		d := race.race(lomid)
		if d < race.record {
			lolo = lomid
		} else {
			lohi = lomid
		}
	}

	hilo := m
	hihi := race.time
	for hihi-hilo > 1 {
		himid := (hilo + hihi) / 2
		d := race.race(himid)
		if d < race.record {
			hihi = himid
		} else {
			hilo = himid
		}
	}
	// expanded := true
	// for expanded {
	// 	expanded = false
	// 	if race.race(lo-1) > race.record {
	// 		lo = lo - 1
	// 		expanded = true
	// 	}
	// 	if race.race(hi+1) > race.record {
	// 		hi = hi + 1
	// 		expanded = true
	// 	}
	// }
	return interval.New(lohi, hilo)
}

func parseInput(input []string) ([]Race, error) {
	races := []Race{}
	reNum := regexp.MustCompile("[0-9]+")

	times := reNum.FindAllString(input[0], -1)
	distances := reNum.FindAllString(input[1], -1)
	if len(times) != len(distances) {
		return nil, fmt.Errorf("different amounts of times and distances")
	}

	for idx, nTime := range times {
		time, _ := strconv.Atoi(nTime)
		record, _ := strconv.Atoi(distances[idx])
		races = append(races, Race{time, record})
	}

	return races, nil
}

func parseInput2(input []string) (Race, error) {
	reNum := regexp.MustCompile("[^0-9]*")

	times := reNum.ReplaceAllString(input[0], "")
	distances := reNum.ReplaceAllString(input[1], "")

	time, err := strconv.Atoi(times)
	if err != nil {
		return Race{}, err
	}
	record, err2 := strconv.Atoi(distances)
	if err2 != nil {
		return Race{}, err2
	}
	return Race{time, record}, nil
}
