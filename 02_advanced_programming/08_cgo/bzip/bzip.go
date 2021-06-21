package bzip

/*
#cgo CFLAGS: -I /usr/include
#cgo LDFLAGS: -I /usr/lib64 -l bz2
#include <bzlib.h>
int bz2compress(bz_stream *s, int action, char *in,
                unsigned *inlen,
                char *out,
                unsigned *outlen);
*/
import "C"

import (
	"io"
	"unsafe"
)

type writer struct {
	w      io.Writer //undedrlying output stream
	stream *C.bz_stream
	outbuf [64 * 1024]byte
}

func NewWriter(out io.Writer) io.WriteCloser {
	const (
		blockSize  = 9
		verbosity  = 0
		workFactor = 30
	)
	w := &writer{w: out, stream: new(C.bz_stream)}
	C.BZ2_bzCompressInit(w.stream, blockSize, verbosity, workFactor)
	return w
}

func (w *writer) Write(data []byte) (int, error) {
	if w.stream == nil {
		panic("closed")
	}
	var total int //uncompressed bytes written
	for len(data) > 0 {
		inlen, outlen := C.uint(len(data)), C.uint(cap(w.outbuf))
		C.bz2compress(w.stream, C.BZ_RUN,
			(*C.char)(unsafe.Pointer(&data[0])), &inlen,
			(*C.char)(unsafe.Pointer(&w.outbuf)), &outlen)
		total += inlen
		data = data[inlen:]
		if _, err := w.w.Write(w.outbuf[:outlen]); err != nil {
			return total, err
		}
	}
	return total, nil
}
