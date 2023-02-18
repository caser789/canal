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
