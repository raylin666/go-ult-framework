# ULT Web API 框架 (基于 GIN)

本框架是基于 `GIN` 进行模块化设计的 API 框架，封装了常用的功能，使用简单，致力于进行快速的业务研发，同时增加了更多限制，约束项目组开发成员，规避混乱无序及自由随意的编码。

提供了方便快捷的 `Makefile` 文件 (帮你快速的生成、构建、执行项目内容)。

当你所需命令不存在时可添加到此文件中, 实现命令统一管理。这也大大的提高了开发者的开发效率, 让开发者更专注于业务代码。

---

## 核心特性

- **分层架构设计**: 采用经典的分层架构 + 依赖注入设计模式，职责清晰
- **严格调用链**: api → service → data 单向调用，避免循环依赖
- **依赖注入**: 使用 Google Wire 实现依赖注入，降低组件耦合
- **中间件管理**: 支持优先级和依赖管理的中间件链，自动验证依赖关系（不满足时 panic）
- **统一错误处理**: 完整的错误码管理系统，支持多语言和告警，BusinessError 实现 error 接口
- **多连接支持**: 支持多数据库和 Redis 连接配置
- **代码生成**: 集成 GORM Gen，自动生成查询器代码
- **优雅关闭**: 完整的服务器生命周期管理和优雅关闭机制
- **链路追踪**: 内置 TraceID 支持，便于分布式追踪，并发安全设计
- **配置管理**: YAML 配置文件，支持多环境配置
- **Docker 支持**: 提供 Dockerfile 和 Docker Compose 配置
- **并发安全**: TraceID 和 RequestContext 使用 sync.Once 确保并发安全
- **性能优化**: Header 和 RequestContext 性能优化，减少内存分配和 GC 压力
- **内存安全**: RawData 返回副本，避免内存泄漏和数据污染
- **对象池**: Context 使用 sync.Pool 复用，减少内存分配和 GC 压力
- **告警通知**: Recovery 中间件支持邮件告警，集成 proposal 通知机制

---

## 整体架构设计

### 架构概览

本框架采用经典的 **分层架构 + 依赖注入** 设计模式，整体架构清晰合理：

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              cmd/                                        │
│                     (入口层 + Wire 依赖注入)                              │
├─────────────────────────────────────────────────────────────────────────┤
│                           internal/                                      │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┬───────┐ │
│  │    router    │      api     │    service   │     data     │  app  │ │
│  │   (路由层)   │   (处理层)   │   (逻辑层)   │   (数据层)   │(工具) │ │
│  └──────────────┴──────────────┴──────────────┴──────────────┴───────┘ │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │                         server                                    │  │
│  │                       (服务器层)                                   │  │
│  └──────────────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────────────┤
│                              pkg/                                        │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┬───────────┐  │
│  │   http   │    db    │  cache   │  logger  │  app     │  types    │  │
│  │(HTTP封装)│ (数据库) │ (缓存)   │  (日志)  │(应用管理)│ (类型)    │  │
│  ├──────────┴──────────┴──────────┴──────────┴──────────┴───────────┤  │
│  │  repositories  │  proposal  │          notify                    │  │
│  │  (数据仓库)    │  (告警提案)│        (告警通知)                  │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │             pkg/http/middleware (中间件管理)                      │  │
│  │     Recovery │ CORS │ Request │ 自定义中间件                      │  │
│  └──────────────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────────────┤
│                            config/                                       │
│                         (配置管理)                                        │
├─────────────────────────────────────────────────────────────────────────┤
│                            errcode/                                      │
│                        (错误码定义)                                       │
└─────────────────────────────────────────────────────────────────────────┘
```

### 分层职责说明

| 层级 | 目录 | 职责说明 |
| --- | --- | --- |
| 入口层 | `cmd/` | 应用启动入口，Wire 依赖注入绑定 |
| 路由层 | `internal/router/` | HTTP 路由注册，将 API 处理器绑定到路由路径 |
| 处理层 | `internal/api/` | HTTP 请求处理，数据校验和响应，调用服务层 |
| 逻辑层 | `internal/service/` | 业务逻辑处理，调用数据层获取数据 |
| 数据层 | `internal/data/` | 数据仓库实例管理，数据库/缓存/RPC 操作 |
| 服务器层 | `internal/server/` | HTTP 服务器创建和配置 |
| 工具层 | `internal/app/` | 公共工具包（日志、日期时间、环境、JWT 等） |
| HTTP封装 | `pkg/http/` | Gin 框架封装，Context、Handler、Response、中间件管理 |
| 中间件 | `pkg/http/middleware/` | Recovery、CORS、Request 中间件及中间件管理器 |
| 数据库 | `pkg/db/` | GORM 数据库连接封装，支持连接重试 |
| 缓存 | `pkg/cache/` | Redis 连接封装，支持连接重试 |
| 日志 | `pkg/logger/` | Zap 日志封装，支持文件轮转和请求/SQL 日志 |
| 应用管理 | `pkg/app/` | 应用生命周期管理，服务器接口定义 |
| 类型定义 | `pkg/types/` | 上下文类型、链路追踪类型、中间件常量 |
| 告警提案 | `pkg/proposal/` | 告警消息定义和通知处理器类型 |
| 告警通知 | `pkg/notify/` | 异常恢复告警通知（邮件等） |
| 配置 | `config/` | YAML 配置文件加载和解析 |
| 错误码 | `errcode/` | 统一错误码定义和管理，支持多语言 |
| 数据仓库 | `pkg/repositories/` | 数据仓库抽象层，管理多连接 |

### 调用链设计

项目定义了严格的单向调用链：

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    api      │ ──> │   service   │ ──> │    data     │
│  (处理层)   │     │  (逻辑层)   │     │  (数据层)   │
└─────────────┘     └─────────────┘     └─────────────┘
```

**核心原则：**
- 单向依赖，避免循环调用
- 职责分离清晰，便于维护和测试
- `api` 层只做数据校验和响应
- `service` 层专注业务逻辑
- `data` 层处理数据操作

### 依赖注入设计

采用 Google Wire 实现依赖注入，降低组件耦合：

```go
// cmd/wire.go
//go:build wireinject
// +build wireinject

package main

import (
    "ult/config"
    "ult/internal/api"
    "ult/internal/app"
    "ult/internal/data"
    "ult/internal/data/repo"
    "ult/internal/router"
    "ult/internal/server"
    "ult/internal/service"
    pkgapp "ult/pkg/app"

    "github.com/google/wire"
)

func initApp(conf *config.Config, tools *app.Tools) (*pkgapp.App, func(), error) {
    panic(wire.Build(
        data.ProviderSet,
        repo.ProviderSet,
        service.ProviderSet,
        api.ProviderSet,
        router.ProviderSet,
        server.ProviderSet,
        newApp))
}
```

