package errcode

import (
	"strings"
)

var registry = &codeRegistry{
	local:     ZhCN,
	codeTexts: zhCNText,
	httpCodes: httpCode,
}

type codeRegistry struct {
	local     string
	codeTexts map[int]string
	httpCodes map[int]int
}

func NewRegistry(local string) {
	local = strings.ToLower(local)
	switch local {
	case ZhCN:
		registry.local = ZhCN
		registry.codeTexts = zhCNText
	case EnUS:
		registry.local = EnUS
		registry.codeTexts = enUSText
	default:
		registry.local = ZhCN
		registry.codeTexts = zhCNText
	}
}

func GetRegistry() *codeRegistry {
	return registry
}

func RegisterHTTPCodes(codes map[int]int) {
	for key, value := range codes {
		registry.httpCodes[key] = value
	}
}

func RegisterTexts(local string, texts map[int]string) {
	local = strings.ToLower(local)
	for key, value := range texts {
		registry.codeTexts[key] = value
	}
}

func GetText(code int) string {
	return registry.codeTexts[code]
}

func GetHTTPCode(code int) int {
	return registry.httpCodes[code]
}

func New(code int) BusinessError {
	return NewError(GetHTTPCode(code), code, GetText(code))
}
