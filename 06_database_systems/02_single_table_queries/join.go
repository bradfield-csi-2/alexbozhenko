package main

type JoinOperator struct {
	left        Operator
	currentLeft Tuple
	right       Operator
	expression  BinaryJoinExpression
	current     Tuple
}

func NewJoinOperator(left Operator, right Operator, expression BinaryJoinExpression) *JoinOperator {

	joinOperator := &JoinOperator{
		left:        left,
		currentLeft: Tuple{}, // we have to save the state of current tuple
		// since Operator only defines execute(), which moves the iterator
		right:      right,
		expression: expression,
		current:    Tuple{},
	}
	getCurrent(joinOperator)
	return joinOperator
}

func (jo *JoinOperator) Next() bool {
	return jo.current.Values != nil

}

func getCurrent(jo *JoinOperator) {
	var current Tuple = Tuple{}
	var currentRight Tuple

	if jo.currentLeft.Values == nil && jo.left.Next() {
		// this is the first invocation of the function,
		// we need to save the tuple on the outer operator
		// for re-use on next invocations
		jo.currentLeft = jo.left.Execute()
	}
out:
	for jo.currentLeft.Values != nil {
		for jo.right.Next() {
			currentRight = jo.right.Execute()
			if jo.expression.Execute(jo.currentLeft, currentRight) {
				current = Tuple{
					Values: append(jo.currentLeft.Values, currentRight.Values...),
				}
				break out
			}
		}
		jo.right.Reset()
		if jo.left.Next() {
			jo.currentLeft = jo.left.Execute()
		} else {
			jo.currentLeft = Tuple{}
		}
	}
	//	fmt.Println(jo.currentLeft)
	jo.current = current
}

func (jo *JoinOperator) Execute() Tuple {
	result := jo.current
	getCurrent(jo)
	//	fmt.Println(current)
	return result
}

func (jo *JoinOperator) Reset() {
	// not implemented
}
