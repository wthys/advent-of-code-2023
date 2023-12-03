package set

import (
	"fmt"
	"strings"
)

type (
	empty struct{}

	Set[T comparable] struct {
		contents map[T]empty
	}

	ForEachFunction[T comparable]           func(value T)
	ForEachStoppingFunction[T comparable]   func(value T) bool
	MapFunction[T comparable, R comparable] func(value T) R
)

func New[T comparable](values ...T) *Set[T] {
	set := Set[T]{map[T]empty{}}
	for _, value := range values {
		set.Add(value)
	}
	return &set
}

func (s Set[T]) String() string {
	str := strings.Builder{}
	fmt.Fprint(&str, "<")
	s.ForEach(func(value T) {
		fmt.Fprintf(&str, " %v", value)
	})
	fmt.Fprint(&str, " >")
	return str.String()
}

func (s Set[T]) Values() []T {
	vals := []T{}
	s.ForEach(func(value T) {
		vals = append(vals, value)
	})
	return vals
}

func (s *Set[T]) Add(value T) *Set[T] {
	(*s).contents[value] = empty{}
	return s
}

func (s Set[T]) Has(value T) bool {
	_, ok := s.contents[value]
	return ok
}

func (s Set[T]) Len() int {
	return len(s.contents)
}

func (s *Set[T]) Remove(value T) *Set[T] {
	delete((*s).contents, value)
	return s
}

func (s Set[T]) Intersect(other *Set[T]) *Set[T] {
	common := New[T]()

	s.ForEach(func(value T) {
		if other.Has(value) {
			common.Add(value)
		}
	})

	return common
}

func (s Set[T]) Union(other *Set[T]) *Set[T] {
	union := New[T]()
	adder := func(value T) {
		union.Add(value)
	}

	s.ForEach(adder)
	other.ForEach(adder)

	return union
}

func (s Set[T]) Subtract(other *Set[T]) *Set[T] {
	sub := New[T]()

	s.ForEach(func(value T) {
		if !other.Has(value) {
			sub.Add(value)
		}
	})

	return sub
}

// Iterates over all elements as long as the result of the forEach function is true.
// Returns true when it was stopped early, false otherwise.
func (s Set[T]) ForEachStopping(forEach ForEachStoppingFunction[T]) bool {
	for value, _ := range s.contents {
		if !forEach(value) {
			return true
		}
	}
	return false
}

func (s Set[T]) ForEach(forEach ForEachFunction[T]) {
	for value, _ := range s.contents {
		forEach(value)
	}
}
