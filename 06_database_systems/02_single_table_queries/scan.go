package main

type ScanOperator struct {
	tableName  string
	child      Operator
	currentRow int
	table      []Tuple
}

func NewScanOperator(tableName string, db *InMemoryDB, child Operator) *ScanOperator {
	table, tableExists := (*db)[tableName]
	if !tableExists {
		panic("Table does not exist")
	}
	return &ScanOperator{
		tableName:  tableName,
		child:      child,
		currentRow: 0,
		table:      table,
	}
}

func (so *ScanOperator) Next() bool {
	return so.currentRow < len(so.table)

}

func (so *ScanOperator) Execute() Tuple {
	tuple := so.table[so.currentRow]
	so.currentRow += 1
	return tuple
}
