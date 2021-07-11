package main

import (
	"testing"
)

func Test_skipListOC_String(t *testing.T) {
	simpleLinkedList := []*skipListNode{
		{
			item: Item{Key: "a", Value: "a_val"},
			forward: []*skipListNode{
				{
					item:    Item{Key: "b", Value: "b_val"},
					forward: []*skipListNode{nil},
				},
			},
		},
	}

	c_node := skipListNode{
		item:    Item{Key: "c", Value: "c_val"},
		forward: []*skipListNode{nil, nil},
	}
	b_node := skipListNode{
		item:    Item{Key: "b", Value: "b_val"},
		forward: []*skipListNode{&c_node},
	}
	a_node := skipListNode{
		item: Item{Key: "a", Value: "a_val"},
		forward: []*skipListNode{
			&b_node,
			&c_node,
			nil,
		},
	}
	threeLevelSkipList := []*skipListNode{&a_node, &a_node, &a_node}

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