各模块通过 `ProviderSet` 注册依赖，Wire 自动生成依赖注入代码。注意实际使用 `panic(wire.Build(...))` 模式，这是 Wire 的标准用法。

Wire 自动生成的依赖链（`cmd/wire_gen.go`）：

```go
func initApp(conf *config.Config, tools *app.Tools) (*app2.App, func(), error) {
    dataData, cleanup := data.NewData(conf, tools)
    dataRepo := data.NewDataRepo(dataData)
    testRepo := repo.NewTestRepo(dataData)
    heartbeatService := service.NewHeartbeatService(dataRepo, testRepo, tools)
    heartbeatInterface := api.NewHeartbeatHandler(heartbeatService, tools)
    httpRouter := router.NewHTTPRouter(heartbeatInterface)
    httpServer := server.NewHTTPServer(conf, tools, httpRouter)
    appApp := newApp(conf, tools, httpServer)
    return appApp, func() { cleanup() }, nil
}
```

### 数据流设计

```
HTTP Request
    │
    ▼
┌─────────────────────────────────────┐
│         pkg/http/server             │
│    (Gin Engine + Middleware)        │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  Recovery 中间件 (PriorityHighest)  │
│  - 异常恢复、告警通知               │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  Request 中间件 (PriorityHigh)      │
│  - Context 初始化 (对象池复用)      │
│  - 验证器设置                       │
│  - 响应处理 (defer)                 │
│  - 404 拦截                         │
│  依赖: Recovery                     │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│  CORS 中间件 (PriorityHigh)         │
│  - 跨域处理                         │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│        internal/router               │
│         (Route Matching)             │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│          internal/api                │
│    (Request Validation + Response)   │
│    - 参数绑定                        │
│    - 数据校验 (Validator)           │
│    - 调用 Service                    │
│    - 响应处理                        │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│        internal/service              │
│       (Business Logic)               │
│    - 业务逻辑处理                    │
│    - 调用 Data Repo                  │
│    - 数据转换                        │
└─────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────┐
│          internal/data               │
│    (Database/Cache Operations)       │
│    - 数据库查询                      │
│    - Redis 操作                      │
│    - 事务管理                        │
└─────────────────────────────────────┘
    │
    ▼
HTTP Response
```

### 错误处理设计

统一的错误处理机制：

```go
// 错误码注册表
type codeRegistry struct {
    mu        sync.RWMutex       // 并发安全
    local     string             // 语言设置
    codeTexts map[int]string     // 错误码消息
    httpCodes map[int]int        // HTTP状态码映射
}

// 业务错误接口
type BusinessError interface {
    error                                    // 实现 error 接口，可作为标准错误使用
    WithStackError(err error) BusinessError  // 堆栈追踪（同时设置描述为 err.Error()）
    StackError() error                       // 获取堆栈错误
    HTTPCode() int                           // HTTP状态码
    BusinessCode() int                       // 业务错误码
    Message() string                         // 错误消息
    Desc() string                            // 错误描述
    Alert() BusinessError                    // 告警标记
    WithDesc(desc string) BusinessError      // 设置错误描述
    IsAlert() bool                           // 是否需要告警
}

// Error() 方法实现
// 返回完整的错误信息：消息 + 描述
func (e *businessError) Error() string {
    if e.desc != "" {
        return e.message + ": " + e.desc
    }
    return e.message
}
```

**关键改进**：
- ✅ BusinessError 实现了 error 接口，可以作为标准错误使用
- ✅ 符合 Go 的错误处理惯例（`if err != nil`）
- ✅ Validator() 方法返回 error 类型，语义更清晰
- ✅ WithStackError() 同时设置描述为 err.Error()，确保堆栈错误有描述信息

**错误码使用方式**：

错误码常量是 `int` 类型，需要通过 `errcode.New()` 创建 BusinessError 实例，或使用 `vars.go` 中预创建的 `Err*` 变量：

```go
// 方式一: 使用 errcode.New() 从错误码常量创建
errcode.New(errcode.DataSelectError).WithDesc("数据库查询失败")

// 方式二: 使用预创建的 Err* 变量（推荐）
errcode.ErrDataSelectError.WithDesc("数据库查询失败")

// 设置堆栈信息
errcode.New(errcode.ServerError).WithStackError(err)
// 或
errcode.ErrServerError.WithStackError(err)

// 设置告警标记
errcode.ErrServerError.Alert()
```

---

## 集成组件

| 名称 | 描述 | 版本 |
| --- | --- | --- |
| Gin | HTTP Web 框架 | 1.11.0 |
| GORM | ORM 数据库组件 | 1.31.1 |
| Wire | 依赖注入 | 0.7.0 |
| Zap | 结构化日志 | 1.27.0 |
| go-redis | Redis 客户端 | 9.0.5 |
| JWT | JWT 认证 | 5.3.1 |
| Validator | 数据校验 | 10.27.0 |
| CORS | 跨域处理 | - |
| Pprof | 性能剖析 | 1.4.0 |
| UUID | 唯一标识生成 | 1.6.0 |
| YAML | 配置文件解析 | 3.0.1 |
| go-utils | 工具库（Server/Middleware/Auth/Validator/Logger/Mail 等） | v2 |

---

## 目录结构

