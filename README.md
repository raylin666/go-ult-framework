# ULT Web API 框架 (基于 Gin 框架)

本框架是基于 `Gin` 进行模块化设计的 API 框架，封装了常用的功能，使用简单，致力于进行快速的业务研发，同时增加了更多限制，约束项目组开发成员，规避混乱无序及自由随意的编码。<br />

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
| qiniu | 上传文件 |
| dingTalk | 钉钉机器人 |
| gomail | 邮件发送 |
| wire | 依赖注入 |
| yaml.v3 | 配置文件解析 |
| RESTFUL API | API 返回值规范 |

### 下载仓库

> git clone git@github.com:raylin666/go-ult-framework.git

### 初始化

> make init

### 下载依赖

> make generate

### 启动服务

> make run

### 编译执行文件 (需要有 .git 提交版本, 你也可以修改 `Makefile` 文件来取消这个限制)

> make build

编译成功后, 可以通过 `./bin/server` 命令运行服务。

<hr />

具体开发使用文档可到 `docs` 目录查看哦 ～
