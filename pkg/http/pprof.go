// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/raylin666/go-utils/v2/server/system"
)

// RegisterPProf 注册 PProf 性能分析路由到 Gin 引擎。
// 自动判断环境，只在非生产环境启用。
// 访问路径: /debug/pprof
//
// 参数:
//   - engine: Gin 引擎实例
//   - environment: 环境标识
func RegisterPProf(engine *gin.Engine, environment string) {
	// 在生产环境不启用 PProf
	if system.NewEnvironment(environment).IsProd() {
		return
	}

	// 注册 PProf 路由
	pprof.Register(engine)
}
