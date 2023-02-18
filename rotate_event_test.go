package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RotateEvent_Dump(t *testing.T) {
	tests := []struct {
		name  string
		event *RotateEvent
		s     string
	}{
		{
			name: "ok",
			event: &RotateEvent{
				Position:    4,
				NextLogName: []byte("mysql-bin.000019"),
			},
			s: `Position: 4
Next log name: mysql-bin.000019

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

func Test_RotateEvent_Decode(t *testing.T) {
	tests := []struct {
		name  string
		b     []byte
		event RotateEvent
		err   error
	}{
		{
			name: "ok",
			b: []byte{
				0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x6d, 0x79, 0x73, 0x71, 0x6c, 0x2d, 0x62, 0x69,
				0x6e, 0x2e, 0x30, 0x30, 0x30, 0x30, 0x31, 0x39,
			},
			event: RotateEvent{
				Position:    4,
				NextLogName: []byte("mysql-bin.000019"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i RotateEvent
			err := i.Decode(tt.b)
			assert.Equal(t, err, tt.err)
			assert.Equal(t, i, tt.event)
		})
	}
}
