package main

import (
	"encoding/binary"
	"fmt"
	"github.com/pingcap/errors"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var (
	checksumVersionSplitMysql   = []int{5, 6, 1}
	checksumVersionProductMysql = (checksumVersionSplitMysql[0]*256+checksumVersionSplitMysql[1])*256 + checksumVersionSplitMysql[2]

	checksumVersionSplitMariaDB   = []int{5, 3, 0}
	checksumVersionProductMariaDB = (checksumVersionSplitMariaDB[0]*256+checksumVersionSplitMariaDB[1])*256 + checksumVersionSplitMariaDB[2]
)

const (
	BINLOG_CHECKSUM_ALG_OFF byte = 0 // Events are without checksum though its generator
	// is checksum-capable New Master (NM).
	BINLOG_CHECKSUM_ALG_CRC32 byte = 1 // CRC32 of zlib algorithm.
	//  BINLOG_CHECKSUM_ALG_ENUM_END,  // the cut line: valid alg range is [1, 0x7f].
	BINLOG_CHECKSUM_ALG_UNDEF byte = 255 // special value to tag undetermined yet checksum
	// or events from checksum-unaware servers
)

// FormatDescriptionEvent https://mariadb.com/kb/en/format_description_event/
//
// This is a descriptor event that is written to the beginning of a binary log file,
// at position 4 (after the 4 magic number bytes)
//
// The whole event written to disk is byte<19> event header + data fields
type FormatDescriptionEvent struct {
	Version uint16 // 2 bytes The binary log format version.
	// 50 bytes
	// The server version (example: 10.2.1-debug-log), padded with 0x00 bytes on the right.
	ServerVersion []byte
	// Timestamp in seconds when this event was created (this is the moment when the binary log was created).
	// This value is redundant; the same value occurs in the timestamp header field.
	// 4 bytes
	CreateTimestamp uint32
	// The header length. 1 byte
	// This length - 19 gives the size of the extra headers field at the end of the header for other events.
	EventHeaderLength uint8
	// Variable-sized. An array that indicates the post-header lengths for all event types.
	// There is one byte per event type that the server knows about.
	EventTypeHeaderLengths []byte

	// 1 byte Checksum Algorithm Type
	// 0 is off, 1 is for CRC32, 255 is undefined
	ChecksumAlgorithm byte
}

func (e *FormatDescriptionEvent) Decode(data []byte) error {
	pos := 0
	e.Version = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	e.ServerVersion = make([]byte, 50)
	copy(e.ServerVersion, data[pos:])
	pos += 50

	e.CreateTimestamp = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	e.EventHeaderLength = data[pos]
	pos++

	if e.EventHeaderLength != byte(EventHeaderSize) {
		return errors.Errorf("invalid event header length %d, must 19", e.EventHeaderLength)
	}

	server := string(e.ServerVersion)
	checksumProduct := checksumVersionProductMysql
	if strings.Contains(strings.ToLower(server), "mariadb") {
		checksumProduct = checksumVersionProductMariaDB
	}

	if calcVersionProduct(string(e.ServerVersion)) >= checksumProduct {
		// here, the last 5 bytes is 1 byte check sum alg type and 4 byte checksum if exists
		e.ChecksumAlgorithm = data[len(data)-5]
		e.EventTypeHeaderLengths = data[pos : len(data)-5]
	} else {
		e.ChecksumAlgorithm = BINLOG_CHECKSUM_ALG_UNDEF
		e.EventTypeHeaderLengths = data[pos:]
	}

	return nil
}

func (e *FormatDescriptionEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "Version: %d\n", e.Version)
	fmt.Fprintf(w, "Server version: %s\n", e.ServerVersion)
	//fmt.Fprintf(w, "Create date: %s\n", time.Unix(int64(e.CreateTimestamp), 0).Format(TimeFormat))
	fmt.Fprintf(w, "Checksum algorithm: %d\n", e.ChecksumAlgorithm)
	//fmt.Fprintf(w, "Event header lengths: \n%s", hex.Dump(e.EventTypeHeaderLengths))
	fmt.Fprintln(w)
}

func calcVersionProduct(server string) int {
	versionSplit := splitServerVersion(server)

	return (versionSplit[0]*256+versionSplit[1])*256 + versionSplit[2]
}

// server version format X.Y.Zabc, a is not . or number
func splitServerVersion(server string) []int {
	seps := strings.Split(server, ".")
	if len(seps) < 3 {
		return []int{0, 0, 0}
	}

	x, _ := strconv.Atoi(seps[0])
	y, _ := strconv.Atoi(seps[1])

	index := 0
	for i, c := range seps[2] {
		if !unicode.IsNumber(c) {
			index = i
			break
		}
	}

	z, _ := strconv.Atoi(seps[2][0:index])

	return []int{x, y, z}
}
