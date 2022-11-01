package http

import (
	goerror "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/raylin666/go-utils/validator"
	cors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/zap"
	nethttp "net/http"
	"runtime/debug"
	"time"
	"ult/internal/constant/errcode"
	"ult/pkg/code"
	"ult/pkg/global"
	"ult/pkg/logger"
)

const (
	_Core_ContextNameKey_ = "_core_context_"
)

// handlerMiddlewares 注册处理中间件
func (srv *HTTPServer) handlerMiddlewares() {
	// 跨域处理
	srv.engine.Use(srv.handlerCORS())

	// 注册数据验证器
	var validator_handler = validator.New(
		validator.WithLocale(srv.config.Validator.Locale),
		validator.WithTagname(srv.config.Validator.Tagname))

	// 请求处理 -> 接口异常及响应将在请求处理的 defer 函数内处理
	srv.engine.Use(srv.handlerRequest(validator_handler))
}

// handlerCORS 跨域处理
func (srv *HTTPServer) handlerCORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(srv.option.cors.domains) > 0 {
			cors.New(cors.Options{
				AllowedOrigins: srv.option.cors.domains,
				AllowedMethods: []string{
					nethttp.MethodHead,
					nethttp.MethodGet,
					nethttp.MethodPost,
					nethttp.MethodPut,
					nethttp.MethodPatch,
					nethttp.MethodDelete,
				},
				AllowedHeaders:     srv.option.cors.domains,
				AllowCredentials:   true,
				OptionsPassthrough: true,
			})
		}

		ctx.Next()
	}
}

// handlerRequest 请求处理
func (srv *HTTPServer) handlerRequest(validator validator.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 拦截 404 请求路由
		if ctx.Writer.Status() == nethttp.StatusNotFound {
			return
		}

		// 请求时间
		ts := time.Now()

		// 初始化核心上下文 Context
		var appctx = newContext(ctx)
		defer recoveryContext(appctx)
		appctx.init()
		appctx.WithValidator(validator)
		ctx.Set(_Core_ContextNameKey_, appctx)

		defer func() {
			// 异常恢复处理
			if err := recover(); err != nil {
				srv.handlerRecovery(ctx, err)
			}

			// 响应处理
			srv.handlerResponse(ts, ctx)
		}()

		ctx.Next()
	}
}

// handlerRecovery 异常/错误处理 (由 handlerRequest 内 defer 函数处理, 请勿单独调用)
func (srv *HTTPServer) handlerRecovery(ctx *gin.Context, err interface{}) {
	// 获取堆栈信息
	var stack = string(debug.Stack())
	srv.logger.UseApp().Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", stack))
	// 获取核心上下文 Context
	appctx := ctx.Value(_Core_ContextNameKey_).(Context)
	if appctx == nil {
		return
	}

	// 设置错误
	appctx.WithAbortError(errcode.NewError(code.ServerError).WithStackError(goerror.New("got panic")))

	// 设置告警提醒 (如发邮件通知、如钉钉告警)

}

// handlerResponse 响应处理 (由 handlerRequest 内 defer 函数处理, 请勿单独调用)
func (srv *HTTPServer) handlerResponse(reqTime time.Time, ctx *gin.Context) {
	var (
		response        interface{}
		httpCode        int
		businessCode    int
		businessMessage string
		stackErr        error
		traceId         string
	)

	// 获取核心上下文 Context
	appctx := ctx.Value(_Core_ContextNameKey_).(Context)
	if appctx == nil {
		return
	}

	// 获取链路追踪 TraceId
	traceId = appctx.TraceID()

	// 发生错误, 进行返回
	if ctx.IsAborted() {
		if err := appctx.getAbortError(); err != nil {
			httpCode = err.HTTPCode()
			businessCode = err.BusinessCode()
			businessMessage = err.Message()
			stackErr = err.StackError()
			// 设置告警提醒 (如发邮件通知、如钉钉告警)
			if err.IsAlert() {
			}

			response = global.ResponseErr{
				TraceId: traceId,
				Code:    businessCode,
				Message: businessMessage,
				Desc:    err.Desc(),
			}
		} else {
			err = errcode.ErrorUnknownError
			httpCode = err.HTTPCode()
			businessCode = err.BusinessCode()
			businessMessage = err.Message()
			stackErr = ctx.Err()
			response = global.ResponseErr{
				TraceId: traceId,
				Code:    businessCode,
				Message: businessMessage,
				Desc:    err.Desc(),
			}
		}

		ctx.JSON(httpCode, response)
	} else {
		// 响应正确返回
		httpCode = nethttp.StatusOK
		businessMessage = "OK"
		response = global.ResponseOK{
			TraceId: traceId,
			Data:    appctx.getPayload(),
		}
		ctx.JSON(httpCode, response)
	}

	costSeconds := time.Since(reqTime).Seconds()

	// 请求日志打印
	srv.logger.RequestLog(&logger.RequestLogFormat{
		TraceId:           traceId,
		ClientIp:          ctx.ClientIP(),
		Method:            appctx.Method(),
		Path:              appctx.URI(),
		RequestProto:      ctx.Request.Proto,
		RequestReferer:    ctx.Request.Referer(),
		RequestUa:         ctx.Request.UserAgent(),
		RequestPostData:   ctx.Request.PostForm.Encode(),
		RequestBodyData:   string(appctx.RawData()),
		RequestHeaderData: appctx.Header(),
		HttpCode:          ctx.Writer.Status(),
		BusinessCode:      businessCode,
		BusinessMessage:   businessMessage,
		RequestTime:       reqTime,
		ResponseTime:      time.Now(),
		CostSeconds:       costSeconds,
	}, stackErr)
}
