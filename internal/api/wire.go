// Package api 提供 API 处理层实现。
package api

import (
	"github.com/google/wire"
)

// ProviderSet Wire 依赖注入提供者集合。
var ProviderSet = wire.NewSet(NewHeartbeatHandler)
