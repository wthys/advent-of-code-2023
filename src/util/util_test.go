package util

import (
	"fmt"
	"testing"
)

type (
	testInt struct {
		Input int
		Want  int
	}

	testFloat struct {
		Input float64
		Want  float64
	}
)

func TestSignInt(t *testing.T) {
	cases := []testInt{
		{int(5), int(1)},
		{int(0), int(0)},
		{int(-245), int(-1)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Sign(cs.Input)
			if s != cs.Want {
				t.Fatalf("Sign(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func TestSignFloat(t *testing.T) {
	cases := []testFloat{
		{float64(3.14), float64(1.0)},
		{float64(0.0), float64(0.0)},
		{float64(-5.256), float64(-1.0)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Sign(cs.Input)
			if s != cs.Want {
				t.Fatalf("Sign(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func TestAbsInt(t *testing.T) {
	cases := []testInt{
		{int(5), int(5)},
		{int(0), int(0)},
		{int(-245), int(245)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Abs(cs.Input)
			if s != cs.Want {
				t.Fatalf("Abs(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func TestAbsFloat(t *testing.T) {
	cases := []testFloat{
		{float64(3.14), float64(3.14)},
		{float64(0.0), float64(0.0)},
		{float64(-5.256), float64(5.256)},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
			s := Abs(cs.Input)
			if s != cs.Want {
				t.Fatalf("Abs(%v) = %v, want %v", cs.Input, s, cs.Want)
			}
		})
	}
}

func hash(array []int) int {
	h := 0
	for _, value := range array {
		h = 13*h + value
	}
	return h
}

func TestPermutationDo(t *testing.T) {
	array := []int{1, 2, 3}

	check := map[int]bool{
		hash([]int{1, 2, 3}): false,
		hash([]int{1, 3, 2}): false,
		hash([]int{2, 1, 3}): false,
		hash([]int{2, 3, 1}): false,
		hash([]int{3, 2, 1}): false,
		hash([]int{3, 1, 2}): false,
	}

	PermutationDo(3, array, func(perm []int) {
		check[hash(perm)] = true
	})

	for value, seen := range check {
		if !seen {
			t.Fatalf("PermutationDo(3, %v, ...) should produce %v, but was not seen", array, value)
		}
	}

}

func TestForEach(t *testing.T) {
	array := []int{1, 2, 3}
	check := map[int]bool{
		1: false,
		2: false,
		3: false,
	}

	ForEach(array, func(value int) {
		check[value] = true
	})

	for value, seen := range check {
		if !seen {
			t.Fatalf("Do(%v, ...) should produce %v, but was not seen", array, value)
		}
	}
}

func TestForEachError(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}
	check := map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
	}

	err := ForEachError(array, func(value int) error {
		if value > 3 {
			return fmt.Errorf("value %v is too large", value)
		}
		check[value] = true
		return nil
	})

	if err == nil {
		t.Fatalf("ForEachError(%v, ...) should produce an error but none was seen", array)
	}

	for value, seen := range check {
		if !seen && value <= 3 {
			t.Fatalf("ForEachError(%v, ...) should produce %v, but was not seen", array, value)
		}
	}
}

func TestForEachStopping(t *testing.T) {
	array := []int{1, 2, 3, 4, 5}
	check := map[int]bool{
		1: false,
		2: false,
		3: false,
		4: false,
		5: false,
	}

	stopped := ForEachStopping(array, func(value int) bool {
		if value > 3 {
			return false
		}
		check[value] = true
		return true
	})

	if !stopped {
		t.Fatalf("ForEachStopping(%v, ...) should produce %v but was %v", array, true, stopped)
	}

	for value, seen := range check {
		if !seen && value <= 3 {
			t.Fatalf("ForEachStopping(%v, ...) should produce %v, but was not seen", array, value)
		}
	}
}
