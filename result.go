// Package mysql
// result protocol https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
package mysql

import "encoding/binary"

type Result struct {
	Status       uint16
	InsertId     uint64
	AffectedRows uint64
}

func (r *Result) GetStatus() uint16 {
	return r.Status
}

func (r *Result) LastInsertId() (int64, error) {
	return int64(r.InsertId), nil
}

func (r *Result) RowsAffected() (int64, error) {
	return int64(r.AffectedRows), nil
}

type Field struct {
	Schema       []byte
	Table        []byte
	OrgTable     []byte
	Name         []byte
	OrgName      []byte
	Charset      uint16
	ColumnLength uint32
	Type         uint8
	Flag         uint16
	Decimal      uint8

	DefaultValueLength uint64
	DefaultValue       []byte

	isFieldList bool
}

func (f *Field) Dump() []byte {
	l := len(f.Schema) + len(f.Table) + len(f.OrgTable) + len(f.Name) + len(f.OrgName) + len(f.DefaultValue) + 48

	data := make([]byte, 0, l)

	data = append(data, PutLengthEncodedString([]byte("def"))...)
	data = append(data, PutLengthEncodedString(f.Schema)...)
	data = append(data, PutLengthEncodedString(f.Table)...)
	data = append(data, PutLengthEncodedString(f.OrgTable)...)
	data = append(data, PutLengthEncodedString(f.Name)...)
	data = append(data, PutLengthEncodedString(f.OrgName)...)

	data = append(data, 0x0c)

	data = append(data, Uint16ToBytes(f.Charset)...)
	data = append(data, Uint32ToBytes(f.ColumnLength)...)
	data = append(data, f.Type)
	data = append(data, Uint16ToBytes(f.Flag)...)
	data = append(data, f.Decimal)
	data = append(data, 0, 0)

	if f.isFieldList {
		data = append(data, Uint64ToBytes(f.DefaultValueLength)...)
		data = append(data, f.DefaultValue...)
	}

	return data
}

func parseField(p []byte) (f Field, err error) {
	var n int
	pos := 0
	// skip catelog, always def
	n, err = SkipLengthEnodedString(p)
	if err != nil {
		return
	}
	pos += n

	// schema
	f.Schema, _, n, err = LengthEnodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// table
	f.Table, _, n, err = LengthEnodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// org_table
	f.OrgTable, _, n, err = LengthEnodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// name
	f.Name, _, n, err = LengthEnodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// org_name
	f.OrgName, _, n, err = LengthEnodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	// skip oc
	pos++

	// charset
	f.Charset = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	// column length
	f.ColumnLength = binary.LittleEndian.Uint32(p[pos:])
	pos += 4

	// type
	f.Type = p[pos]
	pos++

	//flag
	f.Flag = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	// decimals 1
	f.Decimal = p[pos]
	pos++

	// filter [0x00][0x00]
	pos += 2

	f.isFieldList = false
	// if more data, command was field list
	if len(p) > pos {
		f.isFieldList = true
		// length of default value lenenc-int
		f.DefaultValueLength, _, n = LengthEncodedInt(p[pos:])
		pos += n

		if pos+int(f.DefaultValueLength) > len(p) {
			err = ErrMalformPacket
			return
		}

		// default value string[$len
		f.DefaultValue = p[pos:(pos + int(f.DefaultValueLength))]
	}

	return
}
