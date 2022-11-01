package errcode

import (
	pkg_code "ult/pkg/code"
	"ult/pkg/errors"
)

func NewError(code int) errors.BusinessError {
	return errors.NewError(pkg_code.Get().GetHttpCode(code), code, pkg_code.Get().GetText(code))
}

// RegisterNewMerged 注册合并业务状态
func RegisterNewMerged() {
	pkg_code.Get().WithHttpCodes(httpCode)
	pkg_code.Get().WithTexts(pkg_code.EnUS, enUSText)
	pkg_code.Get().WithTexts(pkg_code.ZhCN, zhCNText)
}