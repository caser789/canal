package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/pingcap/errors"
	"github.com/shopspring/decimal"
	"io"
	"strconv"
	"strings"
	"time"
)

var errMissingTableMapEvent = errors.New("invalid table id, no corresponding table map event")

// RowsEventStmtEndFlag is set in the end of the statement.
const RowsEventStmtEndFlag = 0x01

type RowsEvent struct {
	//0, 1, 2
	Version int

	tableIDSize int
	tables      map[uint64]*TableMapEvent
	needBitmap2 bool

	Table *TableMapEvent

	TableID uint64

	Flags uint16

	//if version == 2
	ExtraData []byte

	//lenenc_int
	ColumnCount uint64

	/*
		By default MySQL and MariaDB log the full row image.
		see
			- https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html#sysvar_binlog_row_image
			- https://mariadb.com/kb/en/replication-and-binary-log-system-variables/#binlog_row_image

		ColumnBitmap1, ColumnBitmap2 and SkippedColumns are not set on the full row image.
	*/

	//len = (ColumnCount + 7) / 8
	ColumnBitmap1 []byte

	//if UPDATE_ROWS_EVENTv1 or v2
	//len = (ColumnCount + 7) / 8
	ColumnBitmap2 []byte

	//rows: invalid: int64, float64, bool, []byte, string
	Rows           [][]interface{}
	SkippedColumns [][]int

	parseTime               bool
	timestampStringLocation *time.Location
	useDecimal              bool
	ignoreJSONDecodeErr     bool
}

func (e *RowsEvent) DecodeHeader(data []byte) (int, error) {
	pos := 0
	e.TableID = FixedLengthInt(data[0:e.tableIDSize])
	pos += e.tableIDSize

	e.Flags = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if e.Version == 2 {
		dataLen := binary.LittleEndian.Uint16(data[pos:])
		pos += 2

		e.ExtraData = data[pos : pos+int(dataLen-2)]
		pos += int(dataLen - 2)
	}

	var n int
	e.ColumnCount, _, n = LengthEncodedInt(data[pos:])
	pos += n

	bitCount := bitmapByteSize(int(e.ColumnCount))
	e.ColumnBitmap1 = data[pos : pos+bitCount]
	pos += bitCount

	if e.needBitmap2 {
		e.ColumnBitmap2 = data[pos : pos+bitCount]
		pos += bitCount
	}

	var ok bool
	e.Table, ok = e.tables[e.TableID]
	if !ok {
		if len(e.tables) > 0 {
			return 0, errors.Errorf("invalid table id %d, no corresponding table map event", e.TableID)
		} else {
			return 0, errors.Annotatef(errMissingTableMapEvent, "table id %d", e.TableID)
		}
	}
	return pos, nil
}

func (e *RowsEvent) DecodeData(pos int, data []byte) (err2 error) {
	var (
		n   int
		err error
	)
	// ... repeat rows until event-end
	defer func() {
		if r := recover(); r != nil {
			err2 = errors.Errorf("parse rows event panic %v, data %q, parsed rows %#v, table map %#v", r, data, e, e.Table)
		}
	}()

	// Pre-allocate memory for rows.
	rowsLen := e.ColumnCount
	if e.needBitmap2 {
		rowsLen += e.ColumnCount
	}
	e.SkippedColumns = make([][]int, 0, rowsLen)
	e.Rows = make([][]interface{}, 0, rowsLen)

	for pos < len(data) {
		if n, err = e.decodeRows(data[pos:], e.Table, e.ColumnBitmap1); err != nil {
			return errors.Trace(err)
		}
		pos += n

		if e.needBitmap2 {
			if n, err = e.decodeRows(data[pos:], e.Table, e.ColumnBitmap2); err != nil {
				return errors.Trace(err)
			}
			pos += n
		}
	}

	return nil
}

func (e *RowsEvent) Decode(data []byte) error {
	pos, err := e.DecodeHeader(data)
	if err != nil {
		return err
	}
	return e.DecodeData(pos, data)
}

func isBitSet(bitmap []byte, i int) bool {
	return bitmap[i>>3]&(1<<(uint(i)&7)) > 0
}

