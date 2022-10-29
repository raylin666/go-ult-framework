package router

import (
	"ult/internal/api"
	"ult/internal/service"
	"ult/pkg/http"
)

type Router struct {
	g      http.RouterGroup
	handle struct {
		Heartbeat api.HeartbeatInterface
	}
}

func New(hs *http.HTTPServer) *Router {
	var r = &Router{
		// 创建路由组
		g: hs.CreateRouterGroup(),
		// 注册处理器
		handle: struct {
			Heartbeat api.HeartbeatInterface
		}{
			Heartbeat: api.NewHeartbeatHandler(hs.Logger(), service.NewHeartbeatService(hs.Logger(), hs.DataRepo())),
		},
	}

	// 心跳检测
	r.heartbeat(r.g.Group("/heartbeat"))
	return r
}
