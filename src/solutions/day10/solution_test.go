package day10

import (
	"testing"

	"github.com/wthys/advent-of-code-2023/collections/set"
	"github.com/wthys/advent-of-code-2023/location"
	"github.com/wthys/advent-of-code-2023/util"
)

var UNKNOWN_STARTLOCATION = location.New(-1, -1)

func assertPipelines(t *testing.T, expectedLen int, pipelines []PipeLine) {
	if actual := len(pipelines); actual != expectedLen {
		switch expectedLen {
		case 1:
			t.Fatalf("expected 1 pipeline, got %v => %v", actual, pipelines)
		default:
			t.Fatalf("expected %v pipelines, got %v => %v", expectedLen, actual, pipelines)
		}
	}
}

func assertLocation(t *testing.T, label string, expected location.Location, actual location.Location) {
	if actual != expected {
		t.Fatalf("expected %v to be %v, got %v", label, expected, actual)
	}
}

func assertPipeConnection(t *testing.T, a Pipe, b Pipe, expected Connection) {
	actual := a.FindConnection(b)
	if actual != expected {
		t.Fatalf("expected %v to connect to %v through %v, got %v", a, b, expected, actual)
	}
}

func Test010_Pipe_FindConnection(t *testing.T) {
	center := location.New(0, 0)
	north := location.New(0, -1)
	east := location.New(1, 0)
	south := location.New(0, 1)
	west := location.New(-1, 0)

	cross := NORTH + EAST + WEST + SOUTH

	assertPipeConnection(t, NewPipe(center, cross), NewPipe(east, cross), EAST)
	assertPipeConnection(t, NewPipe(center, cross), NewPipe(west, cross), WEST)
	assertPipeConnection(t, NewPipe(center, cross), NewPipe(south, cross), SOUTH)
	assertPipeConnection(t, NewPipe(center, cross), NewPipe(north, cross), NORTH)
	assertPipeConnection(t, NewPipe(east, cross), NewPipe(west, cross), NONE)
	assertPipeConnection(t, NewPipe(east, cross), NewPipe(south, cross), NONE)

	assertPipeConnection(t, NewPipe(center, cross), NewPipe(east, NORTH+SOUTH), NONE)
	assertPipeConnection(t, NewPipe(center, WEST+NORTH), NewPipe(east, EAST+NORTH), NONE)
}

func Test050_PipeLine_Merge(t *testing.T) {
	pipeA := NewPipe(location.New(0, 0), EAST)
	pipeB := NewPipe(location.New(1, 0), WEST)

	plA := NewPipeLine(pipeA)
	plB := NewPipeLine(pipeB)

	plC, okC := plA.Merge(plB)
	if !okC {
		t.Fatalf("expected %v to merge with %v, did not work", plA, plB)
	}

	plD, okD := plB.Merge(plA)
	if !okD {
		t.Fatalf("expected %v to merge with %v, did not work", plB, plA)
	}

	if plC.Equals(plD) {
		t.Fatalf("expected %v to be the same as %v", plC, plD)
	}
}

func Test100_Pipe_NeejberLocs(t *testing.T) {
	center := location.New(0, 0)
	north := location.New(0, -1)
	east := location.New(1, 0)
	south := location.New(0, 1)
	west := location.New(-1, 0)

	dirs := map[Connection]location.Location{
		NORTH: north,
		EAST:  east,
		WEST:  west,
		SOUTH: south,
	}

	for conn, expected := range dirs {
		pipe := NewPipe(center, conn)
		neejbers := pipe.NeejberLocs()

		if len(neejbers) != 1 {
			t.Fatalf("exptected 1 connection, got %v => %v", len(neejbers), neejbers)
		}
		assertLocation(t, "neejber", expected, neejbers[0])
	}

	util.PermutationDo(2, (NORTH + EAST + SOUTH + WEST).Connections(), func(conns []Connection) {
		conn := NONE
		expected := set.New[location.Location]()
		for _, c := range conns {
			conn += c
			expected.Add(dirs[c])
		}

		pipe := NewPipe(center, conn)
		neejbers := pipe.NeejberLocs()
		actual := set.New(neejbers...)

		if actual.Subtract(expected).Len() > 0 || expected.Subtract(actual).Len() > 0 {
			t.Fatalf("expected %v to have %v neejbers, got %v", pipe, expected, actual)
		}
	})

}
