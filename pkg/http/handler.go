// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

// CoreContextNameKey 用于在 Gin.Context 中存储自定义 Context 的键名。
const (
	CoreContextNameKey = "_core_context_"
)

// Note: 中间件逻辑已迁移到 pkg/http/middleware 包中
// handlerMiddlewares、handlerCORS、handlerRecovery、handlerResponse 等方法已移除
// 请使用 middleware 包中的 CORS、Recovery、Request、Response 等中间件
