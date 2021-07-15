package main

type bstNode struct {
	item  Item
	left  *bstNode
	right *bstNode
}

// bstOC stores items in a binary search tree. Since it doesn't do any sort of
// balancing, performance will degrade badly when items are added in order.
type bstOC struct {
	root *bstNode
}

func newBstOC() *bstOC {
	return &bstOC{}
}

// Finds the first node such that node.item.Key >= key; returns `nil` if no
// such node exists.
func bstFirstGE(node *bstNode, key string) *bstNode {
	if node == nil {
		return nil
	}
	if key < node.item.Key {
		candidate := bstFirstGE(node.left, key)
		if candidate != nil {
			return candidate
		} else {
			return node
		}
	} else if key == node.item.Key {
		return node
	} else {
		return bstFirstGE(node.right, key)
	}
}

func bstPut(node *bstNode, key, value string) (*bstNode, bool) {
	if node == nil {
		return &bstNode{
			item:  Item{key, value},
			left:  nil,
			right: nil,
		}, true
	}
	var ok bool
	if key < node.item.Key {
		node.left, ok = bstPut(node.left, key, value)
	} else if key == node.item.Key {
		node.item.Value = value
	} else {
		node.right, ok = bstPut(node.right, key, value)
	}
	return node, ok
}

func bstDelete(node *bstNode, key string) (*bstNode, bool) {
	if node == nil {
		return nil, false
	}

	var ok bool
	if key < node.item.Key {
		node.left, ok = bstDelete(node.left, key)
		return node, ok
	} else if key > node.item.Key {
		node.right, ok = bstDelete(node.right, key)
		return node, ok
	} else /* key == node.item.Key */ {
		if node.left == nil && node.right == nil {
			return nil, true
		} else if node.left == nil {
			return node.right, true
		} else if node.right == nil {
			return node.left, true
		} else {
			successor := node.right
			for successor.left != nil {
				successor = successor.left
			}
			item := successor.item
			node.right, _ = bstDelete(node.right, item.Key)
			node.item = item
			return node, true
		}
	}
}

func (o *bstOC) Get(key string) (string, bool) {
	node := bstFirstGE(o.root, key)
	if node != nil && node.item.Key == key {
		return node.item.Value, true
	}
	return "", false
}

func (o *bstOC) Put(key, value string) bool {
	var ok bool
	o.root, ok = bstPut(o.root, key, value)
	return ok
}

func (o *bstOC) Delete(key string) bool {
	var ok bool
	o.root, ok = bstDelete(o.root, key)
	return ok
}

func (o *bstOC) RangeScan(startKey, endKey string) Iterator {
	var node *bstNode
	if o.root != nil {
		node = bstFirstGE(o.root, startKey)
	}
	return &bstOCIterator{o, node, startKey, endKey}
}

type bstOCIterator struct {
	o                *bstOC
	node             *bstNode
	startKey, endKey string
}

func (iter *bstOCIterator) Next() {
	// Lazy approach: just search from the root on every iteration
	iter.node = bstFirstGE(iter.o.root, iter.node.item.Key+"\u0000")
}

func (iter *bstOCIterator) Valid() bool {
	return iter.node != nil && iter.node.item.Key <= iter.endKey
}

func (iter *bstOCIterator) Key() string {
	return iter.node.item.Key
}

func (iter *bstOCIterator) Value() string {
	return iter.node.item.Value
}
