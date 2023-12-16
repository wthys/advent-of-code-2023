package day14

import (
	"fmt"

	g "github.com/wthys/advent-of-code-2023/grid"
	l "github.com/wthys/advent-of-code-2023/location"
	"github.com/wthys/advent-of-code-2023/solver"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "14"
}

func RockStringer(rock Rock, _ error) string {
	_, ok := rock.(CubeRock)
	if ok {
		return "#"
	}
	_, ok = rock.(RoundRock)
	if ok {
		return "O"
	}
	return "."
}

func (s solution) Part1(input []string) (string, error) {
	platform, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	platform = Tilt(platform, NORTH)

	bounds, err := platform.Bounds()
	total := 0
	platform.ForEach(func(_ l.Location, rock Rock) {
		_, rounded := rock.(RoundRock)
		if rounded {
			total += bounds.Height() - rock.Pos().Y
		}
	})

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	platform, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	bounds, _ := platform.Bounds()
	loadCalc := func(platform *g.Grid[Rock]) int {
		total := 0
		platform.ForEach(func(_ l.Location, rock Rock) {
			_, rounded := rock.(RoundRock)
			if rounded {
				total += bounds.Height() - rock.Pos().Y
			}
		})
		return total
	}

	limit := 1000000000

	lruCache := map[string]int{}
	lruCache[Hash(platform)] = -1

	loadCache := map[int]int{}

	for i := 0; i < limit; i++ {
		platform = Tilt(platform, NORTH)
		platform = Tilt(platform, WEST)
		platform = Tilt(platform, SOUTH)
		platform = Tilt(platform, EAST)

		hash := Hash(platform)
		lru, ok := lruCache[hash]
		if ok {
			period := i - lru
			likelyIdx := (limit-i)%period + lru - 1
			return solver.Solved(loadCache[likelyIdx])
		}
		lruCache[hash] = i
		loadCache[i] = loadCalc(platform)
	}

	return solver.Error(fmt.Errorf("Could not find a period (I'm sorry)"))
}

const (
	NORTH = 0
	SOUTH = 2
	EAST  = 1
	WEST  = 3
)

var (
	card2dir = map[Cardinal]l.Location{
		NORTH: l.New(0, -1),
		EAST:  l.New(1, 0),
		SOUTH: l.New(0, 1),
		WEST:  l.New(-1, 0),
	}
)

type (
	Cardinal int
	Rock     interface {
		Move(Cardinal) Rock
		Covers(Rock) bool
		Pos() l.Location
	}

	CubeRock struct {
		pos l.Location
	}

	RoundRock struct {
		pos l.Location
	}
	NoRock struct {
		pos l.Location
	}
)

func (rock CubeRock) Move(_ Cardinal) Rock {
	return rock
}

func (rock CubeRock) Covers(other Rock) bool {
	_, empty := other.(NoRock)
	return empty
}

func (rock CubeRock) Pos() l.Location {
	return rock.pos
}

func (rock RoundRock) Move(direction Cardinal) Rock {
	dir, _ := card2dir[direction]
	return RoundRock{rock.pos.Add(dir)}
}

func (rock RoundRock) Covers(other Rock) bool {
	_, empty := other.(NoRock)
	return empty
}

func (rock RoundRock) Pos() l.Location {
	return rock.pos
}

func (rock NoRock) Move(_ Cardinal) Rock {
	return rock
}

func (rock NoRock) Covers(_ Rock) bool {
	return false
}

func (rock NoRock) Pos() l.Location {
	return rock.pos
}

func Hash(platform *g.Grid[Rock]) string {
	bounds, _ := platform.Bounds()
	out := ""
	for y := bounds.Ymin; y <= bounds.Ymax; y++ {
		for x := bounds.Xmin; x <= bounds.Xmax; x++ {
			loc := l.New(x, y)
			rock, err := platform.Get(loc)
			if err != nil {
				continue
			}

			_, round := rock.(RoundRock)
			if round {
				out += rock.Pos().String()
			}
		}
	}
	return out
}

func Tilt(grid *g.Grid[Rock], direction Cardinal) *g.Grid[Rock] {
	platform := g.WithDefaultFunc[Rock](func(loc l.Location) (Rock, error) {
		return NoRock{loc}, nil
	})
	bounds, _ := grid.Bounds()

	mover := func(rock Rock) {
		if rock.Move(direction).Pos() == rock.Pos() {
			platform.Set(rock.Pos(), rock)
			return
		}

		prev := rock
		for {
			candidate := prev.Move(direction)
			if !bounds.Has(candidate.Pos()) {
				platform.Set(prev.Pos(), prev)
				break
			}

			existing, _ := platform.Get(candidate.Pos())

			if !candidate.Covers(existing) {
				platform.Set(prev.Pos(), prev)
				break
			}
			prev = candidate
		}
	}

	switch direction {
	case NORTH:
		for y := bounds.Ymin; y <= bounds.Ymax; y++ {
			for x := bounds.Xmin; x <= bounds.Xmax; x++ {
				loc := l.New(x, y)
				rock, err := grid.Get(loc)
				if err == nil {
					mover(rock)
				}
			}
		}
	case EAST:
		for y := bounds.Ymin; y <= bounds.Ymax; y++ {
			for x := bounds.Xmax; x >= bounds.Xmin; x-- {
				loc := l.New(x, y)
				rock, err := grid.Get(loc)
				if err == nil {
					mover(rock)
				}
			}
		}
	case SOUTH:
		for y := bounds.Ymax; y >= bounds.Ymin; y-- {
			for x := bounds.Xmin; x <= bounds.Xmax; x++ {
				loc := l.New(x, y)
				rock, err := grid.Get(loc)
				if err == nil {
					mover(rock)
				}
			}
		}
	case WEST:
		for y := bounds.Ymin; y <= bounds.Ymax; y++ {
			for x := bounds.Xmin; x <= bounds.Xmax; x++ {
				loc := l.New(x, y)
				rock, err := grid.Get(loc)
				if err == nil {
					mover(rock)
				}
			}
		}
	default:
		panic(fmt.Sprintf("unknown direction '%v'", direction))
	}

	newPlatform := g.New[Rock]()
	bounds.ForEach(func(loc l.Location) {
		rock, err := platform.Get(loc)
		if err != nil {
			rock = NoRock{loc}
		}
		newPlatform.Set(loc, rock)
	})

	return newPlatform
}

func ParseInput(input []string) (*g.Grid[Rock], error) {
	platform := g.New[Rock]()

	for y, line := range input {
		if len(line) == 0 {
			continue
		}

		for x, rockType := range line {
			pos := l.New(x, y)
			switch rockType {
			case '#':
				platform.Set(pos, CubeRock{pos})
			case 'O':
				platform.Set(pos, RoundRock{pos})
			case '.':
				platform.Set(pos, NoRock{pos})
			default:
				return nil, fmt.Errorf("unknown rockType '%v' #%v,%v", rockType, y+1, x+1)
			}
		}
	}

	if platform.Len() == 0 {
		return nil, fmt.Errorf("no rocks found")
	}

	return platform, nil
}
