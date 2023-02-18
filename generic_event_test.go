package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenericEvent_Dump(t *testing.T) {
	tests := []struct {
		name  string
		event *GenericEvent
		s     string
	}{
		{
			name: "ok",
			event: &GenericEvent{
				Data: []byte{'a', 'b', 'c'},
			},
			s: `Event data: 
00000000  61 62 63                                          |abc|

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

func Test_GenericEvent_Decode(t *testing.T) {
	tests := []struct {
		name  string
		b     []byte
		event GenericEvent
		err   error
	}{
		{
			name: "ok",
			b: []byte{
				'a', 'b', 'c',
			},
			event: GenericEvent{
				Data: []byte{0x61, 0x62, 0x63},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i GenericEvent
			err := i.Decode(tt.b)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.event, i)
		})
	}
}