func (e *RowsEvent) decodeRows(data []byte, table *TableMapEvent, bitmap []byte) (int, error) {
	row := make([]interface{}, e.ColumnCount)
	skips := make([]int, 0)

	pos := 0

	// refer: https://github.com/alibaba/canal/blob/c3e38e50e269adafdd38a48c63a1740cde304c67/dbsync/src/main/java/com/taobao/tddl/dbsync/binlog/event/RowsLogBuffer.java#L63
	count := 0
	for i := 0; i < int(e.ColumnCount); i++ {
		if isBitSet(bitmap, i) {
			count++
		}
	}
	count = (count + 7) / 8

	nullBitmap := data[pos : pos+count]
	pos += count

	nullbitIndex := 0

	var n int
	var err error
	for i := 0; i < int(e.ColumnCount); i++ {
		if !isBitSet(bitmap, i) {
			skips = append(skips, i)
			continue
		}

		isNull := (uint32(nullBitmap[nullbitIndex/8]) >> uint32(nullbitIndex%8)) & 0x01
		nullbitIndex++

		if isNull > 0 {
			row[i] = nil
			continue
		}

		row[i], n, err = e.decodeValue(data[pos:], table.ColumnType[i], table.ColumnMeta[i])

		if err != nil {
			return 0, err
		}
		pos += n
	}

	e.Rows = append(e.Rows, row)
	e.SkippedColumns = append(e.SkippedColumns, skips)
	return pos, nil
}

func (e *RowsEvent) parseFracTime(t interface{}) interface{} {
	v, ok := t.(fracTime)
	if !ok {
		return t
	}

	if !e.parseTime {
		// Don't parse time, return string directly
		return v.String()
	}

	// return Golang time directly
	return v.Time
}

