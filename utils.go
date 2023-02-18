package main

import (
	"bytes"
	"reflect"
	"sync"
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
