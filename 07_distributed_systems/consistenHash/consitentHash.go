package consistenHash

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sort"
)

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
	//replicaNumber is used just for debug
	replicaNumber int
}

func (r replicaHashWithNodeName) String() string {
	return fmt.Sprintf("%x => %s:%d\n", r.hash, r.nodeName, r.replicaNumber)
}

type ConsistentHashRing struct {
	// how many times we create points on the hash ring for each node
	replicas int
	// Slice of numNodes*numReplicas hashes, that is sorted by replicaHash
	replicasHashes []replicaHashWithNodeName
	// Map from nodeName to nodeValue, where
	// nodeValue can be anything you want, e.g. an url
	nodes map[string]string
}

func NewConsistentHashRing(replicas int) *ConsistentHashRing {
	rHashes := make([]replicaHashWithNodeName, 0)
	return &ConsistentHashRing{
		replicas:       replicas,
		replicasHashes: rHashes,
		nodes:          make(map[string]string),
	}
}

func replicaHash(nodeName string, replicaNumber int) [md5.Size]byte {
	return md5.Sum([]byte(fmt.Sprintf("%s:%d", nodeName, replicaNumber)))
}

func (c *ConsistentHashRing) GetNodes() map[string]string {
	return c.nodes
}

func (c *ConsistentHashRing) AddNode(nodeName string, nodeValue string) {
	for replicaNumber := 0; replicaNumber < c.replicas; replicaNumber++ {
		replicaHashArray := replicaHash(nodeName, replicaNumber)
		replicaHash := replicaHashWithNodeName{
			hash:          replicaHashArray,
			nodeName:      nodeName,
			replicaNumber: replicaNumber,
		}
		index := sort.Search(len(c.replicasHashes),
			func(i int) bool {
				res := bytes.Compare(c.replicasHashes[i].hash[:], replicaHash.hash[:])
				return res >= 0
			})

		c.replicasHashes = append(c.replicasHashes, replicaHashWithNodeName{})
		copy(c.replicasHashes[index+1:], c.replicasHashes[index:])
		c.replicasHashes[index] = replicaHash
	}
	c.nodes[nodeName] = nodeValue
}

func (c *ConsistentHashRing) DeleteNode(nodeName string) {
	for replicaNumber := 0; replicaNumber < c.replicas; replicaNumber++ {
		replicaHashArray := replicaHash(nodeName, replicaNumber)
		index := sort.Search(len(c.replicasHashes),
			func(i int) bool {
				res := bytes.Compare(c.replicasHashes[i].hash[:], replicaHashArray[:])
				return res >= 0
			})

		if index >= len(c.replicasHashes) {
			panic(fmt.Sprintf("BUG! %s is not present in the hash ring", nodeName))
		}
		c.replicasHashes = append(c.replicasHashes[:index], c.replicasHashes[index+1:]...)
	}
	delete(c.nodes, nodeName)
}

// Given a key, return a node responsible for it
func (c *ConsistentHashRing) GetNode(key string) (nodeName string, nodeValue string) {
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

func (c *ConsistentHashRing) String() string {
	return fmt.Sprint(c.replicasHashes)
}
