package router

import (
	"github.com/google/wire"
	"ult/internal/api"
	"ult/internal/data"
	"ult/internal/service"
	"ult/pkg/http"
)

// ProviderSet is router providers.
var ProviderSet = wire.NewSet(NewHTTPRouter)

type HTTPRouter func(hs *http.HTTPServer)

type httpRouter struct {
	g      http.RouterGroup
	handle struct {
		Heartbeat api.HeartbeatInterface
	}
}

// NewHTTPRouter 创建 HTTP 路由
func NewHTTPRouter() HTTPRouter {
	return func(hs *http.HTTPServer) {
		// 数据仓库
		var _ = data.NewDataRepo(hs.Logger(), hs.DataRepo())
		// HTTP 路由
		var r = &httpRouter{
			// 创建路由组
			g: hs.CreateRouterGroup(),
			// 注册处理器
			handle: struct {
				Heartbeat api.HeartbeatInterface
			}{
				Heartbeat: api.NewHeartbeatHandler(hs.Logger(), service.NewHeartbeatService(hs.Logger())),
			},
		}

		// 心跳检测
		r.heartbeat(r.g.Group("/heartbeat"))
	}
}
