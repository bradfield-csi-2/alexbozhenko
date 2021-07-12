package main

import (
	"fmt"
	"math/rand"
	"time"
)

const maxLevel = 16
const probability = 0.5

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

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
	rand.Seed(time.Now().UnixNano())
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

func getLevelInHumanTerms() (level int) {
	level = 1
	for ; level < maxLevel && rand.Float32() < probability; level++ {
	}
	return
}

// Returns slice of pointes to nodes in the skip list,
// that are less than given node on each level.
// Ordered from the lowest to the highest level
// On the lowest level nil is returned when previous node
// should be located at the head of the list
// Never returns an empty slice, since level of empty list = 1
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

func (skipList *skipListOC) get(key string) (previousNodes []*skipListNode,
	foundNode *skipListNode) {
	foundNode = nil
	previousNodes = skipList.findPreviousNodes(key)
	previousNode := previousNodes[0]
	if previousNode == nil {
		// findPreviousNodes returns nil when previous
		// node is at the head of the list, so we need to check
		// if head of the list is the target node
		headNode := skipList.head[0]
		if headNode != nil && headNode.item.Key == key {
			foundNode = headNode
			return
		}
	} else {
		nextNode := previousNode.forward[0]
		if nextNode != nil && nextNode.item.Key == key {
			foundNode = nextNode
			return
		}
	}
	return
}

func (skipList *skipListOC) Get(key string) (string, bool) {
	_, node := skipList.get(key)
	if node != nil {
		return node.item.Value, true
	}
	return "", false
}

func (skipList *skipListOC) Put(key, value string) bool {
	previousNodes, node := skipList.get(key)
	if node != nil {
		node.item.Value = value
		return false
	}
	newNodeLevel := getLevelInHumanTerms()
	node = &skipListNode{
		item: Item{
			Key:   key,
			Value: value,
		},
		forward: make([]*skipListNode, newNodeLevel),
	}
	// if new node level <= current skipList.Level() we just update
	// pointers to and from it
	for currentLevel := 0; currentLevel <
		min(newNodeLevel, len(previousNodes)); currentLevel++ {
		// we need to update the list header when inserting
		// the first node on this level
		if previousNodes[currentLevel] == nil {
			node.forward[currentLevel] = skipList.head[currentLevel]
			skipList.head[currentLevel] = node
		} else {
			node.forward[currentLevel] = previousNodes[currentLevel].forward[currentLevel]
			previousNodes[currentLevel].forward[currentLevel] = node
		}
	}
	// if new node level > current skipList.Level() we handle pointers
	// from header to the node and from node to nil
	for currentLevel := len(previousNodes); currentLevel <
		newNodeLevel; currentLevel++ {
		node.forward[currentLevel] = nil
		skipList.head = append(skipList.head, node)
	}

	return true
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
