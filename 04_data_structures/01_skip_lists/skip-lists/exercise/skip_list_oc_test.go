package main

import (
	"reflect"
	"testing"
)

var bNodeLinkedList = &skipListNode{
	item:    Item{Key: "b", Value: "b_val"},
	forward: []*skipListNode{nil},
}

var aNodeLinkedList = &skipListNode{
	item:    Item{Key: "a", Value: "a_val"},
	forward: []*skipListNode{bNodeLinkedList},
}

var simpleLinkedList = []*skipListNode{
	aNodeLinkedList,
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
var threeLevelSkipList = []*skipListNode{aNodeSkipList, aNodeSkipList, aNodeSkipList}

func Test_skipListOC_String(t *testing.T) {

	type fields struct {
		head []*skipListNode
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Simple linked list",
			fields: fields{
				head: simpleLinkedList,
			},
			want: `Level 1:
a: a_val -> b: b_val -> nil
`,
		},
		{
			name: "Three level skip list",
			fields: fields{
				head: threeLevelSkipList,
			},
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
			list := skipListOC{
				head: tt.fields.head,
			}
			if got := list.String(); got != tt.want {
				t.Errorf("skipListOC.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_skipListOC_findPrevious(t *testing.T) {
	type fields struct {
		head []*skipListNode
	}
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantNodes []*skipListNode
	}{
		{
			name: "linked list first item",
			fields: fields{
				head: simpleLinkedList,
			},
			args: args{
				key: "a",
			},
			wantNodes: []*skipListNode{nil},
		},
		{
			name: "linked list second item",
			fields: fields{
				head: simpleLinkedList,
			},
			args: args{
				key: "b",
			},
			wantNodes: []*skipListNode{aNodeLinkedList},
		},
		{
			name: "linked list nonexistent item",
			fields: fields{
				head: simpleLinkedList,
			},
			args: args{
				key: "nonexistent",
			},
			wantNodes: []*skipListNode{bNodeLinkedList},
		},
		{
			name: "skipList. Item less than first item",
			fields: fields{
				head: threeLevelSkipList,
			},
			args: args{
				key: "0",
			},
			wantNodes: []*skipListNode{nil, nil, nil},
		},
		{
			name: "skipList. First item",
			fields: fields{
				head: threeLevelSkipList,
			},
			args: args{
				key: "a",
			},
			wantNodes: []*skipListNode{nil, nil, nil},
		},
		{
			name: "skipList. Second item",
			fields: fields{
				head: threeLevelSkipList,
			},
			args: args{
				key: "b",
			},
			wantNodes: []*skipListNode{aNodeSkipList, aNodeSkipList, aNodeSkipList},
		},
		{
			name: "skipList. Third item",
			fields: fields{
				head: threeLevelSkipList,
			},
			args: args{
				key: "c",
			},
			wantNodes: []*skipListNode{bNodeSkipList, aNodeSkipList, aNodeSkipList},
		},
		{
			name: "skipList. Item greater than last item",
			fields: fields{
				head: threeLevelSkipList,
			},
			args: args{
				key: "d",
			},

			wantNodes: []*skipListNode{cNodeSkipList, cNodeSkipList, aNodeSkipList},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &skipListOC{
				head: tt.fields.head,
			}
			if gotNode := list.findPreviousNodes(tt.args.key); !reflect.DeepEqual(gotNode, tt.wantNodes) {
				t.Errorf("skipListOC.findPrevious() = %v, want %v", gotNode, tt.wantNodes)
			}
		})
	}
}
