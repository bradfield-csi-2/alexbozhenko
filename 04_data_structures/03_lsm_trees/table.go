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
	// number of block pointers in the footer (uint32)
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
	totalBlocks := 0
	kvPairsInCurrentBlock := 0

	// Write key_length, value_length, key, value
	// from the beginning of the file
	for _, item := range sortedItems {
		//	fmt.Printf("Writing key(%v):val(%v)= %v:%v\n", len(item.Key), len(item.Value), item.Key, item.Value)
		if bytesWrittenInCurrentBlock > BLOCK_SIZE || // block already filled in
			bytesWrittenInCurrentBlock == 0 { // first block
			fE := footerEntry{
				key:          item.Key,
				block_offset: uint32(dataBytesWrittenTotal),
			}
			footerEntries = append(footerEntries, fE)
			bytesWrittenInCurrentBlock = 0
			totalBlocks++
			//fmt.Printf("Adding footer entry %+v with %v KV pairs\n", fE, kvPairsInCurrentBlock)
			kvPairsInCurrentBlock = 0
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
		kvPairsInCurrentBlock++
	}
	// Write footer consisting of:
	// key_lengh, key, offset
	// for each block
	for _, fE := range footerEntries {
		fE.write(file)
	}

	// write number of block pointers in the footer
	binary.Write(file, binary.BigEndian, uint32(totalBlocks))

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
	footerStartOffset   uint32
	keysToKVPairOffsets *skip_list.SkipListOC
}

// Prepares a Table for efficient access. This will likely involve reading some metadata
// in order to populate the fields of the Table struct.
func LoadTable(filePath string) (*Table, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	keysToKVPairOffsets := skip_list.NewSkipListOC()
	f := bufio.NewReader(file)
	file.Seek(-8, io.SeekEnd)
	buf := make([]byte, 4)
	io.ReadFull(f, buf)
	numberOfBlocks := binary.BigEndian.Uint32(buf)
	io.ReadFull(f, buf)
	footerStartOffset := binary.BigEndian.Uint32(buf)
	file.Seek(int64(footerStartOffset), io.SeekStart)
	fmt.Printf("numberOfBlocks = %v, footerStartOffset = %v\n", numberOfBlocks, footerStartOffset)
	for blockNumber := 0; blockNumber < int(numberOfBlocks); blockNumber++ {
		io.ReadFull(f, buf)
		keyLength := binary.BigEndian.Uint32(buf)
		keyBuf := make([]byte, keyLength)
		io.ReadFull(f, keyBuf)
		key := string(keyBuf)
		io.ReadFull(f, buf)
		blockOffset := binary.BigEndian.Uint32(buf)
		keysToKVPairOffsets.Put(key, blockOffset)
	}

	return &Table{
		file:                file,
		footerStartOffset:   footerStartOffset,
		keysToKVPairOffsets: keysToKVPairOffsets,
	}, nil
}

func (t *Table) Get(key string) (string, bool, error) {
	blockOffset, ok := t.keysToKVPairOffsets.GetLE(key)
	//fmt.Printf("Getting key=%v, found block offset=%v\n", key, blockOffset)
	f := t.file
	if ok {
		// block containing key that is less than or equal was found
		// now let's do a linear scan until we find the matchin key
		buf := make([]byte, 4)
		f.Seek(int64(blockOffset), io.SeekStart)
		// start of the loop
		currentSeek, _ := f.Seek(0, io.SeekCurrent)
		elementsSkipped := 0
		for ; uint32(currentSeek) < t.footerStartOffset; currentSeek, _ = f.Seek(0, io.SeekCurrent) {

			io.ReadFull(f, buf)
			keyLength := binary.BigEndian.Uint32(buf)
			io.ReadFull(f, buf)
			valueLength := binary.BigEndian.Uint32(buf)
			keyBuf := make([]byte, keyLength)
			io.ReadFull(f, keyBuf)
			currentKey := string(keyBuf)
			if currentKey == key {
				// We found a value, return it
				//fmt.Printf("Found a value after doing %v linear skips\n", elementsSkipped)

				valueBuf := make([]byte, valueLength)
				io.ReadFull(f, valueBuf)
				value := string(valueBuf)
				return value, true, nil
			} else if currentKey > key {
				return "", false, nil
			}
			// keep going with linear scan
			f.Seek(int64(valueLength), io.SeekCurrent)
			elementsSkipped++
		}

	}
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
