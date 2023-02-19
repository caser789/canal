package main

import (
	"bytes"
	"fmt"
	"github.com/pingcap/errors"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const (
	TooBigBlockSize = 1024 * 1024 * 4
)

// no copy to change slice to string
// use your own risk
func BytesToString(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

// no copy to change string to slice
// use your own risk
func StringToBytes(s string) (b []byte) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return
}

func Uint64ToInt64(val uint64) int64 {
	return *(*int64)(unsafe.Pointer(&val))
}

func Uint64ToFloat64(val uint64) float64 {
	return *(*float64)(unsafe.Pointer(&val))
}

func Int64ToUint64(val int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&val))
}

func Float64ToUint64(val float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&val))
}

var bytesBufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func GetBytesBuffer() (data *bytes.Buffer) {
	data = bytesBufferPool.Get().(*bytes.Buffer)
	data.Reset()
	return data
}

func PutBytesBuffer(data *bytes.Buffer) {
	if data == nil || data.Len() > TooBigBlockSize {
		return
	}
	bytesBufferPool.Put(data)
}

type ByteSlice struct {
	B []byte
}

var (
	byteSlicePool = sync.Pool{
		New: func() interface{} {
			return new(ByteSlice)
		},
	}
)

func GetByteSlice(length int) *ByteSlice {
	data := byteSlicePool.Get().(*ByteSlice)
	if cap(data.B) < length {
		data.B = make([]byte, length)
	} else {
		data.B = data.B[:length]
	}
	return data
}

func PutByteSlice(data *ByteSlice) {
	data.B = data.B[:0]
	byteSlicePool.Put(data)
}

// fracTime is a help structure wrapping Golang Time.
type fracTime struct {
	time.Time

	// Dec must in [0, 6]
	Dec int

	timestampStringLocation *time.Location
}

var fracTimeFormat []string

func (t fracTime) String() string {
	tt := t.Time
	if t.timestampStringLocation != nil {
		tt = tt.In(t.timestampStringLocation)
	}
	return tt.Format(fracTimeFormat[t.Dec])
}

func formatZeroTime(frac int, dec int) string {
	if dec == 0 {
		return "0000-00-00 00:00:00"
	}

	s := fmt.Sprintf("0000-00-00 00:00:00.%06d", frac)

	// dec must < 6, if frac is 924000, but dec is 3, we must output 924 here.
	return s[0 : len(s)-(6-dec)]
}

func formatBeforeUnixZeroTime(year, month, day, hour, minute, second, frac, dec int) string {
	if dec == 0 {
		return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
	}

	s := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%06d", year, month, day, hour, minute, second, frac)

	// dec must < 6, if frac is 924000, but dec is 3, we must output 924 here.
	return s[0 : len(s)-(6-dec)]
}

func init() {
	fracTimeFormat = make([]string, 7)
	fracTimeFormat[0] = "2006-01-02 15:04:05"

	for i := 1; i <= 6; i++ {
		fracTimeFormat[i] = fmt.Sprintf("2006-01-02 15:04:05.%s", strings.Repeat("0", i))
	}
}

func Pstack() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[0:n])
}

// ErrorEqual returns a boolean indicating whether err1 is equal to err2.
func ErrorEqual(err1, err2 error) bool {
	e1 := errors.Cause(err1)
	e2 := errors.Cause(err2)

	if e1 == e2 {
		return true
	}

	if e1 == nil || e2 == nil {
		return e1 == e2
	}

	return e1.Error() == e2.Error()
}
