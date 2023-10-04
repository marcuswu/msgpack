package router

type Router interface {
	Navigate(string)
	Back()
}
