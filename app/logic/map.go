package logic

import (
	"errors"
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
	newSlice := make([]interface{}, len(a))

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
		items: cloneMap(m.items),
	}
	return newMap
}

func (m *Map) Set(f *Field) error {
	if !f.MapParent {
		return errors.New("Attempt to set array field to map")
	}
	m.items[f.Key()] = f.Value()
	return nil
}

/*func (m *Map) SetField(key string, field *Field) {
	m.items[key] = field.value
}

func (m *Map) SetArray(key string, array *Array) {
	m.items[key] = array.items
}

func (m *Map) SetMap(key string, om *Map) {
	m.items[key] = om.items
}*/

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

/*func (m *Map) GetArray(key string) *Array {
	value, ok := m.items[key]
	if !ok {
		return nil
	}
	arr, ok := value.([]interface{})
	if !ok {
		return nil
	}
	return NewArray(arr)
}

func (m *Map) GetMap(key string) *Map {
	value, ok := m.items[key]
	if !ok {
		return nil
	}
	om, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	return NewMap(om)
}*/

/*func (m *Map) TypeFor(key string) (int, error) {
	value, ok := m.items[key]
	if !ok {
		return int(UnknownType), fmt.Errorf("Key %s does not exist", key)
	}
	return int(TypeOf(value)), nil
}*/

func (m *Map) KeySize() int {
	return len(m.keys)
}

func (m *Map) GetKey(i int) string {
	return m.keys[i]
}
