package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"gonum.org/v1/gonum/stat"
)

const NUM_REPLICAS int = 100
const NUM_KEYS_TO_ADD int = 10_000_000
const MAX_VALUE int = 10_000_000
const NUM_NODES int = 10

// func stdDev(m map[string]int) {
// 	var mean float32
// 	for k, v := range m {
// 		mean +=
// 	}

// }

func TestConsistentHashRing(t *testing.T) {
	hashRing := NewConsistentHashRing(NUM_REPLICAS, NUM_NODES)
	for i := 0; i < NUM_NODES; i++ {
		hashRing.addNode(fmt.Sprintf("name%d", i), fmt.Sprintf("value%d", i))
	}
	nodesToKeysCounter := make(map[string]int)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < NUM_KEYS_TO_ADD; i++ {
		nodeName, _ := hashRing.getNode(fmt.Sprintf("%d", rand.Intn(MAX_VALUE)))
		nodesToKeysCounter[nodeName]++
	}
	// slice where s[i] = n,
	// where n = number of nodes with i keys
	// in other words:
	// this many keys seen on how many nodes
	var nodesCounterFloatValues []float64
	numBuckets := NUM_KEYS_TO_ADD / (NUM_NODES * NUM_REPLICAS)
	sortedNodesDistribution := make([]int, NUM_KEYS_TO_ADD)
	for _, count := range nodesToKeysCounter {
		nodesCounterFloatValues = append(nodesCounterFloatValues, float64(count))
		sortedNodesDistribution[count/numBuckets]++
	}
	for numHits, numberOfNodesWithThisManyHits := range sortedNodesDistribution {
		if numberOfNodesWithThisManyHits != 0 {
			fmt.Printf("%d: %s\n", numHits, strings.Repeat("#", numberOfNodesWithThisManyHits))
		}
	}
	sort.Float64s(nodesCounterFloatValues)
	t.Log(nodesCounterFloatValues)
	stdDev := stat.PopStdDev(nodesCounterFloatValues, nil)
	if stdDev > 0.2 {
		t.Errorf("Nodes deviation: %f, want: <= 0.2", stdDev)
	}
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if got := NewConsistentHashRing(tt.args.replicas); !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("NewConsistentHashRing() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }

}
