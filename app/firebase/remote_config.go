package firebase

// type ActivateCallback interface {
// 	OnActivate(bool)
// }

// func NewActivateCallback(cb func(bool)) ActivateCallback {
// 	return ActivateCallback{OnActivate: cb}
// }

// type ActivateCallback func(bool)
type ActivateCallback struct {
	Callback func(bool)
}

func (cb *ActivateCallback) OnActivate(fetched bool) {
	cb.Callback(fetched)
}

type RemoteConfig interface {
	FetchAndActivate(*ActivateCallback)
	GetBool(string) bool
	GetFloat64(string) float64
	GetInt(string) int
	GetStr(string) string
	GetJson(string, string) error
}
