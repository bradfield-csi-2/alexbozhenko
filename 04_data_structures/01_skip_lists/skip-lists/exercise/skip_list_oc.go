package main

import (
	"fmt"
	"math/rand"
)

const maxLevel = 16
const probability = 0.5

type skipListNode struct {
	item    Item
	forward []*skipListNode
}

func (node skipListNode) Level() int {
	return len(node.forward)
}

type skipListOC struct {
	head []*skipListNode
}

func (list skipListOC) Level() int {
	return len(list.head)
}

func newSkipListOC() *skipListOC {
	return &skipListOC{
		head: []*skipListNode{
			nil,
		},
	}
}

func (list skipListOC) String() string {
	result := ""
	for i := (list.Level() - 1); i >= 0; i-- {
		result += fmt.Sprintf("Level %v:\n", i+1)
		for currentNode := list.head[i]; currentNode !=
			nil && i < currentNode.Level(); currentNode = (*currentNode).forward[i] {
			result += fmt.Sprintf(
				"%s: %s -> ", currentNode.item.Key, currentNode.item.Value)
		}
		result += "nil\n"
	}
	return result
}

func getLevel() (level int) {
	for level := 0; level < maxLevel && rand.Float32() < probability; level++ {
	}
	return
}

// Returns slice of pointes to nodes in the skip list,
// that are less than given node on each level.
// Ordered from the lowest to the highest level
func (list *skipListOC) findPreviousNodes(key string) []*skipListNode {
	previousNodes := make([]*skipListNode, list.Level())
	topLevel := list.Level() - 1
	for level := topLevel; level >= 0; level-- {
		currentNode := list.head[topLevel]
		for currentNode != nil &&
			currentNode.item.Key < key {
			previousNodes[level] = currentNode
			currentNode = currentNode.forward[level]
		}
	}
	return previousNodes
}

func (o *skipListOC) Get(key string) (string, bool) {
	return "", false
}

func (o *skipListOC) Put(key, value string) bool {

	return false
}

func (o *skipListOC) Delete(key string) bool {
	return false
}

func (o *skipListOC) RangeScan(startKey, endKey string) Iterator {
	return &skipListOCIterator{}
}

type skipListOCIterator struct {
}

func (iter *skipListOCIterator) Next() {
}

func (iter *skipListOCIterator) Valid() bool {
	return false
}

func (iter *skipListOCIterator) Key() string {
	return ""
}

func (iter *skipListOCIterator) Value() string {
	return ""
}
