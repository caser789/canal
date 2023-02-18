package main

import (
	"encoding/binary"
	"fmt"
	"github.com/pingcap/errors"
	"io"
	"time"
)

type EventHeader struct {
	Timestamp uint32
	EventType EventType
	ServerID  uint32
	EventSize uint32
	LogPos    uint32
	Flags     uint16
}

func (h *EventHeader) Dump(w io.Writer) {
	fmt.Fprintf(w, "=== %s ===\n", h.EventType)
	fmt.Fprintf(w, "Date: %s\n", time.Unix(int64(h.Timestamp), 0).Format(TimeFormat))
	fmt.Fprintf(w, "Log position: %d\n", h.LogPos)
	fmt.Fprintf(w, "Event size: %d\n", h.EventSize)
}

func (h *EventHeader) Decode(data []byte) error {
	if len(data) < EventHeaderSize {
		return errors.Errorf("header size too short %d, must 19", len(data))
	}

	pos := 0

	// Timestamp 4 bytes
	h.Timestamp = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	// EventType 1 byte
	h.EventType = EventType(data[pos])
	pos++

	// ServerID 4 bytes
	h.ServerID = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	// EventSize 4 bytes
	h.EventSize = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	// LogPos 4 bytes
	h.LogPos = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	// Flags 2 bytes
	h.Flags = binary.LittleEndian.Uint16(data[pos:])
	// pos += 2

	if h.EventSize < uint32(EventHeaderSize) {
		return errors.Errorf("invalid event size %d, must >= 19", h.EventSize)
	}

	return nil
}
