package app

import (
	"github.com/marcuswu/msgpack/app/firebase"
	"github.com/marcuswu/msgpack/mobile/router"
)

type application struct {
	router router.Router
	config firebase.RemoteConfig
}

func SetRouter(router router.Router) {
	app.router = router
}

func Router() router.Router {
	return app.router
}

func SetConfig(config firebase.RemoteConfig) {
	app.config = config
}

func Config() firebase.RemoteConfig {
	return app.config
}

var app = application{}
