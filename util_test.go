package mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomBuf(t *testing.T) {
	x, err := RandomBuf(11)
	assert.NoError(t, err)
	fmt.Println(x)
}

func TestPutLengthEncodedInt(t *testing.T) {
	assert.Equal(t, []byte{0x1}, PutLengthEncodedInt(1))
	assert.Equal(t, []byte{250}, PutLengthEncodedInt(250))
	assert.Equal(t, []byte{252, 251, 0}, PutLengthEncodedInt(251))
	assert.Equal(t, []byte{252, 252, 0}, PutLengthEncodedInt(252))
	assert.Equal(t, []byte{252, 253, 0}, PutLengthEncodedInt(253))
	assert.Equal(t, []byte{253, 0, 0, 1}, PutLengthEncodedInt(0xffff+1))
	assert.Equal(t, []byte{253, 1, 0, 1}, PutLengthEncodedInt(0xffff+2))
	assert.Equal(t, []byte{253, 2, 0, 1}, PutLengthEncodedInt(0xffff+3))
	assert.Equal(t, []byte{253, 3, 0, 1}, PutLengthEncodedInt(0xffff+4))
	assert.Equal(t, []byte{0xfe, 0, 0, 0, 1, 0, 0, 0, 0}, PutLengthEncodedInt(0xffffff+1))
	assert.Equal(t, []byte{0xfe, 1, 0, 0, 1, 0, 0, 0, 0}, PutLengthEncodedInt(0xffffff+2))
	assert.Equal(t, []byte{0xfe, 2, 0, 0, 1, 0, 0, 0, 0}, PutLengthEncodedInt(0xffffff+3))

	var (
		num    uint64
		isNull bool
		n      int
	)
	num, isNull, n = LengthEncodedInt([]byte{0x1})
	assert.Equal(t, num, uint64(1))
	assert.Equal(t, isNull, false)
	assert.Equal(t, n, 1)

	num, isNull, n = LengthEncodedInt([]byte{250})
	assert.Equal(t, num, uint64(250))
	assert.Equal(t, isNull, false)
	assert.Equal(t, n, 1)

	num, isNull, n = LengthEncodedInt([]byte{252, 251, 0})
	assert.Equal(t, num, uint64(251))
	assert.Equal(t, isNull, false)
	assert.Equal(t, n, 3)

	num, isNull, n = LengthEncodedInt([]byte{253, 3, 0, 1})
	assert.Equal(t, num, uint64(0xffff+4))
	assert.Equal(t, isNull, false)
	assert.Equal(t, n, 4)

	num, isNull, n = LengthEncodedInt([]byte{0xfe, 2, 0, 0, 1, 0, 0, 0, 0})
	assert.Equal(t, num, uint64(0xffffff+3))
	assert.Equal(t, isNull, false)
	assert.Equal(t, n, 9)
}

func TestPutLengthEncodedString(t *testing.T) {
	assert.Equal(t, []byte{0}, PutLengthEncodedString([]byte{}))
	assert.Equal(t, []byte{2, 0x61, 0x62}, PutLengthEncodedString([]byte{'a', 'b'}))

	var (
		strGot []byte
		isNull bool
		n      int
		err    error
	)
	strGot, isNull, n, err = LengthEnodedString([]byte{2, 0x61, 0x62})
	assert.Equal(t, strGot, []byte{'a', 'b'})
	assert.Equal(t, isNull, false)
	assert.Equal(t, n, 3)
	assert.NoError(t, err)
}
