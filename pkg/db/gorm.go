// Package db 提供数据库连接封装，基于 GORM 实现。
// 支持多种数据库驱动，提供连接池管理、日志记录和插件扩展功能。
package db

import (
	"context"
	"fmt"
	"time"
	"ult/config/autoload"
	"ult/pkg/logger"

	"github.com/raylin666/go-utils/v2/db/gorm"
	gorm_db "gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

var _ Db = (*db)(nil)

// Predicate 查询条件谓词类型。
type Predicate string

// Db 数据库接口，定义数据库连接的基本操作。
type Db interface {
	Get() gorm.Client // 获取 GORM 客户端
	Close() error     // 关闭数据库连接
	Ping() error      // 测试数据库连接
}

// db 数据库实例，封装 GORM 客户端和日志记录器。
type db struct {
	client gorm.Client
	logger *logger.Logger
}

// Ping 测试数据库连接是否正常。
//
// 返回:
//   - error: 连接测试错误，nil 表示连接正常
func (db *db) Ping() error {
	return db.Get().SqlDB().Ping()
}

// NewDb 创建新的数据库连接实例。
// 支持连接重试机制，可配置最大重试次数和重试间隔。
//
// 参数:
//   - name: 数据库连接名称，用于日志标识
//   - config: 数据库配置（DSN、驱动、连接池等）
//   - logger: 日志记录器实例
//
// 返回:
//   - Db: 数据库实例
//   - error: 连接创建错误
func NewDb(name string, config autoload.DB, logger *logger.Logger) (Db, error) {
	var rdb = new(db)
	rdb.logger = logger

	maxRetries := config.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	retryDelay := config.RetryDelay
	if retryDelay <= 0 {
		retryDelay = 1
	}

	var client gorm.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		client, err = gorm.NewClient(
			gorm.WithDsn(config.Dsn),
			gorm.WithDriver(config.Driver),
			gorm.WithDbName(config.DbName),
			gorm.WithHost(config.Host),
			gorm.WithUserName(config.UserName),
			gorm.WithPassword(config.Password),
			gorm.WithCharset(config.Charset),
			gorm.WithPort(config.Port),
			gorm.WithPrefix(config.Prefix),
			gorm.WithMaxIdleConn(config.MaxIdleConn),
			gorm.WithMaxOpenConn(config.MaxOpenConn),
			gorm.WithMaxLifeTime(time.Duration(config.MaxLifeTime)),
			gorm.WithParseTime(config.ParseTime),
			gorm.WithLoc(config.Loc))
		if err == nil {
			break
		}

		logger.UseApp(context.TODO()).Warn(fmt.Sprintf("数据库连接 %s 尝试 %d/%d 失败：%v，在 %d 秒内重试", name, i+1, maxRetries, err, retryDelay))
		time.Sleep(time.Duration(retryDelay) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("创建数据库连接失败：%v，%d次重试后仍失败", err, maxRetries)
	}

	rdb.client = client

	// 日志处理
	l := rdb.WithLogger()
	if l != nil {
		rdb.client.WithLogger(l)
	}
	// 插件处理
	_ = rdb.client.WithPluginBeforeHandler(rdb.BeforePluginHandler, rdb.AfterPluginHandler)

	return rdb, nil
}

// Get 获取 GORM 客户端实例。
//
// 返回:
//   - gorm.Client: GORM 客户端
func (db *db) Get() gorm.Client {
	return db.client
}

// Close 关闭数据库连接。
//
// 返回:
//   - error: 关闭连接时的错误
func (db *db) Close() error {
	return db.Get().SqlDB().Close()
}

// WithLogger 配置 GORM 日志记录器。
//
// 返回:
//   - gorm_logger.Interface: GORM 日志接口实例
func (db *db) WithLogger() gorm_logger.Interface {
	return NewLogger(
		db.logger,
		WithLoggerLevel(gorm_logger.Info),
		WithLoggerSlowThreshold(time.Second*1),
		WithLoggerIgnoreRecordNotFoundError(true))
}

// BeforePluginHandler DB 插件前置处理方法，在 SQL 执行前调用。
//
// 参数:
//   - rdb: GORM DB 实例
func (db *db) BeforePluginHandler(rdb *gorm_db.DB) {}

// AfterPluginHandler DB 插件后置处理方法，在 SQL 执行后调用。
//
// 参数:
//   - rdb: GORM DB 实例
//   - sql: 执行的 SQL 语句
//   - ts: SQL 执行开始时间
func (db *db) AfterPluginHandler(rdb *gorm_db.DB, sql string, ts time.Time) {}
