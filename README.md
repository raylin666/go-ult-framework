# ULT Web API 框架 (基于 GIN)

本框架是基于 `GIN` 进行模块化设计的 API 框架，封装了常用的功能，使用简单，致力于进行快速的业务研发，同时增加了更多限制，约束项目组开发成员，规避混乱无序及自由随意的编码。

提供了方便快捷的 `Makefile` 文件 (帮你快速的生成、构建、执行项目内容)。

当你所需命令不存在时可添加到此文件中, 实现命令统一管理。这也大大的提高了开发者的开发效率, 让开发者更专注于业务代码。

---

## 核心特性

- **分层架构设计**: 采用经典的分层架构 + 依赖注入设计模式，职责清晰
- **严格调用链**: api → service → data 单向调用，避免循环依赖
- **依赖注入**: 使用 Google Wire 实现依赖注入，降低组件耦合
- **中间件管理**: 支持优先级管理的中间件链，灵活配置
- **统一错误处理**: 完整的错误码管理系统，支持多语言和告警
- **多连接支持**: 支持多数据库和 Redis 连接配置
- **代码生成**: 集成 GORM Gen，自动生成查询器代码
- **优雅关闭**: 完整的服务器生命周期管理和优雅关闭机制
- **链路追踪**: 内置 TraceID 支持，便于分布式追踪
- **配置管理**: YAML 配置文件，支持多环境配置
- **Docker 支持**: 提供 Dockerfile 和 Docker Compose 配置

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
│  │   http   │    db    │  cache   │  logger  │middleware│repositories│ │
│  │(HTTP封装)│ (数据库) │ (缓存)   │  (日志)  │ (中间件) │  (仓库)   │  │
│  └──────────┴──────────┴──────────┴──────────┴──────────┴───────────┘  │
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
| 工具层 | `internal/app/` | 公共工具包（日志、JWT、环境等） |
| HTTP封装 | `pkg/http/` | Gin 框架封装，Context、Handler、Response、中间件管理 |
| 数据库 | `pkg/db/` | GORM 数据库连接封装，支持连接重试 |
| 缓存 | `pkg/cache/` | Redis 连接封装，支持连接重试 |
| 日志 | `pkg/logger/` | Zap 日志封装，支持文件轮转 |
| 配置 | `config/` | YAML 配置文件加载和解析 |
| 错误码 | `errcode/` | 统一错误码定义和管理，支持多语言 |
| 仓库 | `pkg/repositories/` | 数据仓库抽象层，管理多连接 |

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
//+build wireinject

func initApp(conf *config.Config, tools *app.Tools) (*pkgapp.App, func(), error) {
    wire.Build(
        config.ProviderSet,
        logger.ProviderSet,
        db.ProviderSet,
        redis.ProviderSet,
        server.ProviderSet,
        router.ProviderSet,
        api.ProviderSet,
        service.ProviderSet,
        data.ProviderSet,
        app.ProviderSet,
    )
    return nil, nil, nil
}
```

各模块通过 `ProviderSet` 注册依赖，Wire 自动生成依赖注入代码。

### 数据流设计

```
HTTP Request
    │
    ▼
┌─────────────────────────────────────┐
│         pkg/http/server             │
│    (Gin Engine + Middleware)        │
│    - CORS 处理                      │
│    - Recovery 异常恢复              │
│    - Request/Response 处理          │
│    - TraceID 生成                   │
│    - Validator 数据校验             │
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
│    - 数据校验                        │
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
    WithStackError(err error) BusinessError  // 堆栈追踪
    HTTPCode() int                           // HTTP状态码
    BusinessCode() int                       // 业务错误码
    Message() string                         // 错误消息
    Desc() string                            // 错误描述
    Alert() BusinessError                    // 告警标记
    IsAlert() bool                           // 是否需要告警
}
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
│   ├── error.go            # 业务错误接口
│   ├── code.go             # 错误码常量
│   ├── register.go         # 错误码注册表
│   ├── errhttp.go          # HTTP状态码映射
│   ├── vars.go             # 全局变量
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
│   │   ├── tools.go        # 公共工具实例
│   │   └── logo.go         # 启动信息打印
│   ├── data/               # 数据层
│   │   ├── data.go         # 数据仓库管理
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
│   │   └ dbquery/          # GORM Gen 查询器（自动生成）
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
│   │   ├── app.go          # 应用生命周期管理
│   │   └ server.go         # 服务器接口定义
│   ├── http/               # HTTP 封装
│   │   ├── server.go       # Gin 服务器封装
│   │   ├── context.go      # 请求上下文封装
│   │   ├── router.go       # 路由组封装
│   │   ├── response.go     # 响应处理
│   │   ├── option.go       # 服务器选项
│   │   └ middleware/       # 中间件管理
│   │       ├── middleware.go # 中间件接口
│   │       ├── chain.go    # 中间件链
│   │       ├── cors.go     # CORS 中间件
│   │       ├── recovery.go # Recovery 中间件
│   │       └ README.md     # 中间件文档
│   ├── db/                 # 数据库封装
│   │   ├── gorm.go         # GORM 连接封装
│   │   └ logger.go         # 数据库日志
│   ├── cache/              # 缓存封装
│   │   └ redis.go          # Redis 连接封装
│   ├── logger/             # 日志封装
│   │   └ logger.go         # Zap 日志封装
│   ├── repositories/       # 数据仓库封装
│   │   ├── repo.go         # 数据仓库接口
│   │   ├── db.go           # 数据库仓库接口
│   │   └ redis.go          # Redis 仓库接口
│   ├── types/              # 类型定义
│   │   ├── context.go      # 请求上下文类型
│   │   └ trace.go          # 链路追踪类型
│   └ notify/               # 告警通知
│   │   └ recover/          # 异常恢复通知
│   │       └ email/        # 邮件通知
│   │           ├── alert.go # 告警处理
│   │           └ email_template.go # 邮件模板
│   └ proposal/             # 告警提案
│   │   ├── alert.go        # 告警消息定义
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
    database: ult_framework
    max_open_conn: 100
    max_idle_conn: 10
    max_retries: 3
    retry_delay: 2

