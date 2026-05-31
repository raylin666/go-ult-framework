// Package db 提供数据库连接封装，基于 GORM 实现。
package db

import (
	"context"
	"errors"
	"fmt"
	"time"
	"ult/pkg/logger"
	"ult/pkg/types"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var _ gorm_logger.Interface = (*Logger)(nil)

// Logger GORM 日志记录器，实现 gorm_logger.Interface 接口。
// 封装项目日志记录器，提供 SQL 执行日志记录功能。
type Logger struct {
	l *logger.Logger
	*loggerOption
}

// LoggerOption 日志选项函数类型。
type LoggerOption func(*loggerOption)

// loggerOption 日志配置选项。
type loggerOption struct {
	logLevel                  gorm_logger.LogLevel // 日志级别
	slowThreshold             time.Duration        // 慢 SQL 阈值
	ignoreRecordNotFoundError bool                 // 是否忽略记录未找到错误
}

// WithLoggerLevel 设置日志级别选项。
//
// 参数:
//   - level: GORM 日志级别
//
// 返回:
//   - LoggerOption: 日志选项函数
func WithLoggerLevel(level gorm_logger.LogLevel) LoggerOption {
	return func(option *loggerOption) {
		option.logLevel = level
	}
}

// WithLoggerSlowThreshold 设置慢 SQL 阈值选项。
//
// 参数:
//   - slowThreshold: 慢 SQL 时间阈值
//
// 返回:
//   - LoggerOption: 日志选项函数
func WithLoggerSlowThreshold(slowThreshold time.Duration) LoggerOption {
	return func(option *loggerOption) {
		option.slowThreshold = slowThreshold
	}
}

// WithLoggerIgnoreRecordNotFoundError 设置是否忽略 ErrRecordNotFound 错误选项。
//
// 参数:
//   - ignoreRecordNotFoundError: 是否忽略记录未找到错误
//
// 返回:
//   - LoggerOption: 日志选项函数
func WithLoggerIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) LoggerOption {
	return func(option *loggerOption) {
		option.ignoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

// NewLogger 创建新的 GORM 日志记录器。
//
// 参数:
//   - logger: 项目日志记录器实例
//   - opts: 日志选项列表
//
// 返回:
//   - *Logger: GORM 日志记录器实例
func NewLogger(logger *logger.Logger, opts ...LoggerOption) *Logger {
	var l = new(Logger)
	l.loggerOption = new(loggerOption)
	l.l = logger
	for _, opt := range opts {
		opt(l.loggerOption)
	}
	return l
}

// LogMode 设置日志模式并返回新的日志接口实例。
//
// 参数:
//   - level: 日志级别
//
// 返回:
//   - gorm_logger.Interface: 日志接口实例
func (l *Logger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	l.logLevel = level
	return l
}

// Info 记录 Info 级别日志。
//
// 参数:
//   - ctx: 上下文
//   - str: 日志消息
//   - args: 日志参数
func (l *Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.logLevel >= gorm_logger.Info {
		l.l.UseSQL(ctx).Sugar().Infof(str, args...)
	}
}

// Warn 记录 Warn 级别日志。
//
// 参数:
//   - ctx: 上下文
//   - str: 日志消息
//   - args: 日志参数
func (l *Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.logLevel >= gorm_logger.Warn {
		l.l.UseSQL(ctx).Sugar().Warnf(str, args...)
	}
}

// Error 记录 Error 级别日志。
//
// 参数:
//   - ctx: 上下文
//   - str: 日志消息
//   - args: 日志参数
func (l *Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.logLevel >= gorm_logger.Error {
		l.l.UseSQL(ctx).Sugar().Errorf(str, args...)
	}
}

// Trace 记录 SQL 执行追踪日志。
// 根据执行时间和错误情况记录不同级别的日志。
//
// 参数:
//   - ctx: 上下文
//   - begin: SQL 执行开始时间
//   - fc: 获取 SQL 和影响行数的函数
//   - err: SQL 执行错误
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= 0 {
		return
	}

	var (
		traceId string
		sql     string
		rows    int64
	)

	// 请求链路追踪 TraceID
	reqCtx, ok := types.FromRequestContext(ctx)
	if ok {
		traceId = reqCtx.TraceID()
	}

	elapsed := time.Since(begin)
	elapsedStr := zap.String("elapsed", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6))
	fileStr := zap.String("file", utils.FileWithLineNum())
	rowsStr := func(rows int64) zap.Field { return zap.Int64("rows", rows) }
	sqlStr := func(sql string) zap.Field { return zap.String("sql", sql) }
	traceIdStr := func(traceId string) zap.Field { return zap.String("trace_id", traceId) }
	switch {
	case err != nil && l.logLevel >= gorm_logger.Error && (!l.ignoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows = fc()
		l.l.UseSQL(ctx).Error("ERROR SQL", zap.Error(err), fileStr, elapsedStr, rowsStr(rows), sqlStr(sql), traceIdStr(traceId))
	case l.slowThreshold != 0 && elapsed > l.slowThreshold && l.logLevel >= gorm_logger.Warn:
		sql, rows = fc()
		l.l.UseSQL(ctx).Warn(fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold), fileStr, elapsedStr, rowsStr(rows), sqlStr(sql), traceIdStr(traceId))
	case l.logLevel >= gorm_logger.Info:
		sql, rows = fc()
		l.l.UseSQL(ctx).Info("INFO SQL", fileStr, elapsedStr, rowsStr(rows), sqlStr(sql), traceIdStr(traceId))
	}
}
