package main

import (
	"fmt"
	"sort"
)

type sortOrder int

const (
	ASC sortOrder = iota
	DESC
)

type OrderBy struct {
	column string
	order  sortOrder
}

type SortOperator struct {
	child      Operator
	orderBy    []OrderBy
	tuples     []Tuple
	currentRow int
}

func less(a, b Tuple, orderBy []OrderBy) bool {
	var i int
	var foundColumn bool
	for _, ob := range orderBy {
		foundColumn = false
		for i = 0; i < len(a.Values); i++ {
			if a.Values[i].Name == ob.column {
				foundColumn = true
				break
			}
		}
		if !foundColumn {
			panic(fmt.Sprintf("Sorting column %s is not present in the tuple", ob.column))
		}
		left := a.Values[i].StringValue
		right := b.Values[i].StringValue
		if ob.order == DESC {
			left, right = right, left
		}
		if left < right {
			return true
		}
		if left > right {
			return false
		}
	}
	return false
}

func NewSortOperator(orderBy []OrderBy, child Operator) *SortOperator {
	var tuples []Tuple
	for child.Next() {
		tuples = append(tuples, child.Execute())
	}
	sort.SliceStable(tuples, func(i, j int) bool { return less(tuples[i], tuples[j], orderBy) })
	return &SortOperator{
		child:      child,
		orderBy:    orderBy,
		tuples:     tuples,
		currentRow: 0,
	}
}

func (so *SortOperator) Next() bool {
	return so.currentRow < len(so.tuples)

}

func (so *SortOperator) Execute() Tuple {
	result := so.tuples[so.currentRow]
	so.currentRow++
	return result
}
