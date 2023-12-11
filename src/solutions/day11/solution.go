package day11

import (
	"fmt"

	"github.com/wthys/advent-of-code-2023/collections/set"
	"github.com/wthys/advent-of-code-2023/grid"
	"github.com/wthys/advent-of-code-2023/location"
	"github.com/wthys/advent-of-code-2023/solver"
	"github.com/wthys/advent-of-code-2023/util/interval"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "11"
}

func (s solution) Part1(input []string) (string, error) {
	observation, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for idx, a := range observation.galaxies[:len(observation.galaxies)-1] {
		for _, b := range observation.galaxies[idx+1:] {
			dist := a.Subtract(b).Manhattan()
			xd := interval.New(a.X, b.X)
			yd := interval.New(a.Y, b.Y)

			extra := 0
			observation.unusedX.ForEach(func(x int) {
				if xd.Contains(x) {
					extra += 1
				}
			})
			observation.unusedY.ForEach(func(y int) {
				if yd.Contains(y) {
					extra += 1
				}
			})

			total += dist + extra
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	observation, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	expanse := 1_000_000

	total := 0
	for idx, a := range observation.galaxies[:len(observation.galaxies)-1] {
		for _, b := range observation.galaxies[idx+1:] {
			dist := a.Subtract(b).Manhattan()
			xd := interval.New(a.X, b.X)
			yd := interval.New(a.Y, b.Y)

			extra := 0
			observation.unusedX.ForEach(func(x int) {
				if xd.Contains(x) {
					extra += expanse - 1
				}
			})
			observation.unusedY.ForEach(func(y int) {
				if yd.Contains(y) {
					extra += expanse - 1
				}
			})

			total += dist + extra
		}
	}

	return solver.Solved(total)
}

type (
	Observation struct {
		galaxies []location.Location
		unusedX  *set.Set[int]
		unusedY  *set.Set[int]
	}
)

func (obs Observation) String() string {
	return fmt.Sprintf("Observation(g=%v, x=%v, y=%v)", obs.galaxies, obs.unusedX, obs.unusedY)
}

func ParseInput(input []string) (Observation, error) {
	observation := Observation{}
	if len(input) == 0 {
		return observation, fmt.Errorf("not enough input")
	}

	universe := grid.WithDefault(false)
	galaxies := []location.Location{}
	unusedX := set.New[int]()
	unusedY := set.New[int]()

	for y, line := range input {
		for x, space := range line {
			pos := location.New(x, y)
			if space == '#' {
				universe.Set(pos, true)
				galaxies = append(galaxies, pos)
			}
			unusedX.Add(x)
			unusedY.Add(y)
		}
	}

	for _, galaxy := range galaxies {
		unusedX.Remove(galaxy.X)
		unusedY.Remove(galaxy.Y)
	}

	observation.galaxies = galaxies
	observation.unusedX = unusedX
	observation.unusedY = unusedY

	return observation, nil
}
