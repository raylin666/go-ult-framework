package context

import stdCtx "context"

type requestContextKey struct {}

var _ RequestContextInterface = (*RequestContext)(nil)

type RequestContextInterface interface {
	TraceID() string
}

type RequestContext struct {
	traceId string
}

func (ctx *RequestContext) WithTraceID(traceId string) {
	ctx.traceId = traceId
}

func (ctx *RequestContext) TraceID() string {
	return ctx.traceId
}

// NewRequestContext 创建一个新的上下文, 用于业务请求数据传递
func NewRequestContext(ctx stdCtx.Context, reqCtx RequestContextInterface) stdCtx.Context {
	return stdCtx.WithValue(ctx, requestContextKey{}, reqCtx)
}

// FromRequestContext 获取业务请求数据传递上下文
func FromRequestContext(ctx stdCtx.Context) (reqCtx RequestContextInterface, ok bool) {
	reqCtx, ok = ctx.Value(requestContextKey{}).(RequestContextInterface)
	return
}
