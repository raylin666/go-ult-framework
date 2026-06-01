// Package service 提供业务逻辑层实现。
package service

import (
	"github.com/google/wire"
)

// ProviderSet Wire 依赖注入提供者集合。
var ProviderSet = wire.NewSet(NewHeartbeatService)

