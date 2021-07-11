package main

// OC is an Ordered Collection, similar to a map / Python dict / Ruby hash, but
// it additionally stores keys in order (and thus supports range scans).
type OC interface {
	// The second return value will be `false` when the `key` hasn't been
	// associated with any value.
	Get(key string) (string, bool)

	// Put should return `true` if a new key was added, and `false` if an
	// existing key had its value updated.
	Put(key, value string) bool

	// Delete should return whether or not the key was actually deleted, i.e.
	// it should return `true` if the key existed before deletion.
	Delete(key string) bool

	// startKey and endKey are inclusive.
	RangeScan(startKey, endKey string) Iterator
}

type Iterator interface {
	// Advances to the next item in the range. Assumes Valid() == true.
	Next()

	// Indicates whether the iterator is currently pointing to a valid item.
	Valid() bool

	// Returns the Key for the item the iterator is currently pointing to.
	// Assumes Valid() == true.
	Key() string

	// Returns the Value for the item the iterator is currently pointing to.
	// Assumes Valid() == true.
	Value() string
}

type Item struct {
	Key, Value string
}
