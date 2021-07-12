package main

import (
	"math/rand"
	"reflect"
	"testing"
)

var emptyList = newSkipListOC()

var bNodeLinkedList = &skipListNode{
	item:    Item{Key: "b", Value: "b_val"},
	forward: []*skipListNode{nil},
}

var aNodeLinkedList = &skipListNode{
	item:    Item{Key: "a", Value: "a_val"},
	forward: []*skipListNode{bNodeLinkedList},
}

var simpleLinkedList = &skipListOC{
	head: []*skipListNode{
		aNodeLinkedList,
	},
}

var cNodeSkipList = &skipListNode{
	item:    Item{Key: "c", Value: "c_val"},
	forward: []*skipListNode{nil, nil},
}
var bNodeSkipList = &skipListNode{
	item:    Item{Key: "b", Value: "b_val"},
	forward: []*skipListNode{cNodeSkipList},
}
var aNodeSkipList = &skipListNode{
	item: Item{Key: "a", Value: "a_val"},
	forward: []*skipListNode{
		bNodeSkipList,
		cNodeSkipList,
		nil,
	},
}
var threeLevelSkipList = &skipListOC{
	head: []*skipListNode{aNodeSkipList, aNodeSkipList, aNodeSkipList},
}

func Test_skipListOC_String(t *testing.T) {

	tests := []struct {
		name string
		list *skipListOC
		want string
	}{
		{
			name: "Empty linked list",
			list: emptyList,
			want: `Level 1:
nil
`,
		},

		{
			name: "Simple linked list",
			list: simpleLinkedList,
			want: `Level 1:
a: a_val -> b: b_val -> nil
`,
		},
		{
			name: "Three level skip list",
			list: threeLevelSkipList,
			want: `Level 3:
a: a_val -> nil
Level 2:
a: a_val -> c: c_val -> nil
Level 1:
a: a_val -> b: b_val -> c: c_val -> nil
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := tt.list
			if got := list.String(); got != tt.want {
				t.Errorf("skipListOC.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_skipListOC_findPrevious(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		list      *skipListOC
		args      args
		wantNodes []*skipListNode
	}{

		{
			name: "Emtpy list",
			list: emptyList,
			args: args{
				key: "whatever",
			},
			wantNodes: []*skipListNode{nil},
		},
		{
			name: "linked list first item",
			list: simpleLinkedList,
			args: args{
				key: "a",
			},
			wantNodes: []*skipListNode{nil},
		},
		{
			name: "linked list second item",
			list: simpleLinkedList,
			args: args{
				key: "b",
			},
			wantNodes: []*skipListNode{aNodeLinkedList},
		},
		{
			name: "linked list nonexistent item",
			list: simpleLinkedList,
			args: args{
				key: "nonexistent",
			},
			wantNodes: []*skipListNode{bNodeLinkedList},
		},
		{
			name: "skipList. Item less than first item",
			list: threeLevelSkipList,
			args: args{
				key: "0",
			},
			wantNodes: []*skipListNode{nil, nil, nil},
		},
		{
			name: "skipList. First item",
			list: threeLevelSkipList,
			args: args{
				key: "a",
			},
			wantNodes: []*skipListNode{nil, nil, nil},
		},
		{
			name: "skipList. Second item",
			list: threeLevelSkipList,
			args: args{
				key: "b",
			},
			wantNodes: []*skipListNode{aNodeSkipList, aNodeSkipList, aNodeSkipList},
		},
		{
			name: "skipList. Third item",
			list: threeLevelSkipList,
			args: args{
				key: "c",
			},
			wantNodes: []*skipListNode{bNodeSkipList, aNodeSkipList, aNodeSkipList},
		},
		{
			name: "skipList. Item greater than last item",
			list: threeLevelSkipList,
			args: args{
				key: "d",
			},

			wantNodes: []*skipListNode{cNodeSkipList, cNodeSkipList, aNodeSkipList},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := tt.list
			if gotNode := list.findPreviousNodes(tt.args.key); !reflect.DeepEqual(gotNode, tt.wantNodes) {
				t.Errorf("skipListOC.findPrevious() = %v, want %v", gotNode, tt.wantNodes)
			}
		})
	}
}

func Test_skipListOC_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name           string
		list           *skipListOC
		args           args
		wantValue      string
		wantKeyPresent bool
	}{
		{
			name: "empty list",
			list: emptyList,
			args: args{
				key: "whatever",
			},
			wantValue:      "",
			wantKeyPresent: false,
		},
		{
			name: "nonexistent key",
			list: threeLevelSkipList,
			args: args{
				key: "nonexistent",
			},
			wantValue:      "",
			wantKeyPresent: false,
		},
		{
			name: "key at the head",
			list: threeLevelSkipList,
			args: args{
				key: "a",
			},
			wantValue:      "a_val",
			wantKeyPresent: true,
		},
		{
			name: "key in the middle",
			list: threeLevelSkipList,
			args: args{
				key: "b",
			},
			wantValue:      "b_val",
			wantKeyPresent: true,
		},
		{
			name: "key at the end",
			list: threeLevelSkipList,
			args: args{
				key: "c",
			},
			wantValue:      "c_val",
			wantKeyPresent: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := tt.list
			got, got1 := list.Get(tt.args.key)
			if got != tt.wantValue {
				t.Errorf("skipListOC.Get() got = %v, want %v", got, tt.wantValue)
			}
			if got1 != tt.wantKeyPresent {
				t.Errorf("skipListOC.Get() got1 = %v, want %v", got1, tt.wantKeyPresent)
			}
		})
	}
}

func Test_skipListOC_Put(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		list skipListOC
		args []args
		want string
	}{
		{
			name: "one element",
			list: skipListOC{},
			args: []args{
				{
					key:   "a",
					value: "a_val",
				},
			},
			want: `Level 1:
a: a_val -> nil
`,
		},

		{
			name: "five elements",
			list: skipListOC{},
			args: []args{
				{
					key:   "c",
					value: "c_val",
				},
				{
					key:   "a",
					value: "a_val",
				},

				{
					key:   "e",
					value: "e_val",
				},
				{
					key:   "b",
					value: "b_val",
				},
				{
					key:   "d",
					value: "d_val",
				},
			},
			want: `Level 5:
d: d_val -> nil
Level 4:
d: d_val -> nil
Level 3:
b: b_val -> d: d_val -> nil
Level 2:
b: b_val -> d: d_val -> nil
Level 1:
a: a_val -> b: b_val -> c: c_val -> d: d_val -> e: e_val -> nil
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipList := newSkipListOC()
			rand.Seed(1)
			for _, arg := range tt.args {
				skipList.Put(arg.key, arg.value)

			}

			if got := skipList.String(); got != tt.want {
				t.Errorf("skipListOC.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
