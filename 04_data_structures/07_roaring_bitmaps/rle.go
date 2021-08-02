package bitmap

import (
	"encoding/binary"
)

// ==========
// First idea
// ==========
// Each byte contains either compressed or uncompressed data:
// Bytes starting with 0 are uncompressed bytes, with
// leading 0 followed by sequence of literal 0s and 1s,
// 01001011 1000100
//     ^^^^ ^^^^^^^ - up to 12 bits of uncompressed data
//  ^^^             - length of uncompressed data, [1; 15]
// ^                - uncompressed byte flag

// In our scheme we start to make profit only if run is longer than 8,
// thus we can use bits encoding the length to encode values starting from 8
// Bytes starting with 1 are compressed bytes:
// 10100001
//        ^ - value of the run
//  ^^^^^^  - length of the run, starting from 2^3 to to 2^9 - 1 = 511
// ^        - compressed byte flag

// Disadvantages:
// annoying to work with individual bits
//

// ===========
// Second idea
// ===========
// Let's work with whole bytes instead:

// 0_______
//  ^^^^^^^ - length of followed uncompressed bytes[0;127]
// ^        - uncompressed flag

// 1_______
//   ^^^^^^ - length of the run in bytes [0;63]
//  ^       - value, of the run, 0 or 1
// ^        - compressed flag

type state uint

const (
	uncompressed state = iota
	compressed0
	compressed1
)
const maxUncompressedRunLength = 127 // 7 bits
const maxCompressedRunLength = 63    // 6 bits

func getState(b byte) state {
	if b == 0x00 {
		return compressed0
	}
	if b == 0xff {
		return compressed1
	}
	return uncompressed
}

func genHeaderByte(st state, runLength int) byte {
	if st == uncompressed {
		if runLength > maxUncompressedRunLength {
			// TODO abozhenko for robot-dreams:
			// If I know that runLength for uncompressed must be in range [0;127]
			// and length for compressed byte must be in range [0;63],
			// how do I enforce these rules at compile time?
			// In oCaml that can be solved with custom type
			panic("BUG! run length of uncompressed bytes is encoded with just 6 bits")
		}
		return byte(runLength)
	}
	if runLength > maxCompressedRunLength {
		panic("BUG! run length of compressed bytes is encoded with just 5 bits")
	}
	// TODO abozhenko for robot-dreams
	// In ocaml, we can do pattern matching, e.g:
	// https://www2.lib.uchicago.edu/keith/ocaml-class/pattern-matching.html
	// Here, if we introduce a new const, e.g compressedA, nothing would warn
	// us that we forget to add corresponding if statement
	// In ocaml I would introdcue custom enum type,
	// and compiler would warn us about Non-Exhaustive Pattern-Matches
	// How to do something like this in go, to make the code more bulletproof
	// against future modifications like this?
	if st == compressed0 {
		return byte(runLength) | (0b1 << 7)
	} else {
		return byte(runLength) | (0b11 << 6)
	}
}

func compress(b *uncompressedBitmap) []uint64 {
	compressedBytes := []byte{}
	currentRunBytes := []byte{}
	currentBlockBytes := make([]byte, 8)
	var runLength uint8 = 0 // length of consequent bytes of the same type
	currentState := uncompressed
	for _, block := range b.data {
		binary.BigEndian.PutUint64(currentBlockBytes, block)
		for _, b := range currentBlockBytes {
			// When current byte type differs from the previous one,
			// that means end of the run, so we need to write the header byte,
			// followed by the bytes of the run itself.
			// Also, if we accumulated more than enough bytes in the current run,
			// let's create a new run, with new header byte
			if getState(b) != currentState ||
				(currentState == uncompressed && runLength > maxUncompressedRunLength) ||
				(currentState != uncompressed && runLength > maxCompressedRunLength) {
				compressedBytes = append(compressedBytes,
					genHeaderByte(currentState, int(runLength)))
				compressedBytes = append(compressedBytes, currentRunBytes...)
				currentRunBytes = []byte{}
				currentState = getState(b)
				runLength = 0
			}
			runLength++
			currentRunBytes = append(currentRunBytes, b)
		}
	}
	compressedBytes = append(compressedBytes,
		genHeaderByte(currentState, int(runLength)))
	compressedBytes = append(compressedBytes, currentRunBytes...)

	// append 0-bytes to align to uint64
	compressedBytes = append(compressedBytes, make([]byte, len(compressedBytes)%8)...)
	compressedData := make([]uint64, len(compressedBytes)/8)
	for i := range compressedData {
		compressedData[i] = binary.BigEndian.Uint64(compressedBytes[8*i : 8*i+8])
	}

	return compressedData
}

func decompress(compressed []uint64) *uncompressedBitmap {
	var data []uint64
	return &uncompressedBitmap{
		data: data,
	}
}
