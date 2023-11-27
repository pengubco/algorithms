package union_find_test

import (
	"testing"

	"github.com/pengubco/ads/union_find"
	"github.com/stretchr/testify/assert"
)

func TestUnionFind(t *testing.T) {
	assert := assert.New(t)
	uf := union_find.NewUnionFind(5)
	for i := 0; i < 5; i++ {
		assert.Equal(i, uf.Find(i))
	}

	uf.Union(1, 2)
	assert.Equal(uf.Find(1), uf.Find(2))
	assert.Equal(2, uf.Size(1))
	assert.Equal(2, uf.Size(2))

	uf.Union(3, 4)
	assert.Equal(uf.Find(3), uf.Find(4))
	assert.Equal(2, uf.Size(3))
	assert.Equal(2, uf.Size(4))

	uf.Union(0, 3)
	assert.Equal(uf.Find(0), uf.Find(3))
	assert.Equal(uf.Find(0), uf.Find(4))
	assert.Equal(3, uf.Size(0))
	assert.Equal(3, uf.Size(3))
	assert.Equal(3, uf.Size(4))

	uf.Union(4, 1)
	for i := 1; i < 5; i++ {
		assert.Equal(uf.Find(0), uf.Find(i))
		assert.Equal(5, uf.Size(i))
	}
}
