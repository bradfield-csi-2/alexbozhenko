package protocol

import (
	"reflect"
	"testing"
)

func Test_getRequest_Encode(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		want []byte
	}{
		{
			name: "Empty",
			key:  []byte{},
			want: []byte{0x00, 0x01, 0x00},
		},
		{
			name: "Non-empty",
			key:  []byte("aKey"),
			want: []byte{0x00, 0x01, 0x04, 0x61, 0x4b, 0x65, 0x79},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &getRequest{
				key: tt.key,
			}
			var got []byte
			if got = req.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRequest.Encode() = %x, want %x", got, tt.want)
			}
			var roundTripResult getRequest
			roundTripResult.Decode(got)
			if !reflect.DeepEqual(roundTripResult.key, tt.key) {
				t.Errorf("getRequest.Decode() = %x, want %x", roundTripResult.key, tt.key)
			}
		})
	}
}
