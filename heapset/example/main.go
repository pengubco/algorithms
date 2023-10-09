package main

import (
	"fmt"
	"time"

	"github.com/pengubco/ads/heapset"
)

func main() {
	hs := heapset.NewHeapSet[int, *Job](func(v1, v2 *Job) bool {
		return v1.expiration.Before(v2.expiration)
	})
	now := time.Now()
	jobs := []Job{
		{1, now.Add(-1 * time.Minute), "job 1"},
		{2, now.Add(-2 * time.Minute), "job 2"},
		{3, now.Add(-3 * time.Minute), "job 3"},
	}
	for i, _ := range jobs {
		hs.Set(jobs[i].id, &jobs[i])
	}
	id, job, _ := hs.Top()
	fmt.Printf("job with the smallest expiration. id %d, name %s\n", id, job.name)
	job, _ = hs.Get(2)
	fmt.Printf("job 2's expiration is %v\n", job.expiration)
	fmt.Println("taking jobs one by one, in the order of expiration time")
	for hs.Size() > 0 {
		if id, job, ok := hs.Pop(); ok {
			fmt.Printf("job id %d, name %s, expiration %v\n", id, job.name, job.expiration)
		}
	}
}

type Job struct {
	id         int
	expiration time.Time
	name       string
}
