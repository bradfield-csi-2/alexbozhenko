package main

type LimitOperator struct {
	child         Operator
	limit         int
	returnedSoFar int
}

func NewLimitOperator(limit int, child Operator) *LimitOperator {
	return &LimitOperator{
		child:         child,
		limit:         limit,
		returnedSoFar: 0,
	}
}

func (so *LimitOperator) Next() bool {
	return so.returnedSoFar < so.limit

}

func (so *LimitOperator) Execute() Tuple {
	so.returnedSoFar++
	return so.child.Execute()
}