redis:
  default:
    host: 127.0.0.1
    port: 6379
    password: ""
    db: 0
    max_retries: 3
    retry_delay: 2
```

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

---

## Makefile 命令说明

| 命令 | 说明 |
| --- | --- |
| `make init` | 初始化项目，创建配置文件 |
| `make generate` | 下载依赖并生成 Wire 代码 |
| `make wire` | 生成依赖注入文件 |
| `make gormgen` | 生成 GORM Gen 查询器代码 |
| `make run` | 启动开发服务器 |
| `make build` | 编译生产执行文件 |
| `make clean` | 清理编译文件 |
| `make test` | 运行测试 |
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
2. **错误码定义**：使用 `errcode` 包预定义的错误码变量
3. **错误响应**：只能在 `api` 层处理错误响应

```go
// 正确示例
func (s *AccountService) GetByID(ctx context.Context, id int) (*model.Account, BusinessError) {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, errcode.DataSelectError.WithDesc(err.Error())
    }
    if account == nil {
        return nil, errcode.DataNotExistError
    }
    return account, nil
}
```

### 数据响应规范

| 方法 | 说明 | 使用场景 |
| --- | --- | --- |
| `ctx.WithPayload(data)` | 设置成功响应数据 | 业务处理成功 |
| `ctx.WithAbortError(err)` | 设置错误并中断请求 | 业务处理失败 |

**注意事项：**
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

#### 步骤 2：创建 DIY 查询方法

在 `internal/data/dbmethod/` 目录创建方法定义：

```go
// Package dbmethod 提供自定义查询方法定义。
package dbmethod

import "gorm.io/gen"

// Account 自定义查询方法接口。
type Account interface {
    // Where("`username`=@username")
    FindByUserName(username string) (gen.T, error)
    
    // Where("`status`=@status")
    FindByStatus(status int8) ([]gen.T, error)
}
```

#### 步骤 3：生成查询器代码

1. 在 `generate/gormgen/db/default.go` 中注册模型：

```go
var accountModel = model.Account{}

g.ApplyBasic(accountModel)

g.ApplyInterface(func(method dbmethod.Account) {}, accountModel)
```

2. 执行生成命令：

```bash
make gormgen
```

#### 步骤 4：创建数据仓库

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

#### 步骤 5：创建服务层

在 `internal/service/` 目录创建 `account.go`：

```go
// Package service 提供业务逻辑层实现。
package service

import (
    "context"
    "ult/internal/data/repo"
    "ult/errcode"
)

// AccountService 账户服务。
type AccountService struct {
    repo repo.AccountRepo
}

// NewAccountService 创建账户服务实例。
func NewAccountService(repo repo.AccountRepo) *AccountService {
    return &AccountService{repo: repo}
}

// GetByID 根据 ID 获取账户。
func (s *AccountService) GetByID(ctx context.Context, id int) (*AccountResponse, errcode.BusinessError) {
    account, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, errcode.DataSelectError.WithDesc(err.Error())
    }
    if account == nil {
        return nil, errcode.DataNotExistError
    }
    return &AccountResponse{
        ID:       account.ID,
        UserName: account.UserName,
        Avatar:   account.Avatar,
        Status:   account.Status,
    }, nil
}
```

#### 步骤 6：创建 API 处理器

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
        if isErr := ctx.Validator(req); isErr {
            return
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

#### 步骤 7：注册路由

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

#### 步骤 8：注册 Wire ProviderSet

在对应的 `wire.go` 文件中添加 ProviderSet：

```go
// internal/api/wire.go
var ProviderSet = wire.NewSet(NewAccountHandler)

// internal/service/wire.go
var ProviderSet = wire.NewSet(NewAccountService)

