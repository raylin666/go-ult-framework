// Package errcode 提供统一错误码定义和管理。
// 支持多语言错误消息、HTTP 状态码映射、错误堆栈追踪和告警通知。
package errcode

import (
	"github.com/raylin666/go-utils/v2/errors"
)

// BusinessError 接口验证。
var _ BusinessError = (*businessError)(nil)

// BusinessError 业务错误接口。
// 定义统一的错误处理规范，包含 HTTP 状态码、业务错误码、消息、描述等。
// 实现了 error 接口，可以作为标准错误使用。
type BusinessError interface {
	error                                   // 实现 error 接口
	WithStackError(err error) BusinessError // 设置堆栈错误
	StackError() error                      // 获取堆栈错误
	BusinessCode() int                      // 获取业务错误码
	HTTPCode() int                          // 获取 HTTP 状态码
	Message() string                        // 获取错误消息
	Desc() string                           // 获取错误描述
	Alert() BusinessError                   // 设置告警标记
	WithDesc(desc string) BusinessError     // 设置错误描述
	IsAlert() bool                          // 是否需要告警
}

// businessError 业务错误实现结构体。
type businessError struct {
	httpCode     int    // HTTP 状态码
	businessCode int    // 业务错误码
	message      string // 错误消息
	desc         string // 错误描述
	stackError   error  // 堆栈错误
	isAlert      bool   // 是否告警
}

// NewError 创建新的业务错误实例。
//
// 参数:
//   - httpCode: HTTP 状态码
//   - businessCode: 业务错误码
//   - message: 错误消息
//
// 返回:
//   - BusinessError: 业务错误接口
func NewError(httpCode, businessCode int, message string) BusinessError {
	return &businessError{
		httpCode:     httpCode,
		businessCode: businessCode,
		message:      message,
		isAlert:      false,
	}
}

// WithStackError 设置堆栈错误。
func (e *businessError) WithStackError(err error) BusinessError {
	e.WithDesc(err.Error())
	e.stackError = errors.WithStack(err)
	return e
}

// StackError 获取堆栈错误。
func (e *businessError) StackError() error {
	return e.stackError
}

// HTTPCode 获取 HTTP 状态码。
func (e *businessError) HTTPCode() int {
	return e.httpCode
}

// BusinessCode 获取业务错误码。
func (e *businessError) BusinessCode() int {
	return e.businessCode
}

// Message 获取错误消息。
func (e *businessError) Message() string {
	return e.message
}

// Desc 获取错误描述。
func (e *businessError) Desc() string {
	return e.desc
}

// WithDesc 设置错误描述。
func (e *businessError) WithDesc(desc string) BusinessError {
	e.desc = desc
	return e
}

// Alert 设置告警标记。
func (e *businessError) Alert() BusinessError {
	e.isAlert = true
	return e
}

// IsAlert 判断是否需要告警。
func (e *businessError) IsAlert() bool {
	return e.isAlert
}

// Error 实现 error 接口，返回错误消息。
// 如果有描述，返回消息 + 描述；否则只返回消息。
func (e *businessError) Error() string {
	if e.desc != "" {
		return e.message + ": " + e.desc
	}
	return e.message
}
