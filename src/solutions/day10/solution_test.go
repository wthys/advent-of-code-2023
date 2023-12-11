package day10

import (
	"testing"

	"github.com/wthys/advent-of-code-2023/location"
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

func Test100_ParseInput_SimpleLoop(t *testing.T) {
	input := []string{
		"F7",
		"LJ",
	}

	pipelines, startLocation := ParseInput(input)

	assertPipelines(t, 1, pipelines)
	assertLocation(t, "startLocation", UNKNOWN_STARTLOCATION, startLocation)

	if !pipelines[0].IsLoop() {
		t.Fatalf("expected %v to be a loop, it is not", pipelines[0])
	}
}

func Test100_ParseInput_SinglePipe(t *testing.T) {
	input := []string{"-"}

	pipelines, startLocation := ParseInput(input)

	assertPipelines(t, 1, pipelines)
	assertLocation(t, "startLocation", UNKNOWN_STARTLOCATION, startLocation)
}

func Test100_ParseInput_DisjointPipeLines(t *testing.T) {
	input := []string{"-JL-"}

	pipelines, startLocation := ParseInput(input)

	assertPipelines(t, 2, pipelines)
	assertLocation(t, "startLocation", UNKNOWN_STARTLOCATION, startLocation)
}

func Test100_ParseInput__StartLocation(t *testing.T) {
	input := []string{"S"}

	pipelines, startLocation := ParseInput(input)

	assertPipelines(t, 0, pipelines)
	assertLocation(t, "startLocation", location.New(0, 0), startLocation)
}
