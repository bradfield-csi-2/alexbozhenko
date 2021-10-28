package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
)

const GET_PROTOCOL = 1

type getRequest struct {
	key []byte
}

// +------------+------------+------------+----------+
// |            |            |            |          |
// | Reserved   |  Protocol  |  Length    |  Key     |
// | (1 byte)   |  version   |  (UVarint) |  (bytes) |
// |            |  (1 byte)  |            |          |
// |            |            |            |          |
// +------------+------------+------------+----------+
func (req *getRequest) Encode() []byte {
	// allocating maximum possible size according to the
	// structure above
	//result := make([]byte, 1+1+binary.MaxVarintLen64+len(req.key))
	result := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)
	result = append(result, 0x00) //reserved
	result = append(result, GET_PROTOCOL)
	nBytes := binary.PutUvarint(buf, uint64(len(req.key)))
	result = append(result, buf[:nBytes]...)
	result = append(result, req.key...)
	return result
}

func (req *getRequest) Decode(reqBytes []byte) error {
	reader := bytes.NewReader(reqBytes)
	_, err := reader.ReadByte() // ignoring the reserved byte
	if err != nil {
		return err
	}
	_, err = reader.ReadByte() // ignoring protocolVersion for now
	if err != nil {
		return err
	}
	keyLength, err := binary.ReadUvarint(reader)
	if err != nil {
		return err
	}
	keyStartOffset, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	//not using reader io.ReadFull() to avoid copying
	req.key = reqBytes[keyStartOffset : keyStartOffset+int64(keyLength)]

	return nil
}
