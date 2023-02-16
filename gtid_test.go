package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterval_String(t *testing.T) {
	tests := []struct {
		name string
		i    Interval
		want string
	}{
		{
			name: "test length 1",
			i:    Interval{Start: 11, Stop: 12},
			want: "11",
		},
		{
			name: "test length 2",
			i:    Interval{Start: 11, Stop: 21},
			want: "11-20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.i.String())
		})
	}
}

func Test_parseInterval(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		iWant   Interval
		errWant error
	}{
		{
			name:    "if length not 1 or 2 then error",
			s:       "1-2-3",
			errWant: fmt.Errorf("invalid interval format, must be n[-n]"),
		},
		{
			name:    "if end < start then error",
			s:       "11-10",
			iWant:   Interval{Start: 11, Stop: 11},
			errWant: fmt.Errorf("invalid interval format, must be n[-n] and end must >= start"),
		},
		{
			name:  "length 2 happy",
			s:     "11-11",
			iWant: Interval{Start: 11, Stop: 12},
		},
		{
			name:  "length 1 happy",
			s:     "11",
			iWant: Interval{Start: 11, Stop: 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := parseInterval(tt.s)
			assert.Equal(t, i, tt.iWant)
			assert.Equal(t, err, tt.errWant)
		})
	}
}

func Test_IntervalSlice(t *testing.T) {
	tests := []struct {
		name   string
		before IntervalSlice
		after  IntervalSlice
	}{
		{
			name: "before is empty then ok",
		},
		{
			name:   "no overlap",
			before: IntervalSlice{Interval{1, 3}, Interval{5, 7}},
			after:  IntervalSlice{Interval{1, 3}, Interval{5, 7}},
		},
		{
			name:   "has overlap",
			before: IntervalSlice{Interval{1, 5}, Interval{3, 7}},
			after:  IntervalSlice{Interval{1, 7}},
		},
		{
			name:   "has overlap",
			before: IntervalSlice{Interval{3, 7}, Interval{1, 5}},
			after:  IntervalSlice{Interval{1, 7}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := tt.before.Normalize()
			assert.Equal(t, tt.after, o)
		})
	}
}
