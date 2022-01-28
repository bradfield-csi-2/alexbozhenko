package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"hash/fnv"
	"math"
	"os"
	"os/exec"
)

func countTrailingZeroes(b []byte) (i uint) {
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
		res = res + fmt.Sprintf("%08b ", n)
	}
	return
}

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Sprintf("Usage: ./%s file_path", os.Args[0]))
	}
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// key of the map is first eleven bits of the hash
	// value is the max number consequent trailing zeros seen
	// in this bucket
	LogLogMap := make(map[[2]byte]uint)
	scanner := bufio.NewScanner(file)
	fnvHash := fnv.New32()
	maxZeroLength := uint(0)
	for scanner.Scan() {
		word := scanner.Bytes()
		_, err := fnvHash.Write(word)
		if err != nil {
			panic(err)
		}
		wordHash := fnvHash.Sum([]byte{})
		bucket := [2]byte{wordHash[0], wordHash[1] & 0b11100000}
		zeroes := countTrailingZeroes(wordHash)
		zeroesInBucket, ok := LogLogMap[bucket]
		if !ok || zeroesInBucket < zeroes {
			LogLogMap[bucket] = zeroes
		}

		log.Tracef("word=%s hash=%s zeroes=%d\n", string(word), byteSliceToBinString(wordHash), zeroes)
		fnvHash.Reset()
		if zeroes > maxZeroLength {
			maxZeroLength = zeroes
		}
	}
	avgLogLogZeroes := 0.0
	for _, zeroes := range LogLogMap {
		avgLogLogZeroes += float64(zeroes)
	}
	avgLogLogZeroes /= float64(len(LogLogMap))

	log.Tracef("LogLogMap(%d)=%v", len(LogLogMap), LogLogMap)

	wcCardinality, err := exec.Command("wc", "-l", filepath).Output()
	fmt.Printf("Actual cardianlity         = %s", string(wcCardinality))
	log.Tracef("Max zeroes seen=", maxZeroLength)
	fmt.Printf("2^n cardinality estimate   = %d\n", 1<<maxZeroLength)
	log.Tracef("Average LogLog zeroes=", avgLogLogZeroes)
	fmt.Printf("LogLog cardinality estimate= %.0f\n", 0.79402*2048*math.Pow(2, avgLogLogZeroes))
	fmt.Println("=========================")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