```
go-ult-framework/
├── cmd/                    # 应用入口层
│   ├── main.go             # 主程序入口
│   ├── wire.go             # Wire 依赖注入定义
│   └── wire_gen.go         # Wire 自动生成代码
├── config/                 # 配置管理层
│   ├── autoload.go         # 配置加载器
│   ├── autoload/           # 配置模块
│   │   ├── app.go          # 应用配置
│   │   ├── db.go           # 数据库配置
│   │   ├── server.go       # 服务器配置
│   │   ├── redis.go        # Redis 配置
│   │   ├── jwt.go          # JWT 配置
│   │   ├── logger.go       # 日志配置
│   │   ├── language.go     # 语言配置
│   │   ├── datetime.go     # 时间配置
│   │   ├── notify.go       # 告警配置
│   │   └── validator.go    # 校验器配置
│   └── .env.example.yml    # 配置示例文件
├── errcode/                # 错误码定义层
│   ├── error.go            # 业务错误接口和实现
│   ├── code.go             # 错误码常量
│   ├── register.go         # 错误码注册表
│   ├── errhttp.go          # HTTP状态码映射
│   ├── vars.go             # 预创建 BusinessError 变量
│   ├── zh-cn.go            # 中文错误消息
│   └── en-us.go            # 英文错误消息
├── generate/               # 代码生成器
│   └── gormgen/            # GORM Gen 生成器
│       ├── main.go         # 生成器入口
│       └── db/             # 生成器配置
│           └── default.go  # 默认数据库配置
├── internal/               # 内部业务层（不可外部引用）
│   ├── api/                # API 处理层
│   │   ├── heartbeat.go    # 健康检查 API
│   │   └── wire.go         # Wire ProviderSet
│   ├── app/                # 应用工具包
│   │   ├── tools.go        # 公共工具实例（Logger/Datetime/Environment/JWT）
│   │   └── logo.go         # 启动信息打印
│   ├── data/               # 数据层
│   │   ├── data.go         # 数据仓库管理（Data 接口、DataRepo）
│   │   ├── model/          # 数据模型
│   │   │   ├── model.go    # 基础模型
│   │   │   └── test.go     # 测试模型
│   │   ├── repo/           # 数据仓库实现
│   │   │   ├── test.go     # 测试仓库
│   │   │   └── wire.go     # Wire ProviderSet
│   │   ├── redis/          # Redis 操作封装
│   │   │   ├── client.go   # Redis 客户端
│   │   │   └ action/       # Redis 操作
│   │   │       └ lock.go   # 分布式锁
│   │   └── dbquery/        # GORM Gen 查询器（自动生成）
│   │       ├── gen.go      # 生成器代码
│   │       └ api_test.gen.go # 测试查询器
│   ├── router/             # 路由层
│   │   ├── router.go       # 路由注册
│   │   └ heartbeat.go      # 健康检查路由
│   ├── server/             # 服务器层
│   │   ├── http.go         # HTTP 服务器创建
│   │   └ server.go         # Wire ProviderSet
│   └── service/            # 服务层
│       ├── heartbeat.go    # 健康检查服务
│       └ wire.go           # Wire ProviderSet
├── pkg/                    # 通用封装层（可外部引用）
│   ├── app/                # 应用管理
│   │   ├── app.go          # 应用生命周期管理、上下文函数
│   │   └ server.go         # 服务器接口定义、ServerAgreement
│   ├── http/               # HTTP 封装
│   │   ├── server.go       # Gin 服务器封装（UseMiddleware、CreateRequest）
│   │   ├── context.go      # 请求上下文封装（对象池、sync.Once）
│   │   ├── router.go       # 路由组封装
│   │   ├── response.go     # 响应处理（SuccessResponse、ErrorResponse）
│   │   ├── option.go       # 服务器选项（Timeout、Middleware、OpenBrowser）
│   │   ├── pprof.go        # PProf 注册（非生产环境）
│   │   └ middleware/       # 中间件管理
│   │       ├── adapter.go  # 中间件管理器（Manager，依赖验证）
│   │       ├── handler.go  # 中间件接口、函数式中间件、Priority 类型别名
│   │       ├── request.go  # Request 中间件（Context 初始化、响应处理）
│   │       ├── cors.go     # CORS 中间件
│   │       └ recovery.go   # Recovery 中间件（异常恢复、告警通知）
│   ├── db/                 # 数据库封装
│   │   ├── gorm.go         # GORM 连接封装（重试机制、Ping）
│   │   └ logger.go         # 数据库日志（TraceID 集成）
│   ├── cache/              # 缓存封装
│   │   └ redis.go          # Redis 连接封装（重试机制）
│   ├── logger/             # 日志封装
│   │   └ logger.go         # Zap 日志封装（App/SQL/Request 分类日志）
│   ├── repositories/       # 数据仓库封装
│   │   ├── repo.go         # 数据仓库接口（DataRepo）
│   │   ├── db.go           # 数据库仓库接口（DbRepo）
│   │   └ redis.go          # Redis 仓库接口（RedisRepo）
│   ├── types/              # 类型定义
│   │   ├── context.go      # 请求上下文类型、ContextKey 常量
│   │   ├── trace.go        # 链路追踪类型（TraceIdName）
│   │   └ middleware.go      # 中间件名称常量
│   ├── proposal/            # 告警提案
│   │   ├── alert.go        # 告警消息定义（AlertMessage）、NotifyHandler 类型
│   └ notify/               # 告警通知
│   │   └ recover/          # 异常恢复通知
│   │       └ email/        # 邮件通知
│   │           ├── alert.go # 告警处理（异步邮件发送）
│   │           └ email_template.go # 邮件模板（HTML）
├── static/                 # 静态文件
│   └ db/                   # 数据库脚本
│       └ ult.sql           # SQL 初始化脚本
├── .env.example.yml        # 配置示例
├── .gitignore              # Git 忽略文件
├── Dockerfile              # Docker 构建文件
├── docker-compose.yml      # Docker Compose 配置
├── Makefile                # 构建命令
├── go.mod                  # Go 模块定义
└── go.sum                  # Go 模块依赖
```

---

## 快速开始

### 1. 下载仓库

```bash
git clone git@github.com:raylin666/go-ult-framework.git
cd go-ult-framework
```

### 2. 初始化项目

```bash
make init
```

此命令会创建 `.env.yml` 配置文件（从 `.env.example.yml` 复制）。

### 3. 配置数据库和 Redis

编辑 `.env.yml` 文件，配置数据库和 Redis 连接：

```yaml
db:
  default:
    driver: mysql
    host: 127.0.0.1
    port: 3306
    username: root
    password: password
    db_name: ult_framework
    charset: utf8mb4
    prefix: api_
    max_open_conn: 100
    max_idle_conn: 10
    max_life_time: 0
    max_retries: 3
    retry_delay: 2
    parse_time: true
    loc: Local

redis:
  default:
    network: tcp
    addr: 127.0.0.1
    port: 6379
    password: ""
    db: 0
    min_idle_conns: 10
    max_retries: 3
    retry_delay: 2
```

> **注意**: 配置文件具体请参考 `.env.example.yml` 获取完整配置项。

### 4. 下载依赖并生成 Wire 代码

```bash
make generate
```

此命令会：
- 下载 Go 模块依赖
- 安装 Wire 工具
- 生成依赖注入代码

### 5. 启动服务

```bash
make run
```

访问健康检查接口验证服务是否正常：

```bash
curl 127.0.0.1:10001/heartbeat/state
```

