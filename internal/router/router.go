// Package router 提供路由层实现。
// 路由层负责注册 HTTP 路由，将 API 处理器绑定到对应的路由路径。
package router

import (
	"ult/internal/api"
	"ult/pkg/http"

	"github.com/google/wire"
)

// ProviderSet Wire 依赖注入提供者集合。
var ProviderSet = wire.NewSet(NewHTTPRouter)

// HTTPRouter HTTP 路由注册函数类型。
type HTTPRouter func(hs *http.HTTPServer)

// httpRouter HTTP 路由结构体。
// 包含路由组和 API 处理器映射。
type httpRouter struct {
	g      http.RouterGroup // 路由组
	handle struct {
		Heartbeat api.HeartbeatInterface // 健康检查 API 处理器
	}
}

// NewHTTPRouter 创建 HTTP 路由注册函数。
// 注册所有业务模块的路由。
//
// 参数:
//   - heartbeat: 健康检查 API 处理器
//
// 返回:
//   - HTTPRouter: 路由注册函数
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

		// 注册健康检查路由
		r.heartbeat(r.g.Group("/heartbeat"))
	}
}
