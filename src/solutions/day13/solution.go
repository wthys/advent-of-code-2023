package day13

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
	return "13"
}

func (s solution) Part1(input []string) (string, error) {
	patterns, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, pattern := range patterns {
		bounds := g.BoundsFromSlice(pattern)
		corner := l.New(bounds.Xmax-1, bounds.Ymax-1)
		for corner.X >= bounds.Xmin || corner.Y >= bounds.Ymin {
			if pattern.HFold(corner.Y, 0) {
				total += 100 * corner.Y
				break
			}

			if pattern.VFold(corner.X, 0) {
				total += corner.X
				break
			}

			corner = corner.Add(l.New(-1, -1))
		}
	}

	return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
	patterns, err := ParseInput(input)
	if err != nil {
		return solver.Error(err)
	}

	total := 0
	for _, pattern := range patterns {
		bounds := g.BoundsFromSlice(pattern)
		corner := l.New(bounds.Xmax-1, bounds.Ymax-1)
		for corner.X >= bounds.Xmin || corner.Y >= bounds.Ymin {
			if pattern.HFold(corner.Y, 1) {
				total += 100 * corner.Y
				break
			}

			if pattern.VFold(corner.X, 1) {
				total += corner.X
				break
			}

			corner = corner.Add(l.New(-1, -1))
		}
	}

	return solver.Solved(total)
}

type (
	Pattern  []l.Location
	Patterns []Pattern
)

func PrintPattern(pattern Pattern) {
	grid := g.WithDefault(0)
	for _, loc := range pattern {
		grid.Set(loc, 1)
	}

	grid.PrintFunc(func(val int, _ error) string {
		if val > 0 {
			return "#"
		}
		return "."
	})
}

func PrintPatternHFold(pattern Pattern, y int) {
	grid := g.WithDefault(0)
	for _, loc := range pattern {
		if loc.Y > y {
			grid.Set(loc.Add(l.New(0, 1)), 1)
		} else {
			grid.Set(loc, 1)
		}
	}

	b, _ := grid.Bounds()
	for x := b.Xmin; x <= b.Xmax; x++ {
		grid.Set(l.New(x, y+1), 2)
	}

	grid.PrintFunc(func(val int, _ error) string {
		switch val {
		case 1:
			return "#"
		case 2:
			return "-"
		default:
			return "."
		}
	})
}

func PrintPatternVFold(pattern Pattern, x int) {
	grid := g.WithDefault(0)
	for _, loc := range pattern {
		if loc.X > x {
			grid.Set(loc.Add(l.New(1, 0)), 1)
		} else {
			grid.Set(loc, 1)
		}
	}

	b, _ := grid.Bounds()
	for y := b.Ymin; y <= b.Ymax; y++ {
		grid.Set(l.New(x+1, y), 2)
	}

	grid.PrintFunc(func(val int, _ error) string {
		switch val {
		case 1:
			return "#"
		case 2:
			return "|"
		default:
			return "."
		}
	})
}

func MirrorV(loc l.Location, x int) l.Location {
	return l.New(loc.X+2*(x-loc.X)+1, loc.Y)
}

func MirrorH(loc l.Location, y int) l.Location {
	return l.New(loc.X, loc.Y+2*(y-loc.Y)+1)
}

func (pattern Pattern) VFold(x int, smudges int) bool {
	bounds := g.BoundsFromSlice(pattern)

	if x < bounds.Xmin || x >= bounds.Xmax {
		return false
	}

	left := set.New[l.Location]()
	right := set.New[l.Location]()
	mLeft := set.New[l.Location]()
	mRight := set.New[l.Location]()

	size := min(x-bounds.Xmin+1, bounds.Xmax-x)
	xlo := x - size + 1
	xhi := x + size

	for _, loc := range pattern {
		if loc.X <= x {
			if loc.X < xlo {
				continue
			}
			left.Add(loc)
			mLeft.Add(MirrorV(loc, x))
		} else {
			if loc.X > xhi {
				continue
			}
			right.Add(loc)
			mRight.Add(MirrorV(loc, x))
		}
	}

	if util.Abs(left.Len()-right.Len()) != smudges {
		return false
	}

	leftIsect := mRight.Intersect(left)
	rightIsect := mLeft.Intersect(right)
	return util.Abs(leftIsect.Len()-left.Len()) <= smudges && util.Abs(rightIsect.Len()-right.Len()) <= smudges
}

func (pattern Pattern) HFold(y int, smudges int) bool {
	bounds := g.BoundsFromSlice(pattern)

	if y < bounds.Ymin || y >= bounds.Ymax {
		return false
	}

	left := set.New[l.Location]()
	right := set.New[l.Location]()
	mLeft := set.New[l.Location]()
	mRight := set.New[l.Location]()

	size := min(y-bounds.Ymin+1, bounds.Ymax-y)
	ylo := y - size + 1
	yhi := y + size

	for _, loc := range pattern {
		if loc.Y <= y {
			if loc.Y < ylo {
				continue
			}
			left.Add(loc)
			mLeft.Add(MirrorH(loc, y))
		} else {
			if loc.Y > yhi {
				continue
			}
			right.Add(loc)
			mRight.Add(MirrorH(loc, y))
		}
	}

	if util.Abs(left.Len()-right.Len()) != smudges {
		return false
	}

	leftIsect := mRight.Intersect(left)
	rightIsect := mLeft.Intersect(right)
	return util.Abs(leftIsect.Len()-left.Len()) <= smudges && util.Abs(rightIsect.Len()-right.Len()) <= smudges
}

func ParseInput(input []string) (Patterns, error) {
	patterns := Patterns{}

	pattern := Pattern{}
	patternStart := 0
	for y, line := range input {
		if len(line) == 0 {
			patterns = append(patterns, pattern)
			patternStart = -1
			pattern = Pattern{}
		} else {
			if patternStart < 0 {
				patternStart = y
			}

			for x, char := range line {
				if char == '#' {
					loc := l.New(x+1, y-patternStart+1)
					pattern = append(pattern, loc)
				}
			}
		}
	}

	if len(pattern) > 0 {
		patterns = append(patterns, pattern)
	}

	if len(patterns) == 0 {
		return patterns, fmt.Errorf("no patterns found")
	}
	return patterns, nil
}
