package main

import (
	"encoding/binary"
	"encoding/csv"
	"io"
	"os"
)

/*
4KB Page format:
+----------+---------------+-----------------+------------------------+
|          |               |                 |          |             |
|  Tuple   |  Tuple        |    Tuple        |    ...   |             |
|          |               |                 |          |             |
+----------+---------------+-----------------+-----------             |
|   (Tuples of variable length ...)                                   |
|                                                                     |
|                                                                     |
|                                     (4B)             (2B)     (2B)  |
|                            +------+-------+--------+-------+--------|
|                            | ...  |Offset,|Offset, |Start  |  N tup |
|                            | ...  |length |length  |of free|        |
|                            | ...  |       |        |space  |        |
+---------------------------------------------------------------------+
*/

/*
Tuple format:
  +----------+---------------+----------+-----------------+------------+
  |          |               |          |                 |            |                                                                                   |
  | Length   | Value         | Length   |  Value          |   ....     |                                                                                   |
  | (varint) |               | (varint) |                 |            |                                                                                   |
  +----------+---------------+----------+-----------------+------------+
*/

/*
If we wanted to quickly find a page that can fit a tuple,
we could use free space map, with the following format, like
in postgresql:
Free space file:
+-----------+------------+----------------------------------+
| n of      |            |                                  |
| 1/256 th  |            |                                  |
| of a page |            |                                  |
| (1 byte)  |            |                                  |
+-----------+------------+----------------------------------+
*/

/*
But let's implement simpler variant for now:
just tuples of variable length text in the file, with separate file for column names.
*/

type FileWriter struct {
	writer io.Writer
}

func NewFileWriter(filePath string, columnNames []string) *FileWriter {
	metaF, err := os.OpenFile(filePath+".meta", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(metaF)
	csvWriter.Write(columnNames)
	csvWriter.Flush()

	dataF, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	return &FileWriter{
		writer: dataF,
	}
}

func (w *FileWriter) Append(t Tuple) error {
	buf := make([]byte, binary.MaxVarintLen64)

	for _, value := range t.Values {
		varintLen := binary.PutUvarint(buf, uint64(len(value.StringValue)))
		w.writer.Write(buf[:varintLen])
		io.WriteString(w.writer, value.StringValue)
	}
	return nil
}

/*
type DiskValue struct {
	length uint64
	value  string
}

type DiskTuple struct {
	Values []DiskValue
}

type Page struct {
	tuples []DiskTuple
}
*/
