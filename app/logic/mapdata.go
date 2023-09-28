package logic

import (
	"fmt"
	"strconv"
	"strings"
)

// Rewrite this -- expand map to objects which can then collapse back to a map

type Map struct {
	keys  []string
	index map[string]interface{}
	Items map[string]interface{}
}

func (m *Map) buildMapIndex(basePath string, value map[string]interface{}) {
	for k, v := range value {
		key := k
		if len(basePath) > 0 {
			key = fmt.Sprintf("%s/%s", basePath, k)
		}
		_, isMap := v.(map[string]interface{})
		_, isSlice := v.([]interface{})
		if isMap || isSlice {
			m.keys = append(m.keys, key)
			m.index[key] = v
			m.buildIndex(key, v)
		}
	}
}

func (m *Map) buildSliceIndex(basePath string, value []interface{}) {
	for i, v := range value {
		key := strconv.Itoa(i)
		if len(basePath) > 0 {
			key = fmt.Sprintf("%s/%d", basePath, i)
		}
		_, isMap := v.(map[string]interface{})
		_, isSlice := v.([]interface{})
		if isMap || isSlice {
			m.keys = append(m.keys, key)
			m.index[key] = v
			m.buildIndex(key, v)
		}
	}
}

func (m *Map) buildIndex(basePath string, value interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		m.buildMapIndex(basePath, v)
	case []interface{}:
		m.buildSliceIndex(basePath, v)
	default:
		m.keys = append(m.keys, basePath)
	}
}

func (m *Map) rebuildIndex(basePath string, value interface{}) {
	m.keys = make([]string, 0)
	m.index = make(map[string]interface{})
	m.buildIndex("", m.Items)
}

func NewMap(m map[string]interface{}) *Map {
	newMap := &Map{
		keys:  make([]string, 0),
		index: make(map[string]interface{}),
		Items: m,
	}
	newMap.buildIndex("", newMap.Items)
	return newMap
}

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

func (m *Map) Clone() *Map {
	newMap := &Map{
		keys:  make([]string, 0),
		index: make(map[string]interface{}),
		Items: cloneMap(m.Items),
	}
	newMap.buildIndex("", newMap.Items)
	return newMap
}

func (m *Map) Add(destinationPath string, key string, field *Field) error {
	dest, ok := m.index[destinationPath]
	if !ok {
		return fmt.Errorf("Invalid destination %s", destinationPath)
	}

	switch v := dest.(type) {
	case map[string]interface{}:
		v[key] = field.value
	case []interface{}:
		index, err := strconv.Atoi(key)
		if err != nil {
			return fmt.Errorf("Invalid index %s", field)
		}
		if index > len(v) {
			return fmt.Errorf("Invalid index %d", index)
		}
		if index == len(v) {
			v = append(v, field.value)
			return nil
		}

		v = append(v[:index+1], v[index:]...)
		v[index] = field.value
	default:
		return fmt.Errorf("Invalid destination %s", destinationPath)
	}

	return nil
}

func (m *Map) Set(keyPath string, field *Field) error {
	pathParts := strings.Split(keyPath, "/")
	finalKey := pathParts[len(pathParts)-1]
	penultimatePathParts := pathParts[:len(pathParts)-1]
	penultimatePath := strings.Join(penultimatePathParts, "/")
	if len(penultimatePathParts) < 1 {
		m.Items[keyPath] = field.value
		return nil
	}
	dest := m.index[penultimatePath]
	switch v := dest.(type) {
	case map[string]interface{}:
		v[finalKey] = field.value
	case []interface{}:
		index, err := strconv.Atoi(finalKey)
		if err != nil {
			return fmt.Errorf("Invalid index %s", finalKey)
		}
		if index >= len(v) {
			return fmt.Errorf("Invalid index %d", index)
		}
		v[index] = field.value
	default:
		return fmt.Errorf("Invalid path %s", keyPath)
	}

	return nil
}

func (m *Map) Remove(keyPath string) error {
	pathParts := strings.Split(keyPath, "/")
	finalKey := pathParts[len(pathParts)-1]
	penultimatePathParts := pathParts[:len(pathParts)-1]
	penultimatePath := strings.Join(penultimatePathParts, "/")
	if len(penultimatePathParts) < 1 {
		delete(m.Items, keyPath)
		return nil
	}
	dest := m.index[penultimatePath]
	switch v := dest.(type) {
	case map[string]interface{}:
		delete(v, finalKey)
	case []interface{}:
		index, err := strconv.Atoi(finalKey)
		if err != nil {
			return fmt.Errorf("Invalid index %s", finalKey)
		}
		if index >= len(v) {
			return fmt.Errorf("Invalid index %d", index)
		}
		if index < len(v)-1 {
			v = v[:len(v)-1]
			return nil
		}
		v = append(v[:index], v[index+1:]...)
	default:
		return fmt.Errorf("Invalid path %s", keyPath)
	}

	return nil
}

func (m *Map) Get(keyPath string) (*Field, error) {
	pathParts := strings.Split(keyPath, "/")
	finalKey := pathParts[len(pathParts)-1]
	penultimatePathParts := pathParts[:len(pathParts)-1]
	penultimatePath := strings.Join(penultimatePathParts, "/")
	if len(penultimatePathParts) < 1 {
		return &Field{m.Items[keyPath]}, nil
	}
	dest := m.index[penultimatePath]
	switch v := dest.(type) {
	case map[string]interface{}:
		return &Field{v[finalKey]}, nil
	case []interface{}:
		index, err := strconv.Atoi(finalKey)
		if err != nil {
			return nil, fmt.Errorf("Invalid index %s", finalKey)
		}
		if index >= len(v) {
			return nil, fmt.Errorf("Invalid index %d", index)
		}
		return &Field{v[index]}, nil
	default:
		return nil, fmt.Errorf("Invalid path %s", keyPath)
	}
}

func (m *Map) KeySize() int {
	return len(m.keys)
}

func (m *Map) GetKey(i int) string {
	return m.keys[i]
}
