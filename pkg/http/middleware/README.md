# 中间件使用指南

## 概述

新的中间件管理系统提供了更灵活、更强大的中间件配置能力，支持优先级控制、动态添加和条件启用。

## 基本使用

### 1. 使用 WithMiddleware 添加中间件

```go
package main

import (
    "time"
    "ult/config"
    "ult/pkg/http"
    "ult/pkg/http/middleware"
    "ult/pkg/logger"
)

func main() {
    // 创建服务器
    server := http.NewServer(
        config,
        logger,
        []http.ServerOption{},
        // 添加自定义中间件
        http.WithMiddleware(
            middleware.NewRateLimit(&middleware.RateLimitConfig{
                Enabled: true,
                Limit:   100,
                Window:  time.Minute,
            }),
        ),
    )
}
```

### 2. 使用 UseMiddleware 方法

```go
// 在服务器创建后添加中间件
server.UseMiddleware(middleware.NewRateLimit(&middleware.RateLimitConfig{
    Enabled: true,
    Limit:   100,
}))

// 链式添加多个中间件
server.UseMiddleware(middleware1).
       UseMiddleware(middleware2).
       UseMiddleware(middleware3)
```

### 3. 使用 UseMiddlewareFunc 快速添加

```go
// 使用函数方式快速添加中间件
server.UseMiddlewareFunc("custom-auth", middleware.PriorityNormal, func(ctx *gin.Context) {
    // 自定义认证逻辑
    token := ctx.GetHeader("Authorization")
    if token == "" {
        ctx.AbortWithStatus(401)
        return
    }
    ctx.Next()
})
```

## 中间件优先级

中间件按优先级顺序执行，数值越小优先级越高：

```go
const (
    PriorityHighest Priority = iota // 最高优先级（如 Recovery）
    PriorityHigh                    // 高优先级（如 CORS）
    PriorityNormal                  // 正常优先级（如日志、验证）
    PriorityLow                     // 低优先级（如业务中间件）
)
```

### 执行顺序示例

```go
// 中间件执行顺序：
// 1. Recovery (PriorityHighest) - 异常恢复
// 2. CORS (PriorityHigh) - 跨域处理
// 3. Request Handler (内置) - 请求处理
// 4. Custom Auth (PriorityNormal) - 自定义认证
// 5. Rate Limit (PriorityLow) - 限流
```

## 创建自定义中间件

### 方式一：实现 Middleware 接口

```go
package middleware

import (
    "github.com/gin-gonic/gin"
)

// AuthConfig 认证中间件配置
type AuthConfig struct {
    Enabled  bool
    Secret   string
}

// Auth 认证中间件
type Auth struct {
    config *AuthConfig
}

// NewAuth 创建认证中间件
func NewAuth(config *AuthConfig) *Auth {
    return &Auth{config: config}
}

func (a *Auth) Name() string { return "auth" }
func (a *Auth) Priority() Priority { return PriorityNormal }
func (a *Auth) Enabled() bool { return a.config.Enabled }

func (a *Auth) Handler() HandlerFunc {
    return func(ctx *gin.Context) {
        // 认证逻辑
        token := ctx.GetHeader("Authorization")
        if token == "" {
            ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        // 验证 token...
        ctx.Next()
    }
}

// 使用
server.UseMiddleware(NewAuth(&AuthConfig{
    Enabled: true,
    Secret:  "your-secret-key",
}))
```

### 方式二：使用函数式中间件

```go
// 快速创建函数式中间件
server.UseMiddlewareFunc("simple-auth", middleware.PriorityNormal, func(ctx *gin.Context) {
    token := ctx.GetHeader("Authorization")
    if token == "" {
        ctx.AbortWithStatus(401)
        return
    }
    ctx.Next()
})
```

## 内置中间件

### CORS 中间件

```go
// 默认配置
http.EnableCors([]string{"*"})

// 生产环境配置
server.UseMiddleware(middleware.NewCORS(&middleware.CORSConfig{
    Enabled:        true,
    AllowedOrigins: []string{"https://example.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowCredentials: true,
}))
```

### Recovery 中间件

```go
// 自动注册，无需手动添加
// 配置告警通知
http.EnableAlertNotify(alertHandler)
```

## 路由级中间件

```go
// 为特定路由组添加中间件
func setupRoutes(server *http.HTTPServer) {
    apiGroup := server.CreateRouterGroup().Group("/api")
    
    // 为 API 组添加认证中间件
    apiGroup.UseMiddlewareFunc("api-auth", middleware.PriorityNormal, func(ctx *gin.Context) {
        // API 认证逻辑
        ctx.Next()
    })
    
    // 注册路由
    apiGroup.GET("/users", userHandler.List)
    apiGroup.POST("/users", userHandler.Create)
}
```

## 条件启用中间件

```go
// 根据环境启用中间件
if !isProduction {
    server.UseMiddleware(middleware.NewDebug(&middleware.DebugConfig{
        Enabled: true,
    }))
}

// 根据配置启用中间件
if config.Features.RateLimit {
    server.UseMiddleware(middleware.NewRateLimit(&middleware.RateLimitConfig{
        Enabled: true,
        Limit:   config.RateLimit.MaxRequests,
    }))
}
```

## 迁移指南

### 从旧版中间件迁移

**旧版（go-utils 中间件）：**
```go
// 已移除 go-utils 中间件支持
// 请迁移到新的中间件系统
```

**新版（推荐）：**
```go
http.WithMiddleware(newMiddleware)
```

### 主要特性

| 特性 | 说明 |
|------|------|
| 优先级控制 | ✅ 支持 PriorityHighest/High/Normal/Low |
| 条件启用 | ✅ 通过 Config.Enabled() 控制 |
| 动态添加 | ✅ 使用 UseMiddleware() 方法 |
| 链式调用 | ✅ 支持 server.UseMiddleware().UseMiddleware() |
| 路由级中间件 | ✅ 可在路由组中添加 |
| 中间件管理 | ✅ Manager 提供完整的中间件管理功能 |

## 最佳实践

1. **优先级设置**：
   - Recovery、Trace 等基础中间件使用 PriorityHighest
   - CORS、安全检查使用 PriorityHigh
   - 业务中间件使用 PriorityNormal 或 PriorityLow

2. **命名规范**：
   - 中间件名称应简洁明了，如 "auth"、"rate-limit"、"cors"

3. **配置管理**：
   - 使用配置结构体管理中间件参数
   - 提供 Enabled 字段支持条件启用

4. **错误处理**：
   - 中间件中发生错误时使用 ctx.Abort() 中止请求
   - 设置适当的 HTTP 状态码和错误信息

5. **性能考虑**：
   - 避免在中间件中进行耗时操作
   - 使用 sync.Pool 复用对象（如 Recovery 中间件）

## 示例项目

完整的中间件使用示例请参考：
- [pkg/http/middleware/cors.go](../pkg/http/middleware/cors.go)
- [pkg/http/middleware/recovery.go](../pkg/http/middleware/recovery.go)
- [pkg/http/server.go](../pkg/http/server.go)