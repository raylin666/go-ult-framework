// Package errcode 提供统一错误码定义和管理。
package errcode

import (
	"strings"
	"sync"
)

// registry 全局错误码注册表实例。
var registry = &codeRegistry{
	local:     ZhCN,
	codeTexts: zhCNText,
	httpCodes: httpCode,
}

// codeRegistry 错误码注册表结构体。
// 支持多语言错误消息和 HTTP 状态码映射。
type codeRegistry struct {
	mu        sync.RWMutex   // 读写锁，保证并发安全
	local     string         // 当前语言
	codeTexts map[int]string // 错误码消息映射
	httpCodes map[int]int    // 错误码到 HTTP 状态码映射
}

// NewRegistry 创建错误码注册表并设置语言。
//
// 参数:
//   - local: 语言代码 (zh-cn/en-us)
func NewRegistry(local string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

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

// GetRegistry 获取错误码注册表实例。
//
// 返回:
//   - *codeRegistry: 注册表实例
func GetRegistry() *codeRegistry {
	return registry
}

// RegisterHTTPCodes 注册 HTTP 状态码映射。
//
// 参数:
//   - codes: 错误码到 HTTP 状态码的映射
func RegisterHTTPCodes(codes map[int]int) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	for key, value := range codes {
		registry.httpCodes[key] = value
	}
}

// RegisterTexts 注册错误码消息文本。
//
// 参数:
//   - local: 语言代码
//   - texts: 错误码到消息文本的映射
func RegisterTexts(local string, texts map[int]string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	local = strings.ToLower(local)
	for key, value := range texts {
		registry.codeTexts[key] = value
	}
}

// GetText 根据错误码获取消息文本。
//
// 参数:
//   - code: 业务错误码
//
// 返回:
//   - string: 错误消息文本
func GetText(code int) string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return registry.codeTexts[code]
}

// GetHTTPCode 根据业务错误码获取 HTTP 状态码。
//
// 参数:
//   - code: 业务错误码
//
// 返回:
//   - int: HTTP 状态码
func GetHTTPCode(code int) int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return registry.httpCodes[code]
}

// New 根据业务错误码创建业务错误实例。
// 自动从注册表获取 HTTP 状态码和消息文本。
//
// 参数:
//   - code: 业务错误码
//
// 返回:
//   - BusinessError: 业务错误实例
func New(code int) BusinessError {
	return NewError(GetHTTPCode(code), code, GetText(code))
}
