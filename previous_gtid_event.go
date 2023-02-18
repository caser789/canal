package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// PreviousGTIDsEvent https://cloud.tencent.com/developer/article/1396292
//
// PREVIOUS GTID EVENT 是包含在每一个 BINLOG 的开头 用于描述所有以前 BINLOG 所包含的全部 GTID 的一个集合(包括已经删除的 BINLOG )
type PreviousGTIDsEvent struct {
	GTIDSets string
}

func (e *PreviousGTIDsEvent) Decode(data []byte) error {
	var previousGTIDSets []string
	pos := 0
	uuidCount := binary.LittleEndian.Uint16(data[pos : pos+8])
	pos += 8

	for i := uint16(0); i < uuidCount; i++ {
		uuid := e.decodeUuid(data[pos : pos+16])
		pos += 16
		sliceCount := binary.LittleEndian.Uint16(data[pos : pos+8])
		pos += 8
		var intervals []string
		for i := uint16(0); i < sliceCount; i++ {
			start := e.decodeInterval(data[pos : pos+8])
			pos += 8
			stop := e.decodeInterval(data[pos : pos+8])
			pos += 8
			interval := ""
			if stop == start+1 {
				interval = fmt.Sprintf("%d", start)
			} else {
				interval = fmt.Sprintf("%d-%d", start, stop-1)
			}
			intervals = append(intervals, interval)
		}
		previousGTIDSets = append(previousGTIDSets, fmt.Sprintf("%s:%s", uuid, strings.Join(intervals, ":")))
	}
	e.GTIDSets = strings.Join(previousGTIDSets, ",")
	return nil
}

func (e *PreviousGTIDsEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Previous GTID Event: %s\n", e.GTIDSets)
	fmt.Fprintln(w)
}

func (e *PreviousGTIDsEvent) decodeUuid(data []byte) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		hex.EncodeToString(data[0:4]),
		hex.EncodeToString(data[4:6]),
		hex.EncodeToString(data[6:8]),
		hex.EncodeToString(data[8:10]),
		hex.EncodeToString(data[10:]),
	)
}

func (e *PreviousGTIDsEvent) decodeInterval(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}
