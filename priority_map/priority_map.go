// Package priority_map implements a key-value store where key-value pairs can be
// accessed either by key or the order of values.
// The value must be orderable (a < b) and the key must be comparable (a == b).
//
// 1. Get(K) V
// 2. Set(K, V)
// 3. Delete(K)
// 4. Top() (K, V, bool)
// 5. Pop() (K, V, bool)
// 6. Size() int
//
// Usage
//
//	pm, _ := NewPriorityMap[int, string](func (v1, v2 string) bool {
//			return v1 < v2
//	})
//
// pm.Set(1, "a")
// pm.Get(1) // returns "a"
// pm.Set(2, "b")
// pm.Set(3, "c")
// pm.Top() // returns (1, "a", true)
// pm.Pop() // returns (1, "a", true)
// pm.Delete(1)
// pm.Size() // returns 2
//
// See more usage example in the priority_map_test.go.
package priority_map

import (
	"container/heap"
)

// PriorityMap keeps key-value pairs in a hash map and provides access to the pair of
// the minimum value.
type PriorityMap[K comparable, V any] struct {
	// heap
	h heap.Interface

	// hashmap
	m map[K]*Element[K, V]

	emptyK K
	emptyV V
}

// NewPriorityMap returns a PriorityMap where values are ordered by the given less function.
func NewPriorityMap[K comparable, V any](less func(v1, v2 V) bool) *PriorityMap[K, V] {
	hs := PriorityMap[K, V]{
		h: newHeapStruct[K, V](less),
		m: make(map[K]*Element[K, V]),
	}
	heap.Init(hs.h)
	return &hs
}

// Set inserts a k-v pair if the key does not exist. Otherwise, Set updates the value.
func (pm *PriorityMap[K, V]) Set(k K, v V) {
	existingElement, ok := pm.m[k]
	if !ok {
		e := Element[K, V]{
			Key:   k,
			Value: v,
		}
		heap.Push(pm.h, &e)
		pm.m[k] = &e
		return
	}
	existingElement.Value = v
	heap.Fix(pm.h, existingElement.index)
}

// Get returns the value associated with the key
func (pm *PriorityMap[K, V]) Get(k K) (V, bool) {
	e, ok := pm.m[k]
	if !ok {
		return pm.emptyV, false
	}
	return e.Value, true
}

// Delete deletes the key-value pair of the key.
func (pm *PriorityMap[K, V]) Delete(key K) {
	item, ok := pm.m[key]
	if !ok {
		return
	}
	delete(pm.m, key)
	index := item.index
	if index == pm.h.Len()-1 {
		pm.h.Pop()
		return
	}
	pm.h.Swap(index, pm.h.Len()-1)
	pm.h.Pop()
	heap.Fix(pm.h, index)
}

// Top returns the key-value pair of the smallest value. It returns false
// if the set is empty.
func (pm *PriorityMap[K, V]) Top() (K, V, bool) {
	if pm.h.Len() <= 0 {
		return pm.emptyK, pm.emptyV, false
	}

	var e *Element[K, V]
	h := pm.h.(*heapStruct[K, V])
	e = h.e[0]

	return e.Key, e.Value, true
}

// Pop removes and returns the key-value pair of the smallest value. It returns flase
// if the set is empty.
func (pm *PriorityMap[K, V]) Pop() (K, V, bool) {
	if pm.h.Len() == 0 {
		return pm.emptyK, pm.emptyV, false
	}
	e := heap.Pop(pm.h).(*Element[K, V])
	delete(pm.m, e.Key)
	return e.Key, e.Value, true
}

// Size returns the number of key-value pairs.
func (pm *PriorityMap[K, V]) Size() int {
	return pm.h.Len()
}

// Map returns the underlying map. It is here to provide an efficient way of
// iterating over all key-value pairs.
func (pm *PriorityMap[K, V]) Map() map[K]*Element[K, V] {
	return pm.m
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
	e    []*Element[K, V]
	less func(v1, v2 V) bool
}

func newHeapStruct[K comparable, V any](less func(v1, v2 V) bool) *heapStruct[K, V] {
	return &heapStruct[K, V]{
		less: less,
	}
}

func (h *heapStruct[K, V]) Len() int {
	return len(h.e)
}

func (h *heapStruct[K, V]) Less(i, j int) bool {
	return h.less(h.e[i].Value, h.e[j].Value)
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
