// Package repo 提供数据仓库实现。
package repo

import (
	"github.com/google/wire"
)

// ProviderSet Wire 依赖注入提供者集合。
var ProviderSet = wire.NewSet(NewTestRepo)