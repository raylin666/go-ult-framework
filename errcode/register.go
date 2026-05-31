package errcode

import (
	"strings"
	"sync"
)

var registry = &codeRegistry{
	local:     ZhCN,
	codeTexts: zhCNText,
	httpCodes: httpCode,
}

type codeRegistry struct {
	mu        sync.RWMutex
	local     string
	codeTexts map[int]string
	httpCodes map[int]int
}

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

func GetRegistry() *codeRegistry {
	return registry
}

func RegisterHTTPCodes(codes map[int]int) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	for key, value := range codes {
		registry.httpCodes[key] = value
	}
}

func RegisterTexts(local string, texts map[int]string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	local = strings.ToLower(local)
	for key, value := range texts {
		registry.codeTexts[key] = value
	}
}

func GetText(code int) string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return registry.codeTexts[code]
}

func GetHTTPCode(code int) int {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return registry.httpCodes[code]
}

func New(code int) BusinessError {
	return NewError(GetHTTPCode(code), code, GetText(code))
}
