package viewmodels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/marcuswu/msgpack/app/logic"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

type structEdFormat int

const (
	msgpackFormat structEdFormat = iota
	yamlFormat
	jsonFormat
	unknownFormat
)

/*
ViewModel for file viewer screen
viewer actions:
* Change Value
* Add Value
* Remove Value
* Save File
*/
type MsgPackViewerState struct {
	Filename string
	Data     *logic.Field
	Error    error
	format   structEdFormat
}

func (s *MsgPackViewerState) Clone() *MsgPackViewerState {
	data := s.Data
	if data != nil {
		data = s.Data.Clone()
	}
	return &MsgPackViewerState{Data: data, Error: s.Error, format: s.format}
}

type MsgPackStateFunc interface {
	WithState(*MsgPackViewerState)
}
type ViewerStateFunc struct {
	StateFunc func(*MsgPackViewerState)
}

func (sf *ViewerStateFunc) WithState(state *MsgPackViewerState) {
	sf.StateFunc(state)
}

type MsgPackStateObserver interface {
	Update(*MsgPackViewerState)
}

type ViewerViewModel struct {
	state     atomic.Value
	observers map[string]MsgPackStateObserver
}

func NewViewerViewModel(fileData []byte) *ViewerViewModel {
	vm := &ViewerViewModel{observers: make(map[string]MsgPackStateObserver)}
	state := &MsgPackViewerState{Data: nil, Error: nil}

	log.Println("Creating ViewerViewModel")
	// Try MsgPack first
	data, err := vm.readMsgPack(fileData)
	state.format = msgpackFormat
	if err != nil {
		data, err = vm.readJson(fileData)
		state.format = jsonFormat
		if err != nil {
			data, err = vm.readYaml(fileData)
			state.format = yamlFormat
		}
	}

	if err != nil {
		log.Printf("Failed unpack file: %s\n", err.Error())
		state.Error = err
		vm.UpdateState(state)
		return vm
	}

	numKeys := 0
	log.Printf("Unpacked data: %v", data.DebugString())
	if m, _ := data.GetMap(); m != nil {
		numKeys, _ = m.KeySizeAt("")
	}
	if a, _ := data.GetArray(); a != nil {
		numKeys, _ = a.KeySizeAt("")
	}
	log.Printf("Detected encoding format %d", state.format)
	log.Printf("Unpacked and set state data with %d keys", numKeys)
	state.Data = data
	vm.UpdateState(state)
	return vm
}

func (b *ViewerViewModel) readYaml(fileData []byte) (*logic.Field, error) {
	data := make(map[string]interface{})
	err := yaml.Unmarshal(fileData, &data)
	if err != nil {
		log.Printf("Failed unpack file: %s\n", err.Error())
		return nil, err
	}
	return logic.NewFieldWithValue("", data), nil
}

func (b *ViewerViewModel) readJson(fileData []byte) (*logic.Field, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal(fileData, &data)
	if err != nil {
		log.Printf("Failed unpack file: %s\n", err.Error())
		return nil, err
	}
	return logic.NewFieldWithValue("", data), nil
}

func (b *ViewerViewModel) readMsgPack(fileData []byte) (*logic.Field, error) {
	data := make(map[string]interface{})
	err := msgpack.Unmarshal(fileData, &data)
	if err != nil {
		log.Printf("Failed unpack file: %s\n", err.Error())
		return nil, err
	}
	return logic.NewFieldWithValue("", data), nil
}

func (b *ViewerViewModel) UpdateState(newState *MsgPackViewerState) {
	oldState := b.state.Load()
	fmt.Printf("Storing new state %p (old state %p)\n", newState, oldState)
	b.state.Store(newState)
	for _, sub := range b.observers {
		sub.Update(b.state.Load().(*MsgPackViewerState))
	}
}

func (b *ViewerViewModel) CloneState() *MsgPackViewerState {
	state := b.state.Load().(*MsgPackViewerState)
	if state == nil {
		return state
	}
	return state.Clone()
}

func (b *ViewerViewModel) WithState(stateFunc MsgPackStateFunc) {
	stateFunc.WithState(b.state.Load().(*MsgPackViewerState))
}

func (b *ViewerViewModel) Observe(id string, callback MsgPackStateObserver) {
	b.observers[id] = callback
}

func (vm *ViewerViewModel) FileData() []byte {
	byteData := []byte{}
	vm.WithState(&ViewerStateFunc{
		StateFunc: func(state *MsgPackViewerState) {
			var buf bytes.Buffer
			var err error = nil

			fmt.Printf("Encoding format %d\n", state.format)
			switch state.format {
			case msgpackFormat:
				fmt.Println("Encoding data to msgpack")
				enc := msgpack.NewEncoder(&buf)
				enc.SetSortMapKeys(true)
				err = enc.Encode(state.Data.Value())
			case yamlFormat:
				fmt.Println("Encoding data to yaml")
				enc := yaml.NewEncoder(&buf)
				err = enc.Encode(state.Data.Value())
			case jsonFormat:
				fmt.Println("Encoding data to json")
				enc := json.NewEncoder(&buf)
				err = enc.Encode(state.Data.Value())
			}

			if err != nil {
				state.Error = fmt.Errorf("could not convert data: %v", err)
			}

			byteData = buf.Bytes()
		},
	})
	return byteData
}

func (vm *ViewerViewModel) GetPath(path string) *logic.Field {
	state := vm.CloneState()
	a, err := state.Data.GetArray()
	if err == nil {
		val, err := a.GetPath(path)
		if err != nil {
			state.Error = err
			vm.UpdateState(state)
			return nil
		}
		return val
	}
	err = nil
	m, err := state.Data.GetMap()
	if err == nil {
		val, err := m.GetPath(path)
		if err != nil {
			state.Error = err
			vm.UpdateState(state)
			return nil
		}
		return val
	}
	state.Error = err
	vm.UpdateState(state)
	return nil
}

func (vm *ViewerViewModel) SetPath(path string, field *logic.Field) {
	state := vm.CloneState()
	a, err := state.Data.GetArray()
	if err == nil {
		err := a.SetPath(path, field)
		if err != nil {
			state.Error = err
			vm.UpdateState(state)
			return
		}
		vm.UpdateState(state)
		return
	}
	m, err := state.Data.GetMap()
	if err == nil {
		err := m.SetPath(path, field)
		if err != nil {
			state.Error = err
			vm.UpdateState(state)
			return
		}
		vm.UpdateState(state)
		return
	}
	state.Error = err
	vm.UpdateState(state)
}

func (vm *ViewerViewModel) GetFormat() int {
	return int(vm.state.Load().(*MsgPackViewerState).format)
}

func (vm *ViewerViewModel) SetFormat(format int) {
	state := vm.CloneState()
	switch format {
	case int(msgpackFormat):
		state.format = msgpackFormat
	case int(yamlFormat):
		state.format = yamlFormat
	case int(jsonFormat):
		state.format = jsonFormat
	}
	vm.UpdateState(state)
}
