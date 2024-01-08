package maglev_hash

import (
	"fmt"
	"sort"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestBuildPreferences(t *testing.T) {
	cases := []struct {
		slotCnt int
		nodeCnt int
	}{
		{7, 2},
		{7, 3},
		{DefaultSlotCnt, 10},
	}
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			mh := MaglevHash{
				slotCnt: tc.slotCnt,
				nodes: lo.Times(tc.nodeCnt, func(i int) string {
					return fmt.Sprintf("B%d", i)
				}),
			}
			prefList := mh.buildPreferences()
			expectedListAfterSort := lo.Times(tc.slotCnt, func(index int) int { return index })
			assert.Equal(t, len(prefList), tc.nodeCnt)
			for _, l := range prefList {
				newList := lo.Uniq(l)
				sort.Ints(newList)
				assert.Equal(t, expectedListAfterSort, newList)
			}
		})
	}
}

func TestBuildLookup_1(t *testing.T) {
	m := &MaglevHash{
		slotCnt: 7,
		nodeCnt: 3,
		nodes:   []string{"B0", "B1", "B2"},
	}

	lookup := m.buildLookup(
		[][]int{
			[]int{3, 0, 4, 1, 5, 2, 6},
			[]int{0, 2, 4, 6, 1, 3, 5},
			[]int{3, 4, 5, 6, 0, 1, 2},
		})

	// slot assignment: B1, B0, B1, B0, B2, B2, B0
	assert.Equal(t, []int{1, 0, 1, 0, 2, 2, 0}, lookup)
}

func TestBuildLookup_2(t *testing.T) {
	m := &MaglevHash{
		slotCnt: 7,
		nodeCnt: 2,
		nodes: []string{
			"B0", "B2",
		},
	}

	lookup := m.buildLookup([][]int{
		[]int{3, 0, 4, 1, 5, 2, 6},
		[]int{3, 4, 5, 6, 0, 1, 2},
	})

	// slot assignment: B0, B0, B0, B0, B2, B2, B2
	assert.Equal(t, []int{0, 0, 0, 0, 1, 1, 1}, lookup)
}
