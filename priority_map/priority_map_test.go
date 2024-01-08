package priority_map_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"slices"

	"github.com/pengubco/algorithms/priority_map"
	"github.com/stretchr/testify/assert"
)

func TestPriorityMap_Single_Pair(t *testing.T) {
	assert := assert.New(t)
	pm := priority_map.NewPriorityMap[string, int](func(v1, v2 int) bool {
		return v1 < v2
	})
	pm.Set("a", 1)
	k, v, ok := pm.Top()
	assert.Equal(1, pm.Size())
	assert.True(ok)
	assert.Equal("a", k)
	assert.Equal(1, v)

	pm.Set("a", 10)

	v, ok = pm.Get("a")
	assert.True(ok)
	assert.Equal(10, v)

	pm.Delete("a")
	assert.Equal(0, pm.Size())
	v, ok = pm.Get("a")
	assert.False(ok)
	k, v, ok = pm.Pop()
	assert.False(ok)
	k, v, ok = pm.Top()
	assert.False(ok)
}

func TestPriorityMap_DuplicateValues(t *testing.T) {
	assert := assert.New(t)
	pm := priority_map.NewPriorityMap[string, int](func(v1, v2 int) bool {
		return v1 < v2
	})
	pm.Set("a", 20)
	pm.Set("b", 10)
	pm.Set("c", 20)
	k, v, _ := pm.Pop()
	assert.Equal("b", k)
	assert.Equal(10, v)
	k2, v2, _ := pm.Pop()
	k3, v3, _ := pm.Pop()
	assert.Equal(20, v2)
	assert.Equal(20, v3)

	keys := []string{k2, k3}
	slices.Sort(keys)
	assert.Equal([]string{"a", "c"}, keys)
}

func TestPriorityMap_Simple_Pairs(t *testing.T) {
	assert := assert.New(t)

	pm := priority_map.NewPriorityMap[string, int](func(v1, v2 int) bool {
		return v2 < v1
	})
	for i := 0; i < 26; i++ {
		pm.Set(fmt.Sprintf("%c", 'a'+i), i)
	}
	assert.Equal(26, pm.Size())
	k, v, ok := pm.Top()
	assert.Equal("z", k)
	assert.Equal(25, v)
	assert.True(ok)
	pm.Set("a", 100)
	k, v, _ = pm.Top()
	assert.Equal("a", k)
	assert.Equal(100, v)

	pm.Set("a", 0)
	pm.Delete("b")
	pm.Delete("h")
	for i := 0; i < 26; i++ {
		s := fmt.Sprintf("%c", 'a'+25-i)
		if s == "b" || s == "h" {
			continue
		}
		k, v, ok = pm.Pop()
		assert.True(ok)
		assert.Equal(s, k)
	}
}

func TestPriorityMap_CompositeV(t *testing.T) {
	assert := assert.New(t)
	type Job struct {
		expire time.Time
	}
	pm := priority_map.NewPriorityMap[int, *Job](func(v1, v2 *Job) bool {
		return v1.expire.Before(v2.expire)
	})

	for i := 0; i < 100; i++ {
		pm.Set(i, &Job{
			expire: time.Now().Add(time.Duration(i) * time.Minute),
		})
	}
	var prevJobExpire time.Time
	var i int
	for pm.Size() > 0 {
		k, v, ok := pm.Pop()
		assert.True(ok)
		assert.Equal(i, k)
		assert.True(v.expire.After(prevJobExpire))
		prevJobExpire = v.expire
		i++
	}
}

// Add 1M key-value pairs in random values.
func BenchmarkPriorityMap_Add_1M(b *testing.B) {
	n := 1_000_000
	indexes := shuffledIndexes(n)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		pm := priority_map.NewPriorityMap[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})

		b.StartTimer()
		for j := 0; j < n; j++ {
			pm.Set(j, indexes[j])
		}
	}
}

// Updates 1M key-value pairs in random values.
func BenchmarkPriorityMap_Update_1M(b *testing.B) {
	n := 1_000_000
	indexes := shuffledIndexes(n)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		pm := priority_map.NewPriorityMap[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})
		for j := 0; j < n; j++ {
			pm.Set(j, indexes[j])
		}
		rand.Shuffle(n, func(i, j int) {
			indexes[i], indexes[j] = indexes[j], indexes[i]
		})

		b.StartTimer()
		for j := 0; j < n; j++ {
			pm.Set(j, indexes[j])
		}
	}
}

// Delete 1M key-value pairs in random order.
func BenchmarkPriorityMap_Del_1M(b *testing.B) {
	n := 1_000_000
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		pm := priority_map.NewPriorityMap[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})
		for j := 0; j < n; j++ {
			pm.Set(j, r.Int())
		}

		b.StartTimer()
		for j := 0; j < n; j++ {
			pm.Delete(j)
		}
	}
}

// Pop 1M key-value pairs in random order
func BenchmarkPriorityMap_Pop_1M(b *testing.B) {
	n := 1_000_000
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		pm := priority_map.NewPriorityMap[int, int](func(v1, v2 int) bool {
			return v1 < v2
		})
		for j := 0; j < n; j++ {
			pm.Set(j, r.Int())
		}

		b.StartTimer()
		for j := 0; j < n; j++ {
			pm.Pop()
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