成功响应示例：

```json
{
    "trace_id": "b2401c9e-1f6f-4183-952a-b539ddabbb71",
    "data": {
        "status": "healthy",
        "timestamp": "2024-01-01T10:00:00Z",
        "components": {
            "db_default": {"status": "healthy"},
            "redis_default": {"status": "healthy"}
        }
    }
}
```

### 6. 编译执行文件

```bash
make build
```

编译成功后，可通过以下命令运行：

```bash
./bin/server
```

### Docker 部署

支持使用 Docker 和 Docker Compose 部署：

```bash
# 使用 Dockerfile 构建
docker build -t go-ult-framework .

# 使用 Docker Compose
docker-compose up -d
```

> **注意**: 当前 Dockerfile 基础镜像不可用，请建议根据实际需求更换镜像和修改 Dockerfile 配置。

---

## Makefile 命令说明

| 命令 | 说明 |
| --- | --- |
| `make init` | 初始化项目，创建配置文件 |
| `make generate` | 下载依赖并生成 Wire 代码 |
| `make wire` | 生成依赖注入文件 |
| `make gormgen` | 生成 GORM Gen 查询器代码 |
| `make run` | 启动开发服务器 |
| `make build` | 编译生产执行文件（带 git 版本信息） |
| `make help` | 显示帮助信息 |

---

## 开发规范

### 调用链规范

项目定义了严格的单向调用链：

```
api (处理层) → service (逻辑层) → data (数据层)
```

**核心原则：**
- `api` 层只做数据校验和响应，不处理业务逻辑
- `service` 层专注业务逻辑，不直接操作数据库
- `data` 层处理数据操作，不包含业务判断
- 逻辑代码只能下沉，禁止反向调用或互调

### 层级职责规范

| 层级 | 允许使用 | 禁止使用 |
| --- | --- | --- |
| `api` | `service`、`logger`、`validator` | `config`、`dataRepo`、直接数据库操作 |
| `service` | `data`、`logger`、其他 `service` | `config`、直接数据库操作 |
| `data` | `db`、`redis`、`logger` | `service`、`api` |

### 错误处理规范

1. **统一错误类型**：所有业务错误必须返回 `BusinessError` 类型
2. **错误码使用**：使用 `errcode.New(code)` 或 `errcode.Err*` 预创建变量，不要直接对 int 常量调用方法
3. **错误响应**：只能在 `api` 层处理错误响应

```go
// 正确示例 - 使用 errcode.New() 从错误码常量创建
func (s *AccountService) GetByID(ctx context.Context, id int) (*model.Account, errcode.BusinessError) {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, errcode.New(errcode.DataSelectError).WithDesc(err.Error())
    }
    if account == nil {
        return nil, errcode.ErrDataNotExistError
    }
    return account, nil
}

// 也可使用预创建的 Err* 变量（推荐，更简洁）
return nil, errcode.ErrDataSelectError.WithDesc(err.Error())
```

### 数据响应规范

| 方法 | 说明 | 使用场景 |
| --- | --- | --- |
| `ctx.WithPayload(data)` | 设置成功响应数据 | 业务处理成功 |
| `ctx.WithAbortError(err)` | 设置 BusinessError 并中断请求 | 业务处理失败 |

**注意事项：**
- `WithAbortError` 接受 `errcode.BusinessError` 参数，会自动设置 HTTP 状态码并 `AbortWithStatus`
- `WithAbortError` 后必须 `return`，否则会继续执行
- `WithPayload` 一般放在最后，否则需 `return`

---

## 模块开发指南

### 创建新模块

以创建 `account` 模块为例：

#### 步骤 1：创建数据模型

在 `internal/data/model/` 目录创建模型文件：

```go
// Package model 提供数据模型定义。
package model

// Account 账户数据模型。
type Account struct {
    ID       int64  `gorm:"column:id;primaryKey" json:"id"`
    UserName string `gorm:"column:username" json:"username"`
    Password string `gorm:"column:password" json:"password"`
    Avatar   string `gorm:"column:avatar" json:"avatar"`
    Status   int8   `gorm:"column:status" json:"status"`

    BaseModel
}

// TableName 设置表名。
func (Account) TableName() string {
    return "accounts"
}
```

#### 步骤 2：生成查询器代码

1. 在 `generate/gormgen/db/default.go` 中注册模型：

```go
// 添加新模型
var accountModel = model.Account{}

// 在 NewGeneratorDefaultDb 函数中添加
g.ApplyBasic(
    testModel,
    accountModel,  // 新增
)

// 如果需要自定义查询方法，使用 ApplyInterface（可选）
// 注意: 需先定义自定义接口，可参考 GORM Gen 文档
// g.ApplyInterface(func(AccountMethod) {}, accountModel)
```

2. 执行生成命令：

```bash
make gormgen
```

#### 步骤 3：创建数据仓库

在 `internal/data/repo/` 目录创建仓库实现：

```go
// Package repo 提供数据仓库实现。
package repo

import (
    "context"
    "ult/internal/data"
    "ult/internal/data/dbquery"
    "ult/internal/data/model"
)

// AccountRepo 接口验证。
var _ AccountRepo = (*accountRepo)(nil)

// AccountRepo 账户数据仓库接口。
type AccountRepo interface {
    GetByID(ctx context.Context, id int) (*model.Account, error)
    FindByUserName(ctx context.Context, username string) (*model.Account, error)
    Create(ctx context.Context, account *model.Account) error
}

// accountRepo 账户数据仓库实现。
type accountRepo struct {
    data data.Data
}

// NewAccountRepo 创建账户数据仓库实例。
func NewAccountRepo(data data.Data) AccountRepo {
    return &accountRepo{data: data}
}

// GetByID 根据 ID 获取账户。
func (r *accountRepo) GetByID(ctx context.Context, id int) (*model.Account, error) {
    return r.query(ctx).Where(dbquery.Account.ID.Eq(id)).First()
}

// FindByUserName 根据用户名查找账户。
func (r *accountRepo) FindByUserName(ctx context.Context, username string) (*model.Account, error) {
    return r.query(ctx).FindByUserName(username)
}

// query 获取带上下文的查询器。
func (r *accountRepo) query(ctx context.Context) dbquery.IAccountDo {
    return dbquery.Use(r.data.WithContext(ctx).GormDB()).Account.WithContext(ctx)
}
```

#### 步骤 4：创建服务层

在 `internal/service/` 目录创建 `account.go`：

