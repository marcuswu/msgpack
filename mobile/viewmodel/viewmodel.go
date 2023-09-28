package viewmodel

import (
	"sync/atomic"

	"github.com/marcuswu/msgpack/mobile/state"
)

/*
ViewModel holds (produces?) a UI State
UI renders itself based on UI State
ModelView is responsible for producing State based on business logic
This means business logic may simply accept parameters and return a result

ViewModels may:
* Implement simple business logic directly or
* Rely on a Domain Logic layer for complex logic

Business Logic:
  * Loading data
  * Actions with data
  * Transforming data from backend
  * Produces UI state

Responsibility of ViewModel is applying business logic to application data to yield a screen UI state

Data Layer --[Application Data]--> ViewModel --[UI state]--> UI elements
                                            <----[events]----|
ViewModel holds and exposes state to be consumed by the UI
UI state is application data transformed by the ViewModel
ViewModel handles the user actions and updates the state
Updated state is fed back to the UI for rendering
Repeat the above for any event causing state changes
*/

type ViewModel[S state.UIState] interface {
	UpdateState(S)
}

type StateFunc[S state.UIState] func(S)
type StateObserver[S state.UIState] interface {
	Update(S)
}

type BaseViewModel[S state.UIState] struct {
	state     atomic.Value
	observers map[string]StateObserver[S]
}

func (b *BaseViewModel[S]) UpdateState(newState S) {
	b.state.Store(newState)
	for _, sub := range b.observers {
		sub.Update(b.state.Load().(S))
	}
}

func (b *BaseViewModel[S]) CloneState() S {
	return b.state.Load().(S).Clone().(S)
}

func (b *BaseViewModel[S]) WithState(stateFunc StateFunc[S]) {
	stateFunc(b.state.Load().(S))
}

func (b *BaseViewModel[S]) Observe(id string, callback StateObserver[S]) {
	b.observers[id] = callback
}

func (b *BaseViewModel[S]) ReadState() S {
	return b.state.Load().(S)
}
