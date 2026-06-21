// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

import (
	"ult/pkg/types"

	"github.com/gin-gonic/gin"
)

var _ IRouter = (*Router)(nil)

// RouterGroup 路由组接口，定义路由分组功能。
type RouterGroup interface {
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
	IRouter
}

// IRouter 路由接口，定义 HTTP 方法路由注册功能。
type IRouter interface {
	Any(relativePath string, handlers ...HandlerFunc)
	GET(relativePath string, handlers ...HandlerFunc)
	POST(relativePath string, handlers ...HandlerFunc)
	DELETE(relativePath string, handlers ...HandlerFunc)
	PATCH(relativePath string, handlers ...HandlerFunc)
	PUT(relativePath string, handlers ...HandlerFunc)
	OPTIONS(relativePath string, handlers ...HandlerFunc)
	HEAD(relativePath string, handlers ...HandlerFunc)
}

// Router 路由结构体，封装 Gin 路由组。
type Router struct {
	*gin.RouterGroup
}

// NewRouter 创建新的 Router 实例。
//
// 参数:
//   - group: Gin 路由组
//
// 返回:
//   - *Router: 新创建的路由实例
func NewRouter(group *gin.RouterGroup) *Router {
	var r = new(Router)
	r.RouterGroup = group
	return r
}

// Group 创建子路由组。
//
// 参数:
//   - relativePath: 相对路径
//   - handlers: 处理函数列表
//
// 返回:
//   - RouterGroup: 新创建的路由组
func (r *Router) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &Router{r.RouterGroup.Group(relativePath, wrapHandlers(handlers...)...)}
}

// Any 注册所有 HTTP 方法的路由。
func (r *Router) Any(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.Any(relativePath, wrapHandlers(handlers...)...)
}

// GET 注册 GET 方法路由。
func (r *Router) GET(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.GET(relativePath, wrapHandlers(handlers...)...)
}

// POST 注册 POST 方法路由。
func (r *Router) POST(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.POST(relativePath, wrapHandlers(handlers...)...)
}

// DELETE 注册 DELETE 方法路由。
func (r *Router) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.DELETE(relativePath, wrapHandlers(handlers...)...)
}

// PATCH 注册 PATCH 方法路由。
func (r *Router) PATCH(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.PATCH(relativePath, wrapHandlers(handlers...)...)
}

// PUT 注册 PUT 方法路由。
func (r *Router) PUT(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.PUT(relativePath, wrapHandlers(handlers...)...)
}

// OPTIONS 注册 OPTIONS 方法路由。
func (r *Router) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.OPTIONS(relativePath, wrapHandlers(handlers...)...)
}

// HEAD 注册 HEAD 方法路由。
func (r *Router) HEAD(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.HEAD(relativePath, wrapHandlers(handlers...)...)
}

// wrapHandlers 包装处理函数，将自定义 HandlerFunc 转换为 Gin HandlerFunc。
// 复用 Request 中间件已初始化的 Context，避免重复创建。
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		funcs[i] = func(c *gin.Context) {
			// 尝试从 gin.Context 获取 Request 中间件已初始化的 Context 
			// 在 CreateRequest 创建请求中间件的 contextInitializer 中初始化了 Context
			if appCtx, exists := c.Get(types.CoreContextNameKey); exists {
				// 复用已存在的 Context（Request 中间件已初始化）
				ctx := appCtx.(Context)
				handler(ctx)
			} else {
				// 如果不存在（例如未使用 Request 中间件），才创建新的
				ctx, err := newContext(c)
				if err != nil {
					// newContext 已经设置了错误并中止了请求
					// 这里直接返回，不再执行 handler
					return
				}
				defer recoveryContext(ctx)
				handler(ctx)
			}
		}
	}

	return funcs
}
