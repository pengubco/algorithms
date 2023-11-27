// Package rmq implements a range-minimum-query data structure that
// 1. precomputes a sparse table on O(N*logN) time
// 2. answers the minimum value of any subarray of a static array in O(1) time.
//
// Example usage.
//
//	q := NewRMQ[int]([]int{6, 1, 0, 10, 9}, func(v1, v2 int) bool { return v1 < v2 })
//	q.RMQ(0, 0) // 6
//	q.RMQ(0,1) // 1
//
// See https://cp-algorithms.com/sequences/rmq.html for more info.
package rmq

// RMQ is a data structure that precomputes a sparse table on O(N*logN) time and
// answers the minimum value of any subarray of a static array in O(1) time.
type RMQ[V any] struct {
	log2 []int

	// The original array of elements.
	elements []V

	// the sparse table
	// st[i][j]: the min of range [a[j], a[j +2^i - 1]], a sub array of length 2^i.
	st [][]V

	less func(v1, v2 V) bool
}

// NewRMQ creates a RMQ with the elements and less function.
func NewRMQ[V any](elements []V, less func(v1, v2 V) bool) *RMQ[V] {
	n := len(elements)
	q := RMQ[V]{
		log2:     calcLog2n(n),
		elements: elements,
		less:     less,
	}
	q.calcSparseTable()
	return &q
}

// RMQ returns the minimum value between elements[l:r], inclusive at both ends.
func (q *RMQ[V]) RMQ(l, r int) V {
	k := q.log2[r-l+1]
	result := q.st[k][r-(1<<k)+1]
	if q.less(q.st[k][l], result) {
		result = q.st[k][l]
	}
	return result
}

func (q *RMQ[V]) calcSparseTable() {
	n := len(q.elements)
	m := q.log2[n]
	q.st = make([][]V, m+1)
	for i := 0; i <= m; i++ {
		q.st[i] = make([]V, n)
	}
	for j := 0; j < n; j++ {
		q.st[0][j] = q.elements[j]
	}

	for i := 1; i <= m; i++ {
		for j := 0; j+(1<<i)-1 < n; j++ {
			q.st[i][j] = q.st[i-1][j+(1<<(i-1))]
			if q.less(q.st[i-1][j], q.st[i][j]) {
				q.st[i][j] = q.st[i-1][j]
			}
		}
	}
}

// calcLog2n log_2. Can use int(math.Log2(x)) instead.
func calcLog2n(n int) []int {
	log2 := make([]int, n+1)
	log2[1] = 0
	for i := 2; i <= n; i++ {
		log2[i] = log2[i/2] + 1
	}
	return log2
}
