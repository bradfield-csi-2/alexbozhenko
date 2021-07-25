package table

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"sort"
	"testing"
)

const N_WORDS = 80_000

// min and max are inclusive.
func randomWord(min, max int) string {
	n := min + rand.Intn(max-min+1)
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		c := rune(rand.Intn(26))
		buf.WriteRune('a' + c)
	}
	return buf.String()
}

// All items are guaranteed to have unique keys.
func generateSortedItems(n int) []Item {
	m := make(map[string]struct{}, n)
	for len(m) < n {
		key := randomWord(8, 16)
		m[key] = struct{}{}
	}
	keys := make([]string, 0, n)
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]Item, n)
	for i, key := range keys {
		value := randomWord(10, 4100)
		result[i] = Item{key, value}
	}
	return result
}

func TestTable(t *testing.T) {
	dir, err := ioutil.TempDir("", "table")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Building in tmpdir=", dir)
	// Clean up temp directory at end of test; you can remove this for debugging.
	//defer os.RemoveAll(dir)

	tmpfile := filepath.Join(dir, "tmpfile")

	sortedItems := generateSortedItems(N_WORDS)

	toInclude := sortedItems[:N_WORDS/2]
	toExclude := sortedItems[N_WORDS/2:]

	err = Build(tmpfile, toInclude)
	if err != nil {
		t.Fatalf("Error building Table: %v", err)
	}

	table, err := LoadTable(tmpfile)
	if err != nil {
		t.Fatalf("Error loading Table: %v", err)
	}

	for _, item := range toInclude {
		actual, ok, err := table.Get(item.Key)
		if err != nil {
			t.Fatalf("Error performing point read for key %q: %v", item.Key, err)
		}
		if !ok {
			t.Fatalf("Expected key %q to exist", item.Key)
		}
		if actual != item.Value {
			t.Fatalf("Key %q: expected value %q, got %q instead", item.Key, item.Value, actual)
		}
	}

	for _, item := range toExclude {
		_, ok, err := table.Get(item.Key)
		if err != nil {
			t.Fatalf("Error performing point read for key %q: %v", item.Key, err)
		}
		if ok {
			t.Fatalf("Expected key %q not to exist", item.Key)
		}
	}

	// TODO: Uncomment the following to test RangeScan
	/*
		expectedScan := sortedItems[n/4 : n/3]
		startKey := expectedScan[0].Key
		endKey := expectedScan[len(expectedScan)-1].Key
		iter, err := table.RangeScan(startKey, endKey)
		if err != nil {
			t.Fatal(err)
		}
		actualScan := make([]Item, 0, len(expectedScan))
		for ; iter.Valid(); iter.Next() {
			actualScan = append(actualScan, iter.Item())
		}
		if !reflect.DeepEqual(expectedScan, actualScan) {
			t.Fatalf("Unexpected RangeScan result\n\nExpected: %v\n\nActual: %v", expectedScan, actualScan)
		}
	*/
}

func BenchmarkTable(b *testing.B) {

	sortedItems := generateSortedItems(N_WORDS)
	toInclude := sortedItems[:N_WORDS/2]
	toExclude := sortedItems[N_WORDS/2:]

	table, err := LoadTable("/home/alex/large_table")
	if err != nil {
		b.Fatalf("Error loading Table: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, item := range toInclude {
			table.Get(item.Key)
		}

		for _, item := range toExclude {
			table.Get(item.Key)
		}
	}
}