```go
// Package service 提供业务逻辑层实现。
package service

import (
    "context"
    "ult/internal/data"
    "ult/internal/data/repo"
    "ult/errcode"
)

// AccountService 账户服务。
type AccountService struct {
    dataRepo *data.DataRepo
    repo     repo.AccountRepo
    tools    *app.Tools
}

// NewAccountService 创建账户服务实例。
func NewAccountService(dataRepo *data.DataRepo, repo repo.AccountRepo, tools *app.Tools) *AccountService {
    return &AccountService{dataRepo: dataRepo, repo: repo, tools: tools}
}

// GetByID 根据 ID 获取账户。
func (s *AccountService) GetByID(ctx context.Context, id int) (*AccountResponse, errcode.BusinessError) {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, errcode.ErrDataSelectError.WithDesc(err.Error())
    }
    if account == nil {
        return nil, errcode.ErrDataNotExistError
    }
    return &AccountResponse{
        ID:       account.ID,
        UserName: account.UserName,
        Avatar:   account.Avatar,
        Status:   account.Status,
    }, nil
}
```

#### 步骤 5：创建 API 处理器

在 `internal/api/` 目录创建 `account.go`：

```go
// Package api 提供 API 处理层实现。
package api

import (
    "ult/internal/app"
    "ult/internal/service"
    "ult/pkg/http"
)

// AccountInterface 接口验证。
var _ AccountInterface = (*AccountHandler)(nil)

// AccountInterface 账户 API 接口。
type AccountInterface interface {
    GetByID() http.HandlerFunc
    Create() http.HandlerFunc
}

// AccountHandler 账户 API 处理器。
type AccountHandler struct {
    service *service.AccountService
    tools   *app.Tools
}

// NewAccountHandler 创建账户 API 处理器实例。
func NewAccountHandler(service *service.AccountService, tools *app.Tools) AccountInterface {
    return &AccountHandler{
        service: service,
        tools:   tools,
    }
}

// GetByID 根据 ID 获取账户。
func (h *AccountHandler) GetByID() http.HandlerFunc {
    return func(ctx http.Context) {
        var req = new(GetByIDRequest)
        if err := ctx.Validator(req); err != nil {
            return  // 校验失败，Validator 已自动设置错误响应
        }

        account, err := h.service.GetByID(ctx.RequestContext(), req.ID)
        if err != nil {
            ctx.WithAbortError(err)
            return
        }

        ctx.WithPayload(account)
    }
}
```

#### 步骤 6：注册路由

在 `internal/router/router.go` 中添加：

```go
// httpRouter 结构体添加 Account 字段
type httpRouter struct {
    g      http.RouterGroup
    handle struct {
        Heartbeat api.HeartbeatInterface
        Account   api.AccountInterface
    }
}

// NewHTTPRouter 添加 Account 实例化
func NewHTTPRouter(heartbeat api.HeartbeatInterface, account api.AccountInterface) HTTPRouter {
    return func(hs *http.HTTPServer) {
        var r = &httpRouter{
            g: hs.CreateRouterGroup(),
            handle: struct {
                Heartbeat api.HeartbeatInterface
                Account   api.AccountInterface
            }{
                Heartbeat: heartbeat,
                Account:   account,
            },
        }

        r.heartbeat(r.g.Group("/heartbeat"))
        r.account(r.g.Group("/account"))
    }
}

// 新增路由注册方法
func (r *httpRouter) account(group http.RouterGroup) {
    group.GET("/:id", r.handle.Account.GetByID())
    group.POST("/", r.handle.Account.Create())
}
```

#### 步骤 7：注册 Wire ProviderSet

在对应的 `wire.go` 文件中添加 ProviderSet：

```go
// internal/api/wire.go
var ProviderSet = wire.NewSet(NewAccountHandler)

// internal/service/wire.go
var ProviderSet = wire.NewSet(NewAccountService)

// internal/data/repo/wire.go
var ProviderSet = wire.NewSet(NewAccountRepo)
```

#### 步骤 8：重新生成 Wire 代码

```bash
make wire
```

---

## 中间件管理

### 框架内置中间件

框架内置三个核心中间件，按优先级自动排序：

| 中间件 | 名称 | 优先级 | 依赖 | 说明 |
| --- | --- | --- | --- | --- |
| Recovery | `"recovery"` | PriorityHighest (0) | 无 | 异常恢复、告警通知 |
| Request | `"request"` | PriorityHigh (1) | Recovery | Context 初始化、验证器设置、响应处理 |
| CORS | `"cors"` | PriorityHigh (1) | 无 | 跨域处理 |

### 中间件优先级

框架使用 `go-utils/v2/middleware` 包的优先级定义，通过类型别名引入：

```go
// pkg/http/middleware/handler.go
type Priority = utilsMiddleware.Priority

// 优先级值（来自 go-utils/v2/middleware）
const (
    PriorityHighest Priority = iota  // 最高优先级（异常恢复）
    PriorityHigh                     // 高优先级（Request、CORS）
    PriorityNormal                   // 正常优先级（日志、验证）
    PriorityLow                      // 低优先级（权限检查、限流）
)
```

### Recovery 中间件

捕获 panic 异常，记录堆栈日志，支持邮件告警通知：

```go
// 配置
type RecoveryConfig struct {
    Enabled     bool                    // 是否启用
    AlertNotify proposal.NotifyHandler  // 告警通知处理器
    Config      *config.Config          // 应用配置
    PrintStack  bool                    // 是否打印堆栈
}

// 便捷构造
recovery := middleware.NewDefaultRecovery(cfg, logger, alertNotify)
```

Recovery 中间件在捕获 panic 时：
1. 记录 panic 信息和堆栈日志
2. 创建 BusinessError（`errcode.New(errcode.ServerError).WithStackError(goerror.New("got panic"))`）
3. 设置 HTTP 500 状态码并中止请求
4. 如果配置了 AlertNotify，发送告警通知

### Request 中间件

Request 中间件是核心中间件，负责 Context 初始化和响应处理：

```go
// 配置
type RequestConfig struct {
    Enabled            bool               // 是否启用
    Validator          validator.Validator // 数据验证器
    ContextInitializer ContextInitializer // Context 初始化函数
    Response           Response           // 响应处理函数
}

// 函数类型
type ContextInitializer func(ctx *gin.Context) (interface{}, error)
type Response func(reqTime time.Time, ctx *gin.Context)

// 便捷构造
request := middleware.NewDefaultRequest(validatorInst, contextInitializer, responseHandler)
```