// see mysql sql/log_event.cc log_event_print_value
func (e *RowsEvent) decodeValue(data []byte, tp byte, meta uint16) (v interface{}, n int, err error) {
	var length = 0

	if tp == MYSQL_TYPE_STRING {
		if meta >= 256 {
			b0 := uint8(meta >> 8)
			b1 := uint8(meta & 0xFF)

			if b0&0x30 != 0x30 {
				length = int(uint16(b1) | (uint16((b0&0x30)^0x30) << 4))
				tp = b0 | 0x30
			} else {
				length = int(meta & 0xFF)
				tp = b0
			}
		} else {
			length = int(meta)
		}
	}

	switch tp {
	case MYSQL_TYPE_NULL:
		return nil, 0, nil
	case MYSQL_TYPE_LONG:
		n = 4
		v = ParseBinaryInt32(data)
	case MYSQL_TYPE_TINY:
		n = 1
		v = ParseBinaryInt8(data)
	case MYSQL_TYPE_SHORT:
		n = 2
		v = ParseBinaryInt16(data)
	case MYSQL_TYPE_INT24:
		n = 3
		v = ParseBinaryInt24(data)
	case MYSQL_TYPE_LONGLONG:
		n = 8
		v = ParseBinaryInt64(data)
	case MYSQL_TYPE_NEWDECIMAL:
		prec := uint8(meta >> 8)
		scale := uint8(meta & 0xFF)
		v, n, err = decodeDecimal(data, int(prec), int(scale), e.useDecimal)
	case MYSQL_TYPE_FLOAT:
		n = 4
		v = ParseBinaryFloat32(data)
	case MYSQL_TYPE_DOUBLE:
		n = 8
		v = ParseBinaryFloat64(data)
	case MYSQL_TYPE_BIT:
		nbits := ((meta >> 8) * 8) + (meta & 0xFF)
		n = int(nbits+7) / 8

		//use int64 for bit
		v, err = decodeBit(data, int(nbits), n)
	case MYSQL_TYPE_TIMESTAMP:
		n = 4
		t := binary.LittleEndian.Uint32(data)
		if t == 0 {
			v = formatZeroTime(0, 0)
		} else {
			v = e.parseFracTime(fracTime{
				Time:                    time.Unix(int64(t), 0),
				Dec:                     0,
				timestampStringLocation: e.timestampStringLocation,
			})
		}
	case MYSQL_TYPE_TIMESTAMP2:
		v, n, err = decodeTimestamp2(data, meta, e.timestampStringLocation)
		v = e.parseFracTime(v)
	case MYSQL_TYPE_DATETIME:
		n = 8
		i64 := binary.LittleEndian.Uint64(data)
		if i64 == 0 {
			v = formatZeroTime(0, 0)
		} else {
			d := i64 / 1000000
			t := i64 % 1000000
			v = e.parseFracTime(fracTime{
				Time: time.Date(
					int(d/10000),
					time.Month((d%10000)/100),
					int(d%100),
					int(t/10000),
					int((t%10000)/100),
					int(t%100),
					0,
					time.UTC,
				),
				Dec: 0,
			})
		}
	case MYSQL_TYPE_DATETIME2:
		v, n, err = decodeDatetime2(data, meta)
		v = e.parseFracTime(v)
	case MYSQL_TYPE_TIME:
		n = 3
		i32 := uint32(FixedLengthInt(data[0:3]))
		if i32 == 0 {
			v = "00:00:00"
		} else {
			v = fmt.Sprintf("%02d:%02d:%02d", i32/10000, (i32%10000)/100, i32%100)
		}
	case MYSQL_TYPE_TIME2:
		v, n, err = decodeTime2(data, meta)
	case MYSQL_TYPE_DATE:
		n = 3
		i32 := uint32(FixedLengthInt(data[0:3]))
		if i32 == 0 {
			v = "0000-00-00"
		} else {
			v = fmt.Sprintf("%04d-%02d-%02d", i32/(16*32), i32/32%16, i32%32)
		}

	case MYSQL_TYPE_YEAR:
		n = 1
		year := int(data[0])
		if year == 0 {
			v = year
		} else {
			v = year + 1900
		}
	case MYSQL_TYPE_ENUM:
		l := meta & 0xFF
		switch l {
		case 1:
			v = int64(data[0])
			n = 1
		case 2:
			v = int64(binary.LittleEndian.Uint16(data))
			n = 2
		default:
			err = fmt.Errorf("Unknown ENUM packlen=%d", l)
		}
	case MYSQL_TYPE_SET:
		n = int(meta & 0xFF)
		nbits := n * 8

		v, err = littleDecodeBit(data, nbits, n)
	case MYSQL_TYPE_BLOB:
		v, n, err = decodeBlob(data, meta)
	case MYSQL_TYPE_VARCHAR,
		MYSQL_TYPE_VAR_STRING:
		length = int(meta)
		v, n = decodeString(data, length)
	case MYSQL_TYPE_STRING:
		v, n = decodeString(data, length)
	case MYSQL_TYPE_JSON:
		// Refer: https://github.com/shyiko/mysql-binlog-connector-java/blob/master/src/main/java/com/github/shyiko/mysql/binlog/event/deserialization/AbstractRowsEventDataDeserializer.java#L404
		length = int(FixedLengthInt(data[0:meta]))
		n = length + int(meta)
		var d []byte
		d, err = e.decodeJsonBinary(data[meta:n])
		if err == nil {
			v = BytesToString(d)
		}
	case MYSQL_TYPE_GEOMETRY:
		// MySQL saves Geometry as Blob in binlog
		// Seem that the binary format is SRID (4 bytes) + WKB, outer can use
		// MySQL GeoFromWKB or others to create the geometry data.
		// Refer https://dev.mysql.com/doc/refman/5.7/en/gis-wkb-functions.html
		// I also find some go libs to handle WKB if possible
		// see https://github.com/twpayne/go-geom or https://github.com/paulmach/go.geo
		v, n, err = decodeBlob(data, meta)
	default:
		err = fmt.Errorf("unsupport type %d in binlog and don't know how to handle", tp)
	}

	return v, n, err
}

func decodeString(data []byte, length int) (v string, n int) {
	if length < 256 {
		length = int(data[0])

		n = length + 1
		v = BytesToString(data[1:n])
	} else {
		length = int(binary.LittleEndian.Uint16(data[0:]))
		n = length + 2
		v = BytesToString(data[2:n])
	}

	return
}

// ref: https://github.com/mysql/mysql-server/blob/a9b0c712de3509d8d08d3ba385d41a4df6348775/strings/decimal.c#L137
const digitsPerInteger int = 9

var compressedBytes = []int{0, 1, 1, 2, 2, 3, 3, 4, 4, 4}

func decodeDecimalDecompressValue(compIndx int, data []byte, mask uint8) (size int, value uint32) {
	size = compressedBytes[compIndx]
	switch size {
	case 0:
	case 1:
		value = uint32(data[0] ^ mask)
	case 2:
		value = uint32(data[1]^mask) | uint32(data[0]^mask)<<8
	case 3:
		value = uint32(data[2]^mask) | uint32(data[1]^mask)<<8 | uint32(data[0]^mask)<<16
	case 4:
		value = uint32(data[3]^mask) | uint32(data[2]^mask)<<8 | uint32(data[1]^mask)<<16 | uint32(data[0]^mask)<<24
	}
	return
}

