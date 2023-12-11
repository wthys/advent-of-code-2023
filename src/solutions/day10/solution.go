package day10

import (
	"fmt"
	"slices"

	"github.com/wthys/advent-of-code-2023/collections/set"
	"github.com/wthys/advent-of-code-2023/grid"
	"github.com/wthys/advent-of-code-2023/location"
	"github.com/wthys/advent-of-code-2023/solver"
	"github.com/wthys/advent-of-code-2023/util"
)

type solution struct{}

func init() {
	solver.Register(solution{})
}

func (s solution) Day() string {
	return "10"
}

func (s solution) Part1(input []string) (string, error) {
	area, startLocation := ParseInput(input)

	mainLoop, err := findMainLoop(area, startLocation)
	if err != nil {
		return solver.Error(err)
	}

	return solver.Solved(mainLoop.Len() / 2)
}

func (s solution) Part2(input []string) (string, error) {
	pipearea, startLocation := ParseInput(input)

	mainLoop, err := findMainLoop(pipearea, startLocation)
	if err != nil {
		return solver.Error(err)
	}

	area := grid.WithDefault(false)
	east := location.New(1, 0)
	south := location.New(0, 1)
	mainLoop.ForEach(func(pipe Pipe) {
		double := pipe.pos.Scale(2)
		area.Set(double, true)
		if (pipe.connections & EAST).IsConnected() {
			area.Set(double.Add(east), true)
		}
		if (pipe.connections & SOUTH).IsConnected() {
			area.Set(double.Add(south), true)
		}
	})

	b, err := area.Bounds()
	if err != nil {
		return solver.Error(err)
	}
	bounds := grid.Bounds{0, 0, 0, 0}
	bounds.Xmin = b.Xmin - 1
	bounds.Xmax = b.Xmax + 1
	bounds.Ymin = b.Ymin - 1
	bounds.Ymax = b.Ymax + 1

	outside := set.New[location.Location]()
	fringe := set.New[location.Location](location.New(bounds.Xmin, bounds.Ymin))
	for fringe.Len() > 0 {
		newFringe := set.New[location.Location]()
		fringe.ForEach(func(loc location.Location) {
			outside.Add(loc)
			for _, neejber := range loc.OrthoNeejbers() {
				if !bounds.Has(neejber) || outside.Has(neejber) || newFringe.Has(neejber) {
					continue
				}
				pathTile, _ := area.Get(neejber)
				if !pathTile {
					newFringe.Add(neejber)
				}
			}
		})

		fringe = newFringe
	}

	inside := set.New[location.Location]()
	fringe = set.New[location.Location]()
	for _, neejber := range startLocation.Scale(2).Neejbers() {
		onPath, _ := area.Get(neejber)
		if onPath || outside.Has(neejber) || !bounds.Has(neejber) {
			continue
		}
		fringe.Add(neejber)
	}

	for fringe.Len() > 0 {
		newFringe := set.New[location.Location]()
		fringe.ForEach(func(loc location.Location) {
			inside.Add(loc)
			for _, neejber := range loc.OrthoNeejbers() {
				if !bounds.Has(neejber) || inside.Has(neejber) || outside.Has(neejber) || newFringe.Has(neejber) {
					continue
				}
				pathTile, _ := area.Get(neejber)
				if !pathTile {
					newFringe.Add(neejber)
				}
			}
		})

		fringe = newFringe
	}

	realInside := set.New[location.Location]()
	inside.ForEach(func(loc location.Location) {
		realInside.Add(location.New(loc.X/2, loc.Y/2))
	})

	mainLoop.ForEach(func(pipe Pipe) {
		realInside.Remove(pipe.pos)
	})

	return solver.Solved(realInside.Len())
}

const (
	NONE  = Connection(0)
	NORTH = Connection(1)
	EAST  = Connection(2)
	SOUTH = Connection(4)
	WEST  = Connection(8)
)

