package day16

import (
	"fmt"

	"github.com/wthys/advent-of-code-2023/collections/set"
	g "github.com/wthys/advent-of-code-2023/grid"
	l "github.com/wthys/advent-of-code-2023/location"
	"github.com/wthys/advent-of-code-2023/solver"
	"github.com/wthys/advent-of-code-2023/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "16"
}

func energyLevel(cave *g.Grid[Mirror], bounds g.Bounds, beam Beam) int {
	beams := Beams{beam}
	energised := set.New[l.Location]()
	visited := set.New[Beam]()

	for len(beams) > 0 {
		newBeams := Beams{}
		for _, beam := range beams {
			beam = beam.Move()
			if !bounds.Has(beam.pos) {
				continue
			}
			if visited.Has(beam) {
				continue
			}
			visited.Add(beam)
			energised.Add(beam.pos)

			mirror, _ := cave.Get(beam.pos)
			for _, bounce := range mirror.Bounce(beam.dir) {
				newBeams = append(newBeams, Beam{beam.pos, bounce})
			}
		}

		beams = newBeams
	}

	return energised.Len()
}

func (s solution) Part1(input []string) (string, error) {
	cave, bounds, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	beam := Beam{l.New(-1, 0), card2loc[EAST]}

	return solver.Solved(energyLevel(cave, bounds, beam))
}

func (s solution) Part2(input []string) (string, error) {
	cave, bounds, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	maxEnergy := 0
	// SOUTH + NORTH
	for x := 0; x < bounds.Width(); x++ {
		down := energyLevel(cave, bounds, Beam{l.New(x, -1), card2loc[SOUTH]})
		up := energyLevel(cave, bounds, Beam{l.New(x, bounds.Height()), card2loc[NORTH]})
		maxEnergy = max(maxEnergy, down, up)
	}
	// EAST + WEST
	for y := 0; y < bounds.Height(); y++ {
		right := energyLevel(cave, bounds, Beam{l.New(-1, y), card2loc[EAST]})
		left := energyLevel(cave, bounds, Beam{l.New(bounds.Width(), y), card2loc[WEST]})
		maxEnergy = max(maxEnergy, left, right)
	}

	return solver.Solved(maxEnergy)
}

const (
	NONE  = Cardinal(0)
	NORTH = Cardinal(1)
	EAST  = Cardinal(2)
	SOUTH = Cardinal(4)
	WEST  = Cardinal(8)
)

var (
	card2loc = map[Cardinal]l.Location{
		NORTH: {0, -1},
		SOUTH: {0, 1},
		EAST:  {1, 0},
		WEST:  {-1, 0},
	}
	loc2card = map[l.Location]Cardinal{
		{0, 1}:  SOUTH,
		{0, -1}: NORTH,
		{1, 0}:  EAST,
		{-1, 0}: WEST,
	}
)

type (
	Cardinal int
	Mirror   interface {
		Bounce(incoming l.Location) []l.Location
	}

	Empty   struct{}
	Slanted struct {
		top Cardinal
	}
	Straight struct {
		dir Cardinal
	}

	Beams []Beam
	Beam  struct {
		pos l.Location
		dir l.Location
	}
)

func (b Beam) Move() Beam {
	return Beam{b.pos.Add(b.dir), b.dir}
}

func (c Cardinal) String() string {
	return util.IIf(c&EAST > 0, "E", "") + util.IIf(c&NORTH > 0, "N", "") + util.IIf(c&SOUTH > 0, "S", "") + util.IIf(c&WEST > 0, "W", "")
}

func (b Beam) String() string {
	dir, _ := loc2card[b.dir]
	return fmt.Sprintf("[%v %v]", b.pos, dir)
}

func Rotate90(loc l.Location, n int) l.Location {
	switch n % 4 {
	case 0:
		return loc
	case 1, -3:
		return l.New(-loc.Y, loc.X)
	case 2, -2:
		return l.New(-loc.X, -loc.Y)
	case 3, -1:
		return l.New(loc.Y, -loc.X)
	default:
		panic(fmt.Sprintf("there is an unknown value for mod 4, got %v", n%4))
	}
}

func (m Empty) Bounce(incoming l.Location) []l.Location {
	return []l.Location{incoming}
}

func wrapLoc(loc l.Location) []l.Location {
	return []l.Location{loc}
}

func (m Slanted) Bounce(incoming l.Location) []l.Location {
	if m.top&(EAST+WEST) == 0 {
		panic(fmt.Sprintf("got an unexpected value for Slanted.top: %v", m.top))
	}

	idir, _ := loc2card[incoming]

	if m.top&EAST == 0 {
		switch idir {
		case NORTH:
			return wrapLoc(card2loc[WEST])
		case EAST:
			return wrapLoc(card2loc[SOUTH])
		case SOUTH:
			return wrapLoc(card2loc[EAST])
		case WEST:
			return wrapLoc(card2loc[NORTH])
		default:
			panic(fmt.Sprintf("%v is not a cardinal direction", incoming))
		}
	}

	switch idir {
	case NORTH:
		return wrapLoc(card2loc[EAST])
	case EAST:
		return wrapLoc(card2loc[NORTH])
	case SOUTH:
		return wrapLoc(card2loc[WEST])
	case WEST:
		return wrapLoc(card2loc[SOUTH])
	default:
		panic(fmt.Sprintf("%v is not a cardinal direction", incoming))
	}
}

func (m Straight) Bounce(incoming l.Location) []l.Location {
	dir, ok := loc2card[incoming.Unit()]
	if !ok {
		panic(fmt.Sprintf("invalid incoming direction: %v", incoming))
	}
	if dir&m.dir > 0 {
		return []l.Location{incoming}
	}
	return []l.Location{Rotate90(incoming, 1), Rotate90(incoming, -1)}
}

func MirrorFromRune(char rune) Mirror {
	switch char {
	case '-':
		return Straight{EAST + WEST}
	case '|':
		return Straight{NORTH + SOUTH}
	case '\\':
		return Slanted{WEST}
	case '/':
		return Slanted{EAST}
	default:
		return Empty{}
	}
}

func Accomodate(b g.Bounds, loc l.Location) g.Bounds {
	newb := g.Bounds{}
	newb.Xmin = min(loc.X, b.Xmin)
	newb.Xmax = max(loc.X, b.Xmax)
	newb.Ymin = min(loc.Y, b.Ymin)
	newb.Ymax = max(loc.Y, b.Ymax)
	return newb
}

func ParseInput(input []string) (*g.Grid[Mirror], g.Bounds, error) {
	cave := g.WithDefault[Mirror](Empty{})

	bounds := g.Bounds{}
	bounds.Xmin = 1_000_000_000_000
	bounds.Ymin = 1_000_000_000_000
	bounds.Xmax = 0
	bounds.Ymax = 0

	for y, row := range input {
		for x, char := range row {
			loc := l.New(x, y)
			bounds = Accomodate(bounds, loc)
			mirror := MirrorFromRune(char)
			if _, ok := mirror.(Empty); ok {
				continue
			}
			cave.Set(loc, mirror)
		}
	}

	return cave, bounds, nil
}
