package logic

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func pathSlice(path string) []string {
	if len(path) == 0 {
		return []string{}
	}
	return strings.Split(path, "/")
}

func getPathInterface(root interface{}, path []string) (interface{}, error) {
	var current interface{} = root
	for idx, k := range path {
		var value interface{}
		if a, ok := current.([]interface{}); ok {
			index, err := strconv.Atoi(k)
			if err != nil {
				return nil, err
			}
			if index >= len(a) {
				return nil, fmt.Errorf("index out of bounds %d; len: %d", index, len(a))
			}
			value = a[index]
		}
		if m, ok := current.(map[string]interface{}); ok {
			value, ok = m[k]
			if !ok {
				return nil, fmt.Errorf("unknown field %s", k)
			}
		}
		// If we have more path to process, our current value should be an array or map
		if idx < len(path)-1 {
			_, aok := value.([]interface{})
			_, mok := value.(map[string]interface{})
			if !aok && !mok {
				return nil, fmt.Errorf("%s is not an array or dictionary", k)
			}
		}
		current = value
	}
	return current, nil
}

func getPath(root interface{}, path []string) (*Field, error) {
	if len(path) == 0 {
		return NewFieldWithValue("", root), nil
	}
	lastKey := path[len(path)-1]
	value, err := getPathInterface(root, path)
	if err != nil {
		return nil, err
	}
	return NewFieldWithValue(lastKey, value), nil
}

func setPath(root interface{}, path []string, value interface{}) (interface{}, error) {
	if len(path) == 0 {
		return root, nil
	}

	if len(path) == 1 {
		if m, ok := root.(map[string]interface{}); ok {
			m[path[0]] = value
			return m, nil
		}
		if a, ok := root.([]interface{}); ok {
			index, err := strconv.Atoi(path[0])
			if err != nil {
				return nil, err
			}
			if index < len(a) {
				a[index] = value
				return a, nil
			}
			a = append(a, value)
			return a, nil
		}
		return nil, fmt.Errorf("parent field to %s is not a map or an array", path[0])
	}

	key, rest := path[0], path[1:]
	if m, ok := root.(map[string]interface{}); ok {
		result, err := setPath(m[key], rest, value)
		if err != nil {
			return nil, err
		}
		m[key] = result
		return m, nil
	}
	if a, ok := root.([]interface{}); ok {
		index, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}
		if index >= len(a) {
			return nil, fmt.Errorf("index out of bounds %d", index)
		}
		result, err := setPath(a[index], rest, value)
		if err != nil {
			return nil, err
		}
		a[index] = result
		return a, nil
	}
	return nil, fmt.Errorf("parent field to %s is not a map or an array", key)

}

func keySizeAt(root interface{}, path []string) (int, error) {
	value, err := getPathInterface(root, path)
	if err != nil {
		fmt.Printf("Error getting path value %v\n", err)
		return 0, err
	}
	pa, aok := value.([]interface{})
	pm, mok := value.(map[string]interface{})
	if aok {
		return len(pa), nil
	}
	if mok {
		return len(pm), nil
	}
	fmt.Println("Unknown type in keySizeAt")
	return 0, fmt.Errorf("field %s is not a map or an array", path[len(path)-1])
}

func getKeyAt(root interface{}, path []string, index int) (string, error) {
	parent, err := getPathInterface(root, path)
	if err != nil {
		return "", err
	}
	pa, aok := parent.([]interface{})
	pm, mok := parent.(map[string]interface{})
	if aok {
		if index >= len(pa) {
			return "", fmt.Errorf("index %d out of bounds %d", index, len(pa))
		}
		return fmt.Sprintf("%d", index), nil
	}
	if mok {
		if index > len(pm) {
			return "", fmt.Errorf("index %d out of bounds %d", index, len(pm))
		}
		keys := make(sort.StringSlice, 0, len(pm))
		for key := range pm {
			keys = append(keys, key)
		}
		sort.Sort(keys)
		return keys[index], nil
	}
	return "", fmt.Errorf("field %s is not a map or an array", path[len(path)-1])
}

func debugString(root interface{}) string {
	str, _ := json.MarshalIndent(root, "", "\t")
	return string(str)
}
