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

func (list skipListOC) Length() (length int) {
	node := list.head[0]
	length = 0
	for ; node != nil; length++ {
		node = node.forward[0]
	}
	return length

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
		currentNode := list.head[level]
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

func (skipList *skipListOC) Delete(key string) bool {
	previousNodes, node := skipList.get(key)
	if node == nil {
		return false
	}
	foundNodeLevel := node.Level()

	for currentLevel := 0; currentLevel < foundNodeLevel; currentLevel++ {
		if previousNodes[currentLevel] == nil {
			skipList.head[currentLevel] = node.forward[currentLevel]
		} else {
			previousNodes[currentLevel].forward[currentLevel] = node.forward[currentLevel]
		}
	}

	// If after removing the node header has pointer to nil
	// that means that the node had the highest level
	// we need to remove extra levels from the header
	// by search for nil pointers
	// But if node was the only element of the list,
	// we still leave the nil pointer per convention
	for currentLevel := 1; currentLevel < skipList.Level(); currentLevel++ {
		if skipList.head[currentLevel] == nil {
			tmp := make([]*skipListNode, currentLevel+1)
			copy(tmp, skipList.head)
			skipList.head = tmp
			break
		}
	}
	return true
}

func (skipList *skipListOC) RangeScan(startKey, endKey string) Iterator {
	previousNodes, startNode := skipList.get(startKey)
	var currentNode *skipListNode
	if startNode == nil && previousNodes[0] == nil {
		currentNode = skipList.head[0]
	} else if previousNodes[0] != nil {
		currentNode = previousNodes[0].forward[0]
	} else {
		currentNode = startNode
	}
	return &skipListOCIterator{
		skipList:    skipList,
		currentNode: currentNode,
		startKey:    startKey,
		endKey:      endKey,
	}
}

type skipListOCIterator struct {
	skipList         *skipListOC
	currentNode      *skipListNode
	startKey, endKey string
}

func (iter *skipListOCIterator) Next() {
	iter.currentNode = iter.currentNode.forward[0]
}

func (iter *skipListOCIterator) Valid() bool {
	return iter.currentNode != nil && iter.currentNode.item.Key < iter.endKey
}

func (iter *skipListOCIterator) Key() string {
	return iter.currentNode.item.Key
}

func (iter *skipListOCIterator) Value() string {
	return iter.currentNode.item.Value
}
