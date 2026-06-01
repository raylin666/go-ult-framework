// Package logger 提供日志记录封装，基于 Zap 实现。
// 支持 JSON 格式日志输出，提供应用日志、SQL 日志和请求日志分类功能。
package logger

import (
	"context"
	"time"
	"ult/pkg/types"

	"github.com/raylin666/go-utils/v2/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志类型常量。
const (
	LogApp     = "app"     // 应用日志类型
	LogSQL     = "sql"     // SQL 日志类型
	LogRequest = "request" // 请求日志类型
)

// Logger 日志记录器，封装 Zap Logger。
type Logger struct {
	*zap.Logger
}

// NewJSONLogger 创建 JSON 格式的日志记录器。
//
// 参数:
//   - opts: 日志选项列表
//
// 返回:
//   - *Logger: 日志记录器实例
//   - error: 创建错误
func NewJSONLogger(opts ...logger.Option) (*Logger, error) {
	zapLogger, err := logger.NewJSONLogger(opts...)
	return &Logger{zapLogger.Logger}, err
}

// UseApp 获取应用日志记录器，带有 TraceID。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - *zap.Logger: 应用日志记录器
func (log *Logger) UseApp(ctx context.Context) *zap.Logger {
	return log.Logger.Named(LogApp).With(zap.Any("trace_id", ctx.Value(types.TraceIdName)))
}

// UseSQL 获取 SQL 日志记录器，带有 TraceID。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - *zap.Logger: SQL 日志记录器
func (log *Logger) UseSQL(ctx context.Context) *zap.Logger {
	return log.Logger.Named(LogSQL).With(zap.Any("trace_id", ctx.Value(types.TraceIdName)))
}

// UseRequest 获取请求日志记录器，带有 TraceID。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - *zap.Logger: 请求日志记录器
func (log *Logger) UseRequest(ctx context.Context) *zap.Logger {
	return log.Logger.Named(LogRequest).With(zap.Any("trace_id", ctx.Value(types.TraceIdName)))
}

// RequestLogFormat 请求日志格式结构体。
type RequestLogFormat struct {
	ClientIp          string              `json:"client_ip"`           // 客户端 IP
	Method            string              `json:"method"`              // HTTP 方法
	Path              string              `json:"path"`                // 请求路径
	RequestProto      string              `json:"request_proto"`       // 请求协议
	RequestReferer    string              `json:"request_referer"`     // 请求来源
	RequestUa         string              `json:"request_ua"`          // 用户代理
	RequestPostData   string              `json:"request_post_data"`   // POST 数据
	RequestBodyData   string              `json:"request_body_data"`   // 请求体数据
	RequestHeaderData map[string][]string `json:"request_header_data"` // 请求头数据
	HttpCode          int                 `json:"http_code"`           // HTTP 状态码
	BusinessCode      int                 `json:"business_code"`       // 业务错误码
	BusinessMessage   string              `json:"business_message"`    // 业务消息
	RequestTime       time.Time           `json:"request_time"`        // 请求时间
	ResponseTime      time.Time           `json:"response_time"`       // 响应时间
	CostSeconds       float64             `json:"cost_seconds"`        // 耗时（秒）
}

// MarshalLogObject 实现 zapcore.ObjectMarshaler 接口。
// 使用直接字段访问替代反射，提升性能。
//
// 参数:
//   - enc: Zap 对象编码器
//
// 返回:
//   - error: 编码错误
func (rlf *RequestLogFormat) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("client_ip", rlf.ClientIp)
	enc.AddString("method", rlf.Method)
	enc.AddString("path", rlf.Path)
	enc.AddString("request_proto", rlf.RequestProto)
	enc.AddString("request_referer", rlf.RequestReferer)
	enc.AddString("request_ua", rlf.RequestUa)
	enc.AddString("request_post_data", rlf.RequestPostData)
	enc.AddString("request_body_data", rlf.RequestBodyData)
	enc.AddInt("http_code", rlf.HttpCode)
	enc.AddInt("business_code", rlf.BusinessCode)
	enc.AddString("business_message", rlf.BusinessMessage)
	enc.AddTime("request_time", rlf.RequestTime)
	enc.AddTime("response_time", rlf.ResponseTime)
	enc.AddFloat64("cost_seconds", rlf.CostSeconds)

	if rlf.RequestHeaderData != nil {
		enc.AddReflected("request_header_data", rlf.RequestHeaderData)
	}

	return nil
}

// RequestLog 打印请求日志。
// 使用 zap.Object 方式，通过 MarshalLogObject 接口编码日志字段。
//
// 参数:
//   - ctx: 上下文
//   - rlf: 请求日志格式
//   - err: 错误信息（可选）
func (log *Logger) RequestLog(ctx context.Context, rlf *RequestLogFormat, err error) {
	log.UseRequest(ctx).
		With(zap.Object("request", rlf)).
		With(zap.Error(err)).
		Info("请求日志")
}
