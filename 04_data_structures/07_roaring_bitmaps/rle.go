package bitmap

func compress(b *uncompressedBitmap) []uint64 {
	return nil
}

func decompress(compressed []uint64) *uncompressedBitmap {
	var data []uint64
	return &uncompressedBitmap{
		data: data,
	}
}
