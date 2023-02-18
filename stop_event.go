package main

import (
	"io"
)

// StopEvent
// The master server writes the event to the binary log when it shuts down or when resuming after a mysqld process crash.
// A new binary log file is always created but there is no ROTATE_EVENT.
// STOP_EVENT is then the last written event after clean shutdown or resuming a crash.
//
// this event is never sent to slave servers.
//
// Event header with EventType set to STOP_EVENT (0x03).
// Event header NextPos set to EOF
// No special flags added.
//
// The event has no data
type StopEvent struct {
}

func (e *StopEvent) Decode(data []byte) error {
	return nil
}

func (e *StopEvent) Dump(w io.Writer) {
}
