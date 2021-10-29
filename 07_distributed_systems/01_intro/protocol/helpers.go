package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// function appends uvarint-encoded length of item and item itself to *b slice
// buf is a buffer for encoding uvarint
func appendEncodedItemWithLength(b *[]byte, buf, item []byte) {
	nBytes := binary.PutUvarint(buf, uint64(len(item)))
	*b = append(*b, buf[:nBytes]...)
	*b = append(*b, item...)
}

// function decodes uvarint-encoded length from current position of r
// and returns next itemLength bytes of b
// (r should be the reader of b[[]])
// passing both b and r to function is rather dirty,
// but it saves allocations.
func decodeItem(b []byte, r *bytes.Reader) ([]byte, error) {
	itemLength, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	itemStartOffset, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	if len(b)-int(itemStartOffset) < int(itemLength) {
		return nil, fmt.Errorf("expected message of len %v, but got only %v bytes",
			itemLength, r.Size()-itemStartOffset)
	}
	//not using reader io.ReadFull() to avoid copying
	r.Seek(int64(itemLength), io.SeekCurrent)
	return b[itemStartOffset : itemStartOffset+int64(itemLength)], nil
}
