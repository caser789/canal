package main

import "errors"

var (
	ErrBadConn       = errors.New("bad connection")
	ErrMalformPacket = errors.New("Malform packet error")
)
