package day9

import "testing"

type (
	TestCase[T any, R any] struct {
		input    T
		expected R
	}

	TestCase2[T any, R any, S any] struct {
		input     T
		expected1 R
		expected2 S
	}
)

func assertMeasurements(t *testing.T, label string, expected Measurements, actual Measurements) {
	if len(expected) != len(actual) {
		t.Fatalf("%v expected length %v, got %v", label, len(expected), len(actual))
	}

	for i, this := range expected {
		that := actual[i]
		if this != that {
			t.Fatalf("%v expected actual[%v] to be %v, got %v", label, i, this, that)
		}
	}
}

func TestDiff(t *testing.T) {
	cases := []TestCase[Measurements, Measurements]{
		{Measurements{0, 1, 2, 3, 4}, Measurements{1, 1, 1, 1}},
		{Measurements{0, 1, 1, 2, 3, 5}, Measurements{1, 0, 1, 1, 2}},
		{Measurements{0, 1, 3, 6, 10}, Measurements{1, 2, 3, 4}},
		{Measurements{1, -1, 1, -1, 1}, Measurements{-2, 2, -2, 2}},
		{Measurements{0, 1, 0, 1, 0}, Measurements{1, -1, 1, -1}},
	}

	for _, tc := range cases {
		actual := tc.input.Diff()
		assertMeasurements(t, "TestDiff", tc.expected, actual)
	}
}

func TestNext(t *testing.T) {
	cases := []TestCase[Measurements, Measurement]{
		{Measurements{1, 1}, 1},
		{Measurements{1, 2}, 3},
		{Measurements{1, 2, 2}, 1},
		{Measurements{1, 0, 1}, 4},
		{Measurements{1, 0, 1, 0}, -7},
	}

	for _, tc := range cases {
		actual := tc.input.Next()
		if tc.expected != actual {
			t.Fatalf("expected %v.Next() to be %v, got %v", tc.input, tc.expected, actual)
		}
	}
}
