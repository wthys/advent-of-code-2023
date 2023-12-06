package interval

import (
	"fmt"
	"slices"

	"github.com/wthys/advent-of-code-2023/util"
)

type (
	Interval struct {
		low  int
		high int
	}

	Intervals []Interval
)

func New(low, high int) Interval {
	return Interval{min(low, high), max(low, high)}
}

func Empty() Interval {
	return Interval{0, -1}
}

func (i Interval) String() string {
	return fmt.Sprintf("[%v,%v]", i.low, i.high)
}

func (i Interval) Minus(o Interval) Intervals {
	// ---o--hi|  |lo--i--hi|  |lo--o--
	if i.low > o.high || i.high < o.low {
		return Intervals{i}
	}

	//      |lo----------i----------hi|
	// ---o---hi|  |lo---o---hi|  |lo---o---
	return Intervals{Interval{i.low, o.low - 1}, Interval{o.high + 1, i.high}}.Compact()
}

func (i Interval) Plus(o Interval) Intervals {
	if i.IsEmpty() {
		return Intervals{o}.Compact()
	}
	if o.IsEmpty() {
		return Intervals{i}.Compact()
	}
	if i.high+1 < o.low {
		return Intervals{i, o}
	}
	if i.low+1 > o.high {
		return Intervals{o, i}
	}
	return Intervals{Interval{min(i.low, o.low), max(i.high, o.high)}}
}

func (i Interval) Compare(o Interval) int {
	diff := i.low - o.low
	if diff != 0 {
		return util.Sign(diff)
	}
	return util.Sign(i.high - o.high)
}

func (this Intervals) Len() int {
	size := 0
	for _, ivl := range this {
		size += ivl.Len()
	}
	return size
}

func (this Intervals) Contains(val int) bool {
	for _, ivl := range this {
		if ivl.Contains(val) {
			return true
		}
	}
	return false
}

func (i Interval) Contains(val int) bool {
	if i.IsEmpty() {
		return false
	}
	return i.low <= val && val <= i.high
}

func (i Interval) IsEmpty() bool {
	return i.low > i.high
}

func (i Interval) Len() int {
	if i.IsEmpty() {
		return 0
	}
	return i.high - i.low + 1
}

func (i Interval) Lower() int {
	return i.low
}

func (i Interval) Upper() int {
	return i.high
}

func (this Interval) Intersect(that Interval) Interval {
	if this.IsEmpty() || that.IsEmpty() {
		return Empty()
	}

	return Interval{max(this.low, that.low), min(this.high, that.high)}
}

func (this Interval) Equals(that Interval) bool {
	if this.IsEmpty() {
		return that.IsEmpty()
	}
	if that.IsEmpty() {
		return false
	}
	return this.low == that.low && this.high == that.high
}

func (i Interval) ForEach(forEach func(int) bool) bool {
	for val := i.low; val <= i.high; val++ {
		if !forEach(val) {
			return false
		}
	}
	return true
}

func (this Intervals) Equals(that Intervals) bool {
	orderA := slices.DeleteFunc[Intervals, Interval](this, func(a Interval) bool {
		return a.IsEmpty()
	})
	slices.SortFunc[Intervals, Interval](orderA, func(a, b Interval) int {
		return a.Compare(b)
	})
	orderB := slices.DeleteFunc[Intervals, Interval](that, func(a Interval) bool {
		return a.IsEmpty()
	})
	slices.SortFunc[Intervals, Interval](orderB, func(a, b Interval) int {
		return a.Compare(b)
	})

	return slices.Equal(orderA, orderB)
}

func (this Intervals) Compact() Intervals {
	if len(this) == 0 {
		return Intervals{}
	}

	ordered := slices.Clone[Intervals](this)
	slices.SortFunc[Intervals, Interval](ordered, func(a Interval, b Interval) int {
		return a.Compare(b)
	})

	compacted := Intervals{}
	current := ordered[0]
	for _, next := range ordered[1:] {
		merged := current.Plus(next)
		if len(merged) == 1 {
			current = merged[0]
		} else if len(merged) > 1 {
			compacted = append(compacted, current)
			current = next
		}
	}

	if current.IsEmpty() {
		return compacted
	}
	return append(compacted, current)
}

func (this Intervals) Add(that Intervals) Intervals {
	return append(this, that...).Compact()
}

func (this Intervals) ForEach(forEach func(int) bool) bool {
	for _, i := range this {
		if !i.ForEach(forEach) {
			return false
		}
	}
	return true
}
