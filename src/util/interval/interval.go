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

func (i Interval) String() string {
	return fmt.Sprintf("[%v,%v]", i.low, i.high)
}

func (i Interval) Minus(o Interval) Intervals {
	if o.low > i.high || o.high < i.low {
		return Intervals{i}
	}
	if o.low > i.low && o.high < i.high {
		return Intervals{Interval{i.low, o.low - 1}, Interval{o.high + 1, i.high}}
	}

	return Intervals{Interval{max(i.low, min(o.high, i.high)), min(i.high, max(i.low, o.low))}}
}

func (i Interval) Plus(o Interval) Intervals {
	if i.high < o.low {
		return Intervals{i, o}
	}
	if i.low > o.high {
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
	return i.low <= val && val <= i.high
}

func (i Interval) Len() int {
	return i.high - i.low + 1
}

func (this Intervals) Compact() Intervals {
	ordered := slices.Clone[Intervals](this)
	slices.SortFunc[Intervals, Interval](ordered, func(a Interval, b Interval) int {
		return a.Compare(b)
	})

	compacted := Intervals{}
	current := ordered[0]
	for _, next := range ordered[1:] {
		merged := current.Plus(next)
		if len(merged) < 2 {
			current = merged[0]
		} else {
			compacted = append(compacted, current)
			current = next
		}
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

func (i Interval) ForEach(forEach func(int) bool) bool {
	for val := i.low; val <= i.high; val++ {
		if !forEach(val) {
			return false
		}
	}
	return true
}
