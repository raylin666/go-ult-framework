GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
GOVERSION:=$(shell go env GOVERSION)
GIT_VERSION=$(shell git describe --tags --always)

.PHONY: init
# 初始化安装脚本
init:
	cp .env.example.yml .env.yml

.PHONY: generate
# 自动化生成编译前的类库代码
generate:
	go mod download && go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: wire
# 生成依赖注入文件
wire:
	cd ./cmd && $(GOPATH)/bin/wire

.PHONY: run
# 开发环境启动项目
run:
	cd cmd && go run ./...

.PHONY: build
# 编译构建项目 (需要有 .git 提交版本)
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(GIT_VERSION)" -o ./bin/server ./cmd

# 帮助命令
help:
	@echo ''
	@echo 'go version: ' $(GOVERSION)
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
