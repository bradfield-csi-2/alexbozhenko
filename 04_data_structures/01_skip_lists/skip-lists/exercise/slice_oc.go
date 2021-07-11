package main

// sliceOC is not bad for searching, but terrible for insertion and deletion.
// It stores items in a slice and uses binary search to find items of interest
// (whether searching, inserting, or removing). However, inserting and removing
// are both linear time because they involve moving many items in the slice.
//
// Most of the heavy lifting happens in slice_util.go, so that the functions
// can be reused for the "linked block" approach, where each block has its own
// slice of items.
type sliceOC struct {
	items []Item
}

func newSliceOC() *sliceOC {
	return &sliceOC{}
}

func (o *sliceOC) Get(key string) (string, bool) {
	return sliceGet(o.items, key)
}

func (o *sliceOC) Put(key, value string) bool {
	return slicePut(&o.items, key, value)
}

func (o *sliceOC) Delete(key string) bool {
	return sliceDelete(&o.items, key)
}

func (o *sliceOC) RangeScan(startKey, endKey string) Iterator {
	return &sliceOCIterator{o, sliceFirstGE(o.items, startKey), startKey, endKey}
}

type sliceOCIterator struct {
	o                *sliceOC
	index            int
	startKey, endKey string
}

func (iter *sliceOCIterator) Next() {
	iter.index++
}

func (iter *sliceOCIterator) Valid() bool {
	return iter.index < len(iter.o.items) && iter.o.items[iter.index].Key <= iter.endKey
}

func (iter *sliceOCIterator) Key() string {
	return iter.o.items[iter.index].Key
}

func (iter *sliceOCIterator) Value() string {
	return iter.o.items[iter.index].Value
}