Request 中间件功能：
1. **404 拦截**: 检查 `ctx.Writer.Status() == 404`，跳过不匹配的路由
2. **Context 初始化**: 通过 `ContextInitializer` 创建 Context（使用对象池复用）
3. **验证器设置**: 将 Validator 存入 gin.Context，供 `ctx.Validator()` 使用
4. **响应处理**: 通过 defer 调用 Response 函数处理成功/错误响应
5. **依赖 Recovery**: 确保 Recovery 能捕获 Request 中的 panic

服务器自动创建 Request 中间件（`srv.CreateRequest()`），无需手动注册。

### CORS 中间件

```go
// 配置
type CORSConfig struct {
    Enabled            bool
    AllowedOrigins     []string
    AllowedMethods     []string
    AllowedHeaders     []string
    AllowCredentials   bool
    OptionsPassthrough bool
}

// 便捷构造
corsConfig := middleware.DefaultCORSConfig()
```

### 中间件依赖管理

框架支持中间件依赖管理，自动验证依赖关系：

```go
// Middleware 接口
type Middleware interface {
    utilsMiddleware.Middleware      // 继承基础接口（Name、Priority、Enabled）
    Handler() utilsMiddleware.Handler // 返回中间件处理函数
    Dependencies() []string          // 返回依赖的中间件名称列表
}

// 示例：Request 中间件依赖 Recovery
func (r *Request) Dependencies() []string {
    return []string{types.RecoveryMiddlewareName} // 依赖 Recovery 中间件
}
```

**关键特性**：
- ✅ 自动验证中间件依赖是否满足（**不满足时 panic**，确保编译期发现错误）
- ✅ 按优先级和依赖关系自动排序
- ✅ 避免运行时错误，提高系统稳定性

### Middleware Manager

中间件管理器负责注册、排序和构建中间件链：

```go
manager := middleware.NewManager()
manager.Use(recoveryMiddleware)        // 注册 Recovery 中间件
manager.Use(corsMiddleware)            // 注册 CORS 中间件
manager.UseFunc("auth", PriorityLow, authHandler) // 注册函数式中间件

handlers := manager.Build() // 验证依赖 + 构建中间件链
```

### 创建自定义中间件

```go
// 定义中间件结构体
type AuthMiddleware struct {
    config *AuthConfig
}

// 实现 Middleware 接口
func (m *AuthMiddleware) Name() string {
    return "auth"
}

func (m *AuthMiddleware) Priority() utilsMiddleware.Priority {
    return utilsMiddleware.PriorityLow
}

func (m *AuthMiddleware) Enabled() bool {
    return true
}

func (m *AuthMiddleware) Handler() utilsMiddleware.Handler {
    return func(ctx *gin.Context) {
        // 中间件逻辑
        token := ctx.GetHeader("Authorization")
        if token == "" {
            ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        ctx.Next()
    }
}

func (m *AuthMiddleware) Dependencies() []string {
    return []string{} // 无依赖
}

// 注册中间件
srv.UseMiddleware(&AuthMiddleware{config: authConfig})
```

### 使用函数式中间件

```go
// 无依赖的函数式中间件
srv.UseMiddlewareFunc("auth", utilsMiddleware.PriorityLow, func(ctx *gin.Context) {
    token := ctx.GetHeader("Authorization")
    if token == "" {
        ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
        return
    }
    ctx.Next()
})

// 带依赖的函数式中间件
middleware.NewMiddlewareFuncWithDependencies(
    "auth",
    utilsMiddleware.PriorityLow,
    authHandler,
    []string{types.RecoveryMiddlewareName}, // 依赖 Recovery
)
```

---

## 错误码管理

### 错误码分类

| 范围 | 分类 | 说明 |
| --- | --- | --- |
| 100xxx | 服务端错误 | 内部服务器错误 |
| 200xxx | 客户端错误 | 参数错误、认证错误、业务错误、数据操作错误等 |

### 错误码定义

在 `errcode/code.go` 中定义的全部错误码：

| 常量名 | 错误码 | 分类 | 中文消息 |
| --- | --- | --- | --- |
| `ServerError` | 100001 | 服务端 | 内部服务器错误 |
| `AuthorizationError` | 200001 | 客户端 | 签名信息错误 |
| `ParamBindError` | 200002 | 客户端 | 参数信息错误 |
| `RequestError` | 200003 | 客户端 | 请求错误 |
| `ParamValidateError` | 200004 | 客户端 | 参数校验错误 |
| `UnknownError` | 200005 | 客户端 | 未知错误 |
| `DataNotExistError` | 200006 | 客户端 | 数据不存在 |
| `DataExistError` | 200007 | 客户端 | 数据已存在 |
| `RequestNotFoundError` | 200008 | 客户端 | 不存在的请求 |
| `DataDeleteError` | 200009 | 客户端 | 数据删除错误 |
| `ResourceNotExistError` | 200010 | 客户端 | 资源不存在 |
| `DataSelectError` | 200011 | 客户端 | 数据查询失败 |
| `DataCreateError` | 200012 | 客户端 | 数据创建失败 |
| `DataUpdateError` | 200013 | 客户端 | 数据更新失败 |

### 预创建错误变量

`errcode/vars.go` 为每个错误码预创建了 `BusinessError` 实例，推荐直接使用：

```go
var (
    ErrServerError          = New(ServerError)
    ErrAuthorizationError   = New(AuthorizationError)
    ErrParamBindError       = New(ParamBindError)
    ErrRequestError         = New(RequestError)
    ErrParamValidateError   = New(ParamValidateError)
    ErrUnknownError         = New(UnknownError)
    ErrDataNotExistError    = New(DataNotExistError)
    ErrDataExistError       = New(DataExistError)
    ErrRequestNotFoundError = New(RequestNotFoundError)
    ErrDataDeleteError      = New(DataDeleteError)
    ErrDataSelectError      = New(DataSelectError)
    ErrDataCreateError      = New(DataCreateError)
    ErrDataUpdateError      = New(DataUpdateError)
)
```

### 使用错误码

```go
// 返回预定义错误（使用 Err* 变量，推荐）
return errcode.ErrDataNotExistError

// 使用 errcode.New() 从错误码常量创建
return errcode.New(errcode.DataNotExistError)

// 添加错误描述
return errcode.ErrDataSelectError.WithDesc("数据库查询失败")

// 添加堆栈信息（注意：同时会设置描述为 err.Error()）
return errcode.ErrServerError.WithStackError(err)

// 设置告警标记
return errcode.ErrServerError.Alert()
```

### HTTP 状态码映射

