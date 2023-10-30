package logic

type Array struct {
	items []interface{}
}

// Only used w/in Go -- Ok to be skipped by gomobile
func NewArray(items []interface{}) *Array {
	return &Array{items: items}
}

func (a *Array) GetPath(path string) (*Field, error) {
	return getPath(a.items, pathSlice(path))
}

func (a *Array) SetPath(path string, value *Field) error {
	newItems, err := setPath(a.items, pathSlice(path), value.value)
	if err != nil {
		return err
	}
	a.items = newItems.([]interface{})
	return nil
}

func (a *Array) KeySizeAt(path string) (int, error) {
	return keySizeAt(a.items, pathSlice(path))
}

func (a *Array) GetKeyAt(path string, i int) (string, error) {
	return getKeyAt(a.items, pathSlice(path), i)
}

func (a *Array) DebugString() string {
	return debugString(a.items)
}
