package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sort"
)

const PARTITIONS = 640_000

//const NODES_NUMBER = 3

// Terminology:
// Node - a physical server, responsible for subset of keys

type consistentHashRing struct {
	// how many times we create points on the hash ring for each node
	replicas int
	// Sorted array of hashes of all nodes
	replicasHashes [][md5.Size]byte
	// Map from node name to nodeValue, where
	// nodeValue can be anything you want, e.g. an url
	nodes map[string]string
}

func NewConsistentHashRing(replicas, nodeNumber int) *consistentHashRing {
	rHashes := make([][md5.Size]byte, 0)
	return &consistentHashRing{
		replicas:       replicas,
		replicasHashes: rHashes,
		nodes:          make(map[string]string),
	}
}

func (c *consistentHashRing) addNode(nodeName string, nodeValue string) {
	for i := 0; i < c.replicas; i++ {
		replicaHashArray := md5.Sum([]byte(fmt.Sprintf("%s:%d", nodeName, i)))
		replicaHash := replicaHashArray[:]
		index := sort.Search(len(c.replicasHashes),
			func(i int) bool {
				res := bytes.Compare(c.replicasHashes[i][:], replicaHash)
				return res >= 0
			})

		c.replicasHashes = append(append(c.replicasHashes[:index][:], replicaHashArray),
			c.replicasHashes[index:]...)

		c.nodes[string(replicaHash)] = nodeValue
	}
}

func (c *consistentHashRing) deleteNode(nodeName string) {
	for i := 0; i < c.replicas; i++ {
		replicaHashArray := md5.Sum([]byte(fmt.Sprintf("%s:%d", nodeName, i)))
		replicaHash := replicaHashArray[:]
		index := sort.Search(len(c.replicasHashes),
			func(i int) bool {
				return (bytes.Equal(c.replicasHashes[i][:], replicaHash))
			})
		if index == len(c.replicasHashes) {
			panic(fmt.Sprintf("BUG! %s is not present in the hash ring", nodeName))
		}
		c.replicasHashes = append(c.replicasHashes[:index][:], c.replicasHashes[index+1:]...)
	}
	delete(c.nodes, nodeName)
}

// Given a key, return a node responsible for it
func (c *consistentHashRing) getNode(key string) (nodeName string, nodeValue string) {
	keyHash := md5.Sum([]byte(key))
	var index int
	index = sort.Search(len(c.replicasHashes),
		func(i int) bool {
			res := bytes.Compare(c.replicasHashes[i][:], keyHash[:])
			return res >= 0
		})

	if index == len(c.replicasHashes) {
		// Special case to connect end of the hash function range
		// to the beginning to create a "ring"
		index = 0
	}
	replicaHash := c.replicasHashes[index]
	nodeName = string(c.replicasHashes[index][:])
	nodeValue = c.nodes[string(replicaHash[:])]

	return

}
