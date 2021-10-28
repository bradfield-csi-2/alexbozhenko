package protocol

const SET_PROTOCOL = 1

type setRequest struct {
	key []byte
}

func (req *setRequest) Encode() []byte {
	return nil
}

func (req *setRequest) Decode(reqBytes []byte) error {
	return nil
}
