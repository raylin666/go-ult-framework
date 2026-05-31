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

const (
	envFileName = ".env.yml"
)

type Config struct {
	*builder
	ProjectPath string
}

type builder struct {
	Environment string                    `yaml:"environment"`
	App         autoload.App              `yaml:"app"`
	Logger      autoload.Logger           `yaml:"logger"`
	Language    autoload.Language         `yaml:"language"`
	Validator   autoload.Validator        `yaml:"validator"`
	Server      autoload.Server           `yaml:"server"`
	DB          map[string]autoload.DB    `yaml:"db"`
	Redis       map[string]autoload.Redis `yaml:"redis"`
	JWT         autoload.JWT              `yaml:"jwt"`
	Datetime    autoload.Datetime         `yaml:"datetime"`
	Notify      autoload.Notify           `yaml:"notify"`
}

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

// getProjectPath 获取项目根目录
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
