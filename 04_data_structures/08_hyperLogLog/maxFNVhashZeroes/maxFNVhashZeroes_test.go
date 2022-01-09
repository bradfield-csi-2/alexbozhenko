package main

import "testing"

func Test_countTrailingZeroes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI int
	}{
		{
			name:  "Zero",
			args:  args{[]byte{0}},
			wantI: 0,
		},
		{
			name:  "8 zeroes",
			args:  args{[]byte{1, 0}},
			wantI: 8,
		},
		{
			name: "0 zeroes",
			args: args{
				b: []byte{1},
			},
			wantI: 0,
		},
		{
			name: "half byte",
			args: args{
				b: []byte{64, 0},
			},
			wantI: 14,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := countTrailingZeroes(tt.args.b); gotI != tt.wantI {
				t.Errorf("countTrailingZeroes() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}
