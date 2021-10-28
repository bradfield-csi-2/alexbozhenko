package protocol

const getVersion = 1

type getRequest struct {
	key []byte
}

func (req *getRequest) Encode() []byte {
	return nil
}

func (req *getRequest) Decode(reqBytes []byte) error {
	return nil
}
