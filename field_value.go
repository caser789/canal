package main

type FieldValueType uint8

type FieldValue struct {
	Type  FieldValueType
	value uint64 // Also for int64 and float64
	str   []byte
}

const (
	FieldValueTypeNull = iota
	FieldValueTypeUnsigned
	FieldValueTypeSigned
	FieldValueTypeFloat
	FieldValueTypeString
)

func (fv *FieldValue) AsUint64() uint64 {
	return fv.value
}

func (fv *FieldValue) AsInt64() int64 {
	return Uint64ToInt64(fv.value)
}

func (fv *FieldValue) AsFloat64() float64 {
	return Uint64ToFloat64(fv.value)
}

func (fv *FieldValue) AsString() []byte {
	return fv.str
}

func (fv *FieldValue) Value() interface{} {
	switch fv.Type {
	case FieldValueTypeUnsigned:
		return fv.AsUint64()
	case FieldValueTypeSigned:
		return fv.AsInt64()
	case FieldValueTypeFloat:
		return fv.AsFloat64()
	case FieldValueTypeString:
		return fv.AsString()
	default: // FieldValueTypeNull
		return nil
	}
}
