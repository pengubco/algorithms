// Package maglevhash implements Maglev Hash algorithm from the paper
// [Maglev: A Fast and Reliable Software Network Load Balancer](https://static.googleusercontent.com/media/research.google.com/en//pubs/archive/44824.pdf)
// Maglev hash break down the key space to slots and assign slots to nodes.
// Each node is a logic entity that can handle the key. It could be a server,
// a service, a VIP, etc.
// Maglev hash guarantees the following two properties.
//  1. Load Balance. Let N be the number of nodes and M be the number of slots.
//     The difference between any node's assigned slots is less than N/M.
//     For example, if M = 100*N, then the node with most slots has no more than 1% slots
//     than the node with the least slots.
//  2. Minimal disruption. On average, add a new node causes reassigns M/N slots.
//
// Example
// mh, _ := NewMaglev([]string{"B0", "B1"})
// node := mh.Node([]byte("key1"))
package maglev_hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"math"
	"math/big"
	"sort"

	"github.com/samber/lo"
)

const (
	// The default number of slots is the smallest prime number larger than 10,000.
	// It should be OK for less than 100 nodes.
	DefaultSlotCnt = 10007
)

// KeyHashFnType
type KeyHashFnType func([]byte) uint32

// MaglevHash assigns fixed number of slots to a collection of nodes.
// 1. Nodes are identified by string.
// 2. The default hash function hashing key to slot is CRC32.
type MaglevHash struct {
	// Number of slots.
	slotCnt int

	// Number of nodes.
	nodeCnt int

	nodes []string

	// The mapping from slot to node. lookup[i]=j: slot[i] is mapped to node[j]
	lookup []int

	keyHashFn KeyHashFnType
}

// NewMaglev creates a Maglev hash for the given nodes, using default number of
// slots and the CRC32 key hash function.
func NewMaglev(nodes []string) (*MaglevHash, error) {
	return NewMaglevWithTableSize(DefaultSlotCnt, nodes, crc32.ChecksumIEEE)
}

// NewMaglevWithTableSize creates a Maglev hash. Nodes are identified by strings.
// In order to make lookup table stable, the given list of nodes are sorted and
// deduplicated.
func NewMaglevWithTableSize(slotCnt int, nodes []string, keyHashFn KeyHashFnType) (*MaglevHash, error) {
	if !isPrime(slotCnt) {
		return nil, fmt.Errorf("number of slots must be a prime number, %d", slotCnt)
	}
	nodes = lo.Uniq(nodes)
	sort.Strings(nodes)
	if len(nodes) == 0 || len(nodes) > slotCnt {
		return nil, fmt.Errorf("more nodes than slots, %d > %d", len(nodes), slotCnt)
	}
	m := &MaglevHash{
		slotCnt:   slotCnt,
		nodeCnt:   len(nodes),
		nodes:     nodes,
		keyHashFn: keyHashFn,
	}
	m.lookup = m.buildLookup(m.buildPreferences())
	return m, nil
}

// Node returns the assigned node for the given key.
func (m *MaglevHash) Node(key []byte) string {
	return m.nodes[m.lookup[m.keyHashFn(key)%uint32(m.slotCnt)]]
}

// buildLookup calculates the lookup table for slot.
func (m *MaglevHash) buildLookup(preferenceList [][]int) []int {
	lookup := make([]int, m.slotCnt)
	for i := 0; i < m.slotCnt; i++ {
		lookup[i] = -1
	}
	// next[i] indicate the current favorite slot for the i-th node.
	next := make([]int, m.nodeCnt)

	// number of slots that have been assigned to nodes.
	assignedSlotCnt := 0

	for {
		for i := 0; i < m.nodeCnt; i++ {
			c := preferenceList[i][next[i]]
			for lookup[c] >= 0 {
				next[i]++
				c = preferenceList[i][next[i]]
			}
			lookup[c] = i
			next[i]++
			assignedSlotCnt++
			if assignedSlotCnt == m.slotCnt {
				return lookup
			}
		}
	}
}

// buildPreferences builds the preference list of slots for each node.
// Preference list of node. nodePreferences[i][j]=k means that, for node[i],
// the j-th preference is slot[k].
// nodePreferences[i] is a permutation of [0, slotCnt)
func (m *MaglevHash) buildPreferences() [][]int {
	n := len(m.nodes)
	nodePreferences := make([][]int, n)
	for i, node := range m.nodes {
		nodePreferences[i] = make([]int, m.slotCnt)
		offset := md5StringToModulo(fmt.Sprintf("%s:offset", node), m.slotCnt)
		skip := md5StringToModulo(fmt.Sprintf("%s:skip", node), m.slotCnt-1) + 1
		for j := 0; j < m.slotCnt; j++ {
			nodePreferences[i][j] = (offset + j*skip) % m.slotCnt
		}
	}
	return nodePreferences
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	upperBound := math.Sqrt(float64(n)) + 1
	for i := 5; i <= int(upperBound); i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// md5StringToModulo calculate the MD5 of the given string and module the result
// by n.
func md5StringToModulo(s string, n int) int {
	hash := md5.Sum([]byte(s))
	hashString := hex.EncodeToString(hash[:])
	hashInt, _ := big.NewInt(0).SetString(hashString, 16)
	moduloNumber := big.NewInt(int64(n))

	return int(big.NewInt(0).Mod(hashInt, moduloNumber).Int64())
}
