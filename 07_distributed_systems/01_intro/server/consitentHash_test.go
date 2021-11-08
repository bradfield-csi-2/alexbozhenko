package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"gonum.org/v1/gonum/stat"
)

const NUM_REPLICAS int = 100
const NUM_HITS int = 1000
const MAX_VALUE int = 10_000
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
		hashRing.addNode(fmt.Sprintf("%d", i), fmt.Sprintf("%d", i))
	}
	nodesCounter := make(map[string]int)
	for i := 0; i < NUM_HITS; i++ {
		nodeName, _ := hashRing.getNode(fmt.Sprintf("%d", rand.Intn(MAX_VALUE)))
		nodesCounter[nodeName]++
	}
	var nodesCounterFloatValues []float64
	for _, count := range nodesCounter {
		nodesCounterFloatValues = append(nodesCounterFloatValues, float64(count))
		fmt.Printf("%s\n", strings.Repeat("#", count))
	}
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