| 错误码 | HTTP 状态码 |
| --- | --- |
| ServerError (100001) | 500 |
| AuthorizationError (200001) | 401 |
| ParamValidateError (200004) | 422 |
| 其他 200xxx 错误码 | 400 |

---

## 数据校验

### 校验器使用

```go
// 定义请求结构体
type CreateAccountRequest struct {
    UserName string `form:"username" label:"用户名" validate:"required,min=3,max=20"`
    Password string `form:"password" label:"密码" validate:"required,min=6,max=32"`
    Avatar   string `form:"avatar" label:"头像" validate:"omitempty,url"`
}

// 在 API 层调用校验器
// Validator() 返回 error 类型：nil 表示成功，非 nil 表示失败
var req = new(CreateAccountRequest)
if err := ctx.Validator(req); err != nil {
    return  // 校验失败，Validator 已自动设置错误响应并 Abort
}
```

> **注意**: `ctx.Validator()` 返回 `error` 类型（不是 bool），使用 `err != nil` 检查。校验失败时，Validator 会自动调用 `ctx.WithAbortError()` 设置错误响应。

### Validator 方法行为

`Validator()` 方法执行两步操作：
1. **绑定参数**: 调用 `ShouldBindForm(req)` 将表单数据绑定到请求结构体
2. **校验数据**: 使用配置的 Validator 实例调用 `Validate(req)` 进行结构化校验

绑定失败返回 `ParamBindError` 类型错误，校验失败返回 `ParamValidateError` 类型错误，均会自动设置到 Context 中。

### 常用校验规则

| 规则 | 说明 | 示例 |
| --- | --- | --- |
| `required` | 必填字段 | `validate:"required"` |
| `min` | 最小长度/值 | `validate:"min=3"` |
| `max` | 最大长度/值 | `validate:"max=20"` |
| `email` | 箱格式 | `validate:"email"` |
| `url` | URL 格式 | `validate:"url"` |
| `numeric` | 数值类型 | `validate:"numeric"` |
| `omitempty` | 可选字段 | `validate:"omitempty,email"` |

