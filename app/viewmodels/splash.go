package viewmodels

import (
	"sync/atomic"

	"github.com/marcuswu/msgpack/app"
	"github.com/marcuswu/msgpack/app/firebase"
	"github.com/marcuswu/msgpack/mobile/state"
)

/*
ViewModel for splash screen
Startup action:
* load remote config
* proceed to home screen
*/
type StartupState struct {
	HaveConfig bool
}

func (s *StartupState) Clone() state.UIState {
	return &StartupState{HaveConfig: s.HaveConfig}
}

// type SplashStateFunc func(*StartupState)
type SplashStateFunc interface {
	WithState(*StartupState)
}
type StartupStateFunc struct {
	stateFunc func(*StartupState)
}

func (sf *StartupStateFunc) WithState(state *StartupState) {
	sf.stateFunc(state)
}

type SplashStateObserver interface {
	Update(*StartupState)
}

type SplashViewModel struct {
	state     atomic.Value
	observers map[string]SplashStateObserver
}

func NewSplashViewModel() *SplashViewModel {
	vm := &SplashViewModel{}
	state := &StartupState{}
	vm.UpdateState(state)
	return vm
}

func (b *SplashViewModel) UpdateState(newState *StartupState) {
	b.state.Store(newState)
	for _, sub := range b.observers {
		sub.Update(b.state.Load().(*StartupState))
	}
}

func (b *SplashViewModel) CloneState() *StartupState {
	return b.state.Load().(*StartupState).Clone().(*StartupState)
}

func (b *SplashViewModel) WithState(stateFunc SplashStateFunc) {
	stateFunc.WithState(b.state.Load().(*StartupState))
}

func (b *SplashViewModel) Observe(id string, callback SplashStateObserver) {
	b.observers[id] = callback
}

func (s *SplashViewModel) LoadRemoteConfig() {
	app.Config().FetchAndActivate(&firebase.ActivateCallback{Callback: func(bool) {
		newState := s.CloneState()
		newState.HaveConfig = true
		s.UpdateState(newState)
		s.CheckNavigate()
	}})
}

func (s *SplashViewModel) CheckNavigate() {
	s.WithState(&StartupStateFunc{
		stateFunc: func(ss *StartupState) {
			if ss.HaveConfig {
				app.Router().Navigate("home")
			}
		},
	})
}