var zeros = [digitsPerInteger]byte{48, 48, 48, 48, 48, 48, 48, 48, 48}

func decodeDecimal(data []byte, precision int, decimals int, useDecimal bool) (interface{}, int, error) {
	//see python mysql replication and https://github.com/jeremycole/mysql_binlog
	integral := precision - decimals
	uncompIntegral := integral / digitsPerInteger
	uncompFractional := decimals / digitsPerInteger
	compIntegral := integral - (uncompIntegral * digitsPerInteger)
	compFractional := decimals - (uncompFractional * digitsPerInteger)

	binSize := uncompIntegral*4 + compressedBytes[compIntegral] +
		uncompFractional*4 + compressedBytes[compFractional]

	buf := make([]byte, binSize)
	copy(buf, data[:binSize])

	//must copy the data for later change
	data = buf

	// Support negative
	// The sign is encoded in the high bit of the the byte
	// But this bit can also be used in the value
	value := uint32(data[0])
	var res strings.Builder
	res.Grow(precision + 2)
	var mask uint32 = 0
	if value&0x80 == 0 {
		mask = uint32((1 << 32) - 1)
		res.WriteString("-")
	}

	//clear sign
	data[0] ^= 0x80

	zeroLeading := true

	pos, value := decodeDecimalDecompressValue(compIntegral, data, uint8(mask))
	if value != 0 {
		zeroLeading = false
		res.WriteString(strconv.FormatUint(uint64(value), 10))
	}

	for i := 0; i < uncompIntegral; i++ {
		value = binary.BigEndian.Uint32(data[pos:]) ^ mask
		pos += 4
		if zeroLeading {
			if value != 0 {
				zeroLeading = false
				res.WriteString(strconv.FormatUint(uint64(value), 10))
			}
		} else {
			toWrite := strconv.FormatUint(uint64(value), 10)
			res.Write(zeros[:digitsPerInteger-len(toWrite)])
			res.WriteString(toWrite)
		}
	}

	if zeroLeading {
		res.WriteString("0")
	}

	if pos < len(data) {
		res.WriteString(".")

		for i := 0; i < uncompFractional; i++ {
			value = binary.BigEndian.Uint32(data[pos:]) ^ mask
			pos += 4
			toWrite := strconv.FormatUint(uint64(value), 10)
			res.Write(zeros[:digitsPerInteger-len(toWrite)])
			res.WriteString(toWrite)
		}

		if size, value := decodeDecimalDecompressValue(compFractional, data[pos:], uint8(mask)); size > 0 {
			toWrite := strconv.FormatUint(uint64(value), 10)
			padding := compFractional - len(toWrite)
			if padding > 0 {
				res.Write(zeros[:padding])
			}
			res.WriteString(toWrite)
			pos += size
		}
	}

	if useDecimal {
		f, err := decimal.NewFromString(res.String())
		return f, pos, err
	}

	return res.String(), pos, nil
}

func decodeBit(data []byte, nbits int, length int) (value int64, err error) {
	if nbits > 1 {
		switch length {
		case 1:
			value = int64(data[0])
		case 2:
			value = int64(binary.BigEndian.Uint16(data))
		case 3:
			value = int64(BFixedLengthInt(data[0:3]))
		case 4:
			value = int64(binary.BigEndian.Uint32(data))
		case 5:
			value = int64(BFixedLengthInt(data[0:5]))
		case 6:
			value = int64(BFixedLengthInt(data[0:6]))
		case 7:
			value = int64(BFixedLengthInt(data[0:7]))
		case 8:
			value = int64(binary.BigEndian.Uint64(data))
		default:
			err = fmt.Errorf("invalid bit length %d", length)
		}
	} else {
		if length != 1 {
			err = fmt.Errorf("invalid bit length %d", length)
		} else {
			value = int64(data[0])
		}
	}
	return
}

func littleDecodeBit(data []byte, nbits int, length int) (value int64, err error) {
	if nbits > 1 {
		switch length {
		case 1:
			value = int64(data[0])
		case 2:
			value = int64(binary.LittleEndian.Uint16(data))
		case 3:
			value = int64(FixedLengthInt(data[0:3]))
		case 4:
			value = int64(binary.LittleEndian.Uint32(data))
		case 5:
			value = int64(FixedLengthInt(data[0:5]))
		case 6:
			value = int64(FixedLengthInt(data[0:6]))
		case 7:
			value = int64(FixedLengthInt(data[0:7]))
		case 8:
			value = int64(binary.LittleEndian.Uint64(data))
		default:
			err = fmt.Errorf("invalid bit length %d", length)
		}
	} else {
		if length != 1 {
			err = fmt.Errorf("invalid bit length %d", length)
		} else {
			value = int64(data[0])
		}
	}
	return
}

