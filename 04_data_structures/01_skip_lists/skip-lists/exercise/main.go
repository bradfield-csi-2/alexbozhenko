package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

const (
	path     = "/usr/share/dict/words"
	pattern  = "random"
	startKey = "assembly"
	endKey   = "golang"
	stride   = 2
	limit    = 20000
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to file")

func loadWords(limit int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var result []string
	s := bufio.NewScanner(f)
	for s.Scan() && (limit == 0 || len(result) < limit) {
		result = append(result, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

var expectedRangeScanItems int = -1

func runTest(words []string, o OC, name string) {
	fmt.Printf("%-25s", name)

	// Puts
	start := time.Now()
	for _, word := range words {
		ok := o.Put(word, word)
		if !ok {
			log.Fatalf("Word %q shouldn't have been added yet\n", word)
		}
	}
	fmt.Printf("%-20s", time.Since(start))

	// Deletes
	start = time.Now()
	for i := 0; i < len(words); i += stride {
		if !o.Delete(words[i]) {
			log.Fatalf("Couldn't delete %q (not found)\n", words[i])
		}
	}
	fmt.Printf("%-20s", time.Since(start))

	// Gets
	start = time.Now()
	for i, word := range words {
		value, ok := o.Get(word)
		if i%stride == 0 && ok {
			log.Fatalf("Expected deleted word %q to be gone\n", word)
		} else if i%stride != 0 {
			if !ok {
				log.Fatalf("Couldn't find word %q\n", word)
			} else if word != value {
				log.Fatalf("Got wrong value for key %q; expected %q, got %q\n", word, word, value)
			}
		}
	}
	fmt.Printf("%-20s", time.Since(start))

	// RangeScan
	start = time.Now()
	prev := ""
	count := 0
	for iter := o.RangeScan(startKey, endKey); iter.Valid(); iter.Next() {
		curr := iter.Key()
		if curr < prev {
			log.Fatalf("Iterator returned items out of order (%q came before %q)\n", prev, curr)
		}
		if !(startKey <= curr && curr <= endKey) {
			log.Fatalf("Iterator returned item out of range [%q, %q]: %q\n", startKey, endKey, curr)
		}
		prev = curr
		count++
	}
	if expectedRangeScanItems == -1 {
		expectedRangeScanItems = count
	} else if count != expectedRangeScanItems {
		fmt.Println()
		log.Fatalf("Inconsistent number of items from RangeScan: %d vs %d\n", expectedRangeScanItems, count)
	}
	fmt.Printf("%-20s\n", time.Since(start))
}

func main() {
	words, err := loadWords(0)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Set pattern to "sequential" to destroy the performance of the
	// (unbalanced) binary search tree.
	if pattern == "random" {
		rand.Seed(1)
		rand.Shuffle(len(words), func(i, j int) { words[i], words[j] = words[j], words[i] })
	}

	fmt.Printf("Testing using %d words (%s pattern)\n\n", limit, pattern)
	fmt.Printf("%-25s%-20s%-20s%-20s%-20s\n", "Name", "Puts", "Deletes", "Gets", "RangeScan")
	fmt.Printf("----------------------------------------------------------------------------------------------------\n")

	for _, testCase := range []struct {
		o    OC
		name string
	}{
		{newSliceOC(), "Slice"},
		{newLinkedOC(), "Linked List"},
		{newLinkedBlockOC(), "Linked Block"},
		{newBstOC(), "Binary Search Tree"},
		{newRbTreeOC(), "Red Black Tree"},
		{newSkipListOC(), "Skip List"},
	} {
		if len(words) > limit {
			words = words[:limit]
		}
		runTest(words, testCase.o, testCase.name)
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}
}
