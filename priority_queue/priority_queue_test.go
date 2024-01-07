package priorityqueue_test

import (
	"testing"
	"time"

	priorityqueue "github.com/pengubco/ads/priority_queue"
	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue_SimpleType(t *testing.T) {
	pq, err := priorityqueue.NewPriorityQueue[int](func(v1, v2 int) bool {
		return v1 < v2
	})
	assert.NoError(t, err)
	_, err = pq.Top()
	assert.Equal(t, priorityqueue.ErrQueueIsEmpty, err)
	_, err = pq.Pop()
	assert.Equal(t, priorityqueue.ErrQueueIsEmpty, err)

	pq.Push(10)
	v, err := pq.Top()
	assert.NoError(t, err)
	assert.Equal(t, 10, v)

	pq.Push(20)
	pq.Push(5)
	pq.Push(20)
	assert.Equal(t, 4, pq.Size())

	values := []int{5, 10, 20, 20}
	for _, v := range values {
		popped, err := pq.Pop()
		assert.NoError(t, err)
		assert.Equal(t, v, popped)
	}

	assert.Equal(t, 0, pq.Size())
}

func TestPriorityQueue_CustomizedType(t *testing.T) {
	type Job struct {
		name       string
		expiration time.Time
	}
	pq, err := priorityqueue.NewPriorityQueue[*Job](func(v1, v2 *Job) bool {
		return v1.expiration.Before(v2.expiration)
	})
	assert.NoError(t, err)

	baseTime := time.Now()
	pq.Push(&Job{name: "job1", expiration: baseTime})
	pq.Push(&Job{name: "job2", expiration: baseTime.Add(-1 * time.Minute)})
	pq.Push(&Job{name: "job3", expiration: baseTime.Add(-2 * time.Minute)})
	pq.Push(&Job{name: "job4", expiration: baseTime.Add(-2 * time.Minute)})

	job, err := pq.Top()
	assert.NoError(t, err)
	assert.True(t, job.name == "job3" || job.name == "job4")

	pq.Pop()
	pq.Pop()
	job, err = pq.Pop()
	assert.NoError(t, err)
	assert.Equal(t, "job2", job.name)
	job, err = pq.Pop()
	assert.NoError(t, err)
	assert.Equal(t, "job1", job.name)
}
