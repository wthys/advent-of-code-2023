package day16

import (
	"testing"

	"github.com/wthys/advent-of-code-2023/collections/set"
	l "github.com/wthys/advent-of-code-2023/location"
)

type (
	Result[T any, R any] struct {
		input    T
		expected R
	}
)

var (
	east  = l.New(1, 0)
	west  = l.New(-1, 0)
	north = l.New(0, -1)
	south = l.New(0, 1)
)

func TestBounceSlantedWest(t *testing.T) {
	testcases := []Result[l.Location, []l.Location]{
		{north, []l.Location{west}},
		{west, []l.Location{north}},
		{east, []l.Location{south}},
		{south, []l.Location{east}},
	}

	mirror := MirrorFromRune('\\')

	for _, tc := range testcases {
		actual := mirror.Bounce(tc.input)
		if len(actual) != len(tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, actual)
		}
		sact := set.New(actual...)
		sexp := set.New(tc.expected...)
		if sact.Intersect(sexp).Len() < sexp.Len() {
			t.Fatalf("expected %v, got %v", tc.expected, actual)
		}
	}

}

func TestBounceSlantedEast(t *testing.T) {
	testcases := []Result[l.Location, []l.Location]{
		{north, []l.Location{east}},
		{east, []l.Location{north}},
		{west, []l.Location{south}},
		{south, []l.Location{west}},
	}

	mirror := MirrorFromRune('/')

	for _, tc := range testcases {
		actual := mirror.Bounce(tc.input)
		if len(actual) != len(tc.expected) {
			t.Fatalf("expected %v, got %v", tc.expected, actual)
		}
		sact := set.New(actual...)
		sexp := set.New(tc.expected...)
		if sact.Intersect(sexp).Len() < sexp.Len() {
			t.Fatalf("expected %v, got %v", tc.expected, actual)
		}
	}

}
