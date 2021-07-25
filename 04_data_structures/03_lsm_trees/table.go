package table

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	bloom "lsm_trees/bloom_filter"
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
const KEY_LENGTH_BYTES = 4
const VALUE_LENGTH_BYTES = 4
const KEY_VALUE_LENGTH_BYTES = KEY_LENGTH_BYTES + VALUE_LENGTH_BYTES

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
	// bloom filter
	// offset_to_bloomFilter_beginning (uint32)
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
	bloomFilter := bloom.NewMyBloomFilter()

	// Write key_length, value_length, key, value
	// from the beginning of the file
	for _, item := range sortedItems {
		//	fmt.Printf("Writing key(%v):val(%v)= %v:%v\n", len(item.Key), len(item.Value), item.Key, item.Value)
		recordLength := len(item.Key) + len(item.Value) + KEY_VALUE_LENGTH_BYTES
		if (bytesWrittenInCurrentBlock != 0 &&
			bytesWrittenInCurrentBlock+recordLength > BLOCK_SIZE) || // block already filled in
			bytesWrittenInCurrentBlock == 0 { // first block
			fE := footerEntry{
				key:          item.Key,
				block_offset: uint32(dataBytesWrittenTotal),
			}
			//fmt.Printf("Adding footer entry %+v with %v KV pairs and %v bytes\n",
			//	fE, kvPairsInCurrentBlock, bytesWrittenInCurrentBlock)
			footerEntries = append(footerEntries, fE)
			bytesWrittenInCurrentBlock = 0
			totalBlocks++
			kvPairsInCurrentBlock = 0
		}
		bloomFilter.Add(item.Key)
		binary.Write(file, binary.BigEndian, uint32(len(item.Key)))
		binary.Write(file, binary.BigEndian, uint32(len(item.Value)))
		dataBytesWrittenTotal += KEY_VALUE_LENGTH_BYTES
		bytesWrittenInCurrentBlock += KEY_VALUE_LENGTH_BYTES
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

	// write bloom filter content to file
	bloomFilterStartPosition, _ := file.Seek(0, io.SeekCurrent)
	bFilter, _ := bloomFilter.MarshalBinary()
	binary.Write(file, binary.BigEndian, bFilter)

	// write bloom filter starting position
	binary.Write(file, binary.BigEndian, uint32(bloomFilterStartPosition))

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
	bloomFilter         *bloom.MyBloomFilter
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
	buf := make([]byte, 4)

	bloomFilterEndPosition, _ := file.Seek(-12, io.SeekEnd)

	io.ReadFull(f, buf)
	bloomFilterStartPosition := binary.BigEndian.Uint32(buf)

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

	file.Seek(int64(bloomFilterStartPosition), io.SeekStart)
	bloomFilterBuf := make([]byte,
		bloomFilterEndPosition-int64(bloomFilterStartPosition))
	io.ReadFull(file, bloomFilterBuf)

	bloomFilter := bloom.NewMyBloomFilter()
	bloomFilter.UnmarshalBinary(bloomFilterBuf)

	return &Table{
		file:                file,
		footerStartOffset:   footerStartOffset,
		keysToKVPairOffsets: keysToKVPairOffsets,
		bloomFilter:         bloomFilter,
	}, nil
}

func (t *Table) Get(key string) (string, bool, error) {
	blockOffset, ok := t.keysToKVPairOffsets.GetLE(key)
	if !ok {
		return "", false, nil
	}

	//If key is not found in bloomFilter, there is no need to continue
	if !t.bloomFilter.MaybeContains(key) {
		return "", false, nil
	}

	//fmt.Printf("Getting key=%v, found block offset=%v\n", key, blockOffset)
	f := t.file

	// block containing key that is less than or equal to target key was found
	// now let's do a linear scan until we find the matching key
	nextBlockOffset, ok := t.keysToKVPairOffsets.GetG(key)
	if !ok { // current block is last block, just find the end of it
		nextBlockOffset = t.footerStartOffset
	}
	blockSize := nextBlockOffset - blockOffset
	block := make([]byte, blockSize)
	f.Seek(int64(blockOffset), io.SeekStart)
	io.ReadFull(f, block)

	elementsSkipped := 0
	for len(block) > 0 {

		keyLength := binary.BigEndian.Uint32(block[:4])
		block = block[4:]

		valueLength := binary.BigEndian.Uint32(block[:4])
		block = block[4:]

		currentKey := string(block[:keyLength])
		block = block[keyLength:]
		if currentKey == key {
			// We found a value, return it
			//fmt.Printf("Found a value after doing %v linear skips\n", elementsSkipped)

			value := string(block[:valueLength])
			return value, true, nil
		} else if currentKey > key {
			return "", false, nil
		}
		// keep going with linear scan

		block = block[valueLength:]
		elementsSkipped++
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

	// Returns the Item the iterator is currently` pointing to. Assumes Valid() == true.
	Item() Item
}
