// Package server 提供服务器层实现。
package server

import "github.com/google/wire"

// ProviderSet Wire 依赖注入提供者集合。
var ProviderSet = wire.NewSet(NewHTTPServer)
