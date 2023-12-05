package day5

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "5"
}

func (s solution) Part1(input []string) (string, error) {
	seeds, mappers, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	firstLocation := 10000000000000
	for _, seed := range seeds {
		reqs, err := mappers.Gather(seed)
		if err != nil {
			fmt.Println(err)
			continue
		}

		locId, ok := reqs["location"]
		if !ok {
			fmt.Printf("seed %v has no location!! => %v\n", seed, reqs)
			continue
		}

		firstLocation = min(firstLocation, locId)
	}

	return solver.Solved(firstLocation)
}

func (s solution) Part2(input []string) (string, error) {
	seeds, mappers, err := parseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	firstLocation := 100000000000000000

	prev := seeds[0]
	for idx, val := range seeds[1:] {
		if idx%2 == 0 {
			fmt.Printf(" checking %v -> %v (%v)\n", prev, prev+val-1, val)
			for i := 0; i < val; i++ {
				seed := prev + i
				reqs, err := mappers.Gather(seed)
				if err != nil {
					fmt.Println(err)
					continue
				}

				locId, ok := reqs["location"]
				if !ok {
					fmt.Printf("seed %v does not have a location => %v\n", seed, reqs)
					continue
				}

				firstLocation = min(firstLocation, locId)
				if i%100 == 0 {
					fmt.Printf("  checked up until #%v/%v     \r", i, seed)
				}
			}
		}
		prev = val
	}
	fmt.Println()

	return solver.Solved(firstLocation)
}

type (
	MapRange struct {
		source int
		target int
		size   int
	}

	Mapper struct {
		from   string
		to     string
		ranges []MapRange
	}

	Mappers          []Mapper
	SeedRequirements map[string]int
)

func (m Mapper) Map(id int) int {
	for _, maprange := range m.ranges {
		newId, ok := maprange.Map(id)
		if ok {
			return newId
		}
	}
	return id
}

func (m MapRange) Map(id int) (int, bool) {
	diff := id - m.source
	if diff >= 0 && diff < m.size {
		return m.target + diff, true
	}
	return 0, false
}

func (m Mappers) Resolve(from string, id int) (int, string, bool) {
	for _, mapper := range m {
		if mapper.from == from {
			return mapper.Map(id), mapper.to, true
		}
	}
	return 0, "", false
}

func (m Mappers) Gather(seed int) (SeedRequirements, error) {
	reqs := SeedRequirements{"seed": seed}

	from := "seed"
	id := seed
	for from != "location" {
		newId, to, ok := m.Resolve(from, id)
		if !ok {
			return nil, fmt.Errorf("could not map %v from %q", id, from)
		}

		reqs[to] = newId
		id = newId
		from = to
	}

	return reqs, nil
}

func parseInput(input []string) ([]int, Mappers, error) {
	reNum := regexp.MustCompile("[0-9]+")
	reSeeds := regexp.MustCompile("^seeds: ")
	reMapName := regexp.MustCompile("^([a-z]+)-to-([a-z]+) map:")
	reMapRange := regexp.MustCompile("^([0-9]+) ([0-9]+) ([0-9]+)$")
	seeds := []int{}

	var currentMapper *Mapper = nil
	mapRanges := []MapRange{}
	mappers := Mappers{}

	for _, line := range input {
		if len(line) == 0 && currentMapper != nil {
			currentMapper.ranges = mapRanges
			mappers = append(mappers, *currentMapper)
			currentMapper = nil
			mapRanges = []MapRange{}

		} else if reSeeds.MatchString(line) {
			matches := reNum.FindAllString(line, -1)
			for _, num := range matches {
				val, _ := strconv.Atoi(num)
				seeds = append(seeds, val)
			}

		} else if reMapName.MatchString(line) {
			names := reMapName.FindStringSubmatch(line)
			mapper := Mapper{}
			mapper.from = names[1]
			mapper.to = names[2]
			mapper.ranges = []MapRange{}
			currentMapper = &mapper

		} else if reMapRange.MatchString(line) {
			nums := reMapRange.FindStringSubmatch(line)
			maprange := MapRange{}
			maprange.source, _ = strconv.Atoi(nums[2])
			maprange.target, _ = strconv.Atoi(nums[1])
			maprange.size, _ = strconv.Atoi(nums[3])
			mapRanges = append(mapRanges, maprange)
		}
	}

	if currentMapper != nil {
		currentMapper.ranges = mapRanges
		mappers = append(mappers, *currentMapper)
	}

	return seeds, mappers, nil
}
