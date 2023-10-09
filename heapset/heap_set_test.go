package heapset_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pengubco/ads/heapset"
	"github.com/stretchr/testify/assert"
)

func TestHeapSet_Single_Pair(t *testing.T) {
	assert := assert.New(t)
	hs := heapset.NewHeapSet[string, int](func(v1, v2 int) bool {
		return v1 < v2
	})
	hs.Set("a", 1)
	k, v, ok := hs.Top()
	assert.Equal(1, hs.Size())
	assert.True(ok)
	assert.Equal("a", k)
	assert.Equal(1, v)

	hs.Set("a", 10)

	v, ok = hs.Get("a")
	assert.True(ok)
	assert.Equal(10, v)

	hs.Delete("a")
	assert.Equal(0, hs.Size())
	v, ok = hs.Get("a")
	assert.False(ok)
	k, v, ok = hs.Pop()
	assert.False(ok)
	k, v, ok = hs.Top()
	assert.False(ok)
}

func TestHeapSet_Simple_Pairs(t *testing.T) {
	assert := assert.New(t)

	hs := heapset.NewHeapSet[string, int](func(v1, v2 int) bool {
		return v1 > v2
	})
	for i := 0; i < 26; i++ {
		hs.Set(fmt.Sprintf("%c", 'a'+i), i)
	}
	assert.Equal(26, hs.Size())
	k, v, ok := hs.Top()
	assert.Equal("z", k)
	assert.Equal(25, v)
	assert.True(ok)
	hs.Set("a", 100)
	k, v, _ = hs.Top()
	assert.Equal("a", k)
	assert.Equal(100, v)

	hs.Set("a", 0)
	for i := 0; i < 26; i++ {
		k, v, ok = hs.Pop()
		assert.True(ok)
		assert.Equal(fmt.Sprintf("%c", 'a'+25-i), k)
	}
}

func TestHeapSet_CompositeV(t *testing.T) {
	assert := assert.New(t)
	type Job struct {
		expire time.Time
	}
	hs := heapset.NewHeapSet[int, *Job](func(v1, v2 *Job) bool {
		return v1.expire.Before(v2.expire)
	})

	for i := 0; i < 100; i++ {
		hs.Set(i, &Job{
			expire: time.Now().Add(time.Duration(i) * time.Minute),
		})
	}
	var prevJobExpire time.Time
	var i int
	for hs.Size() > 0 {
		k, v, ok := hs.Pop()
		assert.True(ok)
		assert.Equal(i, k)
		assert.True(v.expire.After(prevJobExpire))
		prevJobExpire = v.expire
		i++
	}
}
