package bitmap

import "testing"

func Test_genHeaderByte(t *testing.T) {
	type args struct {
		st        byteType
		runLength uint8
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{name: "uncompressed, 1 length",
			args: args{
				st:        uncompressed,
				runLength: 1,
			},
			want: 0b00000001},
		{name: "uncompressed, 127 length",
			args: args{
				st:        uncompressed,
				runLength: 127,
			},
			want: 0b01111111},
		{name: "compressed0, 1 length",
			args: args{
				st:        compressed0,
				runLength: 1,
			},
			want: 0b10000001},
		{name: "compressed0, 63 length",
			args: args{
				st:        compressed0,
				runLength: 63,
			},
			want: 0b10111111},
		{name: "compressed1, 1 length",
			args: args{
				st:        compressed1,
				runLength: 1,
			},
			want: 0b11000001},
		{name: "compressed1, 63 length",
			args: args{
				st:        compressed1,
				runLength: 63,
			},
			want: 0b11111111},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genHeaderByte(tt.args.st, tt.args.runLength); got != tt.want {
				t.Errorf("genHeaderByte() = %#08b, want %#08b", got, tt.want)
			}
		})
	}
}
