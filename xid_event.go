package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// XIDEvent https://mariadb.com/kb/en/xid_event/
// An XID event is generated for a COMMIT of a transaction that modifies one or more tables of an XA-capable storage engine.
// Event Type is XID_EVENT (0x10)
type XIDEvent struct {
	XID uint64 // The XID transaction number. 8 bytes

	// in fact XIDEvent dosen't have the GTIDSet information, just for beneficial to use
	GSet GTIDSet
}

func (e *XIDEvent) Decode(data []byte) error {
	e.XID = binary.LittleEndian.Uint64(data)
	return nil
}

func (e *XIDEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "XID: %d\n", e.XID)
	if e.GSet != nil {
		fmt.Fprintf(w, "GTIDSet: %s\n", e.GSet.String())
	}
	fmt.Fprintln(w)
}
