package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IntVarEvent_Dump(t *testing.T) {
	tests := []struct {
		name  string
		event *IntVarEvent
		s     string
	}{
		{
			name: "ok",
			event: &IntVarEvent{
				Type:  LAST_INSERT_ID,
				Value: 123,
			},
			s: `Type: 1
Value: 123
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

func Test_IntVarEvent_Decode(t *testing.T) {
	tests := []struct {
		name  string
		b     []byte
		event IntVarEvent
		err   error
	}{
		{
			name: "ok",
			b:    []byte{0x01, 0x01, 0, 0, 0, 0, 0, 0, 0},
			event: IntVarEvent{
				Type:  LAST_INSERT_ID,
				Value: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i IntVarEvent
			err := i.Decode(tt.b)
			assert.Equal(t, err, tt.err)
			assert.Equal(t, i, tt.event)
		})
	}
}
