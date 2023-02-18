package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// RotateEvent https://mariadb.com/kb/en/rotate_event/
// When a binary log file exceeds the configured size limit, a ROTATE_EVENT is written at the end of the file,
// pointing to the next file in the sequence.
// ROTATE_EVENT is generated locally and written to the binary log on the master
// and it's also written when a FLUSH LOGS statement occurs on the master server.
// The ROTATE_EVENT is sent to the connected slave servers.
//
// The Event Type is set ROTATE_EVENT (0x4)
type RotateEvent struct {
	// The position of the first event in the next log file.
	// Note: it always contains the number 4 (meaning the next event starts at position 4 in the next binary log).
	Position uint64

	// The next binary log name. The filename is not null-terminated.
	// filename + EOF
	NextLogName []byte
}

func (e *RotateEvent) Decode(data []byte) error {
	e.Position = binary.LittleEndian.Uint64(data[0:])
	e.NextLogName = data[8:]

	return nil
}

func (e *RotateEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Position: %d\n", e.Position)
	fmt.Fprintf(w, "Next log name: %s\n", e.NextLogName)
	fmt.Fprintln(w)
}
