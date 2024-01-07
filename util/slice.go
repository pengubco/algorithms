package util

import "sort"

func SortAndUniq(a []int) []int {
	m := len(a)
	tmp := make([]int, m)
	copy(tmp, a)
	sort.Ints(tmp)
	deduped := make([]int, 0)
	deduped = append(deduped, tmp[0])
	j := 0
	for i := 1; i < m; i++ {
		if tmp[i] != deduped[j] {
			deduped = append(deduped, tmp[i])
			j++
		}
	}
	return deduped
}
