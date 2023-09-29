package logic

import (
	"errors"
	"time"
)

type FieldType int

const (
	NilType FieldType = iota
	MapType
	ArrayType
	IntType
	Int8Type
	Int16Type
	Int32Type
	Int64Type
	BoolType
	StringType
	Float32Type
	Float64Type
	TimeType
	UnknownType
)

// Only used w/in Go -- Ok to be skipped by gomobile
func TypeOf(value interface{}) FieldType {
	switch value.(type) {
	case map[string]interface{}:
		return MapType
	case []interface{}:
		return ArrayType
	case int, uint:
		return IntType
	case int8, uint8:
		return Int8Type
	case int16, uint16:
		return Int16Type
	case int32, uint32:
		return Int32Type
	case int64, uint64:
		return Int64Type
	case bool:
		return BoolType
	case string:
		return StringType
	case float32:
		return Float32Type
	case float64:
		return Float64Type
	case time.Time:
		return TimeType
	default:
		return UnknownType
	}
}

type Container interface {
	Set(*Field) error
}

type Field struct {
	MapParent bool
	index     int
	key       string
	value     interface{}
}

// Only used w/in Go -- Ok to be skipped by gomobile
func NewArrayField(index int, v interface{}) *Field {
	return &Field{MapParent: false, index: index, value: v}
}

// Only used w/in Go -- Ok to be skipped by gomobile
func NewMapField(key string, v interface{}) *Field {
	return &Field{MapParent: true, key: key, value: v}
}

func (f *Field) Index() int {
	return f.index
}

func (f *Field) Key() string {
	return f.key
}

// For use in Go -- can be skipped by gomobile
func (f *Field) Value() interface{} {
	return f.value
}

func (f *Field) Type() int {
	return int(TypeOf(f.value))
}

func (f *Field) IsNil() bool {
	return f.value == nil
}

func (f *Field) GetInt() (int, error) {
	val, ok := f.value.(int)
	var err error = nil
	if !ok {
		err = errors.New("GetInt() called on a non-int value")
	}
	return val, err
}

func (f *Field) SetInt(v int) {
	f.value = v
}

func (f *Field) GetInt8() (int8, error) {
	val, ok := f.value.(int8)
	var err error = nil
	if !ok {
		err = errors.New("GetInt8() called on a non-int8 value")
	}
	return val, err
}

func (f *Field) SetInt8(v int8) {
	f.value = v
}

func (f *Field) GetInt16() (int16, error) {
	val, ok := f.value.(int16)
	var err error = nil
	if !ok {
		err = errors.New("GetInt16() called on a non-int16 value")
	}
	return val, err
}

func (f *Field) SetInt16(v int16) {
	f.value = v
}

func (f *Field) GetInt32() (int32, error) {
	val, ok := f.value.(int32)
	var err error = nil
	if !ok {
		err = errors.New("GetInt32() called on a non-int32 value")
	}
	return val, err
}

func (f *Field) SetInt32(v int32) {
	f.value = v
}

func (f *Field) GetInt64() (int64, error) {
	val, ok := f.value.(int64)
	var err error = nil
	if !ok {
		err = errors.New("GetInt64() called on a non-int64 value")
	}
	return val, err
}

func (f *Field) SetInt64(v int64) {
	f.value = v
}

/*func (f *Field) GetUint() (uint, error) {
	val, ok := f.value.(uint)
	var err error = nil
	if !ok {
		err = errors.New("GetUint() called on a non-uint value")
	}
	return val, err
}

func (f *Field) SetUInt(v uint) {
	f.value = v
}

func (f *Field) GetUint8() (uint8, error) {
	val, ok := f.value.(uint8)
	var err error = nil
	if !ok {
		err = errors.New("GetUint8() called on a non-uint8 value")
	}
	return val, err
}

func (f *Field) SetUint8(v int8) {
	f.value = v
}

func (f *Field) GetUint16() (uint16, error) {
	val, ok := f.value.(uint16)
	var err error = nil
	if !ok {
		err = errors.New("GetUint16() called on a non-uint16 value")
	}
	return val, err
}

func (f *Field) SetUint16(v uint16) {
	f.value = v
}

func (f *Field) GetUint32() (uint32, error) {
	val, ok := f.value.(uint32)
	var err error = nil
	if !ok {
		err = errors.New("GetUint32() called on a non-uint32 value")
	}
	return val, err
}

func (f *Field) SetUint32(v uint32) {
	f.value = v
}

func (f *Field) GetUint64() (uint64, error) {
	val, ok := f.value.(uint64)
	var err error = nil
	if !ok {
		err = errors.New("GetUint64() called on a non-uint64 value")
	}
	return val, err
}

func (f *Field) SetUint64(v uint64) {
	f.value = v
}*/

func (f *Field) GetBool() (bool, error) {
	val, ok := f.value.(bool)
	var err error = nil
	if !ok {
		err = errors.New("GetBool() called on a non-bool value")
	}
	return val, err
}

func (f *Field) SetBool(v bool) {
	f.value = v
}

func (f *Field) GetString() (string, error) {
	val, ok := f.value.(string)
	var err error = nil
	if !ok {
		err = errors.New("GetString() called on a non-string value")
	}
	return val, err
}

func (f *Field) SetString(v string) {
	f.value = v
}

func (f *Field) GetFloat32() (float32, error) {
	val, ok := f.value.(float32)
	var err error = nil
	if !ok {
		err = errors.New("GetFloat32() called on a non-float32 value")
	}
	return val, err
}

func (f *Field) SetFloat32(v float32) {
	f.value = v
}

func (f *Field) GetFloat64() (float64, error) {
	val, ok := f.value.(float64)
	var err error = nil
	if !ok {
		err = errors.New("GetFloat64() called on a non-float64 value")
	}
	return val, err
}

func (f *Field) SetFloat64(v float64) {
	f.value = v
}

func (f *Field) GetTime() (int64, error) {
	val, ok := f.value.(time.Time)
	if ok {
		return val.UnixMilli(), nil
	}
	return 0, errors.New("GetFloat64() called on a non-float64 value")
}

func (f *Field) SetTime(v int64) {
	f.value = time.UnixMilli(v)
}

func (f *Field) GetMap() (*Map, error) {
	val, ok := f.value.(map[string]interface{})
	if !ok {
		return nil, errors.New("GetMap() called on a non-map value")
	}
	return NewMap(val), nil
}

func (f *Field) SetMap(m *Map) {
	f.value = m.items
}

func (f *Field) GetArray() (*Array, error) {
	val, ok := f.value.([]interface{})
	if !ok {
		return nil, errors.New("GetArray() called on a non-array value")
	}
	return NewArray(val), nil
}

func (f *Field) SetArray(a *Array) {
	f.value = a.items
}
