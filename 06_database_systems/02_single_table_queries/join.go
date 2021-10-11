package main

type JoinOperator struct {
	left       Operator
	right      Operator
	expression BinaryJoinExpression
	current    Tuple
}

func NewJoinOperator(left Operator, right Operator, expression BinaryJoinExpression) *JoinOperator {
	var current Tuple = Tuple{}
	var currentLeft, currentRight Tuple
out:
	for left.Next() {
		currentLeft = left.Execute()
		for right.Next() {
			currentRight = right.Execute()
			if expression.Execute(currentLeft, currentRight) {
				current = Tuple{
					Values: append(currentLeft.Values, currentRight.Values...),
				}
				break out
			}
		}
	}
	return &JoinOperator{
		left:       left,
		right:      right,
		expression: expression,
		current:    current,
	}
}

func (so *JoinOperator) Next() bool {
	return so.current.Values != nil

}

func (so *JoinOperator) Execute() Tuple {
	result := so.current

	var current Tuple = Tuple{}
	var currentLeft, currentRight Tuple
out:
	for so.left.Next() {
		currentLeft = so.left.Execute()
		for so.right.Next() {
			currentRight = so.right.Execute()
			if so.expression.Execute(currentLeft, currentRight) {
				current = Tuple{
					Values: append(currentLeft.Values, currentRight.Values...),
				}
				break out
			}
		}
		// need to do reset here?
		// how, if not all operators implement it?
		// type coercion ?
		//so.right.Reset()
	}
	so.current = current
	return result
}
