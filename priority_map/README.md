# Priority Map 

PriorityMap is a combination of HashMap and Binary Heap. It is useful for workloads that requires
1. A Key-Value store that supports Get, Set, Delete by key. 
2. A Heap that supports access the key-value pair of the smallest value.

Some good use cases:
1. Implementing greedy algorithms like the [Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm) 
and [Prim's algorithm](https://en.wikipedia.org/wiki/Prim%27s_algorithm). 
2. Job Scheduler that schedules job with the highest priority.

Some similar data structures are: 
1. [TreeMap](https://docs.oracle.com/javase/8/docs/api/java/util/TreeMap.html): key-value
  pairs are sorted by key. 
2. [SortedSet](https://redis.io/docs/data-types/sorted-sets/): key-values are sorted by 
  value. 
3. [PriorityQueue](https://docs.oracle.com/javase/8/docs/api/java/util/PriorityQueue.html): pop
  the smallest element out of queue.

Both TreeMap and SortedSet sort **all** key-value pairs. TreeMap keeps order using a balanced 
binary tree, and SortedSet uses a skip list. PriorityMap, however, does not sort all key-value pairs.
Instead, PriorityMap keeps order using a binary heap. Therefore, only the top of Heap is guaranteed
an order, the minimum, among all pairs. 

PriorityQueue does not store key-value pairs, so it does not supports access by key. It does not 
support update an element's priority either.


## How to use it

### Simple Value Type
```go
func main() {
	pm := prioritymap.NewPriorityMap[int, int](func(a, b int) int {
		return a - b
	})
	pm.Set(1, 10)
	pm.Set(2, 10)
	pm.Set(3, 30)
	fmt.Printf("size: %d\n", pm.Size()) // "size: 3"
	if v, ok := pm.Get(1); ok {
		fmt.Printf("key: 1, value: %d\n", v) // "key: 1, value: 10"
	}
	if k, v, ok := pm.Top(); ok {
		fmt.Printf("key: %d, value: %d\n", k, v) // "key: 1, value: 10" or "key: 2, value: 10"
	}
}
```

### Composite Value Type
```go
type Job struct {
	id         int
	expiration time.Time
	name       string
}

func main() {
	pm := prioritymap.NewPriorityMap[int, *Job](func(v1, v2 *Job) bool {
		switch {
		case v1.expire.Before(v2.expire):
			return -1
		case v1.expire.After(v2.expire):
			return 1
		default:
			return 0
		}
	})
	now := time.Now()
	jobs := []Job{
		{1, now.Add(-1 * time.Minute), "job 1"},
		{2, now.Add(-2 * time.Minute), "job 2"},
		{3, now.Add(-3 * time.Minute), "job 3"},
	}
	for i, _ := range jobs {
		pm.Set(jobs[i].id, &jobs[i])
	}
	id, job, _ := pm.Top()
	fmt.Printf("job with the smallest expiration. id %d, name %s\n", id, job.name)
	job, _ = pm.Get(2)
	fmt.Printf("job 2's expiration is %v\n", job.expiration)
	fmt.Println("taking jobs one by one, in the order of expiration time")
	for pm.Size() > 0 {
		if id, job, ok := pm.Pop(); ok {
			fmt.Printf("job id %d, name %s, expiration %v\n", id, job.name, job.expiration)
		}
	}
}
```

## Performance
We benchmark Add, Update, Delete, Pop on a priority map of 1M key-value pairs. The following is a result on my Mac M1 Max. You can run the benchmark with.
```
go test -bench BenchmarkPriorityMap -benchmem  -benchtime 10s
``` 

```text
goos: darwin
goarch: arm64
pkg: github.com/pengubco/algorithms/prioritymap
BenchmarkPriorityMap_Add_1M-10                    66         175442495 ns/op        149403499 B/op       1023415 allocs/op
BenchmarkPriorityMap_Update_1M-10                100         124093262 ns/op           80035 B/op              0 allocs/op
BenchmarkPriorityMap_Del_1M-10                    75         170370006 ns/op               0 B/op              0 allocs/op
BenchmarkPriorityMap_Pop_1M-10                    19         604293805 ns/op               0 B/op              0 allocs/op
PASS
ok      github.com/pengubco/algorithms/prioritymap 94.674s
```

No surprise that `Pop` is most expensive because the heap may need to go from root to a leaf 
to maintain the heap structure. It takes 604ms (~0.6 second) to pop 1M key-value pairs. I think 
this is fast enough for normal production use.

## Correctness 
We run 1M operations on PriorityMap and a SortedSet in Redis, see [redis-compare/main.go](./example/redis-compare/main.go). After each operation, we compare the size and the smallest key-value pair from PriorityMap with corresponding values 
from Redis. This gives us confidence that PriorityMap is correct.

If you found a bug, please open an issue. 
