package main

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileScanOperator struct {
	tableName string
	child     Operator
	columns   []string

	//state
	r *bufio.Reader
	f *os.File
}

func NewFileScanOperator(tableName string, dbPath string, child Operator) *FileScanOperator {
	dataF, err := os.Open(filepath.Join(dbPath, tableName))

	if err != nil {
		panic(err)
	}

	metaF, err := os.Open(filepath.Join(dbPath, tableName+".meta"))
	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}
	defer dataF.Close()

	csvReader := csv.NewReader(metaF)
	csvReader.Read()
	columnNames, err := csvReader.ReadAll()
	columns := columnNames[0]

	return &FileScanOperator{
		tableName: tableName,
		child:     child,
		columns:   columns,
		r:         bufio.NewReader(dataF),
		f:         dataF,
	}
}

func (so *FileScanOperator) Next() bool {
	_, err := so.r.Peek(1)
	if err == io.EOF {
		return false
	} else {
		return true
	}
}

func (so *FileScanOperator) Execute() Tuple {
	var valLength uint64
	var valBytes []byte
	var err error
	var tuple Tuple = Tuple{}
	for _, column := range so.columns {
		valLength, err = binary.ReadUvarint(so.r)
		if err != nil {
			panic(err)
		}
		valBytes = make([]byte, valLength)
		io.ReadFull(so.r, valBytes)
		tuple.Values = append(tuple.Values, Value{
			Name:        column,
			StringValue: string(valBytes),
		})
	}
	return tuple
}

func (so *FileScanOperator) Reset() {
	so.f.Seek(0, io.SeekStart)
	//	so.r = bufio.NewReader(so.f)
}
