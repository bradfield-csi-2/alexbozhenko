package main

type linkedNode struct {
	item Item
	next *linkedNode
	prev *linkedNode
}

// linkedOC is probably one of the worst possible strategies we could use here.
// It stores items in a doubly linked list and does a linear scan through the
// list every time we want to do anything.
//
// Note that `head` and `tail` are dummy nodes that don't store actual items;
// instead, they're used to simplify edge case handling.
type linkedOC struct {
	head *linkedNode
	tail *linkedNode
}

func newLinkedOC() *linkedOC {
	head := &linkedNode{}
	tail := &linkedNode{}
	head.next = tail
	tail.prev = head
	return &linkedOC{head, tail}
}

// Find the first node such that node.item.Key >= key
// returns o.tail if no such node exists
func (o *linkedOC) firstGE(key string) *linkedNode {
	node := o.head.next
	for node != o.tail && node.item.Key < key {
		node = node.next
	}
	return node
}

func (o *linkedOC) Get(key string) (string, bool) {
	node := o.firstGE(key)
	if node != o.tail && node.item.Key == key {
		return node.item.Value, true
	}
	return "", false
}

func (o *linkedOC) Put(key, value string) bool {
	node := o.firstGE(key)
	if node.item.Key == key {
		node.item.Value = value
		return false
	} else {
		newNode := &linkedNode{
			item: Item{key, value},
			next: node,
			prev: node.prev,
		}
		node.prev.next = newNode
		node.prev = newNode
		return true
	}
}

func (o *linkedOC) Delete(key string) bool {
	node := o.firstGE(key)
	if node != o.tail {
		node.prev.next = node.next
		node.next.prev = node.prev
		return true
	}
	return false
}

func (o *linkedOC) RangeScan(startKey, endKey string) Iterator {
	node := o.firstGE(startKey)
	return &linkedOCIterator{o, node, startKey, endKey}
}

type linkedOCIterator struct {
	o                *linkedOC
	node             *linkedNode
	startKey, endKey string
}

func (iter *linkedOCIterator) Next() {
	iter.node = iter.node.next
}

func (iter *linkedOCIterator) Valid() bool {
	return iter.node != iter.o.tail && iter.node.item.Key <= iter.endKey
}

func (iter *linkedOCIterator) Key() string {
	return iter.node.item.Key
}

func (iter *linkedOCIterator) Value() string {
	return iter.node.item.Value
}
