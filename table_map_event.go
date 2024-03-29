package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/pingcap/errors"
	"io"
)

// TableMapEvent https://mariadb.com/kb/en/table_map_event/
type TableMapEvent struct {
	flavor      string
	tableIDSize int

	TableID uint64

	Flags uint16

	Schema []byte // database name
	Table  []byte // table name

	ColumnCount uint64   // number of columns
	ColumnType  []byte   // type of each column, one byte for one column
	ColumnMeta  []uint16 // one byte for each column

	//len = (ColumnCount + 7) / 8
	// Bit-field indicating whether each column can be NULL, one bit per column.
	NullBitmap []byte

	/*
		The following are available only after MySQL-8.0.1 or MariaDB-10.5.0
		By default MySQL and MariaDB do not log the full row metadata.
		see:
			- https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html#sysvar_binlog_row_metadata
			- https://mariadb.com/kb/en/replication-and-binary-log-system-variables/#binlog_row_metadata
	*/

	// SignednessBitmap stores signedness info for numeric columns.
	SignednessBitmap []byte

	// DefaultCharset/ColumnCharset stores collation info for character columns.

	// DefaultCharset[0] is the default collation of character columns.
	// For character columns that have different charset,
	// (character column index, column collation) pairs follows
	DefaultCharset []uint64
	// ColumnCharset contains collation sequence for all character columns
	ColumnCharset []uint64

	// SetStrValue stores values for set columns.
	SetStrValue       [][][]byte
	setStrValueString [][]string

	// EnumStrValue stores values for enum columns.
	EnumStrValue       [][][]byte
	enumStrValueString [][]string

	// ColumnName list all column names.
	ColumnName       [][]byte
	columnNameString []string // the same as ColumnName in string type, just for reuse

	// GeometryType stores real type for geometry columns.
	GeometryType []uint64

	// PrimaryKey is a sequence of column indexes of primary key.
	PrimaryKey []uint64

	// PrimaryKeyPrefix is the prefix length used for each column of primary key.
	// 0 means that the whole column length is used.
	PrimaryKeyPrefix []uint64

	// EnumSetDefaultCharset/EnumSetColumnCharset is similar to DefaultCharset/ColumnCharset but for enum/set columns.
	EnumSetDefaultCharset []uint64
	EnumSetColumnCharset  []uint64
}

func (e *TableMapEvent) Decode(data []byte) error {
	pos := 0
	e.TableID = FixedLengthInt(data[0:e.tableIDSize])
	pos += e.tableIDSize

	e.Flags = binary.LittleEndian.Uint16(data[pos:])
	pos += 2 // not used

	schemaLength := data[pos]
	pos++ // database length

	e.Schema = data[pos : pos+int(schemaLength)]
	pos += int(schemaLength) // database name

	//skip 0x00 database name end will 0x00
	pos++

	tableLength := data[pos]
	pos++ // table name length

	e.Table = data[pos : pos+int(tableLength)]
	pos += int(tableLength) // table name

	//skip 0x00 table name ends with 0xxx
	pos++

	var n int
	e.ColumnCount, _, n = LengthEncodedInt(data[pos:]) // number of columns in the table
	pos += n

	e.ColumnType = data[pos : pos+int(e.ColumnCount)] // one byte for one column type
	pos += int(e.ColumnCount)

	var err error
	var metaData []byte
	if metaData, _, n, err = LengthEncodedString(data[pos:]); err != nil {
		return errors.Trace(err)
	}

	if err = e.decodeMeta(metaData); err != nil {
		return errors.Trace(err)
	}

	pos += n

	nullBitmapSize := bitmapByteSize(int(e.ColumnCount))
	if len(data[pos:]) < nullBitmapSize {
		return io.EOF
	}

	e.NullBitmap = data[pos : pos+nullBitmapSize]

	pos += nullBitmapSize

	if err = e.decodeOptionalMeta(data[pos:]); err != nil {
		return err
	}

	return nil
}

func bitmapByteSize(columnCount int) int {
	return (columnCount + 7) / 8
}

