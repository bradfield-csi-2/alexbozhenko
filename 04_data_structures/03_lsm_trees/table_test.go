package table

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

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
		value := randomWord(10, 20)
		result[i] = Item{key, value}
	}
	return result
}

func TestTable(t *testing.T) {
	dir, err := ioutil.TempDir("", "table")
	if err != nil {
		t.Fatal(err)
	}
	// Clean up temp directory at end of test; you can remove this for debugging.
	defer os.RemoveAll(dir)

	tmpfile := filepath.Join(dir, "tmpfile")

	n := 1000
	sortedItems := generateSortedItems(n)

	toInclude := sortedItems[:n/2]
	toExclude := sortedItems[n/2:]

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
