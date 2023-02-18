package main

import "encoding/binary"

type Field struct {
	Data         FieldData
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
}

type FieldData []byte

func (f *Field) Parse(p FieldData) (err error) {
	f.Data = p

	var n int
	pos := 0
	//skip catelog, always def
	n, err = SkipLengthEncodedString(p)
	if err != nil {
		return
	}
	pos += n

	//schema
	f.Schema, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	//table
	f.Table, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	//org_table
	f.OrgTable, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	//name
	f.Name, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	//org_name
	f.OrgName, _, n, err = LengthEncodedString(p[pos:])
	if err != nil {
		return
	}
	pos += n

	//skip oc
	pos += 1

	//charset
	f.Charset = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	//column length
	f.ColumnLength = binary.LittleEndian.Uint32(p[pos:])
	pos += 4

	//type
	f.Type = p[pos]
	pos++

	//flag
	f.Flag = binary.LittleEndian.Uint16(p[pos:])
	pos += 2

	//decimals 1
	f.Decimal = p[pos]
	pos++

	//filter [0x00][0x00]
	pos += 2

	f.DefaultValue = nil
	//if more data, command was field list
	if len(p) > pos {
		//length of default value lenenc-int
		f.DefaultValueLength, _, n = LengthEncodedInt(p[pos:])
		pos += n

		if pos+int(f.DefaultValueLength) > len(p) {
			err = ErrMalformPacket
			return
		}

		//default value string[$len]
		f.DefaultValue = p[pos:(pos + int(f.DefaultValueLength))]
	}

	return nil
}