// see mysql sql/log_event.h
/*
	0 byte
	MYSQL_TYPE_DECIMAL
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR

	1 byte
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_BLOB
	MYSQL_TYPE_GEOMETRY

	//maybe
	MYSQL_TYPE_TIME2
	MYSQL_TYPE_DATETIME2
	MYSQL_TYPE_TIMESTAMP2

	2 byte
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT
	MYSQL_TYPE_NEWDECIMAL
	MYSQL_TYPE_VAR_STRING
	MYSQL_TYPE_STRING

	This enumeration value is only used internally and cannot exist in a binlog.
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_ENUM
	MYSQL_TYPE_SET
	MYSQL_TYPE_TINY_BLOB
	MYSQL_TYPE_MEDIUM_BLOB
	MYSQL_TYPE_LONG_BLOB
*/
func (e *TableMapEvent) decodeMeta(data []byte) error {
	pos := 0
	e.ColumnMeta = make([]uint16, e.ColumnCount)
	for i, t := range e.ColumnType {
		switch t {
		case MYSQL_TYPE_STRING:
			var x = uint16(data[pos]) << 8 //real type
			x += uint16(data[pos+1])       //pack or field length
			e.ColumnMeta[i] = x
			pos += 2
		case MYSQL_TYPE_NEWDECIMAL:
			var x = uint16(data[pos]) << 8 //precision
			x += uint16(data[pos+1])       //decimals
			e.ColumnMeta[i] = x
			pos += 2
		case MYSQL_TYPE_VAR_STRING,
			MYSQL_TYPE_VARCHAR,
			MYSQL_TYPE_BIT:
			e.ColumnMeta[i] = binary.LittleEndian.Uint16(data[pos:])
			pos += 2
		case MYSQL_TYPE_BLOB,
			MYSQL_TYPE_DOUBLE,
			MYSQL_TYPE_FLOAT,
			MYSQL_TYPE_GEOMETRY,
			MYSQL_TYPE_JSON:
			e.ColumnMeta[i] = uint16(data[pos])
			pos++
		case MYSQL_TYPE_TIME2,
			MYSQL_TYPE_DATETIME2,
			MYSQL_TYPE_TIMESTAMP2:
			e.ColumnMeta[i] = uint16(data[pos])
			pos++
		case MYSQL_TYPE_NEWDATE,
			MYSQL_TYPE_ENUM,
			MYSQL_TYPE_SET,
			MYSQL_TYPE_TINY_BLOB,
			MYSQL_TYPE_MEDIUM_BLOB,
			MYSQL_TYPE_LONG_BLOB:
			return errors.Errorf("unsupport type in binlog %d", t)
		default:
			e.ColumnMeta[i] = 0
		}
	}

	return nil
}

func (e *TableMapEvent) decodeOptionalMeta(data []byte) (err error) {
	pos := 0
	for pos < len(data) {
		// optional metadata fields are stored in Type, Length, Value(TLV) format
		// Type takes 1 byte. Length is a packed integer value. Values takes Length bytes
		t := data[pos]
		pos++

		l, _, n := LengthEncodedInt(data[pos:])
		pos += n

		v := data[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case TABLE_MAP_OPT_META_SIGNEDNESS:
			e.SignednessBitmap = v

		case TABLE_MAP_OPT_META_DEFAULT_CHARSET:
			e.DefaultCharset, err = e.decodeDefaultCharset(v)
			if err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_COLUMN_CHARSET:
			e.ColumnCharset, err = e.decodeIntSeq(v)
			if err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_COLUMN_NAME:
			if err = e.decodeColumnNames(v); err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_SET_STR_VALUE:
			e.SetStrValue, err = e.decodeStrValue(v)
			if err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_ENUM_STR_VALUE:
			e.EnumStrValue, err = e.decodeStrValue(v)
			if err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_GEOMETRY_TYPE:
			e.GeometryType, err = e.decodeIntSeq(v)
			if err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_SIMPLE_PRIMARY_KEY:
			if err = e.decodeSimplePrimaryKey(v); err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_PRIMARY_KEY_WITH_PREFIX:
			if err = e.decodePrimaryKeyWithPrefix(v); err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_ENUM_AND_SET_DEFAULT_CHARSET:
			e.EnumSetDefaultCharset, err = e.decodeDefaultCharset(v)
			if err != nil {
				return err
			}

		case TABLE_MAP_OPT_META_ENUM_AND_SET_COLUMN_CHARSET:
			e.EnumSetColumnCharset, err = e.decodeIntSeq(v)
			if err != nil {
				return err
			}

		default:
			// Ignore for future extension
		}
	}

	return nil
}

