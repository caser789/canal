package main

import (
	"encoding/binary"
	"fmt"
	"github.com/pingcap/errors"
	"math"
)

func ParseBinaryInt8(data []byte) int8 {
	return int8(data[0])
}
func ParseBinaryUint8(data []byte) uint8 {
	return data[0]
}

func ParseBinaryInt16(data []byte) int16 {
	return int16(binary.LittleEndian.Uint16(data))
}
func ParseBinaryUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func ParseBinaryInt24(data []byte) int32 {
	u32 := ParseBinaryUint24(data)
	if u32&0x00800000 != 0 {
		u32 |= 0xFF000000
	}
	return int32(u32)
}
func ParseBinaryUint24(data []byte) uint32 {
	return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16
}

func ParseBinaryInt32(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}
func ParseBinaryUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func ParseBinaryInt64(data []byte) int64 {
	return int64(binary.LittleEndian.Uint64(data))
}
func ParseBinaryUint64(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func ParseBinaryFloat32(data []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data))
}

func ParseBinaryFloat64(data []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(data))
}

func FormatBinaryDate(n int, data []byte) ([]byte, error) {
	switch n {
	case 0:
		return []byte("0000-00-00"), nil
	case 4:
		return []byte(fmt.Sprintf("%04d-%02d-%02d",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3])), nil
	default:
		return nil, errors.Errorf("invalid date packet length %d", n)
	}
}

func FormatBinaryDateTime(n int, data []byte) ([]byte, error) {
	switch n {
	case 0:
		return []byte("0000-00-00 00:00:00"), nil
	case 4:
		return []byte(fmt.Sprintf("%04d-%02d-%02d 00:00:00",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3])), nil
	case 7:
		return []byte(fmt.Sprintf(
			"%04d-%02d-%02d %02d:%02d:%02d",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3],
			data[4],
			data[5],
			data[6])), nil
	case 11:
		return []byte(fmt.Sprintf(
			"%04d-%02d-%02d %02d:%02d:%02d.%06d",
			binary.LittleEndian.Uint16(data[:2]),
			data[2],
			data[3],
			data[4],
			data[5],
			data[6],
			binary.LittleEndian.Uint32(data[7:11]))), nil
	default:
		return nil, errors.Errorf("invalid datetime packet length %d", n)
	}
}

func FormatBinaryTime(n int, data []byte) ([]byte, error) {
	if n == 0 {
		return []byte("0000-00-00"), nil
	}

	var sign byte
	if data[0] == 1 {
		sign = byte('-')
	}

	switch n {
	case 8:
		return []byte(fmt.Sprintf(
			"%c%02d:%02d:%02d",
			sign,
			uint16(data[1])*24+uint16(data[5]),
			data[6],
			data[7],
		)), nil
	case 12:
		return []byte(fmt.Sprintf(
			"%c%02d:%02d:%02d.%06d",
			sign,
			uint16(data[1])*24+uint16(data[5]),
			data[6],
			data[7],
			binary.LittleEndian.Uint32(data[8:12]),
		)), nil
	default:
		return nil, errors.Errorf("invalid time packet length %d", n)
	}
}
