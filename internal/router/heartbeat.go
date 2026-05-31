package router

import (
	"ult/pkg/http"
)

func (r *httpRouter) heartbeat(group http.RouterGroup) {
	group.GET("/state", r.handle.Heartbeat.State())
}