func decodeTimestamp2(data []byte, dec uint16, timestampStringLocation *time.Location) (interface{}, int, error) {
	//get timestamp binary length
	n := int(4 + (dec+1)/2)
	sec := int64(binary.BigEndian.Uint32(data[0:4]))
	usec := int64(0)
	switch dec {
	case 1, 2:
		usec = int64(data[4]) * 10000
	case 3, 4:
		usec = int64(binary.BigEndian.Uint16(data[4:])) * 100
	case 5, 6:
		usec = int64(BFixedLengthInt(data[4:7]))
	}

	if sec == 0 {
		return formatZeroTime(int(usec), int(dec)), n, nil
	}

	return fracTime{
		Time:                    time.Unix(sec, usec*1000),
		Dec:                     int(dec),
		timestampStringLocation: timestampStringLocation,
	}, n, nil
}

const DATETIMEF_INT_OFS int64 = 0x8000000000

func decodeDatetime2(data []byte, dec uint16) (interface{}, int, error) {
	//get datetime binary length
	n := int(5 + (dec+1)/2)

	intPart := int64(BFixedLengthInt(data[0:5])) - DATETIMEF_INT_OFS
	var frac int64 = 0

	switch dec {
	case 1, 2:
		frac = int64(data[5]) * 10000
	case 3, 4:
		frac = int64(binary.BigEndian.Uint16(data[5:7])) * 100
	case 5, 6:
		frac = int64(BFixedLengthInt(data[5:8]))
	}

	if intPart == 0 {
		return formatZeroTime(int(frac), int(dec)), n, nil
	}

	tmp := intPart<<24 + frac
	//handle sign???
	if tmp < 0 {
		tmp = -tmp
	}

	// var secPart int64 = tmp % (1 << 24)
	ymdhms := tmp >> 24

	ymd := ymdhms >> 17
	ym := ymd >> 5
	hms := ymdhms % (1 << 17)

	day := int(ymd % (1 << 5))
	month := int(ym % 13)
	year := int(ym / 13)

	second := int(hms % (1 << 6))
	minute := int((hms >> 6) % (1 << 6))
	hour := int(hms >> 12)

	// DATETIME encoding for nonfractional part after MySQL 5.6.4
	// https://dev.mysql.com/doc/internals/en/date-and-time-data-type-representation.html
	// integer value for 1970-01-01 00:00:00 is
	// year*13+month = 25611 = 0b110010000001011
	// day = 1 = 0b00001
	// hour = 0 = 0b00000
	// minute = 0 = 0b000000
	// second = 0 = 0b000000
	// integer value = 0b1100100000010110000100000000000000000 = 107420450816
	if intPart < 107420450816 {
		return formatBeforeUnixZeroTime(year, month, day, hour, minute, second, int(frac), int(dec)), n, nil
	}

	return fracTime{
		Time: time.Date(year, time.Month(month), day, hour, minute, second, int(frac*1000), time.UTC),
		Dec:  int(dec),
	}, n, nil
}

const TIMEF_OFS int64 = 0x800000000000
const TIMEF_INT_OFS int64 = 0x800000

