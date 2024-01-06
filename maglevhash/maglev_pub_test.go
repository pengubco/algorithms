package maglevhash_test

import (
	"fmt"
	"hash/crc32"
	"math"
	"strconv"
	"testing"

	"github.com/pengubco/ads/maglevhash"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	// slot count must be a prime
	_, err := maglevhash.NewMaglevWithTableSize(10, []string{"B0", "B1"}, crc32.ChecksumIEEE)
	assert.Error(t, err)

	// slot count must be less than node count. Ideally, M >> N.
	_, err = maglevhash.NewMaglevWithTableSize(5, []string{"B0", "B1", "B2", "B3", "B4", "B5"}, crc32.ChecksumIEEE)
	assert.Error(t, err)
}

func TestLoadBalanceAndDisruption(t *testing.T) {
	nodeCnt := 10
	slotCnt := 10_007
	keyspaceCnt := 1_000_000

	nodes := lo.Times(nodeCnt, func(i int) string {
		return fmt.Sprintf("B%d", i)
	})
	mh, err := maglevhash.NewMaglevWithTableSize(slotCnt, nodes, crc32.ChecksumIEEE)
	assert.NoError(t, err)

	// key -> node
	keyToNode1 := make(map[int]string)
	// node -> number of keys
	nodeLoad1 := make(map[string]int)

	for i := 0; i < keyspaceCnt; i++ {
		node := mh.Node([]byte(strconv.FormatInt(int64(i), 10)))
		keyToNode1[i] = node
		nodeLoad1[node]++
	}
	// the min load and max load should be less than N/M of the total key spaces.
	err = verifyEvenLoadBalance(nodeLoad1, int(float64(keyspaceCnt)*1.5*float64(nodeCnt)/float64(slotCnt)))
	assert.NoError(t, err)

	// remove node[5].
	nodeCnt = 9
	nodes = lo.Filter(lo.Times(nodeCnt, func(i int) string {
		return fmt.Sprintf("B%d", i)
	}), func(item string, index int) bool {
		return index != 5
	})
	mh, err = maglevhash.NewMaglevWithTableSize(slotCnt, nodes, crc32.ChecksumIEEE)
	assert.NoError(t, err)

	keyToNode2 := make(map[int]string)
	// node -> number of keys
	nodeLoad2 := make(map[string]int)
	for i := 0; i < keyspaceCnt; i++ {
		node := mh.Node([]byte(strconv.FormatInt(int64(i), 10)))
		keyToNode2[i] = node
		nodeLoad2[node]++
	}
	err = verifyEvenLoadBalance(nodeLoad2, int(float64(keyspaceCnt)*1.5*float64(nodeCnt)/float64(slotCnt)))
	assert.NoError(t, err)
	err = verifyDisruption(keyspaceCnt, keyToNode1, keyToNode2, 2*keyspaceCnt/nodeCnt)
	assert.NoError(t, err)

	// add node[10]
	nodeCnt = 11
	nodes = lo.Times(nodeCnt, func(i int) string {
		return fmt.Sprintf("B%d", i)
	})
	mh, err = maglevhash.NewMaglevWithTableSize(slotCnt, nodes, crc32.ChecksumIEEE)
	assert.NoError(t, err)

	keyToNode3 := make(map[int]string)
	// node -> number of keys
	nodeLoad3 := make(map[string]int)
	for i := 0; i < keyspaceCnt; i++ {
		node := mh.Node([]byte(strconv.FormatInt(int64(i), 10)))
		keyToNode3[i] = node
		nodeLoad3[node]++
	}
	err = verifyEvenLoadBalance(nodeLoad3, int(float64(keyspaceCnt)*math.Ceil(float64(nodeCnt)/float64(slotCnt))))
	assert.NoError(t, err)
	err = verifyDisruption(keyspaceCnt, keyToNode1, keyToNode3, 2*keyspaceCnt/(nodeCnt-1))
	assert.NoError(t, err)
}

func verifyDisruption(keyspaceCnt int, keyToNode1, keyToNode2 map[int]string, threshold int) error {
	// number of keys that have been moved.
	moveCnt := 0
	for i := 0; i < keyspaceCnt; i++ {
		if keyToNode1[i] != keyToNode2[i] {
			moveCnt++
		}
	}
	fmt.Printf("number of keys moved %d, threshold %d\n", moveCnt, threshold)
	if moveCnt > threshold {
		return fmt.Errorf("number of keys moved larger than threshold, %d > %d", moveCnt, threshold)
	}
	return nil
}

func verifyEvenLoadBalance(load map[string]int, threshold int) error {
	minLoad, maxLoad := math.MaxInt, math.MinInt
	for _, v := range load {
		minLoad = min(minLoad, v)
		maxLoad = max(maxLoad, v)
	}
	fmt.Printf("min load: %d max load: %d threshold: %d\n", minLoad, maxLoad, threshold)
	if maxLoad-minLoad > threshold {
		return fmt.Errorf("difference between maxLoad and minLoad larger threshold, %d - %d > %d", maxLoad, minLoad, threshold)
	}
	return nil
}
