package main

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

const memoryUsageInBits = 150_000
const numberOfHashFunctions = 2
const bitsInBucket = 64

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type myBloomFilter struct {
	data []uint64
}

func newMyBloomFilter() *myBloomFilter {
	// ...We will refer to a Bloom filter with
	// k hashes,
	// m bits in the filter,
	// and n elements that have been inserted.

	// Choose a ballpark value for n(number of elements). 123000
	// Choose a value for m. 100_000 bytes = 800_000 bits
	// Calculate the optimal value of k =2
	// Calculate the error rate  = 0.069
	// https://www.wolframalpha.com/input/?i=n+%3D+123115%3B+m+%3D+800000%3B+k%3D2%3B+%281-2.7%5E%28-k*n%2Fm%29%29%5Ek

	return &myBloomFilter{
		data: make([]uint64, memoryUsageInBits/8),
	}
}

func (bloomFilter *myBloomFilter) String() string {
	res := ""
	for i, bucket := range bloomFilter.data {
		if bucket != 0 {
			res += fmt.Sprintf("bucket %v=%064b\n", i, bucket)
		}
	}
	return res
}

func getBitPositions(item string) []uint {
	fnvHash := fnv.New64()
	fnvHash.Write([]byte(item))
	fnvSum := fnvHash.Sum64()
	fnvSum1 := fnvSum & 0xffffffff
	fnvSum2 := fnvSum >> 32
	bitPositions := make([]uint, numberOfHashFunctions)
	// http://willwhim.wpengine.com/2011/09/03/producing-n-hash-functions-by-hashing-only-once/
	// hash(i) = (a + b * i ) % m
	for i := 0; i < numberOfHashFunctions; i++ {
		bitPositions[i] = uint((fnvSum1 + fnvSum2*uint64(i)) % memoryUsageInBits)
	}
	//fmt.Println("bitPositions", bitPositions)
	return bitPositions
}

func (b *myBloomFilter) add(item string) {
	//fmt.Println("Adding", item)
	bitPositions := getBitPositions(item)
	for _, bitPosition := range bitPositions {
		bucketNumber := bitPosition / bitsInBucket
		bitNimberInBucket := (bitsInBucket - 1) - (bitPosition % 8)
		//		fmt.Println("bit number", bitNimberInBucket)
		//		fmt.Println("setting this value ", uint64(1<<bitNimberInBucket))
		b.data[bucketNumber] |= (1 << bitNimberInBucket)
		//fmt.Println("now value is", b.data[bucketNumber])
	}

}

func (b *myBloomFilter) maybeContains(item string) bool {
	//fmt.Println("Getting", item)
	//fmt.Println("data is:", b)
	bitPositions := getBitPositions(item)

	for _, bitPosition := range bitPositions {
		bucketNumber := bitPosition / bitsInBucket
		bitNimberInBucket := (bitsInBucket - 1) - (bitPosition % 8)
		if (b.data[bucketNumber] & (1 << bitNimberInBucket)) == 0 {
			return false
		}
	}
	return true
}

func (b *myBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
