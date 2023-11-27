package binary_search_test

import (
	"testing"

	"github.com/pengubco/ads/binary_search"
	"github.com/stretchr/testify/assert"
)

func TestLowerBoundAndUpperBound(t *testing.T) {
	cases := []struct {
		a          []int
		first      int
		last       int
		val        int
		lowerbound int
		upperbound int
	}{
		{[]int{1, 3}, 0, 2, 0, 0, 0},
		{[]int{1, 3}, 0, 2, 1, 0, 1},
		{[]int{1, 3}, 0, 2, 2, 1, 1},
		{[]int{1, 3}, 0, 2, 3, 1, 2},
		{[]int{1, 3}, 0, 2, 4, 2, 2},
		{[]int{1, 3, 5}, 0, 3, 0, 0, 0},
		{[]int{1, 3, 5}, 0, 3, 3, 1, 2},
		{[]int{1, 3, 5}, 0, 3, 4, 2, 2},
		{[]int{1, 3, 5}, 0, 3, 5, 2, 3},
		{[]int{1, 3, 5}, 0, 3, 6, 3, 3},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			s := binary_search.NewBinarySearcher[int](func(v1, v2 int) int { return v1 - v2 })
			assert.Equal(t, tc.lowerbound, s.LowerBound(tc.a, tc.first, tc.last, tc.val))
			assert.Equal(t, tc.upperbound, s.UpperBound(tc.a, tc.first, tc.last, tc.val))
		})
	}
}
