package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const wordsPath = "/usr/share/dict/words"

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

func main() {
	words, err := loadWords(wordsPath)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	// TODO: Replace trivialBloomFilter with your own implementation
	var b bloomFilter = newTrivialBloomFilter()

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

	falsePositives := 0
	numChecked := 0

	// None of the words at odd indices were added, so whenever
	// maybeContains returns true, it's a false positive
	for i := 1; i < len(words); i += 2 {
		if b.maybeContains(words[i]) {
			falsePositives++
		}
		numChecked++
	}

	falsePositiveRate := float64(falsePositives) / float64(numChecked)

	fmt.Printf("Elapsed time: %s\n", time.Since(start))
	fmt.Printf("Memory usage: %d bytes\n", b.memoryUsage())
	fmt.Printf("False positive rate: %0.2f%%\n", 100*falsePositiveRate)
}
