package state

type UIState interface {
	Clone() UIState
}
