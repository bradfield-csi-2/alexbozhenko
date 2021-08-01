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

func TestBitmap(t *testing.T) {
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
			t.Fatalf("Get should've returned %t for %d\n", ok, x)
		}
	}

	/* TODO: Uncomment this section once you get Union and Intersect working.
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
			t.Fatalf("Union: Get should've returned %t for %d\n", ok1 || ok2, x)
		}
		if (ok1 && ok2) != intersect.Get(x) {
			t.Fatalf("Intersect: Get should've returned %t for %d\n", ok1 && ok2, x)
		}
	}
	*/

	/* TODO: Uncomment this section once you get compression / decompression working
	compressed := compress(b1)
	t.Logf("Uncompressed size: %d words, compressed size: %d words\n", len(b1.data), len(compressed))
	b := decompress(compressed)
	for x := uint32(0); x < start+limit+wordSize; x++ {
		if b1.Get(x) != b.Get(x) {
			t.Fatalf("Compression then decompression produced inconsistent result for %d\n", x)
		}
	}
	*/
}
