package stack_test

import (
	"testing"

	"github.com/pengubco/ads/stack"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := stack.NewStack[int]()
	assert.True(t, s.IsEmpty())
	_, e := s.Top()
	assert.Equal(t, stack.ErrEmpty, e)
	_, e = s.Pop()
	assert.Equal(t, stack.ErrEmpty, e)
	s.Push(1)
	s.Push(2)
	x, e := s.Top()
	assert.NoError(t, e)
	assert.Equal(t, 2, x)
	x, e = s.Top()
	assert.NoError(t, e)
	assert.Equal(t, 2, x)
	x, e = s.Pop()
	assert.NoError(t, e)
	assert.Equal(t, 2, x)
	x, e = s.Pop()
	assert.NoError(t, e)
	assert.Equal(t, 1, x)
	assert.Equal(t, 0, s.Size())
	assert.True(t, s.IsEmpty())
}
