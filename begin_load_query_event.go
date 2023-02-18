package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// BeginLoadQueryEvent
//
// This event is written into the binary log file for LOAD DATA INFILE events
// if the server variable binlog_mode was set to "STATEMENT".
//
// LOAD DATA INFILE Reads rows from a text file into the designated table on the database at a very high speed.
// The file name must be given as a literal string.
//
// Event Type = 0x11
type BeginLoadQueryEvent struct {
	FileID    uint32 // 4 bytes. id of the file
	BlockData []byte // Null terminated data block.
}

func (e *BeginLoadQueryEvent) Decode(data []byte) error {
	pos := 0

	e.FileID = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.BlockData = data[pos:]

	return nil
}

func (e *BeginLoadQueryEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "File ID: %d\n", e.FileID)
	fmt.Fprintf(w, "Block data: %s\n", e.BlockData)
	fmt.Fprintln(w)
}
