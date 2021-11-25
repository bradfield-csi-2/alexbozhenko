package main

import (
	"reflect"
	"testing"
)

func Test_sortedSliceDifference(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: []string{},
		},
		{
			name: "equal",
			args: args{
				a: []string{"a"},
				b: []string{"a"},
			},
			want: []string{},
		},
		{
			name: "equal",
			args: args{
				a: []string{"a", "b"},
				b: []string{"a", "b"},
			},
			want: []string{},
		},
		{
			name: "a longer b",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a", "b"},
			},
			want: []string{"c"},
		},
		{
			name: "b longer a",
			args: args{
				a: []string{"a", "b"},
				b: []string{"a", "b", "c"},
			},
			want: []string{},
		},
		{
			name: "overlaping contents",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"b", "c", "d"},
			},
			want: []string{"a"},
		},
		{
			name: "different contents",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"x", "y", "z"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortedSliceDifference(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortedSliceDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}
