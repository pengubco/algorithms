package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/pengubco/algorithms/maglev_hash"
	"github.com/samber/lo"
)

// This program writes out three lines.
// first line: number of slots assigned to each node
// second line: number of slots that need to be redistributed after removing a node.
// third line: slotCnt/nodeCnt, this is the theoretical minimum number of
// slots need to be moved.
func main() {
	var nodeCnt int
	var slotCnt int
	flag.IntVar(&nodeCnt, "nodeCnt", 0, "number of nodes")
	flag.IntVar(&slotCnt, "slotCnt", 0, "number of slots. must be a prime number and larger than nodeCnt")
	flag.Parse()

	e, err := NewExperiment(nodeCnt, slotCnt)
	if err != nil {
		log.Fatal(err)
	}
	slotToNode := e.getSlotAssignment()
	nodeLoads := lo.Values(calculateNodeLoads(slotToNode))
	b, _ := json.Marshal(nodeLoads)
	fmt.Printf("%s\n", string(b))

	e.removeNode("1")
	newSlotToNode := e.getSlotAssignment()
	slotMoved := calculateSlotMove(slotToNode, newSlotToNode)
	fmt.Printf("%d\n", slotMoved)
	fmt.Printf("%d\n", slotCnt/nodeCnt)
}

func calculateNodeLoads(slotToNode map[int]int) map[int]int {
	load := make(map[int]int)
	for _, v := range slotToNode {
		load[v]++
	}
	return load
}

// calculateSlotMove returns number of slots that are redistributed.
func calculateSlotMove(a, b map[int]int) int {
	result := 0
	for k, _ := range a {
		if a[k] != b[k] {
			result++
		}
	}
	return result
}

func sanityCheck() {
	e, err := NewExperiment(3, 7)
	if err != nil {
		log.Fatal(err)
	}
	slotToNode := e.getSlotAssignment()
	fmt.Printf("%+v\n", slotToNode)

	e.removeNode("1")
	slotToNode2 := e.getSlotAssignment()
	fmt.Printf("%+v\n", slotToNode2)
}

// Experiment has the following settings.
//  1. Nodes: Nodes are number 0, 1, 2, ...
//  2. Use Atoi as the key-hash function which returns the integer form a string.
//     So that we can use Node([]byte("1")) to get the node assigned for slot-1.
type Experiment struct {
	nodeCnt int
	slotCnt int
	nodes   []string

	mh *maglev_hash.MaglevHash
}

func NewExperiment(nodeCnt int, slotCnt int) (*Experiment, error) {
	nodes := lo.Times(nodeCnt, func(i int) string {
		return fmt.Sprintf("%d", i)
	})
	mh, err := maglev_hash.NewMaglevWithTableSize(slotCnt, nodes, func(b []byte) uint32 {
		i, _ := strconv.Atoi(string(b))
		return uint32(i)
	})
	if err != nil {
		return nil, err
	}
	return &Experiment{
		nodeCnt: nodeCnt,
		slotCnt: slotCnt,
		nodes:   nodes,
		mh:      mh,
	}, nil
}

// returns slot -> node mapping.
func (e *Experiment) getSlotAssignment() map[int]int {
	result := make(map[int]int)
	for i := 0; i < e.slotCnt; i++ {
		nodeStr := e.mh.Node([]byte(strconv.Itoa(i)))
		nodeID, _ := strconv.Atoi(nodeStr)
		result[i] = nodeID
	}
	return result
}

func (e *Experiment) removeNode(node string) error {
	found := false
	newNodes := make([]string, 0)
	for _, v := range e.nodes {
		if v != node {
			newNodes = append(newNodes, v)
		} else {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("node does not exist, %s", node)
	}
	e.nodeCnt = len(newNodes)
	e.nodes = newNodes

	mh, err := maglev_hash.NewMaglevWithTableSize(e.slotCnt, e.nodes, func(b []byte) uint32 {
		i, _ := strconv.Atoi(string(b))
		return uint32(i)
	})
	if err != nil {
		return err
	}
	e.mh = mh
	return nil
}
