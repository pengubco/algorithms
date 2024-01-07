// Package binary_search implements binary search algorithms on a slice of
// non decreasing values.
//
// Example usage.
//
// s := NewBinarySearcher[int](func(v1, v2 int) int { return v1 - v2 })
// s.LowerBound([]int{1, 3}, 0, 2, 1) //  0
// s.UpperBound([]int{1, 3}, 0, 2, 1) //  1

package binary_search

// BinarySearcher searches element from a slice of non decreasing elements.
type BinarySearcher[V any] struct {
	// -1: v1 < v2; 0: v1 == v2; >0: v1 > v2.
	compare func(v1, v2 V) int
}

// NewBinarySearcher creates a BinarySearcher using the compare function.
func NewBinarySearcher[V any](compare func(v1, v2 V) int) *BinarySearcher[V] {
	return &BinarySearcher[V]{
		compare: compare,
	}
}

// LowerBound finds the leftmost position such that the element at the position
// is larger than or equal to val. Imagine we'd insert val to the sorted slice and
// keep the slice sorted, the returned position is the lower bound where we can insert.
// The range is half close [first, last).
// Return last if all elements are less than val.
func (s *BinarySearcher[V]) LowerBound(elements []V, first int, last int, val V) int {
	cnt := last - first
	for cnt > 0 {
		half := cnt >> 1
		middle := first + half
		if s.compare(elements[middle], val) < 0 { // This is the only difference with UpperBound
			first = middle
			first++
			cnt -= half + 1
		} else {
			cnt = half
		}
	}
	return first
}

// UpperBound finds the leftmost position such that the element at the position
// is larger than val. Imagine we'd insert val to the sorted slice and keep the slice sorted,
// the returned position is the upper bound where we can insert.
// The range is half close [first, last).
// Return last if all elements are less than or equal to val.
func (s *BinarySearcher[V]) UpperBound(elements []V, first int, last int, val V) int {
	cnt := last - first
	for cnt > 0 {
		half := cnt >> 1
		middle := first + half
		if s.compare(elements[middle], val) <= 0 {
			first = middle
			first++
			cnt -= half + 1
		} else {
			cnt = half
		}
	}
	return first
}
