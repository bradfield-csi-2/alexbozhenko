package protocol

import (
	"bytes"
	"encoding/binary"
)

const SET_PROTOCOL = 1

type setRequest struct {
	key   []byte
	value []byte
}

// +------------+------------+------------+----------+------------+----------+
// |            |            |            |          |            |          |
// | Reserved   |  Protocol  |  Length    |  Key     |  Length    |  Value   |
// | (1 byte)   |  version   |  (UVarint) |  (bytes) |  (UVarint) |  (bytes) |
// |            |  (1 byte)  |            |          |            |          |
// |            |            |            |          |            |          |
// +------------+------------+------------+----------+------------+----------+
func (req *setRequest) Encode() []byte {
	result := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)
	result = append(result, 0x00) //reserved
	result = append(result, SET_PROTOCOL)
	appendEncodedItemWithLength(&result, buf, req.key)
	appendEncodedItemWithLength(&result, buf, req.value)
	return result
}

func (req *setRequest) Decode(reqBytes []byte) error {
	reader := bytes.NewReader(reqBytes)
	_, err := reader.ReadByte() // ignoring the reserved byte
	if err != nil {
		return err
	}
	_, err = reader.ReadByte() // ignoring protocolVersion for now
	if err != nil {
		return err
	}
	key, err := decodeItem(reqBytes, reader)
	if err != nil {
		return err
	}
	req.key = key
	value, err := decodeItem(reqBytes, reader)
	if err != nil {
		return err
	}
	req.value = value

	return nil
}
