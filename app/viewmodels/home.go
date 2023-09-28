package viewmodels

import (
	"errors"
	"os"
	"sync/atomic"

	"github.com/marcuswu/msgpack/app"
	"github.com/marcuswu/msgpack/mobile/state"
)

/*
ViewModel for home screen
Home actions:
* Open file selection
* Select file
*/
type HomeState struct {
	File  string
	Error error
}

func (s *HomeState) Clone() state.UIState {
	return &HomeState{File: s.File, Error: s.Error}
}

type HomeStateFunc func(*HomeState)
type HomeStateObserver interface {
	Update(*HomeState)
}

type HomeViewModel struct {
	state     atomic.Value
	observers map[string]HomeStateObserver
}

func (b *HomeViewModel) UpdateState(newState *HomeState) {
	b.state.Store(newState)
	for _, sub := range b.observers {
		sub.Update(b.state.Load().(*HomeState))
	}
}

func (b *HomeViewModel) CloneState() *HomeState {
	return b.state.Load().(*HomeState).Clone().(*HomeState)
}

func (b *HomeViewModel) WithState(stateFunc HomeStateFunc) {
	stateFunc(b.state.Load().(*HomeState))
}

func (b *HomeViewModel) Observe(id string, callback HomeStateObserver) {
	b.observers[id] = callback
}

func (b *HomeViewModel) ReadState() *HomeState {
	return b.state.Load().(*HomeState)
}

func NewHomeViewModel() *HomeViewModel {
	return &HomeViewModel{}
}

func (vm *HomeViewModel) FileSelected(file string) {
	newState := vm.CloneState()
	newState.File = file

	if _, err := os.Stat(newState.File); errors.Is(err, os.ErrNotExist) {
		// Show error
		newState.Error = err
	}
	vm.UpdateState(newState)

	app.Router().Navigate("viewer")
}
