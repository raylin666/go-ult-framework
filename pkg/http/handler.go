// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

import (
	goerror "errors"
	"fmt"
	nethttp "net/http"
	"runtime/debug"
	"time"
	"ult/errcode"
	"ult/pkg/logger"
	"ult/pkg/proposal"

	"github.com/gin-gonic/gin"
	"github.com/raylin666/go-utils/v2/validator"
	cors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/zap"
)

// CoreContextNameKey 用于在 Gin.Context 中存储自定义 Context 的键名。
const (
	CoreContextNameKey = "_core_context_"
)

// handlerMiddlewares 注册 HTTP 服务器的核心中间件链。
// 设置 CORS 处理和请求处理中间件。
func (srv *HTTPServer) handlerMiddlewares() {
	// 跨域处理
	srv.engine.Use(srv.handlerCORS())

	// 注册数据验证器
	var validatorHandler = validator.New(
		validator.WithLocale(srv.config.Validator.Locale),
		validator.WithTagname(srv.config.Validator.Tagname))

	// 请求处理 -> 接口异常及响应将在请求处理的 defer 函数内处理
	srv.engine.Use(srv.handlerRequest(validatorHandler))
}

// handlerCORS 返回用于 CORS（跨域资源共享）处理的 Gin 中间件。
// 允许配置的域名从浏览器访问 API。
func (srv *HTTPServer) handlerCORS() gin.HandlerFunc {
	if len(srv.option.cors.domains) > 0 {
		return cors.New(cors.Options{
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
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

// handlerRequest 请求处理中间件。
// 初始化核心上下文 Context，处理请求验证和响应处理。
func (srv *HTTPServer) handlerRequest(validator validator.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 拦截 404 请求路由
		if ctx.Writer.Status() == nethttp.StatusNotFound {
			return
		}

		// 请求时间
		ts := time.Now()

		// 初始化核心上下文 Context
		var appCtx = newContext(ctx)
		defer recoveryContext(appCtx)
		appCtx.init()
		appCtx.WithValidator(validator)
		ctx.Set(CoreContextNameKey, appCtx)

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

// handlerRecovery 异常/错误恢复处理。
// 由 handlerRequest 内 defer 函数处理，请勿单独调用。
// 记录 panic 信息和堆栈，设置告警通知。
func (srv *HTTPServer) handlerRecovery(ctx *gin.Context, err interface{}) {
	// 获取堆栈信息
	var stack = string(debug.Stack())
	srv.logger.UseApp(ctx).Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", stack))
	// 获取核心上下文 Context
	appCtx, ok := ctx.Value(CoreContextNameKey).(Context)
	if !ok || appCtx == nil {
		return
	}

	// 设置错误
	appCtx.WithAbortError(errcode.New(errcode.ServerError).WithStackError(goerror.New("got panic")))

	// 设置告警提醒 (如发邮件通知、如钉钉告警)
	if notifyHandler := srv.option.alertNotify; notifyHandler != nil {
		notifyHandler(&proposal.AlertMessage{
			ProjectName:  srv.config.App.Name,
			Environment:  srv.config.Environment,
			TraceID:      appCtx.TraceID(),
			HOST:         appCtx.Host(),
			URI:          appCtx.URI(),
			Method:       appCtx.Method(),
			ErrorMessage: err,
			ErrorStack:   stack,
			Timestamp:    time.Now(),
		})
	}
}

// handlerResponse 响应处理。
// 由 handlerRequest 内 defer 函数处理，请勿单独调用。
// 根据请求状态生成成功或错误响应，记录请求日志。
func (srv *HTTPServer) handlerResponse(reqTime time.Time, ctx *gin.Context) {

	var (
		resp            interface{}
		httpCode        int
		businessCode    int
		businessMessage string
		stackErr        error
		traceId         string
	)

	// 获取核心上下文 Context
	appCtx, ok := ctx.Value(CoreContextNameKey).(Context)
	if !ok || appCtx == nil {
		return
	}

	// 获取链路追踪 TraceId
	traceId = appCtx.TraceID()

	// 发生错误, 进行返回
	if ctx.IsAborted() {
		if err := appCtx.getAbortError(); err != nil {
			httpCode = err.HTTPCode()
			businessCode = err.BusinessCode()
			businessMessage = err.Message()
			stackErr = err.StackError()
			// 设置告警提醒 (如发邮件通知、如钉钉告警)
			if err.IsAlert() {
				if notifyHandler := srv.option.alertNotify; notifyHandler != nil {
					notifyHandler(&proposal.AlertMessage{
						ProjectName:  srv.config.App.Name,
						Environment:  srv.config.Environment,
						TraceID:      traceId,
						HOST:         appCtx.Host(),
						URI:          appCtx.URI(),
						Method:       appCtx.Method(),
						ErrorMessage: err,
						ErrorStack:   fmt.Sprintf("%+v", stackErr),
						Timestamp:    time.Now(),
					})
				}
			}

			resp = NewErrorResponse(traceId, businessCode, businessMessage, err.Desc())
		} else {
			err = errcode.ErrUnknownError
			httpCode = err.HTTPCode()
			businessCode = err.BusinessCode()
			businessMessage = err.Message()
			stackErr = ctx.Err()
			resp = NewErrorResponse(traceId, businessCode, businessMessage, err.Desc())
		}

		ctx.JSON(httpCode, resp)
	} else {
		// 响应正确返回
		httpCode = nethttp.StatusOK
		businessMessage = "OK"
		resp = NewSuccessResponse(traceId, appCtx.getPayload())
		ctx.JSON(httpCode, resp)
	}

	costSeconds := time.Since(reqTime).Seconds()

	// 请求日志打印
	srv.logger.RequestLog(ctx, &logger.RequestLogFormat{
		ClientIp:          ctx.ClientIP(),
		Method:            appCtx.Method(),
		Path:              appCtx.URI(),
		RequestProto:      ctx.Request.Proto,
		RequestReferer:    ctx.Request.Referer(),
		RequestUa:         ctx.Request.UserAgent(),
		RequestPostData:   ctx.Request.PostForm.Encode(),
		RequestBodyData:   string(appCtx.RawData()),
		RequestHeaderData: appCtx.Header(),
		HttpCode:          ctx.Writer.Status(),
		BusinessCode:      businessCode,
		BusinessMessage:   businessMessage,
		RequestTime:       reqTime,
		ResponseTime:      time.Now(),
		CostSeconds:       costSeconds,
	}, stackErr)
}
