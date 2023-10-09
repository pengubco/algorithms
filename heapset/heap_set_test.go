package heapset_test

import (
	"fmt"
	"math/rand"
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
	hs.Delete("b")
	hs.Delete("h")
	for i := 0; i < 26; i++ {
		s := fmt.Sprintf("%c", 'a'+25-i)
		if s == "b" || s == "h" {
			continue
		}
		k, v, ok = hs.Pop()
		assert.True(ok)
		assert.Equal(s, k)
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

// Add 1M key-value pairs in random values.
func BenchmarkHeapSet_Add_1M(b *testing.B) {
	n := 1_000_000
	indexes := shuffledIndexes(n)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		hs := heapset.NewHeapSet[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})

		b.StartTimer()
		for j := 0; j < n; j++ {
			hs.Set(j, indexes[j])
		}
	}
}

// Updates 1M key-value pairs in random values.
func BenchmarkHeapSet_Update_1M(b *testing.B) {
	n := 1_000_000
	indexes := shuffledIndexes(n)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		hs := heapset.NewHeapSet[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})
		for j := 0; j < n; j++ {
			hs.Set(j, indexes[j])
		}
		rand.Shuffle(n, func(i, j int) {
			indexes[i], indexes[j] = indexes[j], indexes[i]
		})

		b.StartTimer()
		for j := 0; j < n; j++ {
			hs.Set(j, indexes[j])
		}
	}
}

// Delete 1M key-value pairs in random order.
func BenchmarkHeapSet_Del_1M(b *testing.B) {
	n := 1_000_000
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		hs := heapset.NewHeapSet[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})
		for j := 0; j < n; j++ {
			hs.Set(j, r.Int())
		}

		b.StartTimer()
		for j := 0; j < n; j++ {
			hs.Delete(j)
		}
	}
}

// Pop 1M key-value pairs in random order
func BenchmarkHeapSet_Pop_1M(b *testing.B) {
	n := 1_000_000
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		hs := heapset.NewHeapSet[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})
		for j := 0; j < n; j++ {
			hs.Set(j, r.Int())
		}

		b.StartTimer()
		for j := 0; j < n; j++ {
			hs.Pop()
		}
	}
}

// returns an array [0, n) in random order
func shuffledIndexes(n int) []int {
	indexes := make([]int, n)
	for i := 0; i < n; i++ {
		indexes[i] = i
	}
	rand.Shuffle(n, func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})
	return indexes
}
