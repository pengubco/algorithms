// Package stack implements the First-In-First-Out (FIFO) stack.
//
// Example usage.
// s := stack.NewStack[int]()
// s.Push(1)
// s.Size(0) // 1
// x, err := s.Top() // 1, nil
// x, err := s.Pop() // 1, nil
// s.IsEmpty() // true
// x, err := s.Top() // 0, ErrEmpty
// x, err := s.Pop() // 0, ErrEmpty

package stack

import "errors"

var (
	ErrEmpty = errors.New("empty stack")
)

// Stack stores elements in a slice and offers First-In-First-Out access.
type Stack[V any] struct {
	elements []V
	size     int

	emptyV V
}

// NewStack creates an empty stack.
func NewStack[V any]() *Stack[V] {
	return &Stack[V]{}
}

// Push pushes a new element to the top of stack.
func (s *Stack[V]) Push(x V) {
	s.elements = append(s.elements, x)
	s.size++
}

// Pop removes the top of stack and returns it, if the stack is not empty.
// Returns ErrEmpty if the stack is empty.
func (s *Stack[V]) Pop() (V, error) {
	if s.size == 0 {
		return s.emptyV, ErrEmpty
	}
	v := s.elements[s.size-1]
	s.size--
	s.elements = s.elements[:s.size]
	return v, nil
}

// Top returns the top of stack, if the stack is not empty.
// Returns ErrEmpty if the stack is empty.
func (s *Stack[V]) Top() (V, error) {
	if s.size > 0 {
		return s.elements[s.size-1], nil
	}
	return s.emptyV, ErrEmpty
}

// IsEmpty returns true iff the stack is empty.
func (s *Stack[V]) IsEmpty() bool {
	return s.size == 0
}

// Size returns the number of elements in the stack.
func (s *Stack[V]) Size() int {
	return s.size
}
