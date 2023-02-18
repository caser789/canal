package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BeginLoadQueryEvent_Dump(t *testing.T) {
	tests := []struct {
		name  string
		event *BeginLoadQueryEvent
		s     string
	}{
		{
			name: "ok",
			event: &BeginLoadQueryEvent{
				FileID:    123,
				BlockData: []byte{'a', 'b', 'c'},
			},
			s: `File ID: 123
Block data: abc

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			tt.event.Dump(b)
			assert.Equal(t, tt.s, b.String())
		})
	}
}

func Test_BeginLoadQueryEvent_Decode(t *testing.T) {
	tests := []struct {
		name  string
		b     []byte
		event BeginLoadQueryEvent
		err   error
	}{
		{
			name: "ok",
			b: []byte{
				0x01, 0x02, 0x03, 0x04,
				0x01, 0x02, 0x03, 0x04,
				0x01, 0x02, 0x03, 0x04,
			},
			event: BeginLoadQueryEvent{
				FileID:    0x4030201,
				BlockData: []byte{0x1, 0x2, 0x3, 0x4, 0x1, 0x2, 0x3, 0x4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i BeginLoadQueryEvent
			err := i.Decode(tt.b)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.event, i)
		})
	}
}
