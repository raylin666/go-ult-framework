package types

import stdCtx "context"

type requestContextKey struct{}

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

func NewRequestContext(ctx stdCtx.Context, reqCtx RequestContextInterface) stdCtx.Context {
	return stdCtx.WithValue(ctx, requestContextKey{}, reqCtx)
}

func FromRequestContext(ctx stdCtx.Context) (reqCtx RequestContextInterface, ok bool) {
	reqCtx, ok = ctx.Value(requestContextKey{}).(RequestContextInterface)
	return
}