type (
	Connection int
	Pipe       struct {
		pos         location.Location
		connections Connection
	}
	Pipes    []Pipe
	PipeLine struct {
		pipes Pipes
		head  Pipe
		tail  Pipe
	}
)

func findMainLoop(area *grid.Grid[Pipe], startLocation location.Location) (PipeLine, error) {
	start := NewPipe(startLocation, NORTH+SOUTH+EAST+WEST)
	checkedLocs := set.New[location.Location]()
	for _, neejber := range startLocation.OrthoNeejbers() {
		if checkedLocs.Has(neejber) {
			continue
		}

		pipeline := NewPipeLine(start)
		pipelineSet := set.New[location.Location](startLocation)
		current := neejber

		for true {
			if checkedLocs.Has(current) {
				break
			}
			pipe, _ := area.Get(current)
			pl, ok := pipeline.Connect(pipe)
			if !ok {
				break
			}
			pipeline = pl
			pipelineSet.Add(current)
			checkedLocs.Add(current)

			nbrs := pipe.NeejberLocs()
			found := false
			for _, nbr := range nbrs {
				if pipelineSet.Has(nbr) {
					continue
				}
				current = nbr
				found = true
				break
			}
			if !found {
				break
			}
		}

		if pipeline.IsLoop() {
			return pipeline, nil
		}
	}

	return PipeLine{}, fmt.Errorf("could not find a looping pipeline")
}

func (con Connection) Connections() []Connection {
	conns := []Connection{}
	if con&NORTH > 0 {
		conns = append(conns, NORTH)
	}
	if con&EAST > 0 {
		conns = append(conns, EAST)
	}
	if con&SOUTH > 0 {
		conns = append(conns, SOUTH)
	}
	if con&WEST > 0 {
		conns = append(conns, WEST)
	}
	return conns
}

var connection2symbol = map[rune]Connection{
	'.': NONE,
	'|': NORTH + SOUTH,
	'-': EAST + WEST,
	'L': NORTH + EAST,
	'J': NORTH + WEST,
	'7': SOUTH + WEST,
	'F': SOUTH + EAST,
}

func ConnectionFromSymbol(symbol rune) Connection {
	cs, ok := connection2symbol[symbol]
	if ok {
		return cs
	}
	return NONE
}

func (con Connection) Invert() Connection {
	conns := NONE
	if con&NORTH > 0 {
		conns += SOUTH
	}
	if con&EAST > 0 {
		conns += WEST
	}
	if con&SOUTH > 0 {
		conns += NORTH
	}
	if con&WEST > 0 {
		conns += EAST
	}
	return conns
}

func (pl PipeLine) Connect(pipe Pipe) (PipeLine, bool) {
	return pl.Merge(NewPipeLine(pipe))
}

func (this PipeLine) Reverse() PipeLine {
	rev := slices.Clone(this.pipes)
	slices.Reverse(rev)
	return PipeLine{rev, this.tail, this.head}
}

func (this PipeLine) Merge(that PipeLine) (PipeLine, bool) {
	pl, merged := this.merge(that)
	if !merged {
		pl, merged = this.merge(that.Reverse())
	}
	return pl, merged
}

func (this PipeLine) merge(that PipeLine) (PipeLine, bool) {
	pl := PipeLine{}
	merged := true
	if conn := this.head.FindConnection(that.tail); conn.IsConnected() {
		// H----that----T=H----this----T
		pl = PipeLine{append(that.pipes, this.pipes...), that.head, this.tail}

	} else if conn := this.tail.FindConnection(that.head); conn.IsConnected() {
		// H----this----T=H----that----T
		pl = PipeLine{append(this.pipes, that.pipes...), this.head, that.tail}

	} else {
		merged = false
	}
	return pl, merged
}

func (pl PipeLine) IsLoop() bool {
	if pl.Len() < 4 {
		return false
	}

	return pl.head.FindConnection(pl.tail).IsConnected()
}

