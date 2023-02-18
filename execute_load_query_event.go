package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExecuteLoadQueryEvent
//
// This event is written into the binary log file for LOAD DATA INFILE events.
// The event format is similar to a QUERY_EVENT except that it has extra static fields.
//
// Event Type = 0x12
type ExecuteLoadQueryEvent struct {
	SlaveProxyID  uint32 // 4 bytes The ID of the thread that issued this statement on the master.
	ExecutionTime uint32 // 4 bytes The time in seconds that the statement took to execute.
	// The length of the name of the database which was the default database when the statement was executed.
	// This name appears later, in the variable data part.
	// It is necessary for statements such as INSERT INTO t VALUES(1) that don't specify the database
	// and rely on the default database previously selected by USE.
	SchemaLength     uint8  // 1 byte
	ErrorCode        uint16 // 2 bytes The error code resulting from execution of the statement on the master.
	StatusVars       uint16 // 2 bytes The length of the status variable block.
	FileID           uint32 // 4 bytes The ID of the loaded file
	StartPos         uint32 // 4 bytes Offset from the start of the statement to the beginning of the filename
	EndPos           uint32 // 4 bytes Offset from the start of the statement to the end of the filename
	DupHandlingFlags uint8  // 1 byte How LOAD DATA INFILE handles duplicates (0x0: error, 0x1: ignore, 0x2: replace).
}

func (e *ExecuteLoadQueryEvent) Decode(data []byte) error {
	pos := 0

	e.SlaveProxyID = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.ExecutionTime = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.SchemaLength = data[pos]
	pos++

	e.ErrorCode = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	e.StatusVars = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	e.FileID = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.StartPos = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.EndPos = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.DupHandlingFlags = data[pos]

	return nil
}

func (e *ExecuteLoadQueryEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Slave proxy ID: %d\n", e.SlaveProxyID)
	fmt.Fprintf(w, "Execution time: %d\n", e.ExecutionTime)
	fmt.Fprintf(w, "Schame length: %d\n", e.SchemaLength)
	fmt.Fprintf(w, "Error code: %d\n", e.ErrorCode)
	fmt.Fprintf(w, "Status vars length: %d\n", e.StatusVars)
	fmt.Fprintf(w, "File ID: %d\n", e.FileID)
	fmt.Fprintf(w, "Start pos: %d\n", e.StartPos)
	fmt.Fprintf(w, "End pos: %d\n", e.EndPos)
	fmt.Fprintf(w, "Dup handling flags: %d\n", e.DupHandlingFlags)
	fmt.Fprintln(w)
}
