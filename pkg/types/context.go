// Package types 提供通用类型定义。
// 定义请求上下文、链路追踪等核心类型。
package types

import stdCtx "context"

// 内部上下文键，用于存储请求数据。
const (
	CoreContextNameKey       = "core_context"        // 存储核心 Context 的键名
	ContextBodyNameKey       = "context_body"        // 存储原始请求体的键
	ContextPayloadNameKey    = "context_payload"     // 存储响应数据的键
	ContextAbortErrorNameKey = "context_abort_error" // 存储中止错误的键
	ContextValidatorNameKey  = "context_validator"   // 存储验证器实例的键
)

// requestContextKey 请求上下文键。
type requestContextKey struct{}

var _ RequestContextInterface = (*RequestContext)(nil)

// RequestContextInterface 请求上下文接口，定义获取 TraceID 的方法。
type RequestContextInterface interface {
	TraceID() string // 获取链路追踪 ID
}

// RequestContext 请求上下文，存储请求相关信息。
type RequestContext struct {
	traceId string // 链路追踪 ID
}

// WithTraceID 设置链路追踪 ID。
//
// 参数:
//   - traceId: 链路追踪 ID
func (ctx *RequestContext) WithTraceID(traceId string) {
	ctx.traceId = traceId
}

// TraceID 获取链路追踪 ID。
//
// 返回:
//   - string: 链路追踪 ID
func (ctx *RequestContext) TraceID() string {
	return ctx.traceId
}

// NewRequestContext 创建带有请求上下文的上下文。
//
// 参数:
//   - ctx: 基础上下文
//   - reqCtx: 请求上下文实例
//
// 返回:
//   - stdCtx.Context: 带有请求上下文的上下文
func NewRequestContext(ctx stdCtx.Context, reqCtx RequestContextInterface) stdCtx.Context {
	return stdCtx.WithValue(ctx, requestContextKey{}, reqCtx)
}

// FromRequestContext 从上下文中获取请求上下文。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - reqCtx: 请求上下文实例
//   - ok: 是否存在
func FromRequestContext(ctx stdCtx.Context) (reqCtx RequestContextInterface, ok bool) {
	reqCtx, ok = ctx.Value(requestContextKey{}).(RequestContextInterface)
	return
}
