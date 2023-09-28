package viewmodels

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/marcuswu/msgpack/app"
	"github.com/marcuswu/msgpack/app/logic"
	"github.com/marcuswu/msgpack/mobile/state"
	"github.com/vmihailenco/msgpack/v5"
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
	data     *logic.Map
	Error    error
}

func (s *MsgPackViewerState) Clone() state.UIState {
	return &MsgPackViewerState{Filename: s.Filename, data: s.data.Clone(), Error: s.Error}
}

type MsgPackStateFunc func(*MsgPackViewerState)
type MsgPackStateObserver interface {
	Update(*MsgPackViewerState)
}

type ViewerViewModel struct {
	state     atomic.Value
	observers map[string]MsgPackStateObserver
}

func NewViewerViewModel(filename string) *ViewerViewModel {
	vm := &ViewerViewModel{}
	state := vm.CloneState()
	state.Filename = filename
	f, err := os.Open(filename)
	if err != nil {
		state.Error = err
		vm.UpdateState(state)
		return vm
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		state.Error = err
		vm.UpdateState(state)
		return vm
	}
	data := make(map[string]interface{})
	err = msgpack.Unmarshal(bytes, &data)
	if err != nil {
		state.Error = err
		vm.UpdateState(state)
		return vm
	}
	state.data = logic.NewMap(data)
	vm.UpdateState(state)
	return vm
}

func (b *ViewerViewModel) UpdateState(newState *MsgPackViewerState) {
	b.state.Store(newState)
	for _, sub := range b.observers {
		sub.Update(b.state.Load().(*MsgPackViewerState))
	}
}

func (b *ViewerViewModel) CloneState() *MsgPackViewerState {
	return b.state.Load().(*MsgPackViewerState).Clone().(*MsgPackViewerState)
}

func (b *ViewerViewModel) WithState(stateFunc MsgPackStateFunc) {
	stateFunc(b.state.Load().(*MsgPackViewerState))
}

func (b *ViewerViewModel) Observe(id string, callback MsgPackStateObserver) {
	b.observers[id] = callback
}

func (b *ViewerViewModel) ReadState() *MsgPackViewerState {
	return b.state.Load().(*MsgPackViewerState)
}

func (vm *ViewerViewModel) AddValue(destinationPath string, key string, value *logic.Field) {
	state := vm.CloneState()
	state.Error = state.data.Add(destinationPath, key, value)
	vm.UpdateState(state)
}

func (vm *ViewerViewModel) ChangeValue(keyPath string, value *logic.Field) {
	state := vm.CloneState()
	state.Error = state.data.Set(keyPath, value)
	vm.UpdateState(state)
}

func (vm *ViewerViewModel) RemoveValue(keyPath string) {
	state := vm.CloneState()
	state.Error = state.data.Remove(keyPath)
	vm.UpdateState(state)
}

func (vm *ViewerViewModel) GetValue(keyPath string) (*logic.Field, error) {
	return vm.ReadState().data.Get(keyPath)
}

func (vm *ViewerViewModel) KeySize() int {
	return vm.ReadState().data.KeySize()
}

func (vm *ViewerViewModel) GetKey(i int) string {
	return vm.ReadState().data.GetKey(i)
}

func (vm *ViewerViewModel) SaveFile(filename string) {
	vm.WithState(func(state *MsgPackViewerState) {
		bytes, err := msgpack.Marshal(state.data.Items)
		if err != nil {
			state.Error = fmt.Errorf("Could not convert data: %v", err)
		}

		f, err := os.Open(filename)
		if err != nil {
			state.Error = fmt.Errorf("Could not open %s: %v", state.Filename, err)
			vm.UpdateState(state)
			return
		}
		_, err = f.Write(bytes)
		if err != nil {
			state.Error = fmt.Errorf("Could not write file: %v", err)
			vm.UpdateState(state)
			return
		}

		app.Router().Navigate("home")
	})
}

func (vm *ViewerViewModel) Save() {
	vm.WithState(func(state *MsgPackViewerState) {
		vm.SaveFile(state.Filename)
	})
}
