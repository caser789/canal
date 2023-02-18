package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

type IntVarEventType byte

const (
	INVALID IntVarEventType = iota
	LAST_INSERT_ID
	INSERT_ID
)

// IntVarEvent
// A INTVAR_EVENT is written every time a statement uses an auto increment column or LAST_INSERT_ID() function.
// Event Type is 5 (0x05)
type IntVarEvent struct {
	Type  IntVarEventType // 1 byte
	Value uint64          // 8 bytes
}

func (i *IntVarEvent) Decode(data []byte) error {
	i.Type = IntVarEventType(data[0])
	i.Value = binary.LittleEndian.Uint64(data[1:])
	return nil
}

func (i *IntVarEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Type: %d\n", i.Type)
	fmt.Fprintf(w, "Value: %d\n", i.Value)
}
