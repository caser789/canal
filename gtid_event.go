package main

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

// GTIDEvent https://mariadb.com/kb/en/gtid_event/
// For global transaction ID, used to start a new transaction event group,
// instead of the old BEGIN query event, and also to mark stand-alone (ddl).
type GTIDEvent struct {
	CommitFlag     uint8
	SID            []byte
	GNO            int64
	LastCommitted  int64
	SequenceNumber int64

	// ImmediateCommitTimestamp/OriginalCommitTimestamp are introduced in MySQL-8.0.1, see:
	// https://mysqlhighavailability.com/replication-features-in-mysql-8-0-1/
	ImmediateCommitTimestamp uint64
	OriginalCommitTimestamp  uint64

	// Total transaction length (including this GTIDEvent), introduced in MySQL-8.0.2, see:
	// https://mysqlhighavailability.com/taking-advantage-of-new-transaction-length-metadata/
	TransactionLength uint64

	// ImmediateServerVersion/OriginalServerVersion are introduced in MySQL-8.0.14, see
	// https://dev.mysql.com/doc/refman/8.0/en/replication-compatibility.html
	ImmediateServerVersion uint32
	OriginalServerVersion  uint32
}

func (e *GTIDEvent) Decode(data []byte) error {
	pos := 0
	e.CommitFlag = data[pos]
	pos++
	e.SID = data[pos : pos+SidLength]
	pos += SidLength
	e.GNO = int64(binary.LittleEndian.Uint64(data[pos:]))
	pos += 8

	if len(data) >= 42 {
		if data[pos] == LogicalTimestampTypeCode {
			pos++
			e.LastCommitted = int64(binary.LittleEndian.Uint64(data[pos:]))
			pos += PartLogicalTimestampLength
			e.SequenceNumber = int64(binary.LittleEndian.Uint64(data[pos:]))
			pos += 8

			// IMMEDIATE_COMMIT_TIMESTAMP_LENGTH = 7
			if len(data)-pos < 7 {
				return nil
			}
			e.ImmediateCommitTimestamp = FixedLengthInt(data[pos : pos+7])
			pos += 7
			if (e.ImmediateCommitTimestamp & (uint64(1) << 55)) != 0 {
				// If the most significant bit set, another 7 byte follows representing OriginalCommitTimestamp
				e.ImmediateCommitTimestamp &= ^(uint64(1) << 55)
				e.OriginalCommitTimestamp = FixedLengthInt(data[pos : pos+7])
				pos += 7
			} else {
				// Otherwise OriginalCommitTimestamp == ImmediateCommitTimestamp
				e.OriginalCommitTimestamp = e.ImmediateCommitTimestamp
			}

			// TRANSACTION_LENGTH_MIN_LENGTH = 1
			if len(data)-pos < 1 {
				return nil
			}
			var n int
			e.TransactionLength, _, n = LengthEncodedInt(data[pos:])
			pos += n

			// IMMEDIATE_SERVER_VERSION_LENGTH = 4
			e.ImmediateServerVersion = UndefinedServerVer
			e.OriginalServerVersion = UndefinedServerVer
			if len(data)-pos < 4 {
				return nil
			}
			e.ImmediateServerVersion = binary.LittleEndian.Uint32(data[pos:])
			pos += 4
			if (e.ImmediateServerVersion & (uint32(1) << 31)) != 0 {
				// If the most significant bit set, another 4 byte follows representing OriginalServerVersion
				e.ImmediateServerVersion &= ^(uint32(1) << 31)
				e.OriginalServerVersion = binary.LittleEndian.Uint32(data[pos:])
				// pos += 4
			} else {
				// Otherwise OriginalServerVersion == ImmediateServerVersion
				e.OriginalServerVersion = e.ImmediateServerVersion
			}
		}
	}
	return nil
}

func (e *GTIDEvent) Dump(w io.Writer) {
	fmtTime := func(t time.Time) string {
		if t.IsZero() {
			return "<n/a>"
		}
		return t.Format(time.RFC3339Nano)
	}

	fmt.Fprintf(w, "Commit flag: %d\n", e.CommitFlag)
	u, _ := uuid.FromBytes(e.SID)
	fmt.Fprintf(w, "GTID_NEXT: %s:%d\n", u.String(), e.GNO)
	fmt.Fprintf(w, "LAST_COMMITTED: %d\n", e.LastCommitted)
	fmt.Fprintf(w, "SEQUENCE_NUMBER: %d\n", e.SequenceNumber)
	fmt.Fprintf(w, "Immediate commmit timestamp: %d (%s)\n", e.ImmediateCommitTimestamp, fmtTime(e.ImmediateCommitTime()))
	fmt.Fprintf(w, "Orignal commmit timestamp: %d (%s)\n", e.OriginalCommitTimestamp, fmtTime(e.OriginalCommitTime()))
	fmt.Fprintf(w, "Transaction length: %d\n", e.TransactionLength)
	fmt.Fprintf(w, "Immediate server version: %d\n", e.ImmediateServerVersion)
	fmt.Fprintf(w, "Orignal server version: %d\n", e.OriginalServerVersion)
	fmt.Fprintln(w)
}

// ImmediateCommitTime returns the commit time of this trx on the immediate server
// or zero time if not available.
func (e *GTIDEvent) ImmediateCommitTime() time.Time {
	return microSecTimestampToTime(e.ImmediateCommitTimestamp)
}

// OriginalCommitTime returns the commit time of this trx on the original server
// or zero time if not available.
func (e *GTIDEvent) OriginalCommitTime() time.Time {
	return microSecTimestampToTime(e.OriginalCommitTimestamp)
}

func microSecTimestampToTime(ts uint64) time.Time {
	if ts == 0 {
		return time.Time{}
	}
	return time.Unix(int64(ts/1000000), int64(ts%1000000)*1000)
}
