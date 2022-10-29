package router

import (
	"ult/pkg/http"
)

func (r *Router) heartbeat(group http.RouterGroup) {
	group.GET("", r.handle.Heartbeat.PONE())
}
