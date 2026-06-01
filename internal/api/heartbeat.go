// Package api 提供 API 处理层实现。
// API 层负责处理 HTTP 请求，进行数据校验和响应，调用服务层处理业务逻辑。
package api

import (
	"ult/internal/app"
	"ult/internal/service"
	"ult/pkg/http"
)

// HeartbeatInterface 接口验证。
var _ HeartbeatInterface = (*HeartbeatHandler)(nil)

// HeartbeatInterface 健康检查 API 接口。
// 定义健康检查相关的 HTTP 处理方法。
type HeartbeatInterface interface {
	State() http.HandlerFunc // 获取系统健康状态
}

// HeartbeatHandler 健康检查 API 处理器。
type HeartbeatHandler struct {
	service *service.HeartbeatService // 健康检查服务
	tools   *app.Tools                // 应用工具包
}

// NewHeartbeatHandler 创建新的健康检查 API 处理器实例。
//
// 参数:
//   - service: 健康检查服务
//   - tools: 应用工具包
//
// 返回:
//   - HeartbeatInterface: 健康检查 API 接口
func NewHeartbeatHandler(service *service.HeartbeatService, tools *app.Tools) HeartbeatInterface {
	return &HeartbeatHandler{
		service: service,
		tools:   tools,
	}
}

// State 获取系统健康状态的 HTTP 处理函数。
// 调用服务层获取数据库和 Redis 连接状态。
//
// 返回:
//   - http.HandlerFunc: HTTP 处理函数
func (h *HeartbeatHandler) State() http.HandlerFunc {
	return func(ctx http.Context) {
		var resp = h.service.State(ctx.RequestContext())
		ctx.WithPayload(resp)
	}
}
