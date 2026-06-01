// Package config 提供配置管理功能。
// 从 YAML 配置文件加载应用配置，支持模块化配置结构。
package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"ult/config/autoload"

	"github.com/raylin666/go-utils/v2/config"
)

// envFileName 配置文件名称。
const (
	envFileName = ".env.yml"
)

// Config 应用配置结构体。
// 包含项目路径和所有模块配置。
type Config struct {
	*builder
	ProjectPath string // 项目根目录路径
}

// builder 配置构建器结构体。
// 定义所有配置模块的字段映射。
type builder struct {
	Environment string                    `yaml:"environment"` // 运行环境 (dev/prod)
	App         autoload.App              `yaml:"app"`         // 应用配置
	Logger      autoload.Logger           `yaml:"logger"`      // 日志配置
	Language    autoload.Language         `yaml:"language"`    // 语言配置
	Validator   autoload.Validator        `yaml:"validator"`   // 验证器配置
	Server      autoload.Server           `yaml:"server"`      // 服务器配置
	DB          map[string]autoload.DB    `yaml:"db"`          // 数据库配置（支持多连接）
	Redis       map[string]autoload.Redis `yaml:"redis"`       // Redis 配置（支持多连接）
	JWT         autoload.JWT              `yaml:"jwt"`         // JWT 认证配置
	Datetime    autoload.Datetime         `yaml:"datetime"`    // 日期时间配置
	Notify      autoload.Notify           `yaml:"notify"`      // 告警通知配置
}

// New 加载配置文件并创建配置实例。
//
// 返回:
//   - *Config: 配置实例
//   - error: 加载错误
func New() (*Config, error) {
	var conf = new(Config)
	conf.ProjectPath = getProjectPath()
	var envFile = fmt.Sprintf("%s/%s", conf.ProjectPath, envFileName)
	err := config.LoadYaml(envFile, &conf.builder)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// getProjectPath 获取项目根目录路径。
func getProjectPath() string {
	_, p, _, ok := runtime.Caller(1)
	if ok {
		p = path.Dir(p)
		var index int
		index = strings.LastIndex(p, string(os.PathSeparator))
		p = p[:index]
	}

	return p
}
