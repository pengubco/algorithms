package util_test

import (
	"testing"

	"github.com/pengubco/ads/util"
	"github.com/stretchr/testify/assert"
)

func TestSortAndUniq(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, util.SortAndUniq([]int{3, 1, 1, 2, 3, 2}))
}