func decodeTime2(data []byte, dec uint16) (string, int, error) {
	//time  binary length
	n := int(3 + (dec+1)/2)

	tmp := int64(0)
	intPart := int64(0)
	frac := int64(0)
	switch dec {
	case 1, 2:
		intPart = int64(BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		frac = int64(data[3])
		if intPart < 0 && frac != 0 {
			/*
			   Negative values are stored with reverse fractional part order,
			   for binary sort compatibility.

			     Disk value  intpart frac   Time value   Memory value
			     800000.00    0      0      00:00:00.00  0000000000.000000
			     7FFFFF.FF   -1      255   -00:00:00.01  FFFFFFFFFF.FFD8F0
			     7FFFFF.9D   -1      99    -00:00:00.99  FFFFFFFFFF.F0E4D0
			     7FFFFF.00   -1      0     -00:00:01.00  FFFFFFFFFF.000000
			     7FFFFE.FF   -1      255   -00:00:01.01  FFFFFFFFFE.FFD8F0
			     7FFFFE.F6   -2      246   -00:00:01.10  FFFFFFFFFE.FE7960

			     Formula to convert fractional part from disk format
			     (now stored in "frac" variable) to absolute value: "0x100 - frac".
			     To reconstruct in-memory value, we shift
			     to the next integer value and then substruct fractional part.
			*/
			intPart++     /* Shift to the next integer value */
			frac -= 0x100 /* -(0x100 - frac) */
		}
		tmp = intPart<<24 + frac*10000
	case 3, 4:
		intPart = int64(BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		frac = int64(binary.BigEndian.Uint16(data[3:5]))
		if intPart < 0 && frac != 0 {
			/*
			   Fix reverse fractional part order: "0x10000 - frac".
			   See comments for FSP=1 and FSP=2 above.
			*/
			intPart++       /* Shift to the next integer value */
			frac -= 0x10000 /* -(0x10000-frac) */
		}
		tmp = intPart<<24 + frac*100

	case 5, 6:
		tmp = int64(BFixedLengthInt(data[0:6])) - TIMEF_OFS
		return timeFormat(tmp, dec, n)
	default:
		intPart = int64(BFixedLengthInt(data[0:3])) - TIMEF_INT_OFS
		tmp = intPart << 24
	}

	if intPart == 0 && frac == 0 {
		return "00:00:00", n, nil
	}

	return timeFormat(tmp, dec, n)
}

func timeFormat(tmp int64, dec uint16, n int) (string, int, error) {
	hms := int64(0)
	sign := ""
	if tmp < 0 {
		tmp = -tmp
		sign = "-"
	}

	hms = tmp >> 24

	hour := (hms >> 12) % (1 << 10) /* 10 bits starting at 12th */
	minute := (hms >> 6) % (1 << 6) /* 6 bits starting at 6th   */
	second := hms % (1 << 6)        /* 6 bits starting at 0th   */
	secPart := tmp % (1 << 24)

	if secPart != 0 {
		s := fmt.Sprintf("%s%02d:%02d:%02d.%06d", sign, hour, minute, second, secPart)
		return s[0 : len(s)-(6-int(dec))], n, nil
	}

	return fmt.Sprintf("%s%02d:%02d:%02d", sign, hour, minute, second), n, nil
}

func decodeBlob(data []byte, meta uint16) (v []byte, n int, err error) {
	var length int
	switch meta {
	case 1:
		length = int(data[0])
		v = data[1 : 1+length]
		n = length + 1
	case 2:
		length = int(binary.LittleEndian.Uint16(data))
		v = data[2 : 2+length]
		n = length + 2
	case 3:
		length = int(FixedLengthInt(data[0:3]))
		v = data[3 : 3+length]
		n = length + 3
	case 4:
		length = int(binary.LittleEndian.Uint32(data))
		v = data[4 : 4+length]
		n = length + 4
	default:
		err = fmt.Errorf("invalid blob packlen = %d", meta)
	}

	return
}

func (e *RowsEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "TableID: %d\n", e.TableID)
	fmt.Fprintf(w, "Flags: %d\n", e.Flags)
	fmt.Fprintf(w, "Column count: %d\n", e.ColumnCount)

	fmt.Fprintf(w, "Values:\n")
	for _, rows := range e.Rows {
		fmt.Fprintf(w, "--\n")
		for j, d := range rows {
			if _, ok := d.([]byte); ok {
				fmt.Fprintf(w, "%d:%q\n", j, d)
			} else {
				fmt.Fprintf(w, "%d:%#v\n", j, d)
			}
		}
	}
	fmt.Fprintln(w)
}

// decodeJsonBinary decodes the JSON binary encoding data and returns
// the common JSON encoding data.
func (e *RowsEvent) decodeJsonBinary(data []byte) ([]byte, error) {
	// Sometimes, we can insert a NULL JSON even we set the JSON field as NOT NULL.
	// If we meet this case, we can return an empty slice.
	if len(data) == 0 {
		return []byte{}, nil
	}
	d := jsonBinaryDecoder{
		useDecimal:      e.useDecimal,
		ignoreDecodeErr: e.ignoreJSONDecodeErr,
	}

	if d.isDataShort(data, 1) {
		return nil, d.err
	}

	v := d.decodeValue(data[0], data[1:])
	if d.err != nil {
		return nil, d.err
	}

	return json.Marshal(v)
}
