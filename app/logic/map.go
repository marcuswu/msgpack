package logic

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

func (m *Map) Remove(key string) {
	delete(m.items, key)
}

func (m *Map) GetPath(path string) (*Field, error) {
	return getPath(m.items, pathSlice(path))
}

func (m *Map) SetPath(path string, field *Field) error {
	newItems, err := setPath(m.items, append(pathSlice(path), field.Key), field.value)
	if err != nil {
		return err
	}
	m.items = newItems.(map[string]interface{})
	return nil
}

func (m *Map) KeySizeAt(path string) (int, error) {
	return keySizeAt(m.items, pathSlice(path))
}

func (m *Map) GetKeyAt(path string, i int) (string, error) {
	return getKeyAt(m.items, pathSlice(path), i)
}

func (m *Map) DebugString() string {
	return debugString(m.items)
}
