package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const NUM_REPLICAS int = 100
const NUM_KEYS_TO_ADD int = 1_000_000
const NUM_NODES int = 10
const ACCEPTABLE_VARIANCE = 0.2
const KEY_LENGTH = 16

func TestConsistentHashRing(t *testing.T) {
	hashRing := NewConsistentHashRing(NUM_REPLICAS, NUM_NODES)
	for i := 0; i < NUM_NODES; i++ {
		nodeName := fmt.Sprintf("name%03d", i)
		nodeValue := fmt.Sprintf("value%d", i)
		hashRing.addNode(nodeName, nodeValue)
	}
	nodesToKeysCounter := make(map[string]float64)
	t.Log(hashRing)
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, KEY_LENGTH)
	for i := 0; i < NUM_KEYS_TO_ADD; i++ {
		rand.Read(key)
		nodeName, _ := hashRing.getNode(string(key))

		nodesToKeysCounter[nodeName]++
	}
	t.Log(fmt.Sprint(nodesToKeysCounter))
	var mean float64
	var observedNodesNumber int = len(nodesToKeysCounter)
	for _, value := range nodesToKeysCounter {
		mean += value
	}
	mean /= float64(observedNodesNumber)
	assert.Equal(t, NUM_NODES, observedNodesNumber,
		"Observed number of nodes must be"+
			"equal to the expected number of nodes")
	var standardDeviation float64
	for _, value := range nodesToKeysCounter {
		standardDeviation += math.Pow(value-mean, 2)
	}
	standardDeviation = math.Pow(standardDeviation/float64(observedNodesNumber), 0.5)
	variationCoefficient := standardDeviation / mean
	assert.Less(t, variationCoefficient, ACCEPTABLE_VARIANCE,
		"Distribution of keys by nodes is far from uniform")
	assert.Equal(t, float64(NUM_KEYS_TO_ADD/NUM_NODES), mean,
		"mean value not equal to expected")

}
