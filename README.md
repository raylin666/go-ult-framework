# ULT Web API 框架 (基于 GIN)

本框架是基于 `GIN` 进行模块化设计的 API 框架，封装了常用的功能，使用简单，致力于进行快速的业务研发，同时增加了更多限制，约束项目组开发成员，规避混乱无序及自由随意的编码。<br />

提供了方便快捷的 `Makefile` 文件 (帮你快速的生成、构建、执行项目内容)。<br />

当你所需命令不存在时可添加到此文件中, 实现命令统一管理。这也大大的提高了开发者的开发效率, 让开发者更专注于业务代码。 <br />

### 集成组件

| 名称 | 描述 | 
| --- | --- |
| cors | 接口跨域 |
| pprof | 性能剖析 |
| errno | 统一定义错误码 |
| zap | 日志收集 |
| gorm | 数据库组件 (支持 `gen` 和 `DIY` 生成文件) |
| go-redis | redis 组件 |
| JWT | 鉴权组件 |
| validator | 数据校验 |
| qiniu | 上传文件 |
| uuid | 唯一值生成 |
| dingTalk | 钉钉机器人 |
| gomail | 邮件发送 |
| wire | 依赖注入 |
| yaml.v3 | 配置文件解析 |
| RESTFUL API | API 返回值规范 |

### 目录介绍

| 目录 | 目录名称 | 目录描述 |
| --- | --- | --- |
| cmd | 项目启动 | 存放项目启动文件及依赖注入绑定 |
| config | 配置文件 |  |
| generate | 代码生成器 | 比如数据库查询器 |
| internal | 内部文件 | 存放项目业务开发文件 |
| pkg | 通用封装包 | 存放项目通用封装逻辑, 代码实现隔离项目内部业务 |
| static | 静态文件 | 比如图片、描述性文件、数据库SQL等 |
| bin | 运行文件 | |
| runtime | 临时/暂存 文件 | 比如日志文件 |

### 下载仓库

> git clone git@github.com:raylin666/go-ult-framework.git

### 初始化

> make init

### 下载依赖

> make generate

### 启动服务

> make run 

访问服务 `curl 127.0.0.1:10001/heartbeat` , 返回 `200` 状态码则表示成功。
```shell
{
    "trace_id": "b2401c9e-1f6f-4183-952a-b539ddabbb71",
    "data": "PONE"
}
```

同时也支持采用 `Dockerfile` 和 `docker-compose` 启动哦 ！

### 编译执行文件 (需要有 .git 提交版本, 你也可以修改 `Makefile` 文件来取消这个限制)

> make build

编译成功后, 可以通过 `./bin/server` 命令运行服务。

<hr />

### 规范约束

> `api` 处理层尽量避免使用 `配置(config)`、`数据仓库(dataRepo)`，职责上它只需要做 `数据校验` 和 `数据响应`。

> `pkg` 通用封装包内逻辑不允许调用 `internal` 内部包代码, 实现代码逻辑隔离, 也避免调用外部代码导致耦合。

> `data` 数据层主要处理业务数据仓库的实例, 数据库逻辑处理、缓存逻辑处理、RPC 远程调用处理等相关操作。

> 逻辑方法的异常情况统一返回 `errors.BusinessError`, 调用异常时统一使用 `internal/constant/errcode/errors.go` 内的变量。

> 调用关系链: (处理层) `api` -> (逻辑层) `service` -> (数据层) `data` , 逻辑代码只能下沉, 注意不要互调哦 ～

### 创建新模块

> 以 `account` 为例: 
1. 在 `internal/router`、`internal/api` 和 `internal/service` 模块分别复制 `heartbeat` 文件, 并依次重命名为新模块名称。
2. 修改 `internal/router/router.go` 文件, 在结构体 `httpRouter.handle` 里添加新模块接口映射；然后在 `NewHTTPRouter` 里的注册处理器添加实例化；最后新增路由注册, 例如: `r.account(r.g.Group("/account"))` 。
3. 此时新模块就创建好了, 运行项目就可以访问对应的路由～

### 数据库模块

> 例如创建个 `account` 模型:
1. 编写模型文件, 在 `internal/repositories/dbrepo/model` 目录创建 `account.go` 文件, 内容如下:

```go
package model

type Account struct {
	UserName string `gorm:"column:username" json:"username"` // 用户名称
	Password string `gorm:"column:password" json:"password"` // 用户密码(加密串)
	Avatar   string `gorm:"column:avatar" json:"avatar"`     // 用户头像
	Status   int8   `gorm:"column:status" json:"status"`     // 用户状态 0:冻结 1:正常 2:暂停

	BaseModel
}
```

2. 如果需要制定 DIY 查询, 可在 `internal/repositories/dbrepo/method` 目录创建 `account.go` 文件, 内容如下:

```go
package method

import (
	"gorm.io/gen"
)

type Account interface {
	// where("`username`=@username")
	FindByUserName(username string) (gen.T, error)
}
```

3. 接下来就是添加到代码生成器中了，很简单。到 `generate/gormgen/db/default.go` 文件中添加 3 行代码即可:
```shell
// 代码 g.UseDB(dbInterface.Get().DB()) 后添加模型定义:
var accountModel = model.Account{}

// 代码 g.ApplyBasic 内添加注册模型:
g.ApplyBasic(
    accountModel,
)
  
// 代码 g.ApplyInterface 内添加注册DIY:
g.ApplyInterface(func(method method.Account) {}, accountModel)
```

4. 生成数据库查询器代码, 执行 `make gormgen` 命令, 成功后会在 `internal/repositories/dbrepo/query` 目录内生成对应的查询器文件。

