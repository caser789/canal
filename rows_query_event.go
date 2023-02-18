package main

import (
	"fmt"
	"io"
)

type RowsQueryEvent struct {
	Query []byte
}

func (e *RowsQueryEvent) Decode(data []byte) error {
	//ignore length byte 1
	e.Query = data[1:]
	return nil
}

func (e *RowsQueryEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Query: %s\n", e.Query)
	fmt.Fprintln(w)
}
