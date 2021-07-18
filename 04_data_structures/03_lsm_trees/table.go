package table

type Item struct {
	Key, Value string
}

// Given a sorted list of key/value pairs, write them out according to the format you designed.
func Build(path string, sortedItems []Item) error {
	return nil
}

// A Table provides efficient access into sorted key/value data that's organized according
// to the format you designed.
//
// Although a Table shouldn't keep all the key/value data in memory, it should contain
// some metadata to help with efficient access (e.g. size, index, optional Bloom filter).
type Table struct {
	// TODO
}

// Prepares a Table for efficient access. This will likely involve reading some metadata
// in order to populate the fields of the Table struct.
func LoadTable(path string) (*Table, error) {
	return nil, nil
}

func (t *Table) Get(key string) (string, bool, error) {
	return "", false, nil
}

func (t *Table) RangeScan(startKey, endKey string) (Iterator, error) {
	return nil, nil
}

type Iterator interface {
	// Advances to the next item in the range. Assumes Valid() == true.
	Next()

	// Indicates whether the iterator is currently pointing to a valid item.
	Valid() bool

	// Returns the Item the iterator is currently pointing to. Assumes Valid() == true.
	Item() Item
}
