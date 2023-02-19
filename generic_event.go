package main

import (
	"encoding/hex"
	"fmt"
	"io"
)

// GenericEvent we don't parse all event, so some we will use GenericEvent instead
type GenericEvent struct {
	Data []byte
}

func (e *GenericEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Event data:\n%s", hex.Dump(e.Data))
	fmt.Fprintln(w)
}

func (e *GenericEvent) Decode(data []byte) error {
	e.Data = data

	return nil
}
