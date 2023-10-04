package logic

import (
	"errors"
	"fmt"
	"reflect"
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
	case nil:
		return NilType
	default:
		return UnknownType
	}
}

func TypeString(t FieldType) string {
	switch t {
	case MapType:
		return "MapType"
	case ArrayType:
		return "ArrayType"
	case IntType:
		return "IntType"
	case Int8Type:
		return "Int8Type"
	case Int16Type:
		return "Int16Type"
	case Int32Type:
		return "Int32Type"
	case Int64Type:
		return "Int64Type"
	case BoolType:
		return "BoolType"
	case StringType:
		return "StringType"
	case Float32Type:
		return "Float32Type"
	case Float64Type:
		return "Float64Type"
	case TimeType:
		return "TimeType"
	case NilType:
		return "NilType"
	case UnknownType:
		return "UnknownType"
	}
	return "WTF"
}

type Container interface {
	Set(*Field) error
}

type Field struct {
	MapParent bool
	Index     int
	Key       string
	value     interface{}
}

// Only used w/in Go -- Ok to be skipped by gomobile
func NewArrayField(index int, v interface{}) *Field {
	return &Field{MapParent: false, Index: index, value: v}
}

// Only used w/in Go -- Ok to be skipped by gomobile
func NewMapField(key string, v interface{}) *Field {
	if TypeOf(v) == UnknownType {
		fmt.Printf("Go Found unknown field with key %s\n", key)
	}
	return &Field{MapParent: true, Key: key, value: v}
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
	uval, uok := f.value.(uint)
	var err error = nil
	if !ok && !uok {
		err = errors.New("GetInt() called on a non-int value")
	}
	if uok {
		return int(uval), err
	}
	return val, err
}

func (f *Field) SetInt(v int) {
	if _, ok := f.value.(uint); ok {
		f.value = uint(v)
	}
	f.value = v
}

func (f *Field) GetInt8() (int8, error) {
	val, ok := f.value.(int8)
	uval, uok := f.value.(uint8)
	var err error = nil
	if !ok && !uok {
		err = errors.New("GetInt8() called on a non-int8 value")
	}
	if uok {
		return int8(uval), err
	}
	return val, err
}

func (f *Field) SetInt8(v int8) {
	if _, ok := f.value.(uint8); ok {
		f.value = uint8(v)
	}
	f.value = v
}

func (f *Field) GetInt16() (int16, error) {
	val, ok := f.value.(int16)
	uval, uok := f.value.(uint16)
	var err error = nil
	if !ok && !uok {
		err = fmt.Errorf("GetInt16() called on a non-int16 (%s) value", reflect.TypeOf(f.value))
	}
	if uok {
		return int16(uval), err
	}
	return val, err
}

func (f *Field) SetInt16(v int16) {
	if _, ok := f.value.(uint16); ok {
		f.value = uint16(v)
	}
	f.value = v
}

func (f *Field) GetInt32() (int32, error) {
	val, ok := f.value.(int32)
	uval, uok := f.value.(uint32)
	var err error = nil
	if !ok && !uok {
		err = errors.New("GetInt32() called on a non-int32 value")
	}
	if uok {
		return int32(uval), err
	}
	return val, err
}

func (f *Field) SetInt32(v int32) {
	if _, ok := f.value.(uint32); ok {
		f.value = uint32(v)
	}
	f.value = v
}

func (f *Field) GetInt64() (int64, error) {
	val, ok := f.value.(int64)
	uval, uok := f.value.(uint64)
	var err error = nil
	if !ok && !uok {
		err = errors.New("GetInt64() called on a non-int64 value")
	}
	if uok {
		return int64(uval), err
	}
	return val, err
}

func (f *Field) SetInt64(v int64) {
	if _, ok := f.value.(uint64); ok {
		f.value = uint64(v)
	}
	f.value = v
}

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

func (f *Field) DebugString(path string) string {
	t := TypeOf(f.value)
	switch t {
	case MapType:
		m, _ := f.GetMap()
		return fmt.Sprintf("%s: Map\n%s", path, m.DebugString(path))
	case ArrayType:
		a, _ := f.GetArray()
		return fmt.Sprintf("%s: Array\n%s", path, a.DebugString(path))
	default:
		return fmt.Sprintf("%s: %s = %v\n", path, TypeString(t), f.value)
	}
}
