package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const NUM_REPLICAS int = 10
const NUM_KEYS_TO_ADD int = 1_000_000
const MAX_VALUE int = 100_000_000
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
		hashRing.addNode(fmt.Sprintf("name%03d", i), fmt.Sprintf("value%d", i))
	}
	nodesToKeysCounter := make(map[string]float64)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < NUM_KEYS_TO_ADD; i++ {
		nodeName, _ := hashRing.getNode(fmt.Sprintf("%d", rand.Intn(MAX_VALUE)))
		nodesToKeysCounter[nodeName]++
	}
	t.Log(fmt.Sprint(nodesToKeysCounter))
	// stdDev := stat.PopStdDev(nodesCounterFloatValues, nil)
	var mean float64
	for _, value := range nodesToKeysCounter {
		mean += value
	}
	var observedNodesNumber int = len(nodesToKeysCounter)
	assert.Equal(t, NUM_NODES, observedNodesNumber,
		"Observed number of nodes must be"+
			"equal to the expected number of nodes")
	mean /= float64(observedNodesNumber)
	var standardDeviation float64
	for _, value := range nodesToKeysCounter {
		standardDeviation += math.Pow(value-mean, 2)
	}
	standardDeviation = math.Pow(standardDeviation/float64(observedNodesNumber), 0.5)
	t.Log("mean = ", mean)
	t.Log("expected = ", NUM_KEYS_TO_ADD/NUM_NODES)
	if standardDeviation > 0.2 {
		t.Errorf("Nodes deviation: %f, want: <= 0.2", standardDeviation)
	}
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		if got := NewConsistentHashRing(tt.args.replicas); !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("NewConsistentHashRing() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }

}
