package logic

import "fmt"

type Array struct {
	items []interface{}
}

// Only used w/in Go -- Ok to be skipped by gomobile
func NewArray(items []interface{}) *Array {
	return &Array{items: items}
}

func (a Array) Set(field *Field) error {
	if field.MapParent {
		return fmt.Errorf("Attempt to set map field to array type")
	}
	i := field.Index
	if i == len(a.items) {
		a.items = append(a.items, field.Value())
		return nil
	}
	if i >= len(a.items) {
		return fmt.Errorf("Index %d out of bounds (len %d)", i, len(a.items))
	}
	a.items[i] = field.value
	return nil
}

func (a Array) Get(i int) *Field {
	return NewArrayField(i, a.items[i])
}

func (a Array) Size() int {
	return len(a.items)
}

func (a Array) DebugString(path string) string {
	out := ""
	for k, _ := range a.items {
		out += a.Get(k).DebugString(fmt.Sprintf("%s/%d", path, k))
	}
	return out
}
