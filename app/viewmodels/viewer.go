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
	Data     *logic.Map
	Error    error
}

func (s *MsgPackViewerState) Clone() state.UIState {
	return &MsgPackViewerState{Filename: s.Filename, Data: s.Data.Clone(), Error: s.Error}
}

// type MsgPackStateFunc func(*MsgPackViewerState)
type MsgPackStateFunc interface {
	WithState(*MsgPackViewerState)
}
type ViewerStateFunc struct {
	stateFunc func(*MsgPackViewerState)
}

func (sf *ViewerStateFunc) WithState(state *MsgPackViewerState) {
	sf.stateFunc(state)
}

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
	state.Data = logic.NewMap(data)
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
	stateFunc.WithState(b.state.Load().(*MsgPackViewerState))
}

func (b *ViewerViewModel) Observe(id string, callback MsgPackStateObserver) {
	b.observers[id] = callback
}

// func (b *ViewerViewModel) ReadState() *MsgPackViewerState {
// 	return b.state.Load().(*MsgPackViewerState)
// }

func (vm *ViewerViewModel) SaveFile(filename string) {
	vm.WithState(&ViewerStateFunc{
		stateFunc: func(state *MsgPackViewerState) {
			bytes, err := msgpack.Marshal(state.Data.Items())
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

		},
	})
}

func (vm *ViewerViewModel) Save() {
	vm.WithState(&ViewerStateFunc{
		stateFunc: func(state *MsgPackViewerState) {
			vm.SaveFile(state.Filename)
		},
	})
}
