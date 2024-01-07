// Package priorityqueue implements Priority Queue with four methods
// 1. Push(V).
// 2. Top() (V, error).
// 3. Pop() (V, error)
// 4. Size() int.
//
// Usage
// pq, err := NewPriorityQueue(func(v1, v2 int){return v1<v2})
// pq.Push(10)
// v, _ := pq.Top()
// v, _ = pq.Pop()
//
// See usage example in the priority_queue_test.go.
//
// Why implement priority queue?
// The container/heap is not easy and intuitive to use because:
// 1. It requires boilerplate code to implement the heap interface.
// 2. Push/pop operations need to use heap.Push and heap.Pop on an array, which
// is not intuitive.
package priorityqueue

import (
	"container/heap"
	"errors"
)

var ErrQueueIsEmpty = errors.New("queue is empty")

// PriorityQueue implements priority queue on any type as long as there is a
// less function to tell the "less-than" relationship between values of the type.
type PriorityQueue[V any] struct {
	hs *heapStruct[V]

	emptyV V
}

// NewPriorityQueue returns a priority queue. Returns error when the less function
// is nil.
func NewPriorityQueue[V any](less func(v1, v2 V) bool) (*PriorityQueue[V], error) {
	if less == nil {
		return nil, errors.New("must provide the compare function")
	}
	pq := &PriorityQueue[V]{
		hs: newHeapStruct[V](less),
	}
	heap.Init(pq.hs)
	return pq, nil
}

// Push inserts a value to the priority queue.
func (pq *PriorityQueue[V]) Push(v V) {
	e := heapElement[V]{
		Value: v,
	}
	heap.Push(pq.hs, &e)
}

// Pop removes and returns the smallest value if the queue is not empty.
func (pq *PriorityQueue[V]) Pop() (V, error) {
	v, err := pq.Top()
	if err != nil {
		return pq.emptyV, err
	}
	heap.Pop(pq.hs)
	return v, nil
}

// Top returns the smallest value if the queue is not empty.
func (pq *PriorityQueue[V]) Top() (V, error) {
	if pq.hs.Len() == 0 {
		return pq.emptyV, ErrQueueIsEmpty
	}
	return pq.hs.e[0].Value, nil
}

// Size returns the number of values in the queue.
func (pq *PriorityQueue[V]) Size() int {
	return pq.hs.Len()
}

// ==== internal types implementing the heap.Interface ====
// heapElement is the unit of data stored in the heap.
type heapElement[V any] struct {
	Value V

	// The array index of the element in the heap.
	index int
}

// heapStruct implements the heap.Interface.
type heapStruct[V any] struct {
	e    []*heapElement[V]
	less func(v1, v2 V) bool
}

func newHeapStruct[V any](less func(v1, v2 V) bool) *heapStruct[V] {
	return &heapStruct[V]{
		less: less,
	}
}

func (h *heapStruct[V]) Len() int {
	return len(h.e)
}

func (h *heapStruct[V]) Less(i, j int) bool {
	return h.less(h.e[i].Value, h.e[j].Value)
}

func (h *heapStruct[V]) Swap(i, j int) {
	h.e[i], h.e[j] = h.e[j], h.e[i]
	h.e[i].index = i
	h.e[j].index = j
}

func (h *heapStruct[V]) Push(x any) {
	n := len(h.e)
	item := x.(*heapElement[V])
	item.index = n
	h.e = append(h.e, item)
}

func (h *heapStruct[V]) Pop() any {
	old := h.e
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	h.e = old[0 : n-1]
	return item
}
