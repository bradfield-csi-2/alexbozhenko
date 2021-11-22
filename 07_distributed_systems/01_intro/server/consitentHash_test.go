package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const NUM_REPLICAS int = 500
const NUM_KEYS_TO_ADD int = 100_000
const NUM_NODES int = 5
const ACCEPTABLE_VARIANCE = 0.2
const KEY_LENGTH = 16

func TestConsistentHashRing(t *testing.T) {
	t.Log("NUM_REPLICAS", NUM_REPLICAS,
		" NUM_NODES ", NUM_NODES,
		" NUM_KEYS_TO_ADD ", NUM_KEYS_TO_ADD)
	hashRing := NewConsistentHashRing(NUM_REPLICAS, NUM_NODES)
	for i := 0; i < NUM_NODES; i++ {
		nodeName := fmt.Sprintf("node%03d", i)
		nodeValue := fmt.Sprintf("value%d", i)
		hashRing.addNode(nodeName, nodeValue)
	}
	nodesToKeysCounter := make(map[string]float64)
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, KEY_LENGTH)

	tmpKeys := make([]replicaHashWithNodeName, NUM_KEYS_TO_ADD)
	for i := 0; i < NUM_KEYS_TO_ADD; i++ {
		rand.Read(key)
		nodeName, _ := hashRing.getNode(string(key))
		nodesToKeysCounter[nodeName]++
		tmpKeys[i] = replicaHashWithNodeName{
			hash:     md5.Sum(key),
			nodeName: nodeName,
		}
	}
	sort.Slice(tmpKeys, func(i, j int) bool {
		return (bytes.Compare(tmpKeys[i].hash[:], tmpKeys[j].hash[:]) == -1)
	})
	for range tmpKeys {
		//	t.Logf("Adding hash(key):%x to %s", rH.hash, rH.nodeName)
	}

	//t.Log(hashRing)
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
	t.Log("mean:", mean)
	t.Log("Variation coefficient:", variationCoefficient)
	assert.Equal(t, float64(NUM_KEYS_TO_ADD)/float64(NUM_NODES), mean,
		"mean value not equal to expected")

}
