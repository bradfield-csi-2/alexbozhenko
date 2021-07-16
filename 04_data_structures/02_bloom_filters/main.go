package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

const wordsPath = "/usr/share/dict/words"

// ngerman is bigger on my machine
// it is 3,9M, and 304k words
//const wordsPath = "/usr/share/dict/ngerman"

func loadWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	words, err := loadWords(wordsPath)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	start := time.Now()

	var b bloomFilter = newMyBloomFilter()
	falsePositives := 0
	numChecked := 0

	for numOfIterations := 0; numOfIterations < 1; numOfIterations++ {
		// Add every other word (even indices)
		for i := 0; i < len(words); i += 2 {
			b.add(words[i])
		}

		// Make sure there are no false negatives
		for i := 0; i < len(words); i += 2 {
			word := words[i]
			if !b.maybeContains(word) {
				log.Fatalf("false negative for word %q\n", word)
			}
		}

		// None of the words at odd indices were added, so whenever
		// maybeContains returns true, it's a false positive
		for i := 1; i < len(words); i += 2 {
			if b.maybeContains(words[i]) {
				falsePositives++
			}
			numChecked++
		}
	}

	falsePositiveRate := float64(falsePositives) / float64(numChecked)

	fmt.Printf("Elapsed time: %s\n", time.Since(start))
	fmt.Printf("Memory usage: %d bytes\n", b.memoryUsage())
	fmt.Printf("False positive rate: %0.2f%%\n", 100*falsePositiveRate)
}
