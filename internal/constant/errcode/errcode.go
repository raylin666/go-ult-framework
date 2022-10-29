package errcode

import "ult/pkg/code"

// RegisterTexts 注册业务状态提示码
func RegisterTexts() {
	code.Get().WithTexts(code.EnUS, enUSText)
	code.Get().WithTexts(code.ZhCN, zhCNText)
}
