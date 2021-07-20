package table

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"lsm_trees/skip_list"
	"os"
	"unsafe"
)

type Item struct {
	Key, Value string
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

type footerEntry struct {
	key          string
	block_offset uint32 // position of the first key_length in the data block
}

func (fE *footerEntry) write(w io.Writer) {
	fmt.Printf("writitng footer entry %+v\n", fE)
	binary.Write(w, binary.BigEndian, uint32(len(fE.key)))
	key := fE.key
	binary.Write(w, binary.BigEndian, *(*[]byte)(unsafe.Pointer(&key)))
	binary.Write(w, binary.BigEndian, uint32(fE.block_offset))
}

const BLOCK_SIZE = 4096 //bytes

// Given a sorted list of key/value pairs, write them out according to the format you designed.
func Build(filePath string, sortedItems []Item) error {
	// Updated file format:
	// data block
	// data block
	// data block
	// ...
	// data block
	// Footer that loaded into memory with metadata that should allow binary search :
	// key_length, key, block_offset
	// key_length, key, block_offset
	// key_length, key, block_offset
	// ...
	// key_length, key, block_offset
	// offset_to_footer_beginning (uint32)

	// Each block contains(we will do linear scan inside the block):
	// key_length varint32 (or just uint32)
	// value_length varint32
	// key
	// value
	// ...
	file, err := os.Create(filePath)
	footerEntries := []footerEntry{}
	if err != nil {
		return err
	}
	defer file.Close()
	bytesWrittenInCurrentBlock := 0
	dataBytesWrittenTotal := 0

	// Write key_length, value_length, key, value
	// from the beginning of the file
	for _, item := range sortedItems {
		fmt.Printf("Writing key(%v):val(%v)= %v:%v\n", len(item.Key), len(item.Value), item.Key, item.Value)
		if bytesWrittenInCurrentBlock > BLOCK_SIZE || // block already filled in
			bytesWrittenInCurrentBlock == 0 { // first block
			footerEntries = append(footerEntries, footerEntry{
				key:          item.Key,
				block_offset: uint32(dataBytesWrittenTotal),
			})
			bytesWrittenInCurrentBlock = 0
		}
		binary.Write(file, binary.BigEndian, uint32(len(item.Key)))
		binary.Write(file, binary.BigEndian, uint32(len(item.Value)))
		dataBytesWrittenTotal += 8
		bytesWrittenInCurrentBlock += 8
		binary.Write(file, binary.BigEndian, *(*[]byte)(unsafe.Pointer(&item.Key)))
		dataBytesWrittenTotal += len(item.Key)
		bytesWrittenInCurrentBlock += len(item.Key)
		binary.Write(file, binary.BigEndian, *(*[]byte)(unsafe.Pointer(&item.Value)))
		dataBytesWrittenTotal += len(item.Value)
		bytesWrittenInCurrentBlock += len(item.Value)
	}
	// Write footer consisting of:
	// write key_lengh, key, offset
	// for each block
	for _, fE := range footerEntries {
		fE.write(file)
	}
	// write position of the beginning of the footer
	binary.Write(file, binary.BigEndian, uint32(dataBytesWrittenTotal))
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

	keyOffsets := make([]uint32, entriesNumber)
	for i := range keyOffsets {
		io.ReadFull(f, buf)
		keyOffsets[i] = binary.BigEndian.Uint32(buf)
	}
	key_ValueOffset_Start := (entriesNumber + 1) * 4
	for i := 0; i < len(keyOffsets)-1; i++ {
		key_ValueOffset_End := keyOffsets[i]
		file.Seek(int64(key_ValueOffset_Start), io.SeekStart)
		buf := make([]byte, key_ValueOffset_End-key_ValueOffset_Start+1)
		io.ReadFull(f, buf)
		key := string(buf[:len(buf)-4])
		valueOffset := binary.BigEndian.Uint32(buf[len(buf)-4:])
		keysToValuesOffsets.Put(key, valueOffset)
		key_ValueOffset_Start = key_ValueOffset_End + 1
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
