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

type replicaHashWithNodeName struct {
	// hash of replicaName, where
	// replicaName constructed like this:
	// nodeName:%d
	// where d is replica number
	hash [md5.Size]byte
	// Name of the related node
	nodeName string
}

type consistentHashRing struct {
	// how many times we create points on the hash ring for each node
	replicas int
	// Slice is sorted by replicaHash
	replicasHashes []replicaHashWithNodeName
	// Map from nodeName to nodeValue, where
	// nodeValue can be anything you want, e.g. an url
	nodes map[string]string
}

func NewConsistentHashRing(replicas, nodeNumber int) *consistentHashRing {
	rHashes := make([]replicaHashWithNodeName, 0)
	return &consistentHashRing{
		replicas:       replicas,
		replicasHashes: rHashes,
		nodes:          make(map[string]string),
	}
}

func replicaHash(nodeName string, replicaNumber int) [md5.Size]byte {
	return md5.Sum([]byte(fmt.Sprintf("%s:%d", nodeName, replicaNumber)))
}

func (c *consistentHashRing) addNode(nodeName string, nodeValue string) {
	for replicaNumber := 0; replicaNumber < c.replicas; replicaNumber++ {
		replicaHashArray := replicaHash(nodeName, replicaNumber)
		replicaHash := replicaHashWithNodeName{
			hash:     replicaHashArray,
			nodeName: nodeName,
		}
		index := sort.Search(len(c.replicasHashes),
			func(i int) bool {
				res := bytes.Compare(c.replicasHashes[i].hash[:], replicaHash.hash[:])
				return res >= 0
			})

		c.replicasHashes = append(append(c.replicasHashes[:index][:],
			replicaHash),
			c.replicasHashes[index:]...)
	}
	c.nodes[nodeName] = nodeValue
}

func (c *consistentHashRing) deleteNode(nodeName string) {
	for replicaNumber := 0; replicaNumber < c.replicas; replicaNumber++ {
		replicaHashArray := replicaHash(nodeName, replicaNumber)
		index := sort.Search(len(c.replicasHashes),
			func(i int) bool {
				return (bytes.Equal(c.replicasHashes[i].hash[:], replicaHashArray[:]))
			})
		if index == len(c.replicasHashes) {
			panic(fmt.Sprintf("BUG! %s is not present in the hash ring", nodeName))
		}
		c.replicasHashes = append(c.replicasHashes[:index], c.replicasHashes[index+1:]...)
	}
	delete(c.nodes, nodeName)
}

// Given a key, return a node responsible for it
func (c *consistentHashRing) getNode(key string) (nodeName string, nodeValue string) {
	keyHash := md5.Sum([]byte(key))
	var index int
	// find index of the first replica hash, where
	// replicaHash >= keyHash
	index = sort.Search(len(c.replicasHashes),
		func(i int) bool {
			res := bytes.Compare(c.replicasHashes[i].hash[:], keyHash[:])
			return res >= 0
		})

	if index == len(c.replicasHashes) {
		// Special case to connect end of the hash function range
		// to the beginning to create a "ring"
		// In other words, is the first replica responsible for
		// the tail of the hash ring
		index = 0
	}
	nodeName = c.replicasHashes[index].nodeName
	nodeValue = c.nodes[nodeName]

	return

}