func (pl PipeLine) Len() int {
	return len(pl.pipes)
}

func (pl PipeLine) ForEach(forEach func(Pipe)) {
	for _, pipe := range pl.pipes {
		forEach(pipe)
	}
}

func (pl PipeLine) HasFunc(hasFunc func(pipe Pipe) bool) bool {
	for _, pipe := range pl.pipes {
		if hasFunc(pipe) {
			return true
		}
	}
	return false
}

func (pipe Pipe) NeejberLocs() []location.Location {
	neejbers := []location.Location{}
	for _, conn := range pipe.connections.Connections() {
		loc, ok := conn2loc[conn]
		if !ok {
			continue
		}
		neejbers = append(neejbers, pipe.pos.Add(loc))
	}
	return neejbers
}

func (con Connection) IsConnected() bool {
	return (con & (NORTH + SOUTH + EAST + WEST)) > 0
}

var loc2conn = map[location.Location]Connection{
	location.New(-1, 0): WEST,
	location.New(1, 0):  EAST,
	location.New(0, -1): NORTH,
	location.New(0, 1):  SOUTH,
}
var conn2loc = map[Connection]location.Location{
	WEST:  location.New(-1, 0),
	EAST:  location.New(1, 0),
	NORTH: location.New(0, -1),
	SOUTH: location.New(0, 1),
}

func (this Pipe) FindConnection(that Pipe) Connection {
	diff := that.pos.Subtract(this.pos)
	if diff.Manhattan() != 1 {
		return NONE
	}

	conn, ok := loc2conn[diff]
	if !ok {
		return NONE
	}

	thisScoped := this.connections & conn
	thatScoped := that.connections.Invert() & conn

	if thisScoped == NONE || thatScoped == NONE {
		return NONE
	}

	return conn
}

func NewPipe(pos location.Location, conn Connection) Pipe {
	return Pipe{pos, conn}
}

func NewPipeLine(pipe Pipe) PipeLine {
	return PipeLine{Pipes{pipe}, pipe, pipe}
}

func (this Pipe) Equals(that Pipe) bool {
	if this.pos != that.pos {
		return false
	}
	return this.connections == that.connections
}

func (this PipeLine) Equals(that PipeLine) bool {
	if this.head != that.head {
		return false
	}
	if this.tail != that.tail {
		return false
	}
	if this.Len() != that.Len() {
		return false
	}
	for idx, pipe := range this.pipes {
		if !pipe.Equals(that.pipes[idx]) {
			return false
		}
	}

	return true
}

func (pl PipeLine) String() string {
	return fmt.Sprintf("PipeLine(%v->%v %v#%v)", pl.head, pl.tail, len(pl.pipes), pl.pipes)
}

func (conn Connection) String() string {
	return fmt.Sprintf("<%v%v%v%v>",
		util.IIf(conn&EAST != NONE, "←", "."),
		util.IIf(conn&NORTH != NONE, "↑", "."),
		util.IIf(conn&SOUTH != NONE, "↓", "."),
		util.IIf(conn&WEST != NONE, "→", "."),
	)
}

func (pipe Pipe) String() string {
	return fmt.Sprintf("Pipe(%v %v)", pipe.pos, pipe.connections)
}

func ParseInput(input []string) (*grid.Grid[Pipe], location.Location) {
	startLocation := location.New(-1, -1)
	area := grid.WithDefaultFunc[Pipe](func(loc location.Location) (Pipe, error) {
		return NewPipe(loc, NONE), nil
	})

	for y, line := range input {
		if len(line) == 0 {
			continue
		}

		for x, sym := range line {
			pos := location.New(x, y)
			if sym == 'S' {
				startLocation = pos
				continue
			}
			conn := ConnectionFromSymbol(sym)
			if !conn.IsConnected() {
				continue
			}
			area.Set(pos, NewPipe(pos, conn))
		}
	}

	return area, startLocation
}
