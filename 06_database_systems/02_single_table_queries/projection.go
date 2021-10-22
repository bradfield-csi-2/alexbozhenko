package main

type ProjectionOperator struct {
	child   Operator
	columns []string
}

func NewProjectionOperator(columns []string, child Operator) *ProjectionOperator {
	return &ProjectionOperator{
		child:   child,
		columns: columns,
	}
}

func (so *ProjectionOperator) Next() bool {
	return so.child.Next()

}

func (so *ProjectionOperator) Execute() Tuple {
	row := so.child.Execute()
	var projected Tuple
	for _, targetColumnName := range so.columns {
		for _, value := range row.Values {
			if targetColumnName == value.Name {
				projected.Values = append(projected.Values, value)
				continue
			}
		}
	}
	return projected
}

func (so *ProjectionOperator) Reset() {
	// Not implemented
}
