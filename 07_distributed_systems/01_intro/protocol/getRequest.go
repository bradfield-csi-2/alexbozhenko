package protocol

import (
	"bytes"
	"encoding/binary"
)

const GET_PROTOCOL = 1

type GetRequest struct {
	Key []byte
}

// +------------+------------+------------+----------+
// |            |            |            |          |
// | Reserved   |  Protocol  |  Length    |  Key     |
// | (1 byte)   |  version   |  (UVarint) |  (bytes) |
// |            |  (1 byte)  |            |          |
// |            |            |            |          |
// +------------+------------+------------+----------+
func (req *GetRequest) Encode() []byte {
	// allocating maximum possible size according to the
	// structure above
	//result := make([]byte, 1+1+binary.MaxVarintLen64+len(req.key))
	result := make([]byte, 0)
	buf := make([]byte, binary.MaxVarintLen64)
	result = append(result, 0x00) //reserved
	result = append(result, GET_PROTOCOL)
	appendEncodedItemWithLength(&result, buf, req.Key)
	return result
}

func (req *GetRequest) Decode(reqBytes []byte) error {
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

	return nil
}
