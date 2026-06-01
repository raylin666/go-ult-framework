// Package app 提供应用工具包。
// 包含日志、日期时间、环境配置、JWT 认证等公共工具实例。
package app

import (
	"ult/pkg/logger"

	"github.com/raylin666/go-utils/v2/auth"
	"github.com/raylin666/go-utils/v2/server/system"
)

// Tools 应用公共工具包结构体。
// 封装日志、日期时间、环境、JWT 等工具实例。
type Tools struct {
	logger      *logger.Logger    // 日志实例
	datetime    *system.Datetime  // 日期时间实例
	environment system.Environment // 环境配置
	jwt         auth.JWT          // JWT 认证实例
}

// NewTools 创建公共工具实例。
//
// 参数:
//   - logger: 日志实例
//   - datetime: 日期时间实例
//   - environment: 环境配置
//   - jwt: JWT 认证实例
//
// 返回:
//   - *Tools: 公共工具实例
func NewTools(
	logger *logger.Logger,
	datetime *system.Datetime,
	environment system.Environment,
	jwt auth.JWT) (tools *Tools) {
	tools = &Tools{logger, datetime, environment, jwt}
	return
}

// Logger 获取日志实例。
func (tools *Tools) Logger() *logger.Logger { return tools.logger }

// Datetime 获取日期时间实例。
func (tools *Tools) Datetime() *system.Datetime { return tools.datetime }

// Environment 获取环境配置。
func (tools *Tools) Environment() system.Environment { return tools.environment }

// JWT 获取 JWT 认证实例。
func (tools *Tools) JWT() auth.JWT { return tools.jwt }
