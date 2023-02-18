package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_EventHeader_Dump(t *testing.T) {
	tests := []struct {
		name   string
		header *EventHeader
		s      string
	}{
		{
			name: "ok",
			header: &EventHeader{
				Timestamp: 123456,
				EventType: QUERY_EVENT,
				ServerID:  11,
				EventSize: 22,
				LogPos:    33,
				Flags:     44,
			},
			s: `=== QueryEvent ===
Date: 1970-01-02 17:47:36
Log position: 33
Event size: 22
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			tt.header.Dump(b)
			assert.Equal(t, tt.s, b.String())
		})
	}
}

func Test_EventHeader_Decode(t *testing.T) {
	tests := []struct {
		name   string
		b      []byte
		header EventHeader
		err    error
	}{
		{
			name: "ok",
			b: []byte{
				0x78, 0xed, 0x1c, 0x5b,
				0x05,
				0x01, 0x00, 0x00, 0x00,
				0x20, 0x00, 0x00, 0x00,
				0x02, 0x03, 0x00, 0x00,
				0x00, 0x00,
				0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2d, 0x3f, 0xa2, 0xf5,
			},
			header: EventHeader{
				Timestamp: 1528622456,
				EventType: INTVAR_EVENT,
				ServerID:  1,
				EventSize: 32,
				LogPos:    770,
				Flags:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h EventHeader
			err := h.Decode(tt.b)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.header, h)
		})
	}
}
