package viewmodels

import (
	"bytes"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/marcuswu/msgpack/app/logic"
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

func (s *MsgPackViewerState) Clone() *MsgPackViewerState {
	data := s.Data
	if data != nil {
		log.Println("Cloning non-nil map data")
		data = s.Data.Clone()
	}
	return &MsgPackViewerState{Data: data, Error: s.Error}
}

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

func NewViewerViewModel(fileData []byte) *ViewerViewModel {
	vm := &ViewerViewModel{observers: make(map[string]MsgPackStateObserver)}
	state := &MsgPackViewerState{Data: nil, Error: nil}

	log.Println("Creating ViewerViewModel")
	data := make(map[string]interface{})
	err := msgpack.Unmarshal(fileData, &data)
	if err != nil {
		log.Printf("Failed unpack file: %s\n", err.Error())
		state.Error = err
		vm.UpdateState(state)
		return vm
	}
	state.Data = logic.NewMap(data)
	log.Printf("Unpacked and set state data with %d keys", len(data))
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
	state := b.state.Load().(*MsgPackViewerState)
	if state == nil {
		log.Println("Cloning nil state")
		return state
	}
	log.Println("Cloning non-nil state")
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
		stateFunc: func(state *MsgPackViewerState) {
			var buf bytes.Buffer
			enc := msgpack.NewEncoder(&buf)
			enc.SetSortMapKeys(true)
			err := enc.Encode(state.Data.Items())

			if err != nil {
				state.Error = fmt.Errorf("Could not convert data: %v", err)
			}

			byteData = buf.Bytes()
		},
	})
	return byteData
}
