package logic

import (
	"errors"
	"fmt"
)

func cloneValue(value interface{}) interface{} {
	switch t := value.(type) {
	case map[string]interface{}:
		return cloneMap(t)
	case []interface{}:
		return cloneArray(t)
	default:
		return value
	}
}

func cloneArray(a []interface{}) []interface{} {
	newSlice := make([]interface{}, 0, len(a))

	for _, v := range a {
		newSlice = append(newSlice, cloneValue(v))
	}

	return newSlice
}

func cloneMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for id, v := range m {
		newMap[id] = cloneValue(v)
	}

	return newMap
}

func mapKeys(m map[string]interface{}) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

type Map struct {
	keys  []string
	items map[string]interface{}
}

// Only used w/in Go -- Ok to be skipped by gomobile
func NewMap(m map[string]interface{}) *Map {
	return &Map{keys: mapKeys(m), items: m}
}

// Only used w/in Go for saving the MsgPack file -- Ok to be skipped by gomobile
func (m *Map) Items() map[string]interface{} {
	return m.items
}

func (m *Map) Clone() *Map {
	newMap := &Map{
		keys:  mapKeys(m.items),
		items: cloneMap(m.items),
	}
	return newMap
}

func (m *Map) Set(f *Field) error {
	if !f.MapParent {
		return errors.New("Attempt to set array field to map")
	}
	fmt.Printf("Change %s from %v to %v\n", f.Key, m.items[f.Key], f.Value())
	m.items[f.Key] = f.Value()
	return nil
}

func (m *Map) Remove(key string) {
	delete(m.items, key)
}

func (m *Map) GetField(key string) *Field {
	value, ok := m.items[key]
	if !ok {
		return nil
	}
	return NewMapField(key, value)
}

func (m *Map) KeySize() int {
	return len(m.keys)
}

func (m *Map) GetKey(i int) string {
	return m.keys[i]
}

func (m *Map) DebugString(path string) string {
	out := ""
	for k := range m.items {
		out += m.GetField(k).DebugString(fmt.Sprintf("%s/%s", path, k))
	}
	return out
}

func (m *Map) String() string {
	return m.DebugString("")
}