// These are TABLE_MAP_EVENT's optional metadata field type, from: libbinlogevents/include/rows_event.h
const (
	TABLE_MAP_OPT_META_SIGNEDNESS byte = iota + 1
	TABLE_MAP_OPT_META_DEFAULT_CHARSET
	TABLE_MAP_OPT_META_COLUMN_CHARSET
	TABLE_MAP_OPT_META_COLUMN_NAME
	TABLE_MAP_OPT_META_SET_STR_VALUE
	TABLE_MAP_OPT_META_ENUM_STR_VALUE
	TABLE_MAP_OPT_META_GEOMETRY_TYPE
	TABLE_MAP_OPT_META_SIMPLE_PRIMARY_KEY
	TABLE_MAP_OPT_META_PRIMARY_KEY_WITH_PREFIX
	TABLE_MAP_OPT_META_ENUM_AND_SET_DEFAULT_CHARSET
	TABLE_MAP_OPT_META_ENUM_AND_SET_COLUMN_CHARSET
)

func (e *TableMapEvent) decodeIntSeq(v []byte) (ret []uint64, err error) {
	p := 0
	for p < len(v) {
		i, _, n := LengthEncodedInt(v[p:])
		p += n
		ret = append(ret, i)
	}
	return
}

func (e *TableMapEvent) decodeDefaultCharset(v []byte) (ret []uint64, err error) {
	ret, err = e.decodeIntSeq(v)
	if err != nil {
		return
	}
	if len(ret)%2 != 1 {
		return nil, errors.Errorf("Expect odd item in DefaultCharset but got %d", len(ret))
	}
	return
}

func (e *TableMapEvent) decodeColumnNames(v []byte) error {
	p := 0
	e.ColumnName = make([][]byte, 0, e.ColumnCount)
	for p < len(v) {
		n := int(v[p])
		p++
		e.ColumnName = append(e.ColumnName, v[p:p+n])
		p += n
	}

	if len(e.ColumnName) != int(e.ColumnCount) {
		return errors.Errorf("Expect %d column names but got %d", e.ColumnCount, len(e.ColumnName))
	}
	return nil
}

func (e *TableMapEvent) decodeStrValue(v []byte) (ret [][][]byte, err error) {
	p := 0
	for p < len(v) {
		nVal, _, n := LengthEncodedInt(v[p:])
		p += n
		vals := make([][]byte, 0, int(nVal))
		for i := 0; i < int(nVal); i++ {
			val, _, n, err := LengthEncodedString(v[p:])
			if err != nil {
				return nil, err
			}
			p += n
			vals = append(vals, val)
		}
		ret = append(ret, vals)
	}
	return
}

func (e *TableMapEvent) decodeSimplePrimaryKey(v []byte) error {
	p := 0
	for p < len(v) {
		i, _, n := LengthEncodedInt(v[p:])
		e.PrimaryKey = append(e.PrimaryKey, i)
		e.PrimaryKeyPrefix = append(e.PrimaryKeyPrefix, 0)
		p += n
	}
	return nil
}

func (e *TableMapEvent) decodePrimaryKeyWithPrefix(v []byte) error {
	p := 0
	for p < len(v) {
		i, _, n := LengthEncodedInt(v[p:])
		e.PrimaryKey = append(e.PrimaryKey, i)
		p += n
		i, _, n = LengthEncodedInt(v[p:])
		e.PrimaryKeyPrefix = append(e.PrimaryKeyPrefix, i)
		p += n
	}
	return nil
}

