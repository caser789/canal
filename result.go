package mysql

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
	return nil
}
