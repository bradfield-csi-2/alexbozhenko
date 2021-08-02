package bitmap

import (
	"math/rand"
	"testing"
)

const (
	start = 1000 * 1000
	limit = 100 * 1000
	items = 10 * 1000
)

func TestBitmap(b *testing.T) {
	//func BenchmarkBitmap(b *testing.B) {
	//for benchNumber := 0; benchNumber < b.N; benchNumber++ {
	b1 := newUncompressedBitmap()
	m1 := make(map[uint32]struct{})

	// Call Set a bunch
	for i := 0; i < items; i++ {
		x := start + uint32(rand.Intn(limit))
		b1.Set(x)
		m1[x] = struct{}{}
	}

	// Make sure subsequent Get works as expected
	for x := uint32(0); x < start+limit+wordSize; x++ {
		_, ok := m1[x]
		if ok != b1.Get(x) {
			b.Fatalf("Get should've returned %t for %d\n", ok, x)
		}
	}

	b2 := newUncompressedBitmap()
	m2 := make(map[uint32]struct{})

	// Call Set a bunch
	for i := 0; i < items; i++ {
		x := uint32(rand.Intn(limit))
		b2.Set(x)
		m2[x] = struct{}{}
	}

	union := b1.Union(b2)
	intersect := b1.Intersect(b2)
	for x := uint32(0); x < start+limit+wordSize; x++ {
		_, ok1 := m1[x]
		_, ok2 := m2[x]
		if (ok1 || ok2) != union.Get(x) {
			b.Fatalf("Union: Get should've returned %t for %d\n", ok1 || ok2, x)
		}
		if (ok1 && ok2) != intersect.Get(x) {
			b.Fatalf("Intersect: Get should've returned %t for %d\n", ok1 && ok2, x)
		}
	}

	compressed := compress(b1)
	b.Logf("Uncompressed size: %d words, compressed size: %d words\n", len(b1.data), len(compressed))
	b_ := decompress(compressed)
	for x := uint32(0); x < start+limit+wordSize; x++ {
		if b1.Get(x) != b_.Get(x) {
			b.Fatalf("Compression then decompression produced inconsistent result for %d\n", x)
		}
	}
	//	}
}
