package table

type Item struct {
	Key, Value string
}

// Given a sorted list of key/value pairs, write them out according to the format you designed.
func Build(filePath string, sortedItems []Item) error {
	// Based on the idea in this answer:
	// https://dba.stackexchange.com/a/11190

	// File format is the fowllowing:
	// +--------------+
	// | 3            | number of entries (4 bytes)
	// +--------------+
	// | 16           | offset of first key data (4 bytes)
	// +--------------+
	// | 24           | offset of second key data (4 bytes)
	// +--------------+
	// | 39           | offset of third key data (4 bytes)
	// +--------------+-----------------+
	// | key one | data offset(4 bytes) |
	// +----------------+----------------------+
	// | key number two | data offset(4 bytes) |
	// +-----------------------+----------------------+
	// | this is the third key | data offset(4 bytes) |
	// +-----------------------+----------------------+
	// +-----------+
	// | value one |
	// +-----------+-------+
	// | value for key two |
	// +-------------------------+
	// | this is the third value |
	// +-------------------------+
	return nil
}

// A Table provides efficient access into sorted key/value data that's organized according
// to the format you designed.
//
// Although a Table shouldn't keep all the key/value data in memory, it should contain
// some metadata to help with efficient access (e.g. size, index, optional Bloom filter).
type Table struct {
	// Design of the file format requires us to load entire list of
	// keys into memory, so we can search quickly, which is suboptiomal.
	// Fix in next iteration?
	// skip list of keys to data offsets
	// file descriptor? pointer? ...
}

// Prepares a Table for efficient access. This will likely involve reading some metadata
// in order to populate the fields of the Table struct.
func LoadTable(filePath string) (*Table, error) {
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
