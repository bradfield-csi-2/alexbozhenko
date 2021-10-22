package main

import "io"

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
*/

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

type FileWriter struct {
	writer io.Writer
}

func (w *FileWriter) Append(t Tuple) error {
	// TODO
	return nil
}

type FileReader struct {
	//...
}

// TODO implement operator interface
