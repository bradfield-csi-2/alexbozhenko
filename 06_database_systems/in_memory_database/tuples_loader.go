package main

import (
	"encoding/csv"
	"log"
	"os"
)

func readCsvFile(filePath string) []Tuple {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	header := records[0]
	var tuples []Tuple
	var values []Value

	for _, row := range records[1:] {
		values = nil
		for i := 0; i < len(row); i++ {
			values = append(values,
				Value{
					Name:        header[i],
					StringValue: row[i],
				})
		}

		tuples = append(tuples, Tuple{
			Values: values,
		})

	}

	return tuples
}
