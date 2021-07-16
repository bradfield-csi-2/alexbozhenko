package main

import (
	"log"
	"testing"
)

func Benchmark_myBloomFilter_add(b *testing.B) {
	words, err := loadWords(wordsPath)
	if err != nil {
		log.Fatal(err)
	}
	var bFilter bloomFilter = newMyBloomFilter()
	b.ResetTimer()

	for j := 0; j < b.N; j++ {
		for i := 0; i < len(words); i += 2 {
			bFilter.add(words[i])
		}
	}
}

func Benchmark_myBloomFilter_maybeContains(b *testing.B) {
	words, err := loadWords(wordsPath)
	if err != nil {
		log.Fatal(err)
	}
	var bFilter bloomFilter = newMyBloomFilter()
	for i := 0; i < len(words); i += 2 {
		bFilter.add(words[i])
	}
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		for i := 0; i < len(words); i += 2 {
			bFilter.maybeContains(words[i])
		}
	}

}
