/* This is taken from the teacher's solution */
package main

import (
	"fmt"
)

// BinaryExpression is the interface for an expression that returns true or false.
type BinaryExpression interface {
	Execute(Tuple) bool
}

// AndExpression is a BinaryExpression that returns the results of AND of its two children.
type AndExpression struct {
	left  BinaryExpression
	right BinaryExpression
}

// NewAndExpression creates a new AndExprsesion.
func NewAndExpression(left, right BinaryExpression) BinaryExpression {
	return &AndExpression{
		left:  left,
		right: right,
	}
}

// Execute returns the result of applying the AND operation to the two children binary expressions
// with the given tuple.
func (a *AndExpression) Execute(tuple Tuple) bool {
	return a.left.Execute(tuple) && a.right.Execute(tuple)
}

// OrExpression is a BinaryExpression that returns the results of OR of its two children.
type OrExpression struct {
	left  BinaryExpression
	right BinaryExpression
}

// NewOrExpression creates a new OrExprsesion.
func NewOrExpression(left, right BinaryExpression) BinaryExpression {
	return &OrExpression{
		left:  left,
		right: right,
	}
}

// Execute returns the result of applying the OR operation to the two children binary expressions
// with the given tuple.
func (a *OrExpression) Execute(tuple Tuple) bool {
	return a.left.Execute(tuple) || a.right.Execute(tuple)
}

// NotExpression is a BinaryExpression that returns the result of applying the NOT operation to its
// child binary expression with the given tuple.
type NotExpression struct {
	child BinaryExpression
}

// NewNotExpression returns a new NotExpression.
func NewNotExpression(child BinaryExpression) BinaryExpression {
	return &NotExpression{
		child: child,
	}
}

// Execute returns the result of applying the NOT operation to the child binary expression with
// the given tuple.
func (n *NotExpression) Execute(tuple Tuple) bool {
	return !n.child.Execute(tuple)
}

// TrueExpression is a BinaryExpression that always returns true regardless of the provided tuple.
type TrueExpression struct {
}

// NewTrueExpression returns a new TrueExpression.
func NewTrueExpression() BinaryExpression {
	return &TrueExpression{}
}

// Execute returns true.
func (n *TrueExpression) Execute(tuple Tuple) bool {
	return true
}

// EQExpression is a BinaryExpression that returns the results of applying an equality check to
// the provided tuple.
type EQExpression struct {
	field string
	value string
}

// NewEQExpression creates a new EQ expression.
func NewEQExpression(field, value string) BinaryExpression {
	return &EQExpression{
		field: field,
		value: value,
	}
}

// Execute returns the results of applying the EQ operation to the provided tuple.
func (e *EQExpression) Execute(tuple Tuple) bool {
	for _, v := range tuple.Values {
		if v.Name == e.field {
			return v.StringValue == e.value
		}
	}

	panic(fmt.Sprintf("tuple: %v did not contain field: %s", tuple, e.field))
}

type BinaryJoinExpression interface {
	Execute(Tuple, Tuple) bool
}

type EQJoinExpression struct {
	field1 string
	field2 string
}

func NewEQJoinExpression(field1, field2 string) *EQJoinExpression {
	return &EQJoinExpression{
		field1: field1,
		field2: field2,
	}
}

func (e *EQJoinExpression) Execute(tuple1, tuple2 Tuple) bool {
	for _, v1 := range tuple1.Values {
		if v1.Name == e.field1 {
			for _, v2 := range tuple2.Values {
				if v2.Name == e.field2 {
					return v1.StringValue == v2.StringValue
				}
			}
			panic(fmt.Sprintf("tuple: %v did not contain field: %s", tuple2, e.field2))
		}
	}
	panic(fmt.Sprintf("tuple: %v did not contain field: %s", tuple1, e.field1))
}
