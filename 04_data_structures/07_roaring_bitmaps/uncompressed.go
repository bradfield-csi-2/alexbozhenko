package bitmap

import (
	"go.uber.org/zap"
)

const wordSize = 64

type uncompressedBitmap struct {
	data []uint64
}

var logger, _ = zap.NewProduction()
var sugar = logger.Sugar()

func newUncompressedBitmap() *uncompressedBitmap {
	// Range for uint32 is [0; 2**32-1], so we need
	// 2**32 bits, or 2**29 bytes, or 2**26 = 67108864 8-bytes blocks
	return &uncompressedBitmap{
		data: []uint64{},
		//	data: make([]uint64, 1<<26),
	}
}

func (b *uncompressedBitmap) Get(x uint32) bool {
	blockNumber := x / wordSize
	dataLen := uint32(len(b.data))
	sugar.Debugw("set", "dataLen", dataLen, "blockNumber", blockNumber)
	if blockNumber >= dataLen {
		return false
	}
	bitPositionInBlock := wordSize - (x % wordSize) - 1
	return (b.data[blockNumber]&(1<<bitPositionInBlock) != 0)
}

func (b *uncompressedBitmap) Set(x uint32) {
	blockNumber := x / wordSize
	dataLen := uint32(len(b.data))
	sugar.Debugw("set", "dataLen", dataLen, "blockNumber", blockNumber)
	if blockNumber >= dataLen {
		emptyBlocks := make([]uint64, blockNumber+1-dataLen)
		b.data = append(b.data, emptyBlocks...)
	}
	bitPositionInBlock := wordSize - (x % wordSize) - 1
	b.data[blockNumber] = b.data[blockNumber] | (1 << bitPositionInBlock)
}

func (b *uncompressedBitmap) Union(other *uncompressedBitmap) *uncompressedBitmap {
	var data []uint64

	return &uncompressedBitmap{
		data: data,
	}
}

func (b *uncompressedBitmap) Intersect(other *uncompressedBitmap) *uncompressedBitmap {
	var data []uint64
	return &uncompressedBitmap{
		data: data,
	}
}
