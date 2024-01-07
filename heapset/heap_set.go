// Package heapset implements a key-value store where key-value pairs can be
// accessed either by key or the order of values.
// The value must be orderable (a < b) and the key must be comparable (a == b).
//
// Heapset combines the hash map and heap. One example is to store
// a set of jobs where each job has a priority. And we want to query jobs by priority
// and update a job's priority.

// If the value is a composite type, the recommended practice is to use pointer type.
// For example,
//
//	type Job struct {
//		expire time.Time
//	  // other fields
//	}
//	hs := heapset.NewHeapSet[int, *Job](func(v1, v2 *Job) int {
//	  switch {
//	  case v1.expiration.Before(v2.expiration):
//	  	return -1
//	  case v1.expiration.After(v2.expiration):
//	  	return 1
//	  default:
//	  	return 0
//	  }
//	})
//
// jobID := 1
//
// job := Job{id: jobID, expire: time.Now(), ....}
//
// hs.Set(jobID, &job)
package heapset

import (
	"container/heap"
)

// HeapSet keeps key-value pairs in a hash map and provides access to the pair of
// the minimum value.
type HeapSet[K comparable, V any] struct {
	// heap
	h heap.Interface

	// hashmap
	s map[K]*Element[K, V]

	emptyK K
	emptyV V
}

// NewHeapSet returns a HeapSet where values are ordered by the given compare function.
func NewHeapSet[K comparable, V any](compare func(v1, v2 V) int) *HeapSet[K, V] {
	hs := HeapSet[K, V]{
		h: newHeapStruct[K, V](compare),
		s: make(map[K]*Element[K, V]),
	}
	heap.Init(hs.h)
	return &hs
}

// Set inserts a k-v pair if the key does not exist. Otherwise, Set updates the value.
func (hs *HeapSet[K, V]) Set(k K, v V) {
	existingElement, ok := hs.s[k]
	if !ok {
		e := Element[K, V]{
			Key:   k,
			Value: v,
		}
		heap.Push(hs.h, &e)
		hs.s[k] = &e
		return
	}
	existingElement.Value = v
	heap.Fix(hs.h, existingElement.index)
}

// Get returns the value associated with the key
func (hs *HeapSet[K, V]) Get(k K) (V, bool) {
	e, ok := hs.s[k]
	if !ok {
		return hs.emptyV, false
	}
	return e.Value, true
}

// Delete deletes the key-value pair of the key.
func (hs *HeapSet[K, V]) Delete(key K) {
	item, ok := hs.s[key]
	if !ok {
		return
	}
	delete(hs.s, key)
	index := item.index
	if index == hs.h.Len()-1 {
		hs.h.Pop()
		return
	}
	hs.h.Swap(index, hs.h.Len()-1)
	hs.h.Pop()
	heap.Fix(hs.h, index)
}

// Top returns the key-value pair of the smallest value. It returns false
// if the set is empty.
func (hs *HeapSet[K, V]) Top() (K, V, bool) {
	if hs.h.Len() <= 0 {
		return hs.emptyK, hs.emptyV, false
	}

	var e *Element[K, V]
	h := hs.h.(*heapStruct[K, V])
	e = h.e[0]

	return e.Key, e.Value, true
}

// Pop removes and returns the key-value pair of the smallest value. It returns flase
// if the set is empty.
func (hs *HeapSet[K, V]) Pop() (K, V, bool) {
	if hs.h.Len() == 0 {
		return hs.emptyK, hs.emptyV, false
	}
	e := heap.Pop(hs.h).(*Element[K, V])
	delete(hs.s, e.Key)
	return e.Key, e.Value, true
}

// Size returns the number of key-value pairs.
func (pq *HeapSet[K, V]) Size() int {
	return pq.h.Len()
}

// Map returns the underlying map. It is here to provide an efficient way of
// iterating over all key-value pairs.
func (hs *HeapSet[K, V]) Map() map[K]*Element[K, V] {
	return hs.s
}

// Element is the unit of data stored in hash map and the heap.
type Element[K comparable, V any] struct {
	Key   K
	Value V

	// The array index of the element in the heap.
	index int
}

// heapStruct implements the heap.Interface.
type heapStruct[K comparable, V any] struct {
	e       []*Element[K, V]
	compare func(v1, v2 V) int
}

func newHeapStruct[K comparable, V any](compare func(v1, v2 V) int) *heapStruct[K, V] {
	return &heapStruct[K, V]{
		compare: compare,
	}
}

func (h *heapStruct[K, V]) Len() int {
	return len(h.e)
}

func (h *heapStruct[K, V]) Less(i, j int) bool {
	return h.compare(h.e[i].Value, h.e[j].Value) < 0
}

func (h *heapStruct[K, V]) Swap(i, j int) {
	h.e[i], h.e[j] = h.e[j], h.e[i]
	h.e[i].index = i
	h.e[j].index = j
}

func (h *heapStruct[K, V]) Push(x any) {
	n := len(h.e)
	item := x.(*Element[K, V])
	item.index = n
	h.e = append(h.e, item)
}

func (h *heapStruct[K, V]) Pop() any {
	old := h.e
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	h.e = old[0 : n-1]
	return item
}
