package main

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_XIDEvent_Dump(t *testing.T) {
	tests := []struct {
		name  string
		event *XIDEvent
		s     string
	}{
		{
			name: "ok",
			event: &XIDEvent{
				XID: 123,
				GSet: &MySQLGTIDSet{
					Sets: map[string]*UUIDSet{
						"de278ad0-2106-11e4-9f8e-6edd0ca20947": {
							SID:       uuid.MustParse("de278ad0-2106-11e4-9f8e-6edd0ca20947"),
							Intervals: IntervalSlice{Interval{1, 3}},
						},
						"de278ad0-2106-11e4-9f8e-6edd0ca20948": {
							SID:       uuid.MustParse("de278ad0-2106-11e4-9f8e-6edd0ca20948"),
							Intervals: IntervalSlice{Interval{1, 3}},
						},
					},
				},
			},
			s: `XID: 123
GTIDSet: de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2,de278ad0-2106-11e4-9f8e-6edd0ca20948:1-2

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

func Test_XIDEvent_Decode(t *testing.T) {
	tests := []struct {
		name  string
		b     []byte
		event XIDEvent
		err   error
	}{
		{
			name:  "ok",
			b:     []byte{0x66, 00, 00, 00, 00, 00, 00, 00},
			event: XIDEvent{XID: 102},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i XIDEvent
			err := i.Decode(tt.b)
			assert.Equal(t, err, tt.err)
			assert.Equal(t, i, tt.event)
		})
	}
}