更多校验规则请参考：[validator 文档](https://github.com/go-playground/validator)

---

## 性能优化

### Header() 方法性能优化

Header() 方法返回请求头的只读引用，性能最优：

```go
// Header() 返回只读引用（性能最优）
headers := ctx.Header()  // 无内存分配，直接返回引用

// CloneHeaders() 返回完整副本（需要修改时使用）
headers := ctx.CloneHeaders()  // 创建完整副本，可以安全修改

// GetHeader() 获取单个请求头（性能最优）
authHeader := ctx.GetHeader("Authorization")  // 单次查找
```

**性能对比**：
- Header()：0.22 ns/op，0 次内存分配，性能提升 **99.96%**
- CloneHeaders()：543.6 ns/op，22 次内存分配（按需使用）
- GetHeader()：2.6 ns/op，0 次内存分配（单次查找）

### RequestContext() 方法性能优化

RequestContext() 方法使用 sync.Once 缓存，多次调用只创建一次：

```go
// 首次调用：创建并缓存 RequestContext
reqCtx1 := ctx.RequestContext()  // 230.5 ns/op，6 次内存分配

// 后续调用：直接返回缓存（性能提升 592倍）
reqCtx2 := ctx.RequestContext()  // 0.39 ns/op，0 次内存分配
```

**性能提升**：
- 首次调用：230.5 ns/op，6 次内存分配
- 后续调用：0.39 ns/op，0 次内存分配，性能提升 **99.83%**

### RawData() 方法内存安全

RawData() 方法返回副本，避免内存泄漏和数据污染：

```go
// 返回副本，可以安全修改
rawData := ctx.RawData()  // 返回副本，不影响原始数据
rawData[0] = 'X'  // 修改副本，不影响原始请求体
```

**内存安全**：
- ✅ 返回副本，避免外部修改影响原始数据
- ✅ 防止 context 对象被放回池中后的数据污染
- ✅ 提高内存安全性

---

## 并发安全性

### TraceID() 方法并发安全

TraceID() 方法使用 sync.Once 确保只生成一次：

```go
// 使用 sync.Once 确保只生成一次 TraceID
func (c *context) TraceID() string {
    c.traceIDOnce.Do(func() {
        // 先检查上下文中是否已存在
        if traceId, ok := c.ctx.Get(pkgtypes.TraceIdName); ok {
            if tid, ok := traceId.(string); ok && len(tid) > 0 {
                return
            }
        }

        // 检查请求头
        var headerTraceId = c.GetHeader(pkgtypes.TraceIdName)
        if len(headerTraceId) <= 0 {
            headerTraceId = uuid.New().String()
        }

        c.ctx.Set(pkgtypes.TraceIdName, headerTraceId)
    })

    // 从上下文中获取 TraceID
    traceId, ok := c.ctx.Get(pkgtypes.TraceIdName)
    // ...
}
```

**并发安全**：
- ✅ 使用 sync.Once 确保只生成一次 UUID
- ✅ 100 个 goroutine 同时调用 TraceID()，都获取到相同的值
- ✅ 避免竞态条件，提高并发安全性

### RequestContext() 方法并发安全

RequestContext() 方法使用 sync.Once 缓存：

```go
// 使用 sync.Once 缓存 RequestContext
func (c *context) RequestContext() stdCtx.Context {
    c.reqContextOnce.Do(func() {
        reqContext := new(pkgtypes.RequestContext)
        reqContext.WithTraceID(c.TraceID())
        c.reqContext = pkgtypes.NewRequestContext(c.ctx.Request.Context(), reqContext)
    })
    return c.reqContext
}
```

**并发安全**：
- ✅ 使用 sync.Once 确保只创建一次 RequestContext
- ✅ 多次调用性能提升 592倍
- ✅ 减少 GC 压力，提高系统性能

### Context 对象池

Context 使用 sync.Pool 复用，减少内存分配：

```go
// contextPool 是用于复用上下文对象的 sync.Pool
var contextPool = &sync.Pool{
    New: func() interface{} {
        return new(context)
    },
}

// newContext 从池中创建或获取上下文对象
func newContext(ctx *gin.Context) (Context, error) {
    context := contextPool.Get().(*context)
    context.ctx = ctx
    if err := context.init(); err != nil {
        return nil, err
    }
    return context, nil
}

// recoveryContext 使用后将上下文归还到池中
func recoveryContext(ctx Context) {
    c, ok := ctx.(*context)
    if !ok {
        return
    }
    c.reset()  // 清空所有字段，包括重置 sync.Once
    contextPool.Put(c)
}
```

**对象池优势**：
- ✅ 减少内存分配，提高性能
- ✅ 减少 GC 压力，提高系统稳定性
- ✅ 对象复用，降低资源消耗
- ✅ reset() 重置 sync.Once 确保池中对象可安全复用

---

## 配置管理

### 配置文件结构

配置文件 `.env.yml` 采用 YAML 格式（注意：字段名与代码中的 YAML struct tag 对应）：

```yaml
environment: dev  # 运行环境: dev/prod

app:
  id: "ult-framework"
  name: "ULT Web API Framework"
  version: "1.0.0"

datetime:
  location: Asia/Shanghai
  cst_layout: 2006-01-02 15:04:05

language:
  local: zh-cn  # 语言设置: zh-cn 或 en-us

validator:
  locale: zh     # 校验器语言
  tagname: label # 字段标签名

logger:           # 注意: 配置键是 "logger"，不是 "log"
  max_size: 128      # 日志文件最大大小(MB)
  max_backups: 5     # 最大备份文件数
  max_age: 7         # 最大保留天数
  local_time: true   # 使用本地时间
  compress: true     # 是否压缩

server:
  http:
    network: "tcp"
    host: "127.0.0.1"
    port: 10001
    cors:
      domains: "all"  # CORS 允许域名: all 或逗号分隔列表

db:
  default:
    driver: "mysql"
    host: "127.0.0.1"
    port: 3306
    username: "root"
    password: "password"
    db_name: "ult_framework"   # 注意: 字段名是 "db_name"，不是 "database"
    charset: "utf8mb4"
    prefix: "api_"
    max_open_conn: 100
    max_idle_conn: 10
    max_life_time: 0           # 最大连接生命周期(秒)
    max_retries: 3
    retry_delay: 2
    parse_time: "true"         # 注意: 类型为 string，值为 "true"
    loc: "Local"

redis:
  default:
    network: "tcp"
    addr: "127.0.0.1"
    port: 6379
    username: ""               # Redis 6.0+ ACL 用户名
    password: ""
    db: 0
    min_idle_conns: 10
    max_retries: 3
    retry_delay: 2
    min_retry_backoff: 0       # 最小重试退避时间(毫秒)
    max_retry_backoff: 0       # 最大重试退避时间(毫秒)
    dial_timeout: 0            # 连接超时(毫秒)
    read_timeout: 0            # 读超时(毫秒)
    write_timeout: 0           # 写超时(毫秒)
    pool_size: 0               # 连接池大小
    max_conn_age: 0            # 最大连接存活时间(毫秒)
    pool_timeout: 0            # 连接池等待超时(毫秒)
    idle_timeout: 0            # 空闲连接超时(毫秒)

jwt:
  app: "ult.service"
  key: "1203822711"
  secret: "Fu83AfHC839F0rTn22V23c"

notify:
  recover:
    email:
      host: "smtp.qq.com"
      port: 465
      user: "xxxxxx@qq.com"
      pass: "123456"
      to: "xxxxxx@qq.com"
```

### 多数据库/Redis 配置

支持配置多个数据库或 Redis 连接：

```yaml
db:
  default:     # 默认连接
    driver: "mysql"
    host: "127.0.0.1"
    ...
  readonly:    # 只读连接
    driver: "mysql"
    host: "127.0.0.1"
    ...

redis:
  default:     # 默认连接
    addr: "127.0.0.1"
    ...
  cache:       # 缓存连接
    addr: "127.0.0.1"
    ...
```

通过 `DataRepo` 接口按名称访问不同连接：

```go
// 获取指定名称的数据库连接
dbConn := dataRepo.DB("default")
dbConn := dataRepo.DB("readonly")

// 获取指定名称的 Redis 连接
redisConn := dataRepo.Redis("default")
redisConn := dataRepo.Redis("cache")
```

---

## 常见问题

### Q: 如何添加新的数据库连接？

A: 在配置文件 `db` 节点下添加新连接配置，框架会自动加载。通过 `DataRepo.DB("连接名")` 访问。

### Q: 如何添加告警通知？

A: Recovery 中间件支持 `proposal.NotifyHandler` 类型的告警处理器。可使用 `pkg/notify/recover/email` 提供的邮件告警功能：

```go
alertNotify := email.NotifyHandler(ctx, config.Notify, logger)
recovery := middleware.NewDefaultRecovery(cfg, logger, alertNotify)
```

### Q: 如何自定义错误码？

A: 在 `errcode/code.go` 定义错误码常量，在 `errcode/zh-cn.go` 定义中文消息，在 `errcode/en-us.go` 定义英文消息，在 `errcode/errhttp.go` 定义 HTTP 状态码映射，在 `errcode/vars.go` 添加预创建变量。

### Q: Wire 生成代码失败怎么办？

A: 确保 `wire.go` 文件头部有 `//go:build wireinject` 和 `// +build wireinject` 注释，然后执行 `make wire`。注意实际使用 `panic(wire.Build(...))` 模式。

### Q: 如何添加自定义中间件？

A: 实现 `middleware.Middleware` 接口（Name、Priority、Enabled、Handler、Dependencies），使用 `srv.UseMiddleware()` 方法注册，或使用 `srv.UseMiddlewareFunc()` 函数式注册。带依赖的函数式中间件使用 `middleware.NewMiddlewareFuncWithDependencies()`。

### Q: 如何处理请求响应？

A: 使用 `ctx.WithPayload(data)` 设置成功响应，使用 `ctx.WithAbortError(err)` 设置 BusinessError 并中止请求。响应处理由 Request 中间件的 Response 函数自动完成。

### Q: Validator() 返回什么类型？

A: `Validator()` 返回 `error` 类型（不是 bool）。使用 `err := ctx.Validator(req); err != nil` 检查。校验失败时，Validator 会自动调用 `WithAbortError()` 设置错误响应。

---

## 版本信息

- Go 版本：1.23.0+
- Gin 版本：1.11.0
- GORM 版本：1.31.1
- Wire 版本：0.7.0
- Zap 版本：1.27.0

---

## 许可证

MIT License

---

## 贡献指南

欢迎提交 Issue 和 Pull Request，请遵循以下规范：

1. 代码风格遵循 Go 客方规范
2. 新功能请添加相应的测试用例
3. 错误码使用 `errcode.New()` 或 `errcode.Err*` 预创建变量，不要直接对 int 常量调用方法

---

## 联系方式

- 作者：raylin666
- GitHub：https://github.com/raylin666/go-ult-framework
