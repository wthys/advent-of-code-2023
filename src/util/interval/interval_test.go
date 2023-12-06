package interval

import "testing"

type (
	BinaryResult[T any, R any] struct {
		left     T
		right    T
		expected R
	}
	UnaryResult[T any, R any] struct {
		left     T
		expected R
	}
)

func TestIntervalEquals(t *testing.T) {
	cases := []BinaryResult[Interval, bool]{
		BinaryResult[Interval, bool]{New(0, 5), New(0, 5), true},
		BinaryResult[Interval, bool]{New(1, 5), New(2, 4), false},
		BinaryResult[Interval, bool]{New(3, 3), New(6, 6), false},
		BinaryResult[Interval, bool]{Interval{3, 2}, Interval{6, 3}, true},
	}

	for _, tc := range cases {
		actualLR := tc.left.Equals(tc.right)
		if actualLR != tc.expected {
			t.Fatalf("%v.Equals(%v) should return %v, got %v", tc.left, tc.right, tc.expected, actualLR)
		}
		actualRL := tc.right.Equals(tc.left)
		if actualRL != tc.expected {
			t.Fatalf("%v.Equals(%v) should return %v, got %v", tc.right, tc.left, tc.expected, actualRL)
		}
	}
}

func TestIntervalIntersect(t *testing.T) {
	cases := []BinaryResult[Interval, Interval]{
		BinaryResult[Interval, Interval]{New(0, 5), New(6, 10), Empty()},
		BinaryResult[Interval, Interval]{New(0, 7), New(5, 10), New(5, 7)},
		BinaryResult[Interval, Interval]{New(3, 7), New(0, 10), New(3, 7)},
		BinaryResult[Interval, Interval]{Empty(), New(3, 5), Empty()},
	}

	for _, tc := range cases {
		if actual := tc.left.Intersect(tc.right); !actual.Equals(tc.expected) {
			t.Fatalf("%v.Intersect(%v) should be %v, got %v", tc.left, tc.right, tc.expected, actual)
		}
		if actual := tc.right.Intersect(tc.left); !actual.Equals(tc.expected) {
			t.Fatalf("%v.Intersect(%v) should be %v, got %v", tc.right, tc.left, tc.expected, actual)
		}
	}
}

func TestIntervalsEquals(t *testing.T) {
	cases := []BinaryResult[Intervals, bool]{
		BinaryResult[Intervals, bool]{Intervals{}, Intervals{}, true},
		BinaryResult[Intervals, bool]{Intervals{}, Intervals{New(1, 2)}, false},
		BinaryResult[Intervals, bool]{Intervals{New(1, 2)}, Intervals{New(1, 2)}, true},
		BinaryResult[Intervals, bool]{Intervals{New(1, 2), New(3, 4)}, Intervals{New(3, 4), New(1, 2)}, true},
	}

	for _, tc := range cases {
		if actual := tc.left.Equals(tc.right); actual != tc.expected {
			t.Fatalf("%v.Equals(%v) should return %v, got %v", tc.left, tc.right, tc.expected, actual)
		}
		if actual := tc.right.Equals(tc.left); actual != tc.expected {
			t.Fatalf("%v.Equals(%v) should return %v, got %v", tc.right, tc.left, tc.expected, actual)
		}
	}
}

func TestIntervalsCompact(t *testing.T) {
	cases := []UnaryResult[Intervals, Intervals]{
		UnaryResult[Intervals, Intervals]{Intervals{}, Intervals{}},
		UnaryResult[Intervals, Intervals]{Intervals{New(1, 2)}, Intervals{New(1, 2)}},
		UnaryResult[Intervals, Intervals]{Intervals{New(1, 2), New(2, 3)}, Intervals{New(1, 3)}},
		UnaryResult[Intervals, Intervals]{Intervals{New(1, 2), New(3, 4)}, Intervals{New(1, 4)}},
		UnaryResult[Intervals, Intervals]{Intervals{New(1, 2), New(4, 5)}, Intervals{New(1, 2), New(4, 5)}},
		UnaryResult[Intervals, Intervals]{Intervals{Empty(), New(1, 2)}, Intervals{New(1, 2)}},
		UnaryResult[Intervals, Intervals]{Intervals{Empty(), New(1, 2), New(1, 2), Empty()}, Intervals{New(1, 2)}},
	}

	for _, tc := range cases {
		if actual := tc.left.Compact(); !actual.Equals(tc.expected) {
			t.Fatalf("%v.Compact() should return %v, got %v", tc.left, tc.expected, actual)
		}
	}
}

func TestIntervalMinus(t *testing.T) {
	cases := []BinaryResult[Interval, Intervals]{
		{Empty(), New(1, 2), Intervals{}},
		{New(1, 2), Empty(), Intervals{New(1, 2)}},
		{New(1, 2), New(1, 1), Intervals{New(2, 2)}},
		{New(1, 3), New(2, 2), Intervals{New(1, 1), New(3, 3)}},
	}

	for _, tc := range cases {
		if actual := tc.left.Minus(tc.right); !actual.Equals(tc.expected) {
			t.Fatalf("%v.Minus(%v) should return %v, got %v", tc.left, tc.right, tc.expected, actual)
		}
	}
}
