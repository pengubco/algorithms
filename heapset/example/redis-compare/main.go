package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"

	"github.com/pengubco/ads/heapset"
	"github.com/redis/go-redis/v9"
)

func main() {
	redisURL := flag.String("redis", "localhost:16379", "redis host:port")
	flag.Parse()
	rdb := redis.NewClient(&redis.Options{
		Addr: *redisURL,
	})

	// Redis sorted set uses string as key and int as value.
	hs := heapset.NewHeapSet[string, float64](func(v1, v2 float64) int {
		return int(v1 - v2)
	})

	err := compareHeapSeatWithRedis(rdb, "ss", hs, 1_000_000)
	if err != nil {
		log.Fatal(err)
	}
}

// Carry out n operations on HeapSet and Redis. After each operation, get the Top() from HeapSet
// and compare it with the minimum value in Redis SortedSet.
func compareHeapSeatWithRedis(rdb *redis.Client, sortedSetName string, hs *heapset.HeapSet[string, float64], n int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rdb.Del(ctx, sortedSetName)

	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("%d", i)
	}
	values := make([]float64, n)
	for i := 0; i < n; i++ {
		values[i] = float64(i)
	}

	rand.Shuffle(n, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	rand.Shuffle(n, func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})

	verify := func() error {
		// check size
		size, err := rdb.ZCard(ctx, sortedSetName).Result()
		if err != nil {
			return err
		}
		if int64(hs.Size()) != size {
			return fmt.Errorf("size is different. Redis: %d, HeapSet: %d", size, hs.Size())
		}

		// check top
		results, err := rdb.ZRangeWithScores(ctx, sortedSetName, 0, 0).Result()
		if err != nil || len(keys) == 0 {
			return err
		}
		key, value, _ := hs.Top()
		redisKey := results[0].Member.(string)
		redisValue := results[0].Score
		if redisKey != key || redisValue != value {
			return fmt.Errorf("top different. redis: key %s, value: %f; heapset: key %s, value %f",
				redisKey, redisValue, key, value)
		}
		return nil
	}

	add := func(k string, v float64) {
		rdb.ZAdd(ctx, sortedSetName, redis.Z{Score: v, Member: k})
		hs.Set(k, v)
	}

	for i := 0; i < n; i++ {
		if (i+1)%1000 == 0 {
			fmt.Printf("progress: %d/%d\n", i+1, n)
		}
		add(keys[i], values[i])
		if err := verify(); err != nil {
			return err
		}
	}
	fmt.Println("complete")
	return nil
}
