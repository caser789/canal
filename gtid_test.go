package main

import (
	"fmt"
	"github.com/google/uuid"
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

func Test_IntervalSlice2(t *testing.T) {
	tests := []struct {
		name           string
		before         IntervalSlice
		afterSort      IntervalSlice
		afterNormalize IntervalSlice
	}{
		{
			name:           "test sort and normalize",
			before:         IntervalSlice{Interval{1, 2}, Interval{2, 4}, Interval{2, 3}},
			afterSort:      IntervalSlice{Interval{1, 2}, Interval{2, 3}, Interval{2, 4}},
			afterNormalize: IntervalSlice{Interval{1, 4}},
		},
		{
			name:           "test sort and normalize",
			before:         IntervalSlice{Interval{1, 2}, Interval{3, 5}, Interval{1, 3}},
			afterSort:      IntervalSlice{Interval{1, 2}, Interval{1, 3}, Interval{3, 5}},
			afterNormalize: IntervalSlice{Interval{1, 5}},
		},
		{
			name:           "test sort and normalize",
			before:         IntervalSlice{Interval{1, 2}, Interval{4, 5}, Interval{1, 3}},
			afterSort:      IntervalSlice{Interval{1, 2}, Interval{1, 3}, Interval{4, 5}},
			afterNormalize: IntervalSlice{Interval{1, 3}, Interval{4, 5}},
		},
		{
			name:           "test sort and normalize",
			before:         IntervalSlice{Interval{1, 4}, Interval{2, 3}},
			afterSort:      IntervalSlice{Interval{1, 4}, Interval{2, 3}},
			afterNormalize: IntervalSlice{Interval{1, 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before.Sort()
			assert.Equal(t, tt.before, tt.afterSort)
			got := tt.before.Normalize()
			assert.Equal(t, got, tt.afterNormalize)
		})
	}
}

func Test_InsertInterval(t *testing.T) {
	tests := []struct {
		name   string
		before IntervalSlice
		insert Interval
		after  IntervalSlice
	}{
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{100, 200}},
			insert: Interval{300, 400},
			after:  IntervalSlice{Interval{100, 200}, Interval{300, 400}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{100, 200}, Interval{300, 400}},
			insert: Interval{50, 70},
			after:  IntervalSlice{Interval{50, 70}, Interval{100, 200}, Interval{300, 400}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 70}, Interval{100, 200}, Interval{300, 400}},
			insert: Interval{101, 201},
			after:  IntervalSlice{Interval{50, 70}, Interval{100, 201}, Interval{300, 400}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 70}, Interval{100, 201}, Interval{300, 400}},
			insert: Interval{99, 202},
			after:  IntervalSlice{Interval{50, 70}, Interval{99, 202}, Interval{300, 400}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 70}, Interval{99, 202}, Interval{300, 400}},
			insert: Interval{102, 302},
			after:  IntervalSlice{Interval{50, 70}, Interval{99, 400}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 70}, Interval{99, 400}},
			insert: Interval{500, 600},
			after:  IntervalSlice{Interval{50, 70}, Interval{99, 400}, Interval{500, 600}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 70}, Interval{99, 400}, Interval{500, 600}},
			insert: Interval{50, 100},
			after:  IntervalSlice{Interval{50, 400}, Interval{500, 600}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 400}, Interval{500, 600}},
			insert: Interval{900, 1000},
			after:  IntervalSlice{Interval{50, 400}, Interval{500, 600}, Interval{900, 1000}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 400}, Interval{500, 600}, Interval{900, 1000}},
			insert: Interval{1010, 1020},
			after:  IntervalSlice{Interval{50, 400}, Interval{500, 600}, Interval{900, 1000}, Interval{1010, 1020}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{50, 400}, Interval{500, 600}, Interval{900, 1000}, Interval{1010, 1020}},
			insert: Interval{49, 1000},
			after:  IntervalSlice{Interval{49, 1000}, Interval{1010, 1020}},
		},
		{
			name:   "test insert interval",
			before: IntervalSlice{Interval{49, 1000}, Interval{1010, 1020}},
			insert: Interval{1, 1012},
			after:  IntervalSlice{Interval{1, 1020}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before.InsertInterval(tt.insert)
			assert.Equal(t, tt.before, tt.after)
		})
	}
}

func Test_UUIDSet_ParseUUIDSet(t *testing.T) {
	tests := []struct {
		name       string
		uUIDSetStr string
		want       *UUIDSet
		errWant    error
	}{
		{
			name:       "happy",
			uUIDSetStr: "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2",
			want: &UUIDSet{
				SID:       uuid.MustParse("de278ad0-2106-11e4-9f8e-6edd0ca20947"),
				Intervals: IntervalSlice{Interval{1, 3}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUUIDSet(tt.uUIDSetStr)
			assert.Equal(t, tt.errWant, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_UUIDSet_Encode(t *testing.T) {
	tests := []struct {
		name string
		uset *UUIDSet
		want []byte
	}{
		{
			name: "happy",
			uset: &UUIDSet{
				SID:       uuid.MustParse("de278ad0-2106-11e4-9f8e-6edd0ca20947"),
				Intervals: IntervalSlice{Interval{1, 3}},
			},
			want: []byte{
				0xde, 0x27, 0x8a, 0xd0, 0x21, 0x6, 0x11, 0xe4, 0x9f, 0x8e, 0x6e, 0xdd, 0xc, 0xa2, 0x9, 0x47,
				0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.uset.Encode()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_UUIDSet_Decode(t *testing.T) {
	tests := []struct {
		name string
		uset *UUIDSet
		b    []byte
		err  error
	}{
		{
			name: "happy",
			uset: &UUIDSet{
				SID:       uuid.MustParse("de278ad0-2106-11e4-9f8e-6edd0ca20947"),
				Intervals: IntervalSlice{Interval{1, 3}},
			},
			b: []byte{
				0xde, 0x27, 0x8a, 0xd0, 0x21, 0x6, 0x11, 0xe4, 0x9f, 0x8e, 0x6e, 0xdd, 0xc, 0xa2, 0x9, 0x47,
				0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UUIDSet{}
			err := a.Decode(tt.b)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.uset, a)
		})
	}
}

func Test_UUIDSet_String(t *testing.T) {
	tests := []struct {
		name string
		uset *UUIDSet
		s    string
	}{
		{
			name: "happy",
			uset: &UUIDSet{
				SID:       uuid.MustParse("de278ad0-2106-11e4-9f8e-6edd0ca20947"),
				Intervals: IntervalSlice{Interval{1, 3}},
			},
			s: "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.uset.String()
			assert.Equal(t, tt.s, s)
		})
	}
}

func Test_ParseMySQLGTIDSet(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    GTIDSet
		errWant error
	}{
		{
			name: "happy",
			str:  "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2,de278ad0-2106-11e4-9f8e-6edd0ca20948:1-2",
			want: &MySQLGTIDSet{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMySQLGTIDSet(tt.str)
			assert.Equal(t, tt.errWant, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
