package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"hash/fnv"
	"os"
)

func countTrailingZeroes(b []byte) (i int) {
	var currentByte byte
	for numByte := len(b) - 1; numByte >= 0; numByte-- {
		currentByte = b[numByte]
		if currentByte != 0 {
			break
		}
		i += 8
	}
	for ; currentByte != 0 && currentByte&1 == 0; currentByte >>= 1 {
		i++
	}

	return
}
func byteSliceToBinString(b []byte) (res string) {
	for _, n := range b {
		//TODO fix format
		res = res + fmt.Sprintf("% 08b", n)
	}
	return
}

func main() {
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fnvHash := fnv.New32()
	maxZeroLength := 0
	log.SetLevel(log.TraceLevel)
	for scanner.Scan() {
		word := scanner.Bytes()
		fnvHash.Write(word)
		wordHash := fnvHash.Sum([]byte{})
		zeroes := countTrailingZeroes(wordHash)
		log.Tracef("word=%s hash=%s zeroes=%d\n", string(word), byteSliceToBinString(wordHash), zeroes)
		fnvHash.Reset()
		if zeroes > maxZeroLength {
			maxZeroLength = zeroes
		}
	}
	fmt.Println("Max zeroes seen=", maxZeroLength)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
