package protocol

import (
	"reflect"
	"testing"
)

func Test_setRequest_Encode_and_Decode(t *testing.T) {
	tests := []struct {
		name  string
		key   []byte
		value []byte
		want  []byte
	}{
		{
			name:  "Empty",
			key:   []byte{},
			value: []byte{},
			want:  []byte{0x00, 0x01, 0x00, 0x00},
		},
		{
			name:  "Non-empty",
			key:   []byte("aKey"),
			value: []byte("aValue"),
			want:  []byte{0x00, 0x01, 0x04, 0x61, 0x4b, 0x65, 0x79, 0x06, 0x61, 0x56, 0x61, 0x6c, 0x75, 0x65},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &setRequest{
				key:   tt.key,
				value: tt.value,
			}
			var got []byte
			if got = req.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRequest.Encode() = %x, want %x", got, tt.want)
			}
			var roundTripResult setRequest
			roundTripResult.Decode(got)
			if !reflect.DeepEqual(roundTripResult.key, tt.key) {
				t.Errorf("setRequest.Decode().key = %x, want %x", roundTripResult.key, tt.key)
			}
			if !reflect.DeepEqual(roundTripResult.value, tt.value) {
				t.Errorf("setRequest.Decode().value = %x, want %x", roundTripResult.value, tt.value)
			}
		})
	}
}

func Test_setRequest_Decode(t *testing.T) {
	tests := []struct {
		name      string
		reqBytes  []byte
		wantKey   []byte
		wantValue []byte
		wantErr   bool
	}{
		{
			name:      "valid",
			reqBytes:  []byte{0x00, 0x01, 0x04, 0x61, 0x4b, 0x65, 0x79, 0x06, 0x61, 0x56, 0x61, 0x6c, 0x75, 0x65},
			wantKey:   []byte("aKey"),
			wantValue: []byte("aValue"),
			wantErr:   false,
		},
		{
			name:      "Invalid message returns error",
			reqBytes:  []byte{0x00, 0x01, 0x04, 0x61, 0x4b, 0x65, 0x79, 0x06, 0x61, 0x56, 0x61, 0x6c, 0x65},
			wantKey:   []byte("aKey"),
			wantValue: []byte("aValue"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &setRequest{
				key:   tt.wantKey,
				value: tt.wantValue,
			}
			if err := req.Decode(tt.reqBytes); (err != nil) != tt.wantErr {
				t.Errorf("setRequest.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
