package day15

import "testing"

type (
	TestCase[I any, E any] struct {
		input    I
		expected E
	}
)

func TestInstructionHash(t *testing.T) {
	testcases := []TestCase[string, int]{
		{"HASH", 52},
		{"rn=1", 30},
		{"cm-", 253},
		{"qp=3", 97},
		{"cm=2", 47},
		{"qp-", 14},
		{"pc=4", 180},
		{"ot=9", 9},
		{"ab=5", 197},
		{"pc-", 48},
		{"pc=6", 214},
		{"ot=7", 231},
	}

	for _, tc := range testcases {
		actual := Instruction(tc.input).Hash()
		if actual != tc.expected {
			t.Fatalf("%v should have hash %v, got %v", tc.input, tc.expected, actual)
		}
	}
}
