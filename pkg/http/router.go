package http

import "github.com/gin-gonic/gin"

var _ IRouter = (*Router)(nil)

type RouterGroup interface {
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
	IRouter
}

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

type Router struct {
	*gin.RouterGroup
}

func NewRouter(group *gin.RouterGroup) *Router {
	var r = new(Router)
	r.RouterGroup = group
	return r
}

func (r *Router) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &Router{r.RouterGroup.Group(relativePath, wrapHandlers(handlers...)...)}
}

func (r *Router) Any(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.Any(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) GET(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.GET(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) POST(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.POST(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.DELETE(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) PATCH(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.PATCH(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) PUT(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.PUT(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.OPTIONS(relativePath, wrapHandlers(handlers...)...)
}

func (r *Router) HEAD(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.HEAD(relativePath, wrapHandlers(handlers...)...)
}

// wrapHandlers 包装处理程序
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler
		funcs[i] = func(c *gin.Context) {
			ctx := newContext(c)
			defer recoveryContext(ctx)
			handler(ctx)
		}
	}

	return funcs
}
