package table

import (
	"bufio"
	"encoding/binary"
	"io"
	"lsm_trees/skip_list"
	"os"
)

type Item struct {
	Key, Value string
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
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
	file                *os.File
	keysToValuesOffsets *skip_list.SkipListOC
}

// Prepares a Table for efficient access. This will likely involve reading some metadata
// in order to populate the fields of the Table struct.
func LoadTable(filePath string) (*Table, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	keysToValuesOffsets := skip_list.NewSkipListOC()
	f := bufio.NewReader(file)
	buf := make([]byte, 4)
	io.ReadFull(f, buf)
	entriesNumber := binary.BigEndian.Uint32(buf)
	keyOffsets := make([]uint32, entriesNumber+1)
	keyOffsets[0] = (entriesNumber + 1) * 4
	for i := range keyOffsets {
		io.ReadFull(f, buf)
		keyOffsets[i] = binary.BigEndian.Uint32(buf)
	}
	for i := 1; i < len(keyOffsets); i++ {
		keyValueStart := keyOffsets[i-1]
		keyValueEnd := keyOffsets[i]
		file.Seek(int64(keyValueStart), io.SeekStart)
		buf := make([]byte, keyValueEnd-keyValueStart+1)
		io.ReadFull(f, buf)
		key := string(buf[:len(buf)-4])
		valueOffset := binary.BigEndian.Uint32(buf[len(buf)-4:])
		keysToValuesOffsets.Put(key, valueOffset)
	}

	return &Table{
		file:                file,
		keysToValuesOffsets: keysToValuesOffsets,
	}, nil
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
