package main

import (
	"reflect"
	"testing"
)

func Test_hostnameToQname(t *testing.T) {
	type args struct {
		hostname string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "example.com",
			args: args{hostname: "example.com"},
			want: []byte{7, 101, 120, 97, 109, 112, 108, 101, 3, 99, 111, 109, 0},
		},
		{
			name: "www.example.com",
			args: args{hostname: "www.example.com"},
			want: []byte{3, 119, 119, 199, 7, 101, 120, 97, 109, 112, 108, 101, 3, 99, 111, 109, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hostnameToQname(tt.args.hostname); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hostnameToQname() = %v, want %v", got, tt.want)
			}
		})
	}
}
