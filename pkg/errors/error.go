package errors

import (
	goerrors "errors"
	"github.com/raylin666/go-utils/errors"
)

var _ BusinessError = (*businessError)(nil)

type BusinessError interface {
	// WithStackError 设置堆栈错误信息
	WithStackError(err error) BusinessError
	// StackError 获取带堆栈的错误信息
	StackError() error
	// BusinessCode 获取业务码
	BusinessCode() int
	// HTTPCode 获取 HTTP 状态码
	HTTPCode() int
	// Message 获取错误描述
	Message() string
	// Desc 获取错误说明
	Desc() string
	// WithDesc 设置错误说明
	WithDesc(desc string) BusinessError
	// WithAlert 设置告警通知
	WithAlert() BusinessError
	// IsAlert 是否开启告警通知
	IsAlert() bool
}

type businessError struct {
	httpCode     int    // HTTP 状态码
	businessCode int    // 业务码
	message      string // 错误描述
	desc         string // 错误说明
	stackError   error  // 含有堆栈信息的错误
	isAlert      bool   // 是否告警通知
}

func NewOriginalError(text string) error {
	return goerrors.New(text)
}

func NewError(httpCode, businessCode int, message string) BusinessError {
	return &businessError{
		httpCode:     httpCode,
		businessCode: businessCode,
		message:      message,
		isAlert:      false,
	}
}

func (e *businessError) WithStackError(err error) BusinessError {
	e.WithDesc(err.Error())
	e.stackError = errors.WithStack(err)
	return e
}

func (e *businessError) StackError() error {
	return e.stackError
}

func (e *businessError) HTTPCode() int {
	return e.httpCode
}

func (e *businessError) BusinessCode() int {
	return e.businessCode
}

func (e *businessError) Message() string {
	return e.message
}

func (e *businessError) Desc() string {
	return e.desc
}

func (e *businessError) WithDesc(desc string) BusinessError {
	e.desc = desc
	return e
}

func (e *businessError) WithAlert() BusinessError {
	e.isAlert = true
	return e
}

func (e *businessError) IsAlert() bool {
	return e.isAlert
}
