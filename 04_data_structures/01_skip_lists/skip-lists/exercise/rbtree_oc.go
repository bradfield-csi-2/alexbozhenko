package main

// rbTreeOC stores items in a Red Black Tree implementation from Github.
type rbTreeOC struct {
	tree *Tree
}

func newRbTreeOC() *rbTreeOC {
	return &rbTreeOC{
		tree: &Tree{},
	}
}

func (o *rbTreeOC) Get(key string) (string, bool) {
	return o.tree.Get(key)
}

func (o *rbTreeOC) Put(key, value string) bool {
	// Note: the 3rd party implementation of a Put doesn't indicate whether
	// a new item was actually added, so here we just assume we're always
	// adding a new key (which is how we use Put in our test).
	o.tree.Put(key, value)
	return true
}

func (o *rbTreeOC) Delete(key string) bool {
	// Note: the 3rd party implementation of a Delete doesn't indicate whether
	// the key was actually deleted, so here we just assume we're always
	// deleting an existing key (which is how we use Delete in our test).
	o.tree.Remove(key)
	return true
}

func (o *rbTreeOC) RangeScan(startKey, endKey string) Iterator {
	var rbIter *RBIterator
	node, ok := o.tree.Ceiling(startKey)
	if ok {
		rbIter = o.tree.IteratorAt(node)
	}
	return &rbTreeOCIterator{o, rbIter, startKey, endKey}
}

type rbTreeOCIterator struct {
	o                *rbTreeOC
	rbIter           *RBIterator
	startKey, endKey string
}

func (iter *rbTreeOCIterator) Next() {
	ok := iter.rbIter.Next()
	if !ok {
		iter.rbIter = nil
	}
}

func (iter *rbTreeOCIterator) Valid() bool {
	return iter.rbIter != nil && iter.Key() <= iter.endKey
}

func (iter *rbTreeOCIterator) Key() string {
	return iter.rbIter.Key()
}

func (iter *rbTreeOCIterator) Value() string {
	return iter.rbIter.Value()
}
