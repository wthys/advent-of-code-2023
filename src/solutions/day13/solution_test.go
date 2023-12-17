package day13

import (
	"strings"
	"testing"

	l "github.com/wthys/advent-of-code-2023/location"
)

type (
	TestCase[I any, E any] struct {
		input    I
		expected E
	}
	FoldInput struct {
		label string
		str   string
		fold  int
	}
	MirrorInput struct {
		loc    l.Location
		offset int
	}
)

func toPattern(t *testing.T, input string) Pattern {
	patterns, err := ParseInput(strings.Split(input, "\n"))
	if err != nil {
		t.Fatalf("did not expect err %v", err)
	}
	return patterns[0]
}

const (
	pCircle = ".##.\n#..#\n#..#\n.##."
	pVert   = "#..#\n.##.\n....\n####"
	pHoriz  = "#..#\n#.#.\n#.#.\n#..#"
	pAssymH = "#.\n#.\n.#\n.#\n#.\n#.\n..\n#."
	pAssymV = "#....##..##\n#.....####."
	pFull   = "##\n##"
)

func TestPatternHFold(t *testing.T) {
	testcases := []TestCase[FoldInput, bool]{
		{FoldInput{"circle", pCircle, 0}, false},
		{FoldInput{"circle", pCircle, 1}, false},
		{FoldInput{"circle", pCircle, 2}, true},
		{FoldInput{"circle", pCircle, 3}, false},
		{FoldInput{"circle", pCircle, 4}, false},
		{FoldInput{"vert", pVert, 0}, false},
		{FoldInput{"vert", pVert, 1}, false},
		{FoldInput{"vert", pVert, 2}, false},
		{FoldInput{"vert", pVert, 3}, false},
		{FoldInput{"vert", pVert, 4}, false},
		{FoldInput{"horiz", pHoriz, 0}, false},
		{FoldInput{"horiz", pHoriz, 1}, false},
		{FoldInput{"horiz", pHoriz, 2}, true},
		{FoldInput{"horiz", pHoriz, 3}, false},
		{FoldInput{"horiz", pHoriz, 4}, false},
		{FoldInput{"assymH", pAssymH, 0}, false},
		{FoldInput{"assymH", pAssymH, 1}, true},
		{FoldInput{"assymH", pAssymH, 2}, false},
		{FoldInput{"assymH", pAssymH, 3}, true},
		{FoldInput{"assymH", pAssymH, 4}, false},
		{FoldInput{"assymH", pAssymH, 5}, false},
		{FoldInput{"assymH", pAssymH, 6}, false},
		{FoldInput{"assymH", pAssymH, 7}, false},
		{FoldInput{"assymV", pAssymV, 1}, false},
		{FoldInput{"assymV", pAssymV, 2}, false},
		{FoldInput{"full", pFull, 0}, false},
		{FoldInput{"full", pFull, 1}, true},
		{FoldInput{"full", pFull, 2}, false},
	}

	for _, tc := range testcases {
		pattern := toPattern(t, tc.input.str)
		actual := pattern.HFold(tc.input.fold, 0)
		if actual != tc.expected {
			PrintPatternHFold(pattern, tc.input.fold)
			t.Fatalf("exected %q HFold@%v to be %v, got %v", tc.input.label, tc.input.fold, tc.expected, actual)
		}
	}
}

func TestMirrorV(t *testing.T) {
	testcases := []TestCase[MirrorInput, l.Location]{
		{MirrorInput{l.New(0, 0), 0}, l.New(1, 0)},
		{MirrorInput{l.New(0, 1), 1}, l.New(3, 1)},
		{MirrorInput{l.New(0, 2), 2}, l.New(5, 2)},
		{MirrorInput{l.New(0, 3), 3}, l.New(7, 3)},
		{MirrorInput{l.New(0, 4), 4}, l.New(9, 4)},
		{MirrorInput{l.New(1, 0), 0}, l.New(0, 0)},
		{MirrorInput{l.New(2, 1), 0}, l.New(-1, 1)},
		{MirrorInput{l.New(3, 2), 0}, l.New(-2, 2)},
		{MirrorInput{l.New(4, 3), 0}, l.New(-3, 3)},
		{MirrorInput{l.New(5, 4), 0}, l.New(-4, 4)},
	}

	for _, tc := range testcases {
		actual := MirrorV(tc.input.loc, tc.input.offset)
		if actual != tc.expected {
			t.Fatalf("Expected MirrorH(%v, %v) to be %v, got %v", tc.input.loc, tc.input.offset, tc.expected, actual)
		}
	}
}

func TestMirrorH(t *testing.T) {
	testcases := []TestCase[MirrorInput, l.Location]{
		{MirrorInput{l.New(0, 0), 0}, l.New(0, 1)},
		{MirrorInput{l.New(1, 0), 1}, l.New(1, 3)},
		{MirrorInput{l.New(2, 0), 2}, l.New(2, 5)},
		{MirrorInput{l.New(3, 0), 3}, l.New(3, 7)},
		{MirrorInput{l.New(4, 0), 4}, l.New(4, 9)},
		{MirrorInput{l.New(0, 1), 0}, l.New(0, 0)},
		{MirrorInput{l.New(1, 2), 0}, l.New(1, -1)},
		{MirrorInput{l.New(2, 3), 0}, l.New(2, -2)},
		{MirrorInput{l.New(3, 4), 0}, l.New(3, -3)},
		{MirrorInput{l.New(4, 5), 0}, l.New(4, -4)},
	}

	for _, tc := range testcases {
		actual := MirrorH(tc.input.loc, tc.input.offset)
		if actual != tc.expected {
			t.Fatalf("Expected MirrorH(%v, %v) to be %v, got %v", tc.input.loc, tc.input.offset, tc.expected, actual)
		}
	}
}
