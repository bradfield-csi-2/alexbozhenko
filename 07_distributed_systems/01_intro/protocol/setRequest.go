package protocol

import (
	"bytes"
	"encoding/binary"
)

const SET_PROTOCOL = 1

type SetRequest struct {
	Key   []byte
	Value []byte
}

// +------------+------------+------------+----------+------------+----------+
// |            |            |            |          |            |          |
// | Reserved   |  Protocol  |  Length    |  Key     |  Length    |  Value   |
// | (1 byte)   |  version   |  (UVarint) |  (bytes) |  (UVarint) |  (bytes) |
// |            |  (1 byte)  |            |          |            |          |
// |            |            |            |          |            |          |
// +------------+------------+------------+----------+------------+----------+
func (req *SetRequest) Encode() []byte {
	result := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)
	result = append(result, 0x00) //reserved
	result = append(result, SET_PROTOCOL)
	appendEncodedItemWithLength(&result, buf, req.Key)
	appendEncodedItemWithLength(&result, buf, req.Value)
	return result
}

func (req *SetRequest) Decode(reqBytes []byte) error {
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
	req.Key = key
	value, err := decodeItem(reqBytes, reader)
	if err != nil {
		return err
	}
	req.Value = value

	return nil
}
