package main

import (
	"fmt"
	"time"

	"github.com/pengubco/ads/heapset"
)

func main() {
	fmt.Println("example using simple value type")
	withSimpleValue()

	fmt.Println()
	fmt.Println("example using composite value type")
	withCompositeValue()
}

func withSimpleValue() {
	hs := heapset.NewHeapSet[string, int](func(v1, v2 int) bool {
		return v1 < v2
	})

	hs.Set("apple", 1)
	hs.Set("banana", 2)
	hs.Set("cherry", 3)
	k, _, _ := hs.Top()
	fmt.Printf("my favorite fruit is %s\n", k) // apple
	hs.Pop()
	k, _, _ = hs.Top()
	fmt.Printf("my 2nd favorite fruit is %s\n", k) // banana
	hs.Pop()
	k, _, _ = hs.Top()
	fmt.Printf("my 3rd favorite fruit is %s\n", k) // cherry
	hs.Pop()
}

func withCompositeValue() {
	type Job struct {
		id         int
		expiration time.Time
		name       string
	}

	hs := heapset.NewHeapSet[int, *Job](func(v1, v2 *Job) bool {
		return v1.expiration.Before(v2.expiration)
	})
	now, err := time.Parse("2006-01-02", "2022-12-30")
	if err != nil {
		fmt.Println(err)
		return
	}
	jobs := []Job{
		{1, now, "job 1"},
		{2, now.Add(-1 * time.Minute), "job 2"},
	}
	for i := range jobs {
		hs.Set(jobs[i].id, &jobs[i])
	}
	id, job, _ := hs.Top()
	fmt.Printf("job with the earliest expiration. id: %d, name: %s\n", id, job.name) // id 2, name: job 2
	job, _ = hs.Get(2)
	fmt.Printf("job 2's expiration is %v\n", job.expiration) // 2022-12-29 23:59:00 +0000 UTC
}
