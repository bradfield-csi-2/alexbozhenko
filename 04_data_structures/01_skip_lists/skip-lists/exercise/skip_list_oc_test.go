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

var cNodeSkipList = skipListNode{
	item:    Item{Key: "c", Value: "c_val"},
	forward: []*skipListNode{nil, nil},
}
var bNodeSkipList = skipListNode{
	item:    Item{Key: "b", Value: "b_val"},
	forward: []*skipListNode{&cNodeSkipList},
}
var aNodeSkipList = skipListNode{
	item: Item{Key: "a", Value: "a_val"},
	forward: []*skipListNode{
		&bNodeSkipList,
		&cNodeSkipList,
		nil,
	},
}
var threeLevelSkipList = []*skipListNode{&aNodeSkipList, &aNodeSkipList, &aNodeSkipList}

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
		name     string
		fields   fields
		args     args
		wantNode *skipListNode
	}{
		{
			name: "linked list first item",
			fields: fields{
				head: simpleLinkedList,
			},
			args: args{
				key: "a",
			},
			wantNode: nil,
		},
		{
			name: "linked list second item",
			fields: fields{
				head: simpleLinkedList,
			},
			args: args{
				key: "b",
			},
			wantNode: aNodeLinkedList,
		},
		{
			name: "linked list nonexistent item",
			fields: fields{
				head: simpleLinkedList,
			},
			args: args{
				key: "nonexistent",
			},
			wantNode: bNodeLinkedList,
		},
		{
			name: "skipList. Item less than first item",
			fields: fields{
				head: simpleLinkedList,
			},
			args: args{
				key: "0",
			},
			wantNode: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &skipListOC{
				head: tt.fields.head,
			}
			if gotNode := list.findPrevious(tt.args.key); !reflect.DeepEqual(gotNode, tt.wantNode) {
				t.Errorf("skipListOC.findPrevious() = %v, want %v", gotNode, tt.wantNode)
			}
		})
	}
}
