package router

import (
	"ult/internal/api"
	"ult/pkg/http"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewHTTPRouter)

type HTTPRouter func(hs *http.HTTPServer)

type httpRouter struct {
	g      http.RouterGroup
	handle struct {
		Heartbeat api.HeartbeatInterface
	}
}

func NewHTTPRouter(heartbeat api.HeartbeatInterface) HTTPRouter {
	return func(hs *http.HTTPServer) {
		var r = &httpRouter{
			g: hs.CreateRouterGroup(),
			handle: struct {
				Heartbeat api.HeartbeatInterface
			}{
				Heartbeat: heartbeat,
			},
		}
		r.heartbeat(r.g.Group("/heartbeat"))
	}
}
