package main

import (
	"fmt"
	"time"

	"github.com/pengubco/ads/prioritymap"
)

func main() {
	fmt.Println("example 1: simple key value types")
	simpleKV()

	fmt.Println()
	fmt.Println("example 2: advanced key value types: job scheduler")
	jobScheduler()
}

func simpleKV() {
	hs := prioritymap.NewPriorityMap[int, int](func(a, b int) int {
		return a - b
	})
	hs.Set(1, 10)
	hs.Set(2, 10)
	hs.Set(3, 30)
	fmt.Printf("size: %d\n", hs.Size()) // "size: 3"
	if v, ok := hs.Get(1); ok {
		fmt.Printf("key: 1, value: %d\n", v) // "key: 1, value: 10"
	}
	if k, v, ok := hs.Top(); ok {
		fmt.Printf("key: %d, value: %d\n", k, v) // "key: 1, value: 10" or "key: 2, value: 10"
	}
}

type Job struct {
	id         int
	expiration time.Time
	name       string
}

func jobScheduler() {
	hs := prioritymap.NewPriorityMap[int, *Job](func(v1, v2 *Job) int {
		switch {
		case v1.expiration.Before(v2.expiration):
			return -1
		case v1.expiration.After(v2.expiration):
			return 1
		default:
			return 0
		}
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
