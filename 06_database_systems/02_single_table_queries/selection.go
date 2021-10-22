package main

type SelectionOperator struct {
	child      Operator
	expression BinaryExpression
	current    Tuple
}

func NewSelectionOperator(expression BinaryExpression, child Operator) *SelectionOperator {
	var current Tuple = Tuple{}
	for child.Next() {
		current = child.Execute()
		if expression.Execute(current) {
			break
		}
	}
	return &SelectionOperator{
		child:      child,
		current:    current,
		expression: expression,
	}
}

func (so *SelectionOperator) Next() bool {
	return so.current.Values != nil

}

func (so *SelectionOperator) Execute() Tuple {
	record := Tuple{}
	result := so.current
	for so.child.Next() {
		tmp := so.child.Execute()
		if so.expression.Execute(tmp) {
			record = tmp
			break
		}
	}
	so.current = record
	return result
}

func (so *SelectionOperator) Reset() {
	// Not implemented
}
