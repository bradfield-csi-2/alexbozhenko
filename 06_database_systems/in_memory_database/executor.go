package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

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

func executor(root RootOperator) {
	for root.child.Next() {
		fmt.Printf("%s\n", root.child.Execute())
	}
}

type InMemoryDB map[string][]Tuple

var DB InMemoryDB = make(InMemoryDB)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to file")

func main() {

	//DB["movies"] = readCsvFile("movies.csv")
	//DB["tags"] = readCsvFile("tags.csv")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/*
		moviesFileWriter := NewFileWriter(
			"./movies",
			[]string{"movieId", "title", "genres"},
		)
		for _, tuple := range DB["movies"] {
			moviesFileWriter.Append(tuple)
		}

		tagsFileWriter := NewFileWriter(
			"./tags",
			[]string{"userId", "movieId", "tag", "timestamp"},
		)
		for _, tuple := range DB["tags"] {
			fmt.Println(tuple)
			tagsFileWriter.Append(tuple)
		}
	*/

	// root := RootOperator{
	// 	child: NewLimitOperator(5,
	// 		NewProjectionOperator([]string{"title", "genres", "movieId"},
	// 			NewSortOperator(
	// 				[]OrderBy{
	// 					{
	// 						column: "genres",
	// 						order:  ASC,
	// 					},
	// 					{
	// 						column: "title",
	// 						order:  DESC,
	// 					},
	// 					{
	// 						column: "movieId",
	// 						order:  DESC,
	// 					},
	// 				},
	// 				NewSelectionOperator(
	// 					NewOrExpression(
	// 						NewEQExpression("genres", "Action|Adventure|Thriller"),
	// 						NewEQExpression("genres", "Adventure|Animation|Children|Comedy|Fantasy"),
	// 					),
	// 					NewScanOperator("movies", &DB, nil),
	// 				)))),
	// }
	//fmt.Println(strings.Join(executor(root), "\n"))

	/*
		root := RootOperator{
			child: NewJoinOperator(
				NewFileScanOperator("movies", ".", nil),
				NewFileScanOperator("tags", ".", nil),
				NewEQJoinExpression("movieId", "movieId"),
			),
		}
	*/

	root := RootOperator{
		child: NewFileScanOperator("tags", ".", nil),
	}

	executor(root)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		defer f.Close()
	}
}
