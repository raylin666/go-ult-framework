# ULT Web API 框架 (基于 GIN 框架)

本框架是基于 `GIN` 进行模块化设计的 API 框架，封装了常用的功能，使用简单，致力于进行快速的业务研发，同时增加了更多限制，约束项目组开发成员，规避混乱无序及自由随意的编码。<br />

提供了方便快捷的 `Makefile` 文件 (帮你快速的生成、构建、执行项目内容)。<br />

当你所需命令不存在时可添加到此文件中, 实现命令统一管理。<br />

这也大大的提高了开发者的开发效率, 让开发者更专注于业务代码。 <br />

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

### 编译执行文件 (需要有 .git 提交版本, 你也可以修改 `Makefile` 文件来取消这个限制)

> make build

编译成功后, 可以通过 `./bin/server` 命令运行服务。

<hr />

### 规范约束

> `api` 处理层尽量避免使用 `配置(config)`、`数据仓库(dataRepo)`，职责上它只需要做 `数据校验` 和 `数据响应`。
> `pkg` 通用封装包内逻辑不允许调用 `internal` 内部包代码, 实现代码逻辑隔离, 也避免调用外部代码导致耦合。

### 创建新模块

> 以 `heartbeat` 为例: 
1. 在 `internal/router`、`internal/api` 和 `internal/service` 模块分别复制 `heartbeat` 文件, 并依次重命名为新模块名称。
2. 修改 `internal/router/router.go` 文件, 在结构体 `httpRouter.handle` 里添加新模块接口映射；然后在 `NewHTTPRouter` 里的注册处理器添加实例化；最后新增路由注册, 例如: `r.heartbeat(r.g.Group("/heartbeat"))` 。
3. 此时新模块就创建好了, 运行项目就可以访问对应的路由～

### 数据库模块