// internal/data/repo/wire.go
var ProviderSet = wire.NewSet(NewAccountRepo)
```

#### 步骤 9：重新生成 Wire 代码

```bash
make wire
```

---

## 中间件管理

### 中间件优先级

框架支持中间件优先级管理，按优先级顺序执行：

```go
const (
    PriorityHighest Priority = iota  // 最高优先级（异常恢复、链路追踪）
    PriorityHigh                     // 高优先级（CORS、安全检查）
    PriorityNormal                   // 正常优先级（日志、验证）
    PriorityLow                      // 低优先级（权限检查、限流）
)
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

func (m *AuthMiddleware) Priority() Priority {
    return PriorityLow
}

func (m *AuthMiddleware) Handler() HandlerFunc {
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

// 注册中间件
srv.UseMiddleware(&AuthMiddleware{config: authConfig})
```

### 使用函数式中间件

```go
srv.UseMiddlewareFunc("auth", PriorityLow, func(ctx *gin.Context) {
    token := ctx.GetHeader("Authorization")
    if token == "" {
        ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
        return
    }
    ctx.Next()
})
```

---

## 错误码管理

### 错误码分类

| 范围 | 分类 | 说明 |
| --- | --- | --- |
| 100xxx | 服务端错误 | 内部服务器错误、数据库错误等 |
| 200xxx | 客户端错误 | 参数错误、认证错误、业务错误等 |

### 定义错误码

在 `errcode/code.go` 中定义：

```go
const (
    // 服务端错误 (100xxx)
    ServerError       = 100001  // 内部服务器错误
    DataSelectError   = 100002  // 数据查询错误
    DataInsertError   = 100003  // 数据插入错误
    DataUpdateError   = 100004  // 数据更新错误
    DataDeleteError   = 100005  // 数据删除错误
    
    // 客户端错误 (200xxx)
    AuthorizationError = 200001  // 签名信息错误
    ParamBindError     = 200002  // 参数绑定错误
    ParamValidateError = 200004  // 参数校验错误
    DataNotExistError  = 200006  // 数据不存在
)
```

### 使用错误码

```go
// 返回预定义错误
return errcode.DataNotExistError

// 添加错误描述
return errcode.DataSelectError.WithDesc("数据库查询失败")

// 添加堆栈信息
return errcode.ServerError.WithStackError(err)

// 设置告警标记
return errcode.ServerError.Alert()
```

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
var req = new(CreateAccountRequest)
if isErr := ctx.Validator(req); isErr {
    return  // 校验失败，自动返回错误响应
}
```

### 常用校验规则

| 规则 | 说明 | 示例 |
| --- | --- | --- |
| `required` | 必填字段 | `validate:"required"` |
| `min` | 最小长度/值 | `validate:"min=3"` |
| `max` | 最大长度/值 | `validate:"max=20"` |
| `email` | 邮箱格式 | `validate:"email"` |
| `url` | URL 格式 | `validate:"url"` |
| `numeric` | 数值类型 | `validate:"numeric"` |
| `omitempty` | 可选字段 | `validate:"omitempty,email"` |

更多校验规则请参考：[validator 文档](https://github.com/go-playground/validator)

---

## 配置管理

### 配置文件结构

配置文件 `.env.yml` 采用 YAML 格式：

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
  local: zh-cn  # 语言设置

validator:
  locale: zh     # 校验器语言
  tagname: label # 字段标签名

logger:
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
    database: "ult_framework"
    charset: "utf8mb4"
    prefix: "api_"
    max_open_conn: 100
    max_idle_conn: 10
    max_retries: 3
    retry_delay: 2
    parse_time: true
    loc: Local

redis:
  default:
    network: "tcp"
    addr: "127.0.0.1"
    port: 6379
    password: ""
    db: 0
    min_idle_conns: 10
    max_retries: 3
    retry_delay: 2

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
    host: "127.0.0.1"
    ...
  cache:       # 缓存连接
    host: "127.0.0.1"
    ...
```

---

## 常见问题

### Q: 如何添加新的数据库连接？

A: 在配置文件 `db` 节点下添加新连接配置，然后在 `internal/data/data.go` 中注册。

### Q: 如何添加告警通知？

A: 使用 `pkg/http` 的 `EnableAlertNotify` 选项配置告警处理器，或配置邮件通知。

### Q: 如何自定义错误码？

A: 在 `errcode/code.go` 定义错误码常量，在 `errcode/zh-cn.go` 定义中文消息，在 `errcode/errhttp.go` 定义 HTTP 状态码映射。

### Q: Wire 生成代码失败怎么办？

A: 确保 `wire.go` 文件头部有 `//+build wireinject` 注释，然后执行 `make generate`。

### Q: 如何添加自定义中间件？

A: 实现 `Middleware` 接口，使用 `srv.UseMiddleware()` 方法注册，或使用 `srv.UseMiddlewareFunc()` 函数式注册。

### Q: 如何处理请求响应？

A: 使用 `ctx.WithPayload(data)` 设置成功响应，使用 `ctx.WithAbortError(err)` 设置错误响应。

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

1. 代码风格遵循 Go 官方规范
2. 新功能请添加相应的测试用例

---

## 联系方式

- 作者：raylin666
- GitHub：https://github.com/raylin666/go-ult-framework