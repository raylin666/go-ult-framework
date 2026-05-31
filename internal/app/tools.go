package app

import (
	"ult/pkg/logger"

	"github.com/raylin666/go-utils/v2/auth"
	"github.com/raylin666/go-utils/v2/server/system"
)

type Tools struct {
	logger      *logger.Logger
	datetime    *system.Datetime
	environment system.Environment
	jwt         auth.JWT
}

// NewTools 创建公共工具实例
func NewTools(
	logger *logger.Logger,
	datetime *system.Datetime,
	environment system.Environment,
	jwt auth.JWT) (tools *Tools) {
	tools = &Tools{logger, datetime, environment, jwt}
	return
}

// Logger 日志
func (tools *Tools) Logger() *logger.Logger { return tools.logger }

// Datetime 日期时间
func (tools *Tools) Datetime() *system.Datetime { return tools.datetime }

// Environment 开发环境
func (tools *Tools) Environment() system.Environment { return tools.environment }

// JWT 鉴权认证
func (tools *Tools) JWT() auth.JWT { return tools.jwt }
