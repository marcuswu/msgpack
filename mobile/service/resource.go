package service

type LoadingState int

const (
	Loading LoadingState = iota
	Success
	Error
)

/*
A Service is just a go routine
A Resource, maintains the state of the go routine (the service call state)
On completion, the Resource sets Error or Result
*/
type Resource[T any] struct {
	State  LoadingState
	Error  error
	Result T
}
