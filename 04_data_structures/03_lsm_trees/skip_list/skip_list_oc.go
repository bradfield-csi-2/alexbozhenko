package skip_list

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

type skipListNode struct {
	item    Item
	forward []*skipListNode
}

func (node skipListNode) Level() int {
	return len(node.forward)
}

// head is a dummy node with nil item, that is used to
// simplify handiling edge cases. Head never contains an item
type SkipListOC struct {
	head *skipListNode
}

func (list SkipListOC) Level() int {
	return len(list.head.forward)
}

func NewSkipListOC() *SkipListOC {
	rand.Seed(time.Now().UnixNano())
	return &SkipListOC{
		head: &skipListNode{
			// item is nil
			forward: []*skipListNode{nil},
		},
	}
}

func (list SkipListOC) String() string {
	result := ""
	for i := (list.Level() - 1); i >= 0; i-- {
		result += fmt.Sprintf("Level %v:\n", i+1)
		for currentNode := list.head.forward[i]; currentNode !=
			nil && i < currentNode.Level(); currentNode = (*currentNode).forward[i] {
			result += fmt.Sprintf(
				"%s: %v -> ", currentNode.item.Key, currentNode.item.Value)
		}
		result += "nil\n"
	}
	return result
}

func (list SkipListOC) Length() (length int) {
	node := list.head.forward[0]
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
// dummy head node is returned when "previous" node at this Level
// should be the first in the list
// Never returns an empty slice, since level of empty list = 1
func (list *SkipListOC) findPreviousNodes(key string) []*skipListNode {
	previousNodes := make([]*skipListNode, list.Level())
	topLevel := list.Level() - 1
	currentNode := list.head
	for level := topLevel; level >= 0; level-- {
		for currentNode != nil &&
			currentNode.forward[level] != nil &&
			currentNode.forward[level].item.Key < key {
			currentNode = currentNode.forward[level]
		}
		previousNodes[level] = currentNode
	}
	return previousNodes
}

// Return previous nodes, and nil for node, if node was not found in the list
func (skipList *SkipListOC) get(key string) (previousNodes []*skipListNode,
	foundNode *skipListNode) {
	foundNode = nil
	previousNodes = skipList.findPreviousNodes(key)
	previousNode := previousNodes[0]
	nextNode := previousNode.forward[0]
	if nextNode != nil && nextNode.item.Key == key {
		foundNode = nextNode
		return
	}
	return
}

func (skipList *SkipListOC) Get(key string) (uint32, bool) {
	_, node := skipList.get(key)
	if node != nil {
		return node.item.Value, true
	}
	return 0, false
}

func (skipList *SkipListOC) Put(key string, value uint32) bool {
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
		node.forward[currentLevel] = previousNodes[currentLevel].forward[currentLevel]
		previousNodes[currentLevel].forward[currentLevel] = node
	}
	// if new node level > current skipList.Level() we handle pointers
	// from header to the node and from node to nil
	for currentLevel := len(previousNodes); currentLevel <
		newNodeLevel; currentLevel++ {
		node.forward[currentLevel] = nil
		skipList.head.forward = append(skipList.head.forward, node)
	}

	return true
}

func (skipList *SkipListOC) Delete(key string) bool {
	previousNodes, node := skipList.get(key)
	if node == nil {
		return false
	}
	foundNodeLevel := node.Level()

	for currentLevel := 0; currentLevel < foundNodeLevel; currentLevel++ {
		previousNodes[currentLevel].forward[currentLevel] = node.forward[currentLevel]
	}

	// When after removing the node, header has pointer to nil,
	// that means that the deleted node had the highest level
	// we need to remove extra levels from the header
	// by search for nil pointers
	// But if node was the only element of the list,
	// we still leave the nil pointer per convention
	for currentLevel := 1; currentLevel < skipList.Level(); currentLevel++ {
		if skipList.head.forward[currentLevel] == nil {
			tmp := make([]*skipListNode, currentLevel+1)
			copy(tmp, skipList.head.forward)
			skipList.head.forward = tmp
			break
		}
	}
	return true
}

func (skipList *SkipListOC) RangeScan(startKey, endKey string) Iterator {
	previousNodes, startNode := skipList.get(startKey)
	var currentNode *skipListNode
	if startNode != nil {
		// when node was found, start from it
		currentNode = startNode
	} else {
		// When node was not found in the list, start from the first greater element
		currentNode = previousNodes[0].forward[0]
	}
	return &skipListOCIterator{
		skipList:    skipList,
		currentNode: currentNode,
		startKey:    startKey,
		endKey:      endKey,
	}
}

type skipListOCIterator struct {
	skipList         *SkipListOC
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

func (iter *skipListOCIterator) Value() uint32 {
	return iter.currentNode.item.Value
}
