package main

import "fmt"

type EventError struct {
	Header *EventHeader

	//Error message
	Err string

	//Event data
	Data []byte
}

func (e *EventError) Error() string {
	return fmt.Sprintf("Header %#v, Data %q, Err: %v", e.Header, e.Data, e.Err)
}