5. 在 `internal/repositories/dbrepo/query.go` 文件中增加如下代码, 方便默认查询调用:
```go
// NewDefaultDbQuery 创建默认数据库查询
func NewDefaultDbQuery(dbInterface db.Db) *query.Query {
	return query.Use(dbInterface.Get().DB())
}
```

### 数据层处理 

该层的设计目的是 <b>解藕数据与业务逻辑</b> 代码, 使层级更清晰, 服务逻辑层不再需要引入 `DataRepo` 来处理数据逻辑，只需要专心处理业务逻辑即可。当业务复杂、多人协作开发、功能模块多的项目强烈建议采用数据层来降低后期维护成本。

> 依照如上 `account` 模型为例:
1. 在服务逻辑层添加数据层调用代码, 打开 `internal/service/account.go` 文件, 修改内容如下：
```shell
// 在 import 之后添加数据层接口定义
type AccountRepo interface {
	ID(ctx context.Context, username string) int
}

// 在 AccountService 内添加数据层实例, 例如
type AccountService struct {
    logger *logger.Logger
	repo   AccountRepo
}

// 在 NewAccountService 方法参数添加数据层实例及设置入服务逻辑层结构体的 repo, 例如
func NewAccountService(logger *logger.Logger, repo AccountRepo) *AccountService {
	return &AccountService{
		logger: logger,
		repo:   repo,
	}
}
```

2. 编写数据处理文件, 在 `internal/data` 目录创建 `account.go` 文件, 内容如下:
```go
package data

import (
	"context"
	"ult/internal/constant/defined"
	"ult/internal/repositories/dbrepo"
	"ult/internal/service"
	"ult/pkg/global"
	"ult/pkg/logger"
)

type accountRepo struct {
	repo   global.DataRepo
	logger *logger.Logger
}

func NewAccountRepo(logger *logger.Logger, repo global.DataRepo) service.AccountRepo {
	return &accountRepo{
		logger: logger,
		repo:   repo,
	}
}

func (data *accountRepo) ID(ctx context.Context, username string) int {
	var q = dbrepo.NewDefaultDbQuery(data.repo.DB(defined.DB_CONNECTION_DEFAULT_NAME))
	res, err := q.Account.WithContext(ctx).FindByUserName(username)
	if err != nil {
		return 0
	}

	return res.ID
}
```

3. 数据层定义各仓库实例化, 在 `internal/data/data/go` 文件中新增 `Account` , 内容如下：
```shell
// DataRepo 结构体新增 `Account`
type DataRepo struct {
	Account service.AccountRepo
}

// NewDataRepo 方法实例化 `Account`
func NewDataRepo(logger *logger.Logger, repo global.DataRepo) *DataRepo {
	return &DataRepo{
		Account: NewAccountRepo(logger, repo),
	}
}
```

4. 服务逻辑层调用数据层, 在 `internal/router/router.go` 文件找到代码 `data.NewDataRepo(hs.Logger(), hs.DataRepo())`，修改为:
```go
var repo = data.NewDataRepo(hs.Logger(), hs.DataRepo())

// 在注册处理器里的 `NewAccountService` 调用中添加参数 `repo.Account` 即可
var r = &httpRouter{
    // 创建路由组
    g: hs.CreateRouterGroup(),
    // 注册处理器
    handle: struct {
        Account api.AccountInterface
    }{
        Account: api.NewAccountHandler(hs.Logger(), service.NewAccountService(hs.Logger(), repo.Account)),
    },
}
```

接下来就可以在 `service` 服务逻辑层直接调用数据层咯～

### 数据响应

> 主要分为 `正常响应` 和 `异常响应`。

1. 只能在 `api` 层处理, 否则将会导致项目规范凌乱, 无法管理。
2. `ctx.WithAbortError` 抛出异常后记得 `return` 断言, 否则是会往下执行的。
3. `ctx.WithPayload` 正常响应数据, 一般放在最后, 否则记得 `return` 断言, 否则也是会往下执行的。

### 数据校验

验证器组件文档: `https://github.com/go-playground/validator`
数据验证几乎不需要你做啥, 只要定义好参数结构体, 然后在 `api` 层调用验证器即可。例如：
```go
type Account struct {
    Id int `form:"id" label:"唯一值" validate:"numeric,min=1"`
}

var req = new(Account)
if isErr := ctx.Validator(req); isErr {
    return
}

// 如上的响应HTTP状态码为 422, 响应内容如下:
// {"trace_id":"eab08d57-bd3f-48e6-99de-9dbcd14ecc33","code":100005,"message":"参数校验错误","desc":"唯一值必须是一个有效的数值"}
```

### 异常处理/错误状态码

1. 为了统一规范, 基础错误状态码放置在 `pkg/code/code.go` 文件中, HTTP 状态码设置在 `pkg/code/errhttp.go` 文件中 (未设置的状态码 响应时默认是 `400`), 状态码提示分为中英文 `zh-cn` 和 `en-us`。
2. 业务状态码放置在 `internal/constant/errcode/errcode.go` 文件中, HTTP 状态码设置在 `internal/constant/errcode/errhttp.go` 文件中 (未设置的状态码 响应时默认是 `400`), 状态码提示也分为中英文 `zh-cn` 和 `en-us`。
3. 规范化的错误状态管理, 通过 `internal/constant/errcode/errors.go` 文件来定义错误, 抛出异常时只需要调用该变量即可实现, 避免了 普通错误码需要引用 `pkg` 包内的错误码, 而业务错误码则需要引用 `errcode` 包内的错误码, 导致调用极不统一的现象。
4. 逻辑方法的异常情况统一为 `errors.BusinessError` 对象, 例如: 
```go
func (h *AccountService) ID(ctx context.Context) errors.BusinessError {
    if h.repo.ID(ctx) <= 0 {
        return errcode.ErrorDataSelectError
    }   
    
    return nil 
}
```

