package main

const (
	maxBlockSize = 500
	splitSize    = maxBlockSize / 2
)

type linkedBlockNode struct {
	items []Item
	next  *linkedBlockNode
	prev  *linkedBlockNode
}

// linkedBlockOC is similar to linkedListOC, but each node stores a block of
// items, so that we can skip blocks of items at a time when searching.
type linkedBlockOC struct {
	head *linkedBlockNode
	tail *linkedBlockNode
}

func newLinkedBlockOC() *linkedBlockOC {
	head := &linkedBlockNode{}
	tail := &linkedBlockNode{}
	head.next = tail
	tail.prev = head
	return &linkedBlockOC{head, tail}
}

// Find the first block such that the last item in block.items satisfies
// item.Key >= key; returns o.tail if no such block exists.
func (o *linkedBlockOC) firstGE(key string) *linkedBlockNode {
	b := o.head.next
	for b != o.tail && b.items[len(b.items)-1].Key < key {
		b = b.next
	}
	return b
}

func (o *linkedBlockOC) Get(key string) (string, bool) {
	b := o.firstGE(key)
	if b == o.tail {
		return "", false
	}
	return sliceGet(b.items, key)
}

func (o *linkedBlockOC) Put(key, value string) bool {
	b := o.firstGE(key)
	if b == o.tail {
		b = b.prev
	}

	// If the collection was empty, add the first (non-sentinel) block.
	if b == o.head {
		newBlock := &linkedBlockNode{
			next: o.tail,
			prev: o.head,
		}
		o.head.next = newBlock
		o.tail.prev = newBlock
		b = newBlock
	}

	ok := slicePut(&b.items, key, value)

	// Split the current block if it got too large.
	if len(b.items) > maxBlockSize {
		newBlock := &linkedBlockNode{
			items: b.items[splitSize:],
			next:  b.next,
			prev:  b,
		}
		b.items = b.items[:splitSize]
		b.next.prev = newBlock
		b.next = newBlock
	}
	return ok
}

func (o *linkedBlockOC) Delete(key string) bool {
	b := o.firstGE(key)
	if b == o.tail {
		return false
	}
	ok := sliceDelete(&b.items, key)
	if len(b.items) == 0 {
		b.prev.next = b.next
		b.next.prev = b.prev
	}
	return ok
}

func (o *linkedBlockOC) RangeScan(startKey, endKey string) Iterator {
	b := o.firstGE(startKey)
	index := 0
	if b != o.tail {
		index = sliceFirstGE(b.items, startKey)
	}
	return &linkedBlockOCIterator{o, b, index, startKey, endKey}
}

type linkedBlockOCIterator struct {
	o                *linkedBlockOC
	b                *linkedBlockNode
	index            int
	startKey, endKey string
}

func (iter *linkedBlockOCIterator) Next() {
	iter.index++
	if iter.index == len(iter.b.items) {
		iter.b = iter.b.next
		iter.index = 0
	}
}

func (iter *linkedBlockOCIterator) Valid() bool {
	return iter.b != iter.o.tail && iter.b.items[iter.index].Key <= iter.endKey
}

func (iter *linkedBlockOCIterator) Key() string {
	return iter.b.items[iter.index].Key
}

func (iter *linkedBlockOCIterator) Value() string {
	return iter.b.items[iter.index].Value
}
