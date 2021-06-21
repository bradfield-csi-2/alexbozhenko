package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
type CounterWithFilename struct {
	counter int
	filename string
}
*/

func main() {
	counts := make(map[string]map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, filenameCounterMap := range counts {
	        var count int
		for _, thisFileCount := range filenameCounterMap  {
			count += thisFileCount 
		}
		if count > 1 {
			fmt.Printf("%s\n", line )

			for filename, thisFileCount := range filenameCounterMap  {
				fmt.Printf("\t\t%s\t%d\n", filename, thisFileCount)
			}
		}
	}
}

func countLines(f *os.File, counts  map[string]map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		var counterWithFilename  map[string]int
		counterWithFilename = counts[input.Text()]
		if len(counterWithFilename) == 0 {
			counterWithFilename  = make(map[string]int)
		}
		counterWithFilename[f.Name()]++ 
		counts[input.Text()] = counterWithFilename 
	}
}
