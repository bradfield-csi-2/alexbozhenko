package protocol

import (
	"reflect"
	"testing"
)

func Test_getRequest_Encode_and_Decode(t *testing.T) {
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
			req := &GetRequest{
				Key: tt.key,
			}
			var got []byte
			if got = req.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRequest.Encode() = %x, want %x", got, tt.want)
			}
			var roundTripResult GetRequest
			roundTripResult.Decode(got)
			if !reflect.DeepEqual(roundTripResult.Key, tt.key) {
				t.Errorf("getRequest.Decode() = %x, want %x", roundTripResult.Key, tt.key)
			}
		})
	}
}

func Test_getRequest_Decode(t *testing.T) {
	tests := []struct {
		name     string
		reqBytes []byte
		wantKey  []byte
		wantErr  bool
	}{
		{
			name:     "valid",
			reqBytes: []byte{0x00, 0x01, 0x04, 0x61, 0x4b, 0x65, 0x79},
			wantKey:  []byte("aKey"),
			wantErr:  false,
		},
		{
			name:     "Invalid message returns error",
			reqBytes: []byte{0x00, 0x01, 0x04, 0x61, 0x4b, 0x65},
			wantKey:  []byte("aKey"),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &GetRequest{
				Key: tt.wantKey,
			}
			if err := req.Decode(tt.reqBytes); (err != nil) != tt.wantErr {
				t.Errorf("getRequest.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
