package main

import "fmt"

/*
For example query like this:
select genres from movies where title ilike '%golden%' order by movieid limit 5;

Looks like order will be the following:

Limit which limits the number of output rows.
  (int limit)

Projection which yields only a subset of the columns in
  the underlying rows, possibly renaming some of the columns.
  ([]column_name)

Sort which yields records in sorted order.
  (column_name, order)

Selection which yields only records that return true for a
  predicate function, or more interestingly, arbitrary “expressions”
  of predicates (e.g. A == B OR C == D).
  (column_name string, OP WHAT_TYPE?, value string)

Scan which yields each row for the table as needed. In this
  initial implementation your Scan operator can return
  rows from a predefined list in memory.
  (table_name)
*/

func executor(root RootOperator) []Tuple {
	var result []Tuple
	for root.child.Next() {
		result = append(result, root.child.Execute())
	}
	return result
}

type InMemoryDB map[string][]Tuple

var DB InMemoryDB = make(InMemoryDB)

func main() {
	DB["movies"] = readCsvFile("movies.csv")
	DB["tags"] = readCsvFile("tags.csv")

	root := RootOperator{
		child: NewScanOperator("movies", &DB, nil),
	}
	fmt.Println(executor(root))
}
