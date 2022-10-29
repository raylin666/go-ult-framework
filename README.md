# ULT Web 框架 (基于 Gin 框架)

本框架致力于 Web 项目, 如果你的是纯 RPC 或 gRPC 服务, 建议使用 `Kratos`、`go-micro` 等框架 ~
本框架提供了方便快捷的 `Makefile` 文件 (帮你快速的生成、构建、执行项目内容), 当你所需命令不存在时可添加到此文件中, 实现命令统一管理。这也大大的提高了开发者的开发效率, 让开发者更专注于业务代码。 本框架是基于B站(
bilibili)开源产品 `Kratos` 框架的基础模版改造而成, 脚本架为原始的 `kratos`, 意味着你可以直接使用 (前提是要初始化项目后)。

### 下载仓库

> git clone git@github.com:raylin666/go-ult-framework.git

### 初始化

> make init

### 下载依赖

> make generate

### 启动服务

> make run