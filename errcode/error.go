package errcode

import (
	"github.com/raylin666/go-utils/v2/errors"
)

var _ BusinessError = (*businessError)(nil)

type BusinessError interface {
	WithStackError(err error) BusinessError
	StackError() error
	BusinessCode() int
	HTTPCode() int
	Message() string
	Desc() string
	Alert() BusinessError
	WithDesc(desc string) BusinessError
	IsAlert() bool
}

type businessError struct {
	httpCode     int
	businessCode int
	message      string
	desc         string
	stackError   error
	isAlert      bool
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

func (e *businessError) Alert() BusinessError {
	e.isAlert = true
	return e
}

func (e *businessError) IsAlert() bool {
	return e.isAlert
}