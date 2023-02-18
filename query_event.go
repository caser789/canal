package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// QueryEvent https://mariadb.com/kb/en/query_event/
// This event is written into the binary log file for:
//     STATEMENT based replication (updating statements)
//     DDLs
//     COMMIT related to non-transactional engines (MyISAM, BLACKHOLE etc)
//
// Event Type = 0x02
type QueryEvent struct {
	SlaveProxyID  uint32
	ExecutionTime uint32
	ErrorCode     uint16
	StatusVars    []byte
	Schema        []byte
	Query         []byte

	// in fact QueryEvent dosen't have the GTIDSet information, just for beneficial to use
	GSet GTIDSet
}

func (e *QueryEvent) Decode(data []byte) error {
	pos := 0

	e.SlaveProxyID = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.ExecutionTime = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	schemaLength := data[pos]
	pos++

	e.ErrorCode = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	statusVarsLength := binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	e.StatusVars = data[pos : pos+int(statusVarsLength)]
	pos += int(statusVarsLength)

	e.Schema = data[pos : pos+int(schemaLength)]
	pos += int(schemaLength)

	//skip 0x00
	pos++

	e.Query = data[pos:]
	return nil
}

func (e *QueryEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Slave proxy ID: %d\n", e.SlaveProxyID)
	fmt.Fprintf(w, "Execution time: %d\n", e.ExecutionTime)
	fmt.Fprintf(w, "Error code: %d\n", e.ErrorCode)
	//fmt.Fprintf(w, "Status vars: \n%s", hex.Dump(e.StatusVars))
	fmt.Fprintf(w, "Schema: %s\n", e.Schema)
	fmt.Fprintf(w, "Query: %s\n", e.Query)
	if e.GSet != nil {
		fmt.Fprintf(w, "GTIDSet: %s\n", e.GSet.String())
	}
	fmt.Fprintln(w)
}
