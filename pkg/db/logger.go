package db

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
	"ult/pkg/global"
	"ult/pkg/logger"
)

var _ gorm_logger.Interface = (*Logger)(nil)

type Logger struct {
	l *logger.Logger
	*loggerOption
}

type LoggerOption func(*loggerOption)

type loggerOption struct {
	logLevel                  gorm_logger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

// WithLoggerLevel 日志级别
func WithLoggerLevel(level gorm_logger.LogLevel) LoggerOption {
	return func(option *loggerOption) {
		option.logLevel = level
	}
}

// WithLoggerSlowThreshold 慢 SQL 阈值
func WithLoggerSlowThreshold(slowThreshold time.Duration) LoggerOption {
	return func(option *loggerOption) {
		option.slowThreshold = slowThreshold
	}
}

// WithLoggerIgnoreRecordNotFoundError 忽略 ErrRecordNotFound（记录未找到）错误
func WithLoggerIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) LoggerOption {
	return func(option *loggerOption) {
		option.ignoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

func NewLogger(logger *logger.Logger, opts ...LoggerOption) *Logger {
	var l = new(Logger)
	l.loggerOption = new(loggerOption)
	l.l = logger
	for _, opt := range opts {
		opt(l.loggerOption)
	}
	return l
}

func (l *Logger) LogMode(level gorm_logger.LogLevel) gorm_logger.Interface {
	l.logLevel = level
	return l
}

func (l *Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.logLevel >= gorm_logger.Info {
		l.l.UseSQL(ctx).Sugar().Infof(str, args...)
	}
}

func (l *Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.logLevel >= gorm_logger.Warn {
		l.l.UseSQL(ctx).Sugar().Warnf(str, args...)
	}
}

func (l *Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.logLevel >= gorm_logger.Error {
		l.l.UseSQL(ctx).Sugar().Errorf(str, args...)
	}
}

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
	reqCtx, ok := global.FromRequestContext(ctx)
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
