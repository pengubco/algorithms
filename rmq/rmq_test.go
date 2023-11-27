package rmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRMA(t *testing.T) {
	r := NewRMQ[int]([]int{6, 1, 0, 10, 9}, func(v1, v2 int) bool { return v1 < v2 })
	assert.Equal(t, 6, r.RMQ(0, 0))
	assert.Equal(t, 1, r.RMQ(0, 1))
	assert.Equal(t, 0, r.RMQ(0, 2))
	assert.Equal(t, 0, r.RMQ(0, 3))
	assert.Equal(t, 0, r.RMQ(0, 4))

	assert.Equal(t, 9, r.RMQ(3, 4))
}
