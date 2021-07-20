package bloom

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"unsafe"
)

// ...We will refer to a Bloom filter with
// k hashes,
// m bits in the filter,
// and n elements that have been inserted.

// Choose a ballpark value for n(number of elements). 123000
// Choose a value for m. 100_000 bytes = 800_000 bits
// Calculate the optimal value of k =2
// Calculate the error rate  = 0.069
// https://www.wolframalpha.com/input/?i=n+%3D+123115%3B+m+%3D+800000%3B+k%3D2%3B+%281-2.7%5E%28-k*n%2Fm%29%29%5Ek

const memoryUsageInBits = 80_000_000 // 10MB
const numberOfHashFunctions = 9
const bitsInBucket = 64

var bitPositions = make([]uint64, numberOfHashFunctions)

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type MyBloomFilter struct {
	data []uint64
}

func NewMyBloomFilter() *MyBloomFilter {
	return &MyBloomFilter{
		data: make([]uint64, int64(math.Ceil(
			float64(memoryUsageInBits)/float64(bitsInBucket)))),
	}
}

func (bloomFilter *MyBloomFilter) String() string {
	res := ""
	for i, bucket := range bloomFilter.data {
		if bucket != 0 {
			res += fmt.Sprintf("bucket %v=%064b\n", i, bucket)
		}
	}
	return res
}

// Produce multiple hash functions using the trick described in
// http://willwhim.wpengine.com/2011/09/03/producing-n-hash-functions-by-hashing-only-once/
// hash(i) = (a + b * i ) % m
func getBitPositions(item string) []uint64 {
	fnvHash := fnv.New64()

	// The hack below was copied from the solution:
	// next time I should be smart enough to come up with it myself:)
	// since I saw in pprof that stringtobyteslice takes 27% of total cpu time

	// The basic way to convert `item` into a byte array is as follows:
	//
	//     bytes := []byte(item)
	//
	// However, this `unsafe.Pointer` hack lets us avoid copying data and
	// get a slight speed-up by directly referencing the underlying bytes
	// of `item`.

	fnvHash.Write(*(*[]byte)(unsafe.Pointer(&item)))
	fnvSum := fnvHash.Sum64()
	fnvSum1 := fnvSum & 0xffffffff
	fnvSum2 := fnvSum >> 32
	for i := uint64(0); i < numberOfHashFunctions; i++ {
		bitPositions[i] = (fnvSum1 + fnvSum2*i) % memoryUsageInBits
	}
	return bitPositions
}

func (b *MyBloomFilter) Add(item string) {
	bitPositions := getBitPositions(item)
	for _, bitPosition := range bitPositions {
		bucketNumber := bitPosition / bitsInBucket
		bitNumberInBucket := (bitsInBucket - 1) - (bitPosition % bitsInBucket)
		b.data[bucketNumber] |= (1 << bitNumberInBucket)
	}
}

func (b *MyBloomFilter) MaybeContains(item string) bool {
	bitPositions := getBitPositions(item)
	for _, bitPosition := range bitPositions {
		bucketNumber := bitPosition / bitsInBucket
		bitNumberInBucket := (bitsInBucket - 1) - (bitPosition % bitsInBucket)
		if (b.data[bucketNumber] & (1 << bitNumberInBucket)) == 0 {
			return false
		}
	}
	return true
}

func (b *MyBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
