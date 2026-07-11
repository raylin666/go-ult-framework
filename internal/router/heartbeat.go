// Package router 提供路由层实现。
package router

import (
	"ult/pkg/http"
)

// heartbeat 注册健康检查路由。
//
// 参数:
//   - group: 路由组
func (r *httpRouter) heartbeat(group http.RouterGroup) {
	group.GET("/state", r.handle.Heartbeat.State)
}