func (e *TableMapEvent) Dump(w io.Writer) {
	fmt.Fprintf(w, "TableID: %d\n", e.TableID)
	fmt.Fprintf(w, "TableID size: %d\n", e.tableIDSize)
	fmt.Fprintf(w, "Flags: %d\n", e.Flags)
	fmt.Fprintf(w, "Schema: %s\n", e.Schema)
	fmt.Fprintf(w, "Table: %s\n", e.Table)
	fmt.Fprintf(w, "Column count: %d\n", e.ColumnCount)
	fmt.Fprintf(w, "Column type: \n%s", hex.Dump(e.ColumnType))
	fmt.Fprintf(w, "NULL bitmap: \n%s", hex.Dump(e.NullBitmap))

	fmt.Fprintf(w, "Signedness bitmap: \n%s", hex.Dump(e.SignednessBitmap))
	fmt.Fprintf(w, "Default charset: %v\n", e.DefaultCharset)
	fmt.Fprintf(w, "Column charset: %v\n", e.ColumnCharset)
	fmt.Fprintf(w, "Set str value: %v\n", e.SetStrValueString())
	fmt.Fprintf(w, "Enum str value: %v\n", e.EnumStrValueString())
	fmt.Fprintf(w, "Column name: %v\n", e.ColumnNameString())
	fmt.Fprintf(w, "Geometry type: %v\n", e.GeometryType)
	fmt.Fprintf(w, "Primary key: %v\n", e.PrimaryKey)
	fmt.Fprintf(w, "Primary key prefix: %v\n", e.PrimaryKeyPrefix)
	fmt.Fprintf(w, "Enum/set default charset: %v\n", e.EnumSetDefaultCharset)
	fmt.Fprintf(w, "Enum/set column charset: %v\n", e.EnumSetColumnCharset)

	unsignedMap := e.UnsignedMap()
	fmt.Fprintf(w, "UnsignedMap: %#v\n", unsignedMap)

	collationMap := e.CollationMap()
	fmt.Fprintf(w, "CollationMap: %#v\n", collationMap)

	enumSetCollationMap := e.EnumSetCollationMap()
	fmt.Fprintf(w, "EnumSetCollationMap: %#v\n", enumSetCollationMap)

	enumStrValueMap := e.EnumStrValueMap()
	fmt.Fprintf(w, "EnumStrValueMap: %#v\n", enumStrValueMap)

	setStrValueMap := e.SetStrValueMap()
	fmt.Fprintf(w, "SetStrValueMap: %#v\n", setStrValueMap)

	geometryTypeMap := e.GeometryTypeMap()
	fmt.Fprintf(w, "GeometryTypeMap: %#v\n", geometryTypeMap)

	nameMaxLen := 0
	for _, name := range e.ColumnName {
		if len(name) > nameMaxLen {
			nameMaxLen = len(name)
		}
	}
	nameFmt := "  %s"
	if nameMaxLen > 0 {
		nameFmt = fmt.Sprintf("  %%-%ds", nameMaxLen)
	}

	primaryKey := map[int]struct{}{}
	for _, pk := range e.PrimaryKey {
		primaryKey[int(pk)] = struct{}{}
	}

	fmt.Fprintf(w, "Columns: \n")
	for i := 0; i < int(e.ColumnCount); i++ {
		if len(e.ColumnName) == 0 {
			fmt.Fprintf(w, nameFmt, "<n/a>")
		} else {
			fmt.Fprintf(w, nameFmt, e.ColumnName[i])
		}

		fmt.Fprintf(w, "  type=%-3d", e.realType(i))

		if e.IsNumericColumn(i) {
			if len(unsignedMap) == 0 {
				fmt.Fprintf(w, "  unsigned=<n/a>")
			} else if unsignedMap[i] {
				fmt.Fprintf(w, "  unsigned=yes")
			} else {
				fmt.Fprintf(w, "  unsigned=no ")
			}
		}
		if e.IsCharacterColumn(i) {
			if len(collationMap) == 0 {
				fmt.Fprintf(w, "  collation=<n/a>")
			} else {
				fmt.Fprintf(w, "  collation=%d ", collationMap[i])
			}
		}
		if e.IsEnumColumn(i) {
			if len(enumSetCollationMap) == 0 {
				fmt.Fprintf(w, "  enum_collation=<n/a>")
			} else {
				fmt.Fprintf(w, "  enum_collation=%d", enumSetCollationMap[i])
			}

			if len(enumStrValueMap) == 0 {
				fmt.Fprintf(w, "  enum=<n/a>")
			} else {
				fmt.Fprintf(w, "  enum=%v", enumStrValueMap[i])
			}
		}
		if e.IsSetColumn(i) {
			if len(enumSetCollationMap) == 0 {
				fmt.Fprintf(w, "  set_collation=<n/a>")
			} else {
				fmt.Fprintf(w, "  set_collation=%d", enumSetCollationMap[i])
			}

			if len(setStrValueMap) == 0 {
				fmt.Fprintf(w, "  set=<n/a>")
			} else {
				fmt.Fprintf(w, "  set=%v", setStrValueMap[i])
			}
		}
		if e.IsGeometryColumn(i) {
			if len(geometryTypeMap) == 0 {
				fmt.Fprintf(w, "  geometry_type=<n/a>")
			} else {
				fmt.Fprintf(w, "  geometry_type=%v", geometryTypeMap[i])
			}
		}

		available, nullable := e.Nullable(i)
		if !available {
			fmt.Fprintf(w, "  null=<n/a>")
		} else if nullable {
			fmt.Fprintf(w, "  null=yes")
		} else {
			fmt.Fprintf(w, "  null=no ")
		}

		if _, ok := primaryKey[i]; ok {
			fmt.Fprintf(w, "  pri")
		}

		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintln(w)
}

// Nullable returns the nullablity of the i-th column.
// If null bits are not available, available is false.
// i must be in range [0, ColumnCount).
func (e *TableMapEvent) Nullable(i int) (available, nullable bool) {
	if len(e.NullBitmap) == 0 {
		return
	}
	return true, e.NullBitmap[i/8]&(1<<uint(i%8)) != 0
}

// SetStrValueString returns values for set columns as string slices.
// nil is returned if not available or no set columns at all.
func (e *TableMapEvent) SetStrValueString() [][]string {
	if e.setStrValueString == nil {
		if len(e.SetStrValue) == 0 {
			return nil
		}
		e.setStrValueString = make([][]string, 0, len(e.SetStrValue))
		for _, vals := range e.SetStrValue {
			e.setStrValueString = append(
				e.setStrValueString,
				e.bytesSlice2StrSlice(vals),
			)
		}
	}
	return e.setStrValueString
}

// EnumStrValueString returns values for enum columns as string slices.
// nil is returned if not available or no enum columns at all.
func (e *TableMapEvent) EnumStrValueString() [][]string {
	if e.enumStrValueString == nil {
		if len(e.EnumStrValue) == 0 {
			return nil
		}
		e.enumStrValueString = make([][]string, 0, len(e.EnumStrValue))
		for _, vals := range e.EnumStrValue {
			e.enumStrValueString = append(
				e.enumStrValueString,
				e.bytesSlice2StrSlice(vals),
			)
		}
	}
	return e.enumStrValueString
}

// ColumnNameString returns column names as string slice.
// nil is returned if not available.
func (e *TableMapEvent) ColumnNameString() []string {
	if e.columnNameString == nil {
		e.columnNameString = e.bytesSlice2StrSlice(e.ColumnName)
	}
	return e.columnNameString
}

func (e *TableMapEvent) bytesSlice2StrSlice(src [][]byte) []string {
	if src == nil {
		return nil
	}
	ret := make([]string, 0, len(src))
	for _, item := range src {
		ret = append(ret, string(item))
	}
	return ret
}

// UnsignedMap returns a map: column index -> unsigned.
// Note that only numeric columns will be returned.
// nil is returned if not available or no numeric columns at all.
func (e *TableMapEvent) UnsignedMap() map[int]bool {
	if len(e.SignednessBitmap) == 0 {
		return nil
	}
	p := 0
	ret := make(map[int]bool)
	for i := 0; i < int(e.ColumnCount); i++ {
		if !e.IsNumericColumn(i) {
			continue
		}
		ret[i] = e.SignednessBitmap[p/8]&(1<<uint(7-p%8)) != 0
		p++
	}
	return ret
}

// CollationMap returns a map: column index -> collation id.
// Note that only character columns will be returned.
// nil is returned if not available or no character columns at all.
func (e *TableMapEvent) CollationMap() map[int]uint64 {
	return e.collationMap(e.IsCharacterColumn, e.DefaultCharset, e.ColumnCharset)
}

// EnumSetCollationMap returns a map: column index -> collation id.
// Note that only enum or set columns will be returned.
// nil is returned if not available or no enum/set columns at all.
func (e *TableMapEvent) EnumSetCollationMap() map[int]uint64 {
	return e.collationMap(e.IsEnumOrSetColumn, e.EnumSetDefaultCharset, e.EnumSetColumnCharset)
}

func (e *TableMapEvent) collationMap(includeType func(int) bool, defaultCharset, columnCharset []uint64) map[int]uint64 {
	if len(defaultCharset) != 0 {
		defaultCollation := defaultCharset[0]

		// character column index -> collation
		collations := make(map[int]uint64)
		for i := 1; i < len(defaultCharset); i += 2 {
			collations[int(defaultCharset[i])] = defaultCharset[i+1]
		}

		p := 0
		ret := make(map[int]uint64)
		for i := 0; i < int(e.ColumnCount); i++ {
			if !includeType(i) {
				continue
			}

			if collation, ok := collations[p]; ok {
				ret[i] = collation
			} else {
				ret[i] = defaultCollation
			}
			p++
		}

		return ret
	}

	if len(columnCharset) != 0 {
		p := 0
		ret := make(map[int]uint64)
		for i := 0; i < int(e.ColumnCount); i++ {
			if !includeType(i) {
				continue
			}

			ret[i] = columnCharset[p]
			p++
		}

		return ret
	}

	return nil
}

// EnumStrValueMap returns a map: column index -> enum string value.
// Note that only enum columns will be returned.
// nil is returned if not available or no enum columns at all.
func (e *TableMapEvent) EnumStrValueMap() map[int][]string {
	return e.strValueMap(e.IsEnumColumn, e.EnumStrValueString())
}

// SetStrValueMap returns a map: column index -> set string value.
// Note that only set columns will be returned.
// nil is returned if not available or no set columns at all.
func (e *TableMapEvent) SetStrValueMap() map[int][]string {
	return e.strValueMap(e.IsSetColumn, e.SetStrValueString())
}

func (e *TableMapEvent) strValueMap(includeType func(int) bool, strValue [][]string) map[int][]string {
	if len(strValue) == 0 {
		return nil
	}
	p := 0
	ret := make(map[int][]string)
	for i := 0; i < int(e.ColumnCount); i++ {
		if !includeType(i) {
			continue
		}
		ret[i] = strValue[p]
		p++
	}
	return ret
}

// GeometryTypeMap returns a map: column index -> geometry type.
// Note that only geometry columns will be returned.
// nil is returned if not available or no geometry columns at all.
func (e *TableMapEvent) GeometryTypeMap() map[int]uint64 {
	if len(e.GeometryType) == 0 {
		return nil
	}
	p := 0
	ret := make(map[int]uint64)
	for i := 0; i < int(e.ColumnCount); i++ {
		if !e.IsGeometryColumn(i) {
			continue
		}

		ret[i] = e.GeometryType[p]
		p++
	}
	return ret
}

// Below realType and IsXXXColumn are base from:
//   table_def::type in sql/rpl_utility.h
//   Table_map_log_event::print_columns in mysql-8.0/sql/log_event.cc and mariadb-10.5/sql/log_event_client.cc

func (e *TableMapEvent) realType(i int) byte {
	typ := e.ColumnType[i]

	switch typ {
	case MYSQL_TYPE_STRING:
		rtyp := byte(e.ColumnMeta[i] >> 8)
		if rtyp == MYSQL_TYPE_ENUM || rtyp == MYSQL_TYPE_SET {
			return rtyp
		}

	case MYSQL_TYPE_DATE:
		return MYSQL_TYPE_NEWDATE
	}

	return typ
}

func (e *TableMapEvent) IsNumericColumn(i int) bool {
	switch e.realType(i) {
	case MYSQL_TYPE_TINY,
		MYSQL_TYPE_SHORT,
		MYSQL_TYPE_INT24,
		MYSQL_TYPE_LONG,
		MYSQL_TYPE_LONGLONG,
		MYSQL_TYPE_NEWDECIMAL,
		MYSQL_TYPE_FLOAT,
		MYSQL_TYPE_DOUBLE:
		return true

	default:
		return false
	}
}

// IsCharacterColumn returns true if the column type is considered as character type.
// Note that JSON/GEOMETRY types are treated as character type in mariadb.
// (JSON is an alias for LONGTEXT in mariadb: https://mariadb.com/kb/en/json-data-type/)
func (e *TableMapEvent) IsCharacterColumn(i int) bool {
	switch e.realType(i) {
	case MYSQL_TYPE_STRING,
		MYSQL_TYPE_VAR_STRING,
		MYSQL_TYPE_VARCHAR,
		MYSQL_TYPE_BLOB:
		return true

	case MYSQL_TYPE_GEOMETRY:
		if e.flavor == "mariadb" {
			return true
		}
		return false

	default:
		return false
	}
}

func (e *TableMapEvent) IsEnumColumn(i int) bool {
	return e.realType(i) == MYSQL_TYPE_ENUM
}

func (e *TableMapEvent) IsSetColumn(i int) bool {
	return e.realType(i) == MYSQL_TYPE_SET
}

func (e *TableMapEvent) IsGeometryColumn(i int) bool {
	return e.realType(i) == MYSQL_TYPE_GEOMETRY
}

func (e *TableMapEvent) IsEnumOrSetColumn(i int) bool {
	rtyp := e.realType(i)
	return rtyp == MYSQL_TYPE_ENUM || rtyp == MYSQL_TYPE_SET
}
