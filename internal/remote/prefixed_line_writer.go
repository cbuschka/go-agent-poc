package remote

import (
	"io"
	"log"
)

type PrefixedLineWriter struct {
	prefix string
	buf    []byte
}

func NewPrefixedLineWriter(prefix string) io.Writer {
	wr := PrefixedLineWriter{prefix: prefix, buf: make([]byte, 1024)}
	return io.Writer(&wr)
}

func (w *PrefixedLineWriter) Write(bs []byte) (int, error) {
	w.buf = append(w.buf, bs...)
	for i := 0; i < len(w.buf); i++ {
		if w.buf[i] == '\n' {
			log.Printf("%s %s", w.prefix, string(w.buf[0:i+1]))
			w.buf = w.buf[i+1 : len(w.buf)]
			i = 0
		}
	}

	return len(bs), nil
}

func (w *PrefixedLineWriter) Close() error {
	return nil
}
