package bitmap

import (
	"encoding/binary"
	"fmt"
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
//  ^^^^^^^ - length of followed uncompressed bytes[1;127]
// ^        - uncompressed flag

// 1_______
//   ^^^^^^ - length of the run in bytes [1;63]
//  ^       - value, of the run, 0 or 1
// ^        - compressed flag

type state uint

const (
	compressed0 state = iota
	compressed1
	uncompressed
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

func genHeaderByte(st state, runLength uint8) byte {
	if runLength == 0 {
		panic("BUG! run length must be non-zero to produce valid encoding")
	}
	if st == uncompressed {
		if runLength > maxUncompressedRunLength {
			// TODO abozhenko for robot-dreams:
			// If I know that runLength for uncompressed must be in range [1;127]
			// and length for compressed byte must be in range [1;63],
			// how do I enforce these rules at compile time?
			// In oCaml that can be solved with custom type
			panic(fmt.Sprintf("BUG! run length(%v) of uncompressed bytes must be encoded with just 6 bits",
				runLength))
		}
		return byte(runLength)
	}
	if runLength > maxCompressedRunLength {
		panic(fmt.Sprintf("BUG! run length(%v) of compressed bytes must be encoded with just 5 bits",
			runLength))
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
	//atom.SetLevel(zapcore.DebugLevel)
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
			if getState(b) != currentState && runLength > 0 ||
				(currentState == uncompressed && runLength >= maxUncompressedRunLength) ||
				(currentState != uncompressed && runLength >= maxCompressedRunLength) {

				sugar.Debugw(
					"Detected new run",
					"previous runLength", runLength,
					"previous state", currentState,
					"new state", getState(b),
					"len compressedBytes", len(compressedBytes),
				)
				compressedBytes = append(compressedBytes,
					genHeaderByte(currentState, runLength))
				compressedBytes = append(compressedBytes, currentRunBytes...)
				currentRunBytes = []byte{}
				currentState = getState(b)
				runLength = 0
			}
			runLength++
			if currentState == uncompressed {
				currentRunBytes = append(currentRunBytes, b)
			}
		}
	}
	compressedBytes = append(compressedBytes,
		genHeaderByte(currentState, runLength))
	compressedBytes = append(compressedBytes, currentRunBytes...)

	// append 0-bytes to align to uint64
	sugar.Debugw(
		"before adding 0-bytes to align to uint64",
		"len was", len(compressedBytes),
		"len/8", len(compressedBytes)/8,
		"len%8", len(compressedBytes)%8,
		"adding zero bytes: ", len(compressedBytes)%8,
	)
	compressedBytes = append(compressedBytes, make([]byte, 8-len(compressedBytes)%8)...)
	compressedData := make([]uint64, len(compressedBytes)/8)
	for i := range compressedData {
		compressedData[i] = binary.BigEndian.Uint64(compressedBytes[8*i : 8*i+8])
	}
	sugar.Debugw(
		"Finished compression",
		"compressed len",
		len(compressedData),
		"uncompressed len",
		len(b.data),
	)

	return compressedData
}

func parseCompressedHeader(header byte) (st state, runLength uint8) {
	compressedMark := header & (1 << 7)
	if compressedMark != 0 {
		st = state((header & (1 << 6)) >> 6)
	} else {
		st = uncompressed
	}
	if st == uncompressed {
		runLength = uint8(header & 0b01111111)
	} else {
		runLength = uint8(header & 0b00111111)
	}
	return
}

func decompress(compressed []uint64) *uncompressedBitmap {

	compressedBytesBuffer := make([]byte, 8)
	uncompressedBytes := []byte{}
	var runLength uint8 = 0
	var byteToAppend byte
	var st state
	for _, block := range compressed {
		binary.BigEndian.PutUint64(compressedBytesBuffer, block)
		for _, b := range compressedBytesBuffer {
			// runLength == 0 means that we are at the header byte
			if runLength == 0 {
				if b == 0 {
					// if we reached zero-bytes we had added to the end
					// to align to uint64, we need to stop processing
					break
				}
				st, runLength = parseCompressedHeader(b)
				if st == uncompressed {
					continue
				}
			}
			if st == uncompressed {
				uncompressedBytes = append(uncompressedBytes, b)
				runLength--
			} else {
				if st == compressed0 {
					byteToAppend = 0b00000000
				} else {
					byteToAppend = 0b11111111
				}
				for i := uint8(0); i < runLength; i++ {
					uncompressedBytes = append(uncompressedBytes, byteToAppend)
				}
				runLength = 0
			}
		}
	}
	data := make([]uint64, len(uncompressedBytes)/8)
	for i := range data {
		data[i] = binary.BigEndian.Uint64(uncompressedBytes[8*i : 8*i+8])
	}

	return &uncompressedBitmap{
		data: data,
	}
}